// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package super_node

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/eth/core"
	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

const (
	PayloadChanBufferSize = 2000
)

// SuperNode is the top level interface for streaming, converting to IPLDs, publishing,
// and indexing all chain data; screening this data; and serving it up to subscribed clients
// This service is compatible with the Ethereum service interface (node.Service)
type SuperNode interface {
	// APIs(), Protocols(), Start() and Stop()
	node.Service
	// Data processing event loop
	Sync(wg *sync.WaitGroup, forwardPayloadChan chan<- shared.ConvertedData) error
	// Pub-Sub handling event loop
	Serve(wg *sync.WaitGroup, screenAndServePayload <-chan shared.ConvertedData)
	// Method to subscribe to the service
	Subscribe(id rpc.ID, sub chan<- SubscriptionPayload, quitChan chan<- bool, params shared.SubscriptionSettings)
	// Method to unsubscribe from the service
	Unsubscribe(id rpc.ID)
	// Method to access the node info for the service
	Node() *core.Node
	// Method to access chain type
	Chain() shared.ChainType
}

// Service is the underlying struct for the super node
type Service struct {
	// Used to sync access to the Subscriptions
	sync.Mutex
	// Interface for streaming payloads over an rpc subscription
	Streamer shared.PayloadStreamer
	// Interface for converting raw payloads into IPLD object payloads
	Converter shared.PayloadConverter
	// Interface for publishing the IPLD payloads to IPFS
	Publisher shared.IPLDPublisher
	// Interface for indexing the CIDs of the published IPLDs in Postgres
	Indexer shared.CIDIndexer
	// Interface for filtering and serving data according to subscribed clients according to their specification
	Filterer shared.ResponseFilterer
	// Interface for fetching IPLD objects from IPFS
	IPLDFetcher shared.IPLDFetcher
	// Interface for searching and retrieving CIDs from Postgres index
	Retriever shared.CIDRetriever
	// Chan the processor uses to subscribe to payloads from the Streamer
	PayloadChan chan shared.RawChainData
	// Used to signal shutdown of the service
	QuitChan chan bool
	// A mapping of rpc.IDs to their subscription channels, mapped to their subscription type (hash of the StreamFilters)
	Subscriptions map[common.Hash]map[rpc.ID]Subscription
	// A mapping of subscription params hash to the corresponding subscription params
	SubscriptionTypes map[common.Hash]shared.SubscriptionSettings
	// Info for the Geth node that this super node is working with
	NodeInfo *core.Node
	// Number of publishAndIndex workers
	WorkerPoolSize int
	// chain type for this service
	chain shared.ChainType
	// Path to ipfs data dir
	ipfsPath string
	// Underlying db
	db *postgres.DB
}

// NewSuperNode creates a new super_node.Interface using an underlying super_node.Service struct
func NewSuperNode(settings *Config) (SuperNode, error) {
	sn := new(Service)
	var err error
	// If we are syncing, initialize the needed interfaces
	if settings.Sync {
		sn.Streamer, sn.PayloadChan, err = NewPayloadStreamer(settings.Chain, settings.WSClient)
		if err != nil {
			return nil, err
		}
		sn.Converter, err = NewPayloadConverter(settings.Chain)
		if err != nil {
			return nil, err
		}
		sn.Publisher, err = NewIPLDPublisher(settings.Chain, settings.IPFSPath, settings.SyncDBConn, settings.IPFSMode)
		if err != nil {
			return nil, err
		}
		sn.Indexer, err = NewCIDIndexer(settings.Chain, settings.SyncDBConn, settings.IPFSMode)
		if err != nil {
			return nil, err
		}
		sn.Filterer, err = NewResponseFilterer(settings.Chain)
		if err != nil {
			return nil, err
		}
	}
	// If we are serving, initialize the needed interfaces
	if settings.Serve {
		sn.Retriever, err = NewCIDRetriever(settings.Chain, settings.ServeDBConn)
		if err != nil {
			return nil, err
		}
		sn.IPLDFetcher, err = NewIPLDFetcher(settings.Chain, settings.IPFSPath, settings.ServeDBConn, settings.IPFSMode)
		if err != nil {
			return nil, err
		}
		sn.db = settings.ServeDBConn
	}
	sn.QuitChan = settings.Quit
	sn.Subscriptions = make(map[common.Hash]map[rpc.ID]Subscription)
	sn.SubscriptionTypes = make(map[common.Hash]shared.SubscriptionSettings)
	sn.WorkerPoolSize = settings.Workers
	sn.NodeInfo = &settings.NodeInfo
	sn.ipfsPath = settings.IPFSPath
	sn.chain = settings.Chain
	return sn, nil
}

