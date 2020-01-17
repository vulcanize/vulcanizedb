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
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/config"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

const (
	PayloadChanBufferSize = 20000
)

// SuperNode is the top level interface for streaming, converting to IPLDs, publishing,
// and indexing all Ethereum data; screening this data; and serving it up to subscribed clients
// This service is compatible with the Ethereum service interface (node.Service)
type SuperNode interface {
	// APIs(), Protocols(), Start() and Stop()
	node.Service
	// Main event loop for syncAndPublish processes
	SyncAndPublish(wg *sync.WaitGroup, forwardPayloadChan chan<- interface{}, forwardQuitchan chan<- bool) error
	// Main event loop for handling client pub-sub
	ScreenAndServe(wg *sync.WaitGroup, screenAndServePayload <-chan interface{}, screenAndServeQuit <-chan bool)
	// Method to subscribe to receive state diff processing output
	Subscribe(id rpc.ID, sub chan<- Payload, quitChan chan<- bool, params SubscriptionSettings)
	// Method to unsubscribe from state diff processing
	Unsubscribe(id rpc.ID)
	// Method to access the node info for this service
	Node() core.Node
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
	// Interface for resolving IPLDs to their data types
	Resolver shared.IPLDResolver
	// Chan the processor uses to subscribe to payloads from the Streamer
	PayloadChan chan interface{}
	// Used to signal shutdown of the service
	QuitChan chan bool
	// A mapping of rpc.IDs to their subscription channels, mapped to their subscription type (hash of the StreamFilters)
	Subscriptions map[common.Hash]map[rpc.ID]Subscription
	// A mapping of subscription params hash to the corresponding subscription params
	SubscriptionTypes map[common.Hash]SubscriptionSettings
	// Info for the Geth node that this super node is working with
	NodeInfo core.Node
	// Number of publishAndIndex workers
	WorkerPoolSize int
	// chain type for this service
	chain config.ChainType
	// Path to ipfs data dir
	ipfsPath string
	// Underlying db
	db *postgres.DB
}

// NewSuperNode creates a new super_node.Interface using an underlying super_node.Service struct
func NewSuperNode(settings *config.SuperNode) (SuperNode, error) {
	if err := ipfs.InitIPFSPlugins(); err != nil {
		return nil, err
	}
	sn := new(Service)
	var err error
	// If we are syncing, initialize the needed interfaces
	if settings.Sync {
		sn.Streamer, sn.PayloadChan, err = NewPayloadStreamer(settings.Chain, settings.WSClient)
		if err != nil {
			return nil, err
		}
		sn.Converter, err = NewPayloadConverter(settings.Chain, params.MainnetChainConfig)
		if err != nil {
			return nil, err
		}
		sn.Publisher, err = NewIPLDPublisher(settings.Chain, settings.IPFSPath)
		if err != nil {
			return nil, err
		}
		sn.Indexer, err = NewCIDIndexer(settings.Chain, settings.DB)
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
		sn.Retriever, err = NewCIDRetriever(settings.Chain, settings.DB)
		if err != nil {
			return nil, err
		}
		sn.IPLDFetcher, err = NewIPLDFetcher(settings.Chain, settings.IPFSPath)
		if err != nil {
			return nil, err
		}
		sn.Resolver, err = NewIPLDResolver(settings.Chain)
		if err != nil {
			return nil, err
		}
	}
	sn.QuitChan = settings.Quit
	sn.Subscriptions = make(map[common.Hash]map[rpc.ID]Subscription)
	sn.SubscriptionTypes = make(map[common.Hash]SubscriptionSettings)
	sn.WorkerPoolSize = settings.Workers
	sn.NodeInfo = settings.NodeInfo
	sn.ipfsPath = settings.IPFSPath
	sn.chain = settings.Chain
	sn.db = settings.DB
	return sn, nil
}

// Protocols exports the services p2p protocols, this service has none
func (sap *Service) Protocols() []p2p.Protocol {
	return []p2p.Protocol{}
}

