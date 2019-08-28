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

package seed_node

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/statediff"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/streamer"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
)

const payloadChanBufferSize = 20000 // the max eth sub buffer size

// NodeInterface is the top level interface for streaming, converting to IPLDs, publishing,
// and indexing all Ethereum data; screening this data; and serving it up to subscribed clients
// This service is compatible with the Ethereum service interface (node.Service)
type NodeInterface interface {
	// APIs(), Protocols(), Start() and Stop()
	node.Service
	// Main event loop for syncAndPublish processes
	SyncAndPublish(wg *sync.WaitGroup, forwardPayloadChan chan<- ipfs.IPLDPayload, forwardQuitchan chan<- bool) error
	// Main event loop for handling client pub-sub
	ScreenAndServe(screenAndServePayload <-chan ipfs.IPLDPayload, screenAndServeQuit <-chan bool)
	// Method to subscribe to receive state diff processing output
	Subscribe(id rpc.ID, sub chan<- streamer.SeedNodePayload, quitChan chan<- bool, streamFilters config.Subscription)
	// Method to unsubscribe from state diff processing
	Unsubscribe(id rpc.ID)
}

// Service is the underlying struct for the SyncAndPublish interface
type Service struct {
	// Used to sync access to the Subscriptions
	sync.Mutex
	// Interface for streaming statediff payloads over a geth rpc subscription
	Streamer streamer.IStateDiffStreamer
	// Interface for converting statediff payloads into ETH-IPLD object payloads
	Converter ipfs.PayloadConverter
	// Interface for publishing the ETH-IPLD payloads to IPFS
	Publisher ipfs.IPLDPublisher
	// Interface for indexing the CIDs of the published ETH-IPLDs in Postgres
	Repository CIDRepository
	// Interface for filtering and serving data according to subscribed clients according to their specification
	Filterer ResponseFilterer
	// Interface for fetching ETH-IPLD objects from IPFS
	Fetcher ipfs.IPLDFetcher
	// Interface for searching and retrieving CIDs from Postgres index
	Retriever CIDRetriever
	// Interface for resolving ipfs blocks to their data types
	Resolver ipfs.IPLDResolver
	// Chan the processor uses to subscribe to state diff payloads from the Streamer
	PayloadChan chan statediff.Payload
	// Used to signal shutdown of the service
	QuitChan chan bool
	// A mapping of rpc.IDs to their subscription channels, mapped to their subscription type (hash of the StreamFilters)
	Subscriptions map[common.Hash]map[rpc.ID]Subscription
	// A mapping of subscription hash type to the corresponding StreamFilters
	SubscriptionTypes map[common.Hash]config.Subscription
	// Number of workers
	WorkerPoolSize int
}

// NewSeedNode creates a new seed_node.Interface using an underlying seed_node.Service struct
func NewSeedNode(ipfsPath string, db *postgres.DB, rpcClient core.RpcClient, qc chan bool, workers int) (NodeInterface, error) {
	publisher, err := ipfs.NewIPLDPublisher(ipfsPath)
	if err != nil {
		return nil, err
	}
	fetcher, err := ipfs.NewIPLDFetcher(ipfsPath)
	if err != nil {
		return nil, err
	}
	return &Service{
		Streamer:          streamer.NewStateDiffStreamer(rpcClient),
		Repository:        NewCIDRepository(db),
		Converter:         ipfs.NewPayloadConverter(params.MainnetChainConfig),
		Publisher:         publisher,
		Filterer:          NewResponseFilterer(),
		Fetcher:           fetcher,
		Retriever:         NewCIDRetriever(db),
		Resolver:          ipfs.NewIPLDResolver(),
		PayloadChan:       make(chan statediff.Payload, payloadChanBufferSize),
		QuitChan:          qc,
		Subscriptions:     make(map[common.Hash]map[rpc.ID]Subscription),
		SubscriptionTypes: make(map[common.Hash]config.Subscription),
		WorkerPoolSize:    workers,
	}, nil
}

// Protocols exports the services p2p protocols, this service has none
func (sap *Service) Protocols() []p2p.Protocol {
	return []p2p.Protocol{}
}

// APIs returns the RPC descriptors the StateDiffingService offers
func (sap *Service) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: APIName,
			Version:   APIVersion,
			Service:   NewPublicSeedNodeAPI(sap),
			Public:    true,
		},
	}
}