// Protocols exports the services p2p protocols, this service has none
func (sap *Service) Protocols() []p2p.Protocol {
	return []p2p.Protocol{}
}

// APIs returns the RPC descriptors the super node service offers
func (sap *Service) APIs() []rpc.API {
	ifnoAPI := NewInfoAPI()
	apis := []rpc.API{
		{
			Namespace: APIName,
			Version:   APIVersion,
			Service:   NewPublicSuperNodeAPI(sap),
			Public:    true,
		},
		{
			Namespace: "rpc",
			Version:   APIVersion,
			Service:   ifnoAPI,
			Public:    true,
		},
		{
			Namespace: "net",
			Version:   APIVersion,
			Service:   ifnoAPI,
			Public:    true,
		},
		{
			Namespace: "admin",
			Version:   APIVersion,
			Service:   ifnoAPI,
			Public:    true,
		},
	}
	chainAPI, err := NewPublicAPI(sap.chain, sap.db, sap.ipfsPath)
	if err != nil {
		log.Error(err)
		return apis
	}
	return append(apis, chainAPI)
}

// Sync streams incoming raw chain data and converts it for further processing
// It forwards the converted data to the publishAndIndex process(es) it spins up
// If forwards the converted data to a ScreenAndServe process if it there is one listening on the passed screenAndServePayload channel
// This continues on no matter if or how many subscribers there are
func (sap *Service) Sync(wg *sync.WaitGroup, screenAndServePayload chan<- shared.ConvertedData) error {
	sub, err := sap.Streamer.Stream(sap.PayloadChan)
	if err != nil {
		return err
	}
	wg.Add(1)

	// Channels for forwarding data to the publishAndIndex workers
	publishAndIndexPayload := make(chan shared.ConvertedData, PayloadChanBufferSize)
	// publishAndIndex worker pool to handle publishing and indexing concurrently, while
	// limiting the number of Postgres connections we can possibly open so as to prevent error
	for i := 0; i < sap.WorkerPoolSize; i++ {
		sap.publishAndIndex(i, publishAndIndexPayload)
	}
	go func() {
		for {
			select {
			case payload := <-sap.PayloadChan:
				ipldPayload, err := sap.Converter.Convert(payload)
				if err != nil {
					log.Errorf("super node conversion error for chain %s: %v", sap.chain.String(), err)
					continue
				}
				log.Infof("%s data streamed at head height %d", sap.chain.String(), ipldPayload.Height())
				// If we have a ScreenAndServe process running, forward the iplds to it
				select {
				case screenAndServePayload <- ipldPayload:
				default:
				}
				// Forward the payload to the publishAndIndex workers
				// this channel acts as a ring buffer
				select {
				case publishAndIndexPayload <- ipldPayload:
				default:
					<-publishAndIndexPayload
					publishAndIndexPayload <- ipldPayload
				}
			case err := <-sub.Err():
				log.Errorf("super node subscription error for chain %s: %v", sap.chain.String(), err)
			case <-sap.QuitChan:
				log.Infof("quiting %s SyncAndPublish process", sap.chain.String())
				wg.Done()
				return
			}
		}
	}()
	log.Infof("%s Sync goroutine successfully spun up", sap.chain.String())
	return nil
}

// publishAndIndex is spun up by SyncAndConvert and receives converted chain data from that process
// it publishes this data to IPFS and indexes their CIDs with useful metadata in Postgres
func (sap *Service) publishAndIndex(id int, publishAndIndexPayload <-chan shared.ConvertedData) {
	go func() {
		for {
			select {
			case payload := <-publishAndIndexPayload:
				log.Debugf("publishing %s data streamed at head height %d", sap.chain.String(), payload.Height())
				cidPayload, err := sap.Publisher.Publish(payload)
				if err != nil {
					log.Errorf("super node publishAndIndex worker %d publishing error for chain %s: %v", id, sap.chain.String(), err)
					continue
				}
				log.Debugf("indexing %s data streamed at head height %d", sap.chain.String(), payload.Height())
				if err := sap.Indexer.Index(cidPayload); err != nil {
					log.Errorf("super node publishAndIndex worker %d indexing error for chain %s: %v", id, sap.chain.String(), err)
				}
			}
		}
	}()
	log.Debugf("%s publishAndIndex goroutine successfully spun up", sap.chain.String())
}

// Serve listens for incoming converter data off the screenAndServePayload from the SyncAndConvert process
// It filters and sends this data to any subscribers to the service
// This process can be stood up alone, without an screenAndServePayload attached to a SyncAndConvert process
// and it will hang on the WaitGroup indefinitely, allowing the Service to serve historical data requests only
func (sap *Service) Serve(wg *sync.WaitGroup, screenAndServePayload <-chan shared.ConvertedData) {
	wg.Add(1)
	go func() {
		for {
			select {
			case payload := <-screenAndServePayload:
				sap.filterAndServe(payload)
			case <-sap.QuitChan:
				log.Infof("quiting %s ScreenAndServe process", sap.chain.String())
				wg.Done()
				return
			}
		}
	}()
	log.Infof("%s Serve goroutine successfully spun up", sap.chain.String())
}