// APIs returns the RPC descriptors the super node service offers
func (sap *Service) APIs() []rpc.API {
	apis := []rpc.API{
		{
			Namespace: APIName,
			Version:   APIVersion,
			Service:   NewPublicSuperNodeAPI(sap),
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

// SyncAndPublish is the backend processing loop which streams data from geth, converts it to iplds, publishes them to ipfs, and indexes their cids
// This continues on no matter if or how many subscribers there are, it then forwards the data to the ScreenAndServe() loop
// which filters and sends relevant data to client subscriptions, if there are any
func (sap *Service) SyncAndPublish(wg *sync.WaitGroup, screenAndServePayload chan<- interface{}, screenAndServeQuit chan<- bool) error {
	sub, err := sap.Streamer.Stream(sap.PayloadChan)
	if err != nil {
		return err
	}
	wg.Add(1)

	// Channels for forwarding data to the publishAndIndex workers
	publishAndIndexPayload := make(chan interface{}, PayloadChanBufferSize)
	publishAndIndexQuit := make(chan bool, sap.WorkerPoolSize)
	// publishAndIndex worker pool to handle publishing and indexing concurrently, while
	// limiting the number of Postgres connections we can possibly open so as to prevent error
	for i := 0; i < sap.WorkerPoolSize; i++ {
		sap.publishAndIndex(i, publishAndIndexPayload, publishAndIndexQuit)
	}
	go func() {
		for {
			select {
			case payload := <-sap.PayloadChan:
				ipldPayload, err := sap.Converter.Convert(payload)
				if err != nil {
					log.Error(err)
					continue
				}
				// If we have a ScreenAndServe process running, forward the payload to it
				select {
				case screenAndServePayload <- ipldPayload:
				default:
				}
				// Forward the payload to the publishAndIndex workers
				select {
				case publishAndIndexPayload <- ipldPayload:
				default:
				}
			case err := <-sub.Err():
				log.Error(err)
			case <-sap.QuitChan:
				// If we have a ScreenAndServe process running, forward the quit signal to it
				select {
				case screenAndServeQuit <- true:
				default:
				}
				// Also forward a quit signal for each of the publishAndIndex workers
				for i := 0; i < sap.WorkerPoolSize; i++ {
					select {
					case publishAndIndexQuit <- true:
					default:
					}
				}
				log.Info("quiting SyncAndPublish process")
				wg.Done()
				return
			}
		}
	}()
	log.Info("syncAndPublish goroutine successfully spun up")
	return nil
}

func (sap *Service) publishAndIndex(id int, publishAndIndexPayload <-chan interface{}, publishAndIndexQuit <-chan bool) {
	go func() {
		for {
			select {
			case payload := <-publishAndIndexPayload:
				cidPayload, err := sap.Publisher.Publish(payload)
				if err != nil {
					log.Errorf("worker %d error: %v", id, err)
					continue
				}
				if err := sap.Indexer.Index(cidPayload); err != nil {
					log.Errorf("worker %d error: %v", id, err)
				}
			case <-publishAndIndexQuit:
				log.Infof("quiting publishAndIndex worker %d", id)
				return
			}
		}
	}()
	log.Info("publishAndIndex goroutine successfully spun up")
}

// ScreenAndServe is the loop used to screen data streamed from the state diffing eth node
// and send the appropriate portions of it to a requesting client subscription, according to their subscription configuration
func (sap *Service) ScreenAndServe(wg *sync.WaitGroup, screenAndServePayload <-chan interface{}, screenAndServeQuit <-chan bool) {
	wg.Add(1)
	go func() {
		for {
			select {
			case payload := <-screenAndServePayload:
				sap.sendResponse(payload)
			case <-screenAndServeQuit:
				log.Info("quiting ScreenAndServe process")
				wg.Done()
				return
			}
		}
	}()
	log.Info("screenAndServe goroutine successfully spun up")
}

func (sap *Service) sendResponse(payload interface{}) {
	sap.Lock()
	for ty, subs := range sap.Subscriptions {
		// Retrieve the subscription parameters for this subscription type
		subConfig, ok := sap.SubscriptionTypes[ty]
		if !ok {
			log.Errorf("subscription configuration for subscription type %s not available", ty.Hex())
			sap.closeType(ty)
			continue
		}
		response, err := sap.Filterer.Filter(subConfig, payload)
		if err != nil {
			log.Error(err)
			sap.closeType(ty)
			continue
		}
		for id, sub := range subs {
			select {
			case sub.PayloadChan <- Payload{response, ""}:
				log.Infof("sending super node payload to subscription %s", id)
			default:
				log.Infof("unable to send payload to subscription %s; channel has no receiver", id)
			}
		}
	}
	sap.Unlock()
}

// Subscribe is used by the API to subscribe to the service loop
// The params must be rlp serializable and satisfy the Params() interface
func (sap *Service) Subscribe(id rpc.ID, sub chan<- Payload, quitChan chan<- bool, params SubscriptionSettings) {
	log.Info("Subscribing to the super node service")
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
	// Subscription type is defined as the hash of the subscription settings
	by, err := rlp.EncodeToBytes(params)
	if err != nil {
		sendNonBlockingErr(subscription, err)
		sendNonBlockingQuit(subscription)
		return
	}
	subscriptionType := crypto.Keccak256Hash(by)
	// If the subscription requests a backfill, use the Postgres index to lookup and retrieve historical data
	// Otherwise we only filter new data as it is streamed in from the state diffing geth node
	if params.HistoricalData() || params.HistoricalDataOnly() {
		if err := sap.backFill(subscription, id, params); err != nil {
			sendNonBlockingErr(subscription, err)
			sendNonBlockingQuit(subscription)
			return
		}
	}
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
}

func (sap *Service) backFill(sub Subscription, id rpc.ID, params SubscriptionSettings) error {
	log.Debug("sending historical data for subscriber", id)
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
	log.Debug("historical data starting block:", params.StartingBlock())
	log.Debug("histocial data ending block:", endingBlock)
	go func() {
		for i := startingBlock; i <= endingBlock; i++ {
			cidWrapper, empty, err := sap.Retriever.Retrieve(params, i)
			if err != nil {
				sendNonBlockingErr(sub, fmt.Errorf("CID Retrieval error at block %d\r%s", i, err.Error()))
				continue
			}
			if empty {
				continue
			}
			blocksWrapper, err := sap.IPLDFetcher.Fetch(cidWrapper)
			if err != nil {
				sendNonBlockingErr(sub, fmt.Errorf("IPLD Fetching error at block %d\r%s", i, err.Error()))
				continue
			}
			backFillIplds, err := sap.Resolver.Resolve(blocksWrapper)
			if err != nil {
				sendNonBlockingErr(sub, fmt.Errorf("IPLD Resolving error at block %d\r%s", i, err.Error()))
				continue
			}
			select {
			case sub.PayloadChan <- Payload{backFillIplds, ""}:
				log.Infof("sending super node historical data payload to subscription %s", id)
			default:
				log.Infof("unable to send back-fill payload to subscription %s; channel has no receiver", id)
			}
		}
	}()
	return nil
}

// Unsubscribe is used to unsubscribe to the StateDiffingService loop
func (sap *Service) Unsubscribe(id rpc.ID) {
	log.Info("Unsubscribing from the super node service")
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
func (sap *Service) Start(*p2p.Server) error {
	log.Info("Starting super node service")
	wg := new(sync.WaitGroup)
	payloadChan := make(chan interface{}, PayloadChanBufferSize)
	quitChan := make(chan bool, 1)
	if err := sap.SyncAndPublish(wg, payloadChan, quitChan); err != nil {
		return err
	}
	sap.ScreenAndServe(wg, payloadChan, quitChan)
	return nil
}

// Stop is used to close down the service
func (sap *Service) Stop() error {
	log.Info("Stopping super node service")
	sap.Lock()
	close(sap.QuitChan)
	sap.close()
	sap.Unlock()
	return nil
}

// Node returns the node info for this service
func (sap *Service) Node() core.Node {
	return sap.NodeInfo
}

// close is used to close all listening subscriptions
// close needs to be called with subscription access locked
func (sap *Service) close() {
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
	subs := sap.Subscriptions[subType]
	for _, sub := range subs {
		sendNonBlockingQuit(sub)
	}
	delete(sap.Subscriptions, subType)
	delete(sap.SubscriptionTypes, subType)
}