// SyncAndPublish is the backend processing loop which streams data from geth, converts it to iplds, publishes them to ipfs, and indexes their cids
// This continues on no matter if or how many subscribers there are, it then forwards the data to the ScreenAndServe() loop
// which filters and sends relevant data to client subscriptions, if there are any
func (sap *Service) SyncAndPublish(wg *sync.WaitGroup, screenAndServePayload chan<- ipfs.IPLDPayload, screenAndServeQuit chan<- bool) error {
	sub, err := sap.Streamer.Stream(sap.PayloadChan)
	if err != nil {
		return err
	}
	wg.Add(1)

	// Channels for forwarding data to the publishAndIndex workers
	publishAndIndexPayload := make(chan ipfs.IPLDPayload, payloadChanBufferSize)
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
				case screenAndServePayload <- *ipldPayload:
				default:
				}
				// Forward the payload to the publishAndIndex workers
				select {
				case publishAndIndexPayload <- *ipldPayload:
				default:
				}
			case err = <-sub.Err():
				log.Error(err)
			case <-sap.QuitChan:
				// If we have a ScreenAndServe process running, forward the quit signal to it
				select {
				case screenAndServeQuit <- true:
				default:
				}
				// Also forward a quit signal for each of the workers
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

	return nil
}

func (sap *Service) publishAndIndex(id int, publishAndIndexPayload <-chan ipfs.IPLDPayload, publishAndIndexQuit <-chan bool) {
	go func() {
		for {
			select {
			case payload := <-publishAndIndexPayload:
				cidPayload, err := sap.Publisher.Publish(&payload)
				if err != nil {
					log.Errorf("worker %d error: %v", id, err)
					continue
				}
				err = sap.Repository.Index(cidPayload)
				if err != nil {
					log.Errorf("worker %d error: %v", id, err)
				}
			case <-publishAndIndexQuit:
				log.Infof("quiting publishAndIndex worker %d", id)
				return
			}
		}
	}()
}

// ScreenAndServe is the loop used to screen data streamed from the state diffing eth node
// and send the appropriate portions of it to a requesting client subscription, according to their subscription configuration
func (sap *Service) ScreenAndServe(screenAndServePayload <-chan ipfs.IPLDPayload, screenAndServeQuit <-chan bool) {
	go func() {
		for {
			select {
			case payload := <-screenAndServePayload:
				err := sap.sendResponse(payload)
				if err != nil {
					log.Error(err)
				}
			case <-screenAndServeQuit:
				log.Info("quiting ScreenAndServe process")
				return
			}
		}
	}()
}

func (sap *Service) sendResponse(payload ipfs.IPLDPayload) error {
	sap.Lock()
	for ty, subs := range sap.Subscriptions {
		// Retrieve the subscription parameters for this subscription type
		subConfig, ok := sap.SubscriptionTypes[ty]
		if !ok {
			log.Errorf("subscription configuration for subscription type %s not available", ty.Hex())
			continue
		}
		response, err := sap.Filterer.FilterResponse(subConfig, payload)
		if err != nil {
			log.Error(err)
			continue
		}
		for id, sub := range subs {
			select {
			case sub.PayloadChan <- response:
				log.Infof("sending seed node payload to subscription %s", id)
			default:
				log.Infof("unable to send payload to subscription %s; channel has no receiver", id)
			}
		}
	}
	sap.Unlock()
	return nil
}

// Subscribe is used by the API to subscribe to the service loop
func (sap *Service) Subscribe(id rpc.ID, sub chan<- streamer.SeedNodePayload, quitChan chan<- bool, streamFilters config.Subscription) {
	log.Info("Subscribing to the seed node service")
	// Subscription type is defined as the hash of its content
	// Group subscriptions by type and screen payloads once for subs of the same type
	by, err := rlp.EncodeToBytes(streamFilters)
	if err != nil {
		log.Error(err)
	}
	subscriptionHash := crypto.Keccak256(by)
	subscriptionType := common.BytesToHash(subscriptionHash)
	subscription := Subscription{
		PayloadChan: sub,
		QuitChan:    quitChan,
	}
	// If the subscription requests a backfill, use the Postgres index to lookup and retrieve historical data
	// Otherwise we only filter new data as it is streamed in from the state diffing geth node
	if streamFilters.BackFill || streamFilters.BackFillOnly {
		sap.backFill(subscription, id, streamFilters)
	}
	if !streamFilters.BackFillOnly {
		sap.Lock()
		if sap.Subscriptions[subscriptionType] == nil {
			sap.Subscriptions[subscriptionType] = make(map[rpc.ID]Subscription)
		}
		sap.Subscriptions[subscriptionType][id] = subscription
		sap.SubscriptionTypes[subscriptionType] = streamFilters
		sap.Unlock()
	}
}