// filterAndServe filters the payload according to each subscription type and sends to the subscriptions
func (sap *Service) filterAndServe(payload shared.ConvertedData) {
	log.Debugf("Sending %s payload to subscriptions", sap.chain.String())
	sap.Lock()
	for ty, subs := range sap.Subscriptions {
		// Retrieve the subscription parameters for this subscription type
		subConfig, ok := sap.SubscriptionTypes[ty]
		if !ok {
			log.Errorf("super node %s subscription configuration for subscription type %s not available", sap.chain.String(), ty.Hex())
			sap.closeType(ty)
			continue
		}
		if subConfig.EndingBlock().Int64() > 0 && subConfig.EndingBlock().Int64() < payload.Height() {
			// We are not out of range for this subscription type
			// close it, and continue to the next
			sap.closeType(ty)
			continue
		}
		response, err := sap.Filterer.Filter(subConfig, payload)
		if err != nil {
			log.Errorf("super node filtering error for chain %s: %v", sap.chain.String(), err)
			sap.closeType(ty)
			continue
		}
		responseRLP, err := rlp.EncodeToBytes(response)
		if err != nil {
			log.Errorf("super node rlp encoding error for chain %s: %v", sap.chain.String(), err)
			continue
		}
		for id, sub := range subs {
			select {
			case sub.PayloadChan <- SubscriptionPayload{Data: responseRLP, Err: "", Flag: EmptyFlag, Height: response.Height()}:
				log.Debugf("sending super node %s payload to subscription %s", sap.chain.String(), id)
			default:
				log.Infof("unable to send %s payload to subscription %s; channel has no receiver", sap.chain.String(), id)
			}
		}
	}
	sap.Unlock()
}

// Subscribe is used by the API to remotely subscribe to the service loop
// The params must be rlp serializable and satisfy the SubscriptionSettings() interface
func (sap *Service) Subscribe(id rpc.ID, sub chan<- SubscriptionPayload, quitChan chan<- bool, params shared.SubscriptionSettings) {
	log.Infof("New %s subscription %s", sap.chain.String(), id)
	subscription := Subscription{
		ID:          id,
		PayloadChan: sub,
		QuitChan:    quitChan,
	}
	if params.ChainType() != sap.chain {
		sendNonBlockingErr(subscription, fmt.Errorf("subscription %s is for chain %s, service supports chain %s", id, params.ChainType().String(), sap.chain.String()))
		sendNonBlockingQuit(subscription)
		return
	}
	// Subscription type is defined as the hash of the rlp-serialized subscription settings
	by, err := rlp.EncodeToBytes(params)
	if err != nil {
		sendNonBlockingErr(subscription, err)
		sendNonBlockingQuit(subscription)
		return
	}
	subscriptionType := crypto.Keccak256Hash(by)
	if !params.HistoricalDataOnly() {
		// Add subscriber
		sap.Lock()
		if sap.Subscriptions[subscriptionType] == nil {
			sap.Subscriptions[subscriptionType] = make(map[rpc.ID]Subscription)
		}
		sap.Subscriptions[subscriptionType][id] = subscription
		sap.SubscriptionTypes[subscriptionType] = params
		sap.Unlock()
	}
	// If the subscription requests a backfill, use the Postgres index to lookup and retrieve historical data
	// Otherwise we only filter new data as it is streamed in from the state diffing geth node
	if params.HistoricalData() || params.HistoricalDataOnly() {
		if err := sap.sendHistoricalData(subscription, id, params); err != nil {
			sendNonBlockingErr(subscription, fmt.Errorf("super node subscriber backfill error for chain %s: %v", sap.chain.String(), err))
			sendNonBlockingQuit(subscription)
			return
		}
	}
}