func (sap *Service) backFill(sub Subscription, id rpc.ID, con config.Subscription) {
	log.Debug("back-filling data for id", id)
	// Retrieve cached CIDs relevant to this subscriber
	var endingBlock int64
	var startingBlock int64
	var err error
	startingBlock, err = sap.Retriever.RetrieveFirstBlockNumber()
	if err != nil {
		sub.PayloadChan <- streamer.SeedNodePayload{
			ErrMsg: "unable to set block range; error: " + err.Error(),
		}
	}
	if startingBlock < con.StartingBlock.Int64() {
		startingBlock = con.StartingBlock.Int64()
	}
	if con.EndingBlock.Int64() <= 0 || con.EndingBlock.Int64() <= startingBlock {
		endingBlock, err = sap.Retriever.RetrieveLastBlockNumber()
		if err != nil {
			sub.PayloadChan <- streamer.SeedNodePayload{
				ErrMsg: "unable to set block range; error: " + err.Error(),
			}
		}
	}
	log.Debug("backfill starting block:", con.StartingBlock)
	log.Debug("backfill ending block:", endingBlock)
	// Backfilled payloads are sent concurrently to the streamed payloads, so the receiver needs to pay attention to
	// the blocknumbers in the payloads they receive to keep things in order
	// TODO: separate backfill into a different rpc subscription method altogether?
	go func() {
		for i := con.StartingBlock.Int64(); i <= endingBlock; i++ {
			cidWrapper, err := sap.Retriever.RetrieveCIDs(con, i)
			if err != nil {
				sub.PayloadChan <- streamer.SeedNodePayload{
					ErrMsg: "CID retrieval error: " + err.Error(),
				}
				continue
			}
			if ipfs.EmptyCIDWrapper(*cidWrapper) {
				continue
			}
			blocksWrapper, err := sap.Fetcher.FetchCIDs(*cidWrapper)
			if err != nil {
				log.Error(err)
				sub.PayloadChan <- streamer.SeedNodePayload{
					ErrMsg: "IPLD fetching error: " + err.Error(),
				}
				continue
			}
			backFillIplds, err := sap.Resolver.ResolveIPLDs(*blocksWrapper)
			if err != nil {
				log.Error(err)
				sub.PayloadChan <- streamer.SeedNodePayload{
					ErrMsg: "IPLD resolving error: " + err.Error(),
				}
				continue
			}
			select {
			case sub.PayloadChan <- backFillIplds:
				log.Infof("sending seed node back-fill payload to subscription %s", id)
			default:
				log.Infof("unable to send back-fill payload to subscription %s; channel has no receiver", id)
			}
		}
	}()
}

// Unsubscribe is used to unsubscribe to the StateDiffingService loop
func (sap *Service) Unsubscribe(id rpc.ID) {
	log.Info("Unsubscribing from the seed node service")
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
	log.Info("Starting seed node service")
	wg := new(sync.WaitGroup)
	payloadChan := make(chan ipfs.IPLDPayload, payloadChanBufferSize)
	quitChan := make(chan bool, 1)
	if err := sap.SyncAndPublish(wg, payloadChan, quitChan); err != nil {
		return err
	}
	sap.ScreenAndServe(payloadChan, quitChan)
	return nil
}

// Stop is used to close down the service
func (sap *Service) Stop() error {
	log.Info("Stopping seed node service")
	close(sap.QuitChan)
	return nil
}

// close is used to close all listening subscriptions
func (sap *Service) close() {
	sap.Lock()
	for ty, subs := range sap.Subscriptions {
		for id, sub := range subs {
			select {
			case sub.QuitChan <- true:
				log.Infof("closing subscription %s", id)
			default:
				log.Infof("unable to close subscription %s; channel has no receiver", id)
			}
		}
		delete(sap.Subscriptions, ty)
		delete(sap.SubscriptionTypes, ty)
	}
	sap.Unlock()
}