// sendHistoricalData sends historical data to the requesting subscription
func (sap *Service) sendHistoricalData(sub Subscription, id rpc.ID, params shared.SubscriptionSettings) error {
	log.Infof("Sending %s historical data to subscription %s", sap.chain.String(), id)
	// Retrieve cached CIDs relevant to this subscriber
	var endingBlock int64
	var startingBlock int64
	var err error
	startingBlock, err = sap.Retriever.RetrieveFirstBlockNumber()
	if err != nil {
		return err
	}
	if startingBlock < params.StartingBlock().Int64() {
		startingBlock = params.StartingBlock().Int64()
	}
	endingBlock, err = sap.Retriever.RetrieveLastBlockNumber()
	if err != nil {
		return err
	}
	if endingBlock > params.EndingBlock().Int64() && params.EndingBlock().Int64() > 0 && params.EndingBlock().Int64() > startingBlock {
		endingBlock = params.EndingBlock().Int64()
	}
	log.Debugf("%s historical data starting block: %d", sap.chain.String(), params.StartingBlock().Int64())
	log.Debugf("%s historical data ending block: %d", sap.chain.String(), endingBlock)
	go func() {
		for i := startingBlock; i <= endingBlock; i++ {
			cidWrappers, empty, err := sap.Retriever.Retrieve(params, i)
			if err != nil {
				sendNonBlockingErr(sub, fmt.Errorf("super node %s CID Retrieval error at block %d\r%s", sap.chain.String(), i, err.Error()))
				continue
			}
			if empty {
				continue
			}
			for _, cids := range cidWrappers {
				response, err := sap.IPLDFetcher.Fetch(cids)
				if err != nil {
					sendNonBlockingErr(sub, fmt.Errorf("super node %s IPLD Fetching error at block %d\r%s", sap.chain.String(), i, err.Error()))
					continue
				}
				responseRLP, err := rlp.EncodeToBytes(response)
				if err != nil {
					log.Error(err)
					continue
				}
				select {
				case sub.PayloadChan <- SubscriptionPayload{Data: responseRLP, Err: "", Flag: EmptyFlag, Height: response.Height()}:
					log.Debugf("sending super node historical data payload to %s subscription %s", sap.chain.String(), id)
				default:
					log.Infof("unable to send back-fill payload to %s subscription %s; channel has no receiver", sap.chain.String(), id)
				}
			}
		}
		// when we are done backfilling send an empty payload signifying so in the msg
		select {
		case sub.PayloadChan <- SubscriptionPayload{Data: nil, Err: "", Flag: BackFillCompleteFlag}:
			log.Debugf("sending backfill completion notice to %s subscription %s", sap.chain.String(), id)
		default:
			log.Infof("unable to send backfill completion notice to %s subscription %s", sap.chain.String(), id)
		}
	}()
	return nil
}

// Unsubscribe is used by the API to remotely unsubscribe to the StateDiffingService loop
func (sap *Service) Unsubscribe(id rpc.ID) {
	log.Infof("Unsubscribing %s from the %s super node service", id, sap.chain.String())
	sap.Lock()
	for ty := range sap.Subscriptions {
		delete(sap.Subscriptions[ty], id)
		if len(sap.Subscriptions[ty]) == 0 {
			// If we removed the last subscription of this type, remove the subscription type outright
			delete(sap.Subscriptions, ty)
			delete(sap.SubscriptionTypes, ty)
		}
	}
	sap.Unlock()
}

// Start is used to begin the service
// This is mostly just to satisfy the node.Service interface
func (sap *Service) Start(*p2p.Server) error {
	log.Infof("Starting %s super node service", sap.chain.String())
	wg := new(sync.WaitGroup)
	payloadChan := make(chan shared.ConvertedData, PayloadChanBufferSize)
	if err := sap.Sync(wg, payloadChan); err != nil {
		return err
	}
	sap.Serve(wg, payloadChan)
	return nil
}

// Stop is used to close down the service
// This is mostly just to satisfy the node.Service interface
func (sap *Service) Stop() error {
	log.Infof("Stopping %s super node service", sap.chain.String())
	sap.Lock()
	close(sap.QuitChan)
	sap.close()
	sap.Unlock()
	return nil
}

// Node returns the node info for this service
func (sap *Service) Node() *core.Node {
	return sap.NodeInfo
}

// Chain returns the chain type for this service
func (sap *Service) Chain() shared.ChainType {
	return sap.chain
}

// close is used to close all listening subscriptions
// close needs to be called with subscription access locked
func (sap *Service) close() {
	log.Infof("Closing all %s subscriptions", sap.chain.String())
	for subType, subs := range sap.Subscriptions {
		for _, sub := range subs {
			sendNonBlockingQuit(sub)
		}
		delete(sap.Subscriptions, subType)
		delete(sap.SubscriptionTypes, subType)
	}
}

// closeType is used to close all subscriptions of given type
// closeType needs to be called with subscription access locked
func (sap *Service) closeType(subType common.Hash) {
	log.Infof("Closing all %s subscriptions of type %s", sap.chain.String(), subType.String())
	subs := sap.Subscriptions[subType]
	for _, sub := range subs {
		sendNonBlockingQuit(sub)
	}
	delete(sap.Subscriptions, subType)
	delete(sap.SubscriptionTypes, subType)
}
