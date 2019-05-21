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

package ipfs

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/statediff"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const payloadChanBufferSize = 800 // 1/10th max eth sub buffer size

// SyncPublishScreenAndServe is an interface for streaming, converting to IPLDs, publishing,
// indexing all Ethereum data screening this data, and serving it up to subscribed clients
type SyncPublishScreenAndServe interface {
	// APIs(), Protocols(), Start() and Stop()
	node.Service
	// Main event loop for syncAndPublish processes
	SyncAndPublish(wg *sync.WaitGroup, forwardPayloadChan chan<- IPLDPayload, forwardQuitchan chan<- bool) error
	// Main event loop for handling client pub-sub
	ScreenAndServe(wg *sync.WaitGroup, receivePayloadChan <-chan IPLDPayload, receiveQuitchan <-chan bool)
	// Method to subscribe to receive state diff processing output
	Subscribe(id rpc.ID, sub chan<- ResponsePayload, quitChan chan<- bool, streamFilters *StreamFilters)
	// Method to unsubscribe from state diff processing
	Unsubscribe(id rpc.ID) error
}

// Service is the underlying struct for the SyncAndPublish interface
type Service struct {
	// Used to sync access to the Subscriptions
	sync.Mutex
	// Interface for streaming statediff payloads over a geth rpc subscription
	Streamer StateDiffStreamer
	// Interface for converting statediff payloads into ETH-IPLD object payloads
	Converter PayloadConverter
	// Interface for publishing the ETH-IPLD payloads to IPFS
	Publisher IPLDPublisher
	// Interface for indexing the CIDs of the published ETH-IPLDs in Postgres
	Repository CIDRepository
	// Interface for filtering and serving data according to subscribed clients according to their specification
	Screener ResponseScreener
	// Interface for fetching ETH-IPLD objects from IPFS
	Fetcher IPLDFetcher
	// Interface for searching and retrieving CIDs from Postgres index
	Retriever CIDRetriever
	// Interface for resolving ipfs blocks to their data types
	Resolver IPLDResolver
	// Chan the processor uses to subscribe to state diff payloads from the Streamer
	PayloadChan chan statediff.Payload
	// Used to signal shutdown of the service
	QuitChan chan bool
	// A mapping of rpc.IDs to their subscription channels
	Subscriptions map[rpc.ID]Subscription
}

// NewIPFSProcessor creates a new Processor interface using an underlying Processor struct
func NewIPFSProcessor(ipfsPath string, db *postgres.DB, ethClient core.EthClient, rpcClient core.RpcClient, qc chan bool) (SyncPublishScreenAndServe, error) {
	publisher, err := NewIPLDPublisher(ipfsPath)
	if err != nil {
		return nil, err
	}
	fetcher, err := NewIPLDFetcher(ipfsPath)
	if err != nil {
		return nil, err
	}
	return &Service{
		Streamer:      NewStateDiffStreamer(rpcClient),
		Repository:    NewCIDRepository(db),
		Converter:     NewPayloadConverter(ethClient),
		Publisher:     publisher,
		Screener:      NewResponseScreener(),
		Fetcher:       fetcher,
		Retriever:     NewCIDRetriever(db),
		Resolver:      NewIPLDResolver(),
		PayloadChan:   make(chan statediff.Payload, payloadChanBufferSize),
		QuitChan:      qc,
		Subscriptions: make(map[rpc.ID]Subscription),
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
// It then forwards the data to the Serve() loop which filters and sends relevent data to client subscriptions
func (sap *Service) SyncAndPublish(wg *sync.WaitGroup, forwardPayloadChan chan<- IPLDPayload, forwardQuitchan chan<- bool) error {
	sub, err := sap.Streamer.Stream(sap.PayloadChan)
	if err != nil {
		return err
	}
	wg.Add(1)
	go func() {
		for {
			select {
			case payload := <-sap.PayloadChan:
				if payload.Err != nil {
					log.Error(err)
					continue
				}
				ipldPayload, err := sap.Converter.Convert(payload)
				if err != nil {
					log.Error(err)
					continue
				}
				// If we have a ScreenAndServe process running, forward the payload to it
				// If the ScreenAndServe process loop is slower than this one, will it miss some incoming payloads??
				select {
				case forwardPayloadChan <- *ipldPayload:
				default:
				}
				cidPayload, err := sap.Publisher.Publish(ipldPayload)
				if err != nil {
					log.Error(err)
					continue
				}
				err = sap.Repository.Index(cidPayload)
				if err != nil {
					log.Error(err)
				}
			case err = <-sub.Err():
				log.Error(err)
			case <-sap.QuitChan:
				// If we have a ScreenAndServe process running, forward the quit signal to it
				select {
				case forwardQuitchan <- true:
				default:
				}
				log.Info("quiting SyncAndPublish process")
				wg.Done()
				return
			}
		}
	}()

	return nil
}

// ScreenAndServe is the processing loop used to screen data streamed from the state diffing eth node and send the appropriate data to a requesting client subscription
func (sap *Service) ScreenAndServe(wg *sync.WaitGroup, receivePayloadChan <-chan IPLDPayload, receiveQuitchan <-chan bool) {
	wg.Add(1)
	go func() {
		for {
			select {
			case payload := <-receivePayloadChan:
				err := sap.processResponse(payload)
				if err != nil {
					log.Error(err)
				}
			case <-receiveQuitchan:
				log.Info("quiting ScreenAndServe process")
				wg.Done()
				return
			}
		}
	}()
}

func (sap *Service) processResponse(payload IPLDPayload) error {
	for id, sub := range sap.Subscriptions {
		response, err := sap.Screener.ScreenResponse(sub.StreamFilters, payload)
		if err != nil {
			return err
		}
		sap.serve(id, *response)
	}
	return nil
}

// Subscribe is used by the API to subscribe to the service loop
func (sap *Service) Subscribe(id rpc.ID, sub chan<- ResponsePayload, quitChan chan<- bool, streamFilters *StreamFilters) {
	log.Info("Subscribing to the statediff service")
	sap.Lock()
	sap.Subscriptions[id] = Subscription{
		PayloadChan:   sub,
		QuitChan:      quitChan,
		StreamFilters: streamFilters,
	}
	sap.Unlock()
	// If the subscription requests a backfill, use the Postgres index to lookup and retrieve historical data
	// Otherwise we only filter new data as it is streamed in from the state diffing geth node
	if streamFilters.BackFill {
		// Retrieve cached CIDs relevant to this subscriber
		cids, err := sap.Retriever.RetrieveCIDs(*streamFilters)
		if err != nil {
			log.Error(err)
			sap.serve(id, ResponsePayload{
				Err: err,
			})
			return
		}
		for _, cid := range cids {
			blocksWrapper, err := sap.Fetcher.FetchCIDs(cid)
			if err != nil {
				log.Error(err)
				sap.serve(id, ResponsePayload{
					Err: err,
				})
				return
			}
			backFillIplds, err := sap.Resolver.ResolveIPLDs(*blocksWrapper)
			if err != nil {
				log.Error(err)
				sap.serve(id, ResponsePayload{
					Err: err,
				})
				return
			}
			sap.serve(id, *backFillIplds)
		}
	}
}

// Unsubscribe is used to unsubscribe to the StateDiffingService loop
func (sap *Service) Unsubscribe(id rpc.ID) error {
	log.Info("Unsubscribing from the statediff service")
	sap.Lock()
	_, ok := sap.Subscriptions[id]
	if !ok {
		return fmt.Errorf("cannot unsubscribe; subscription for id %s does not exist", id)
	}
	delete(sap.Subscriptions, id)
	sap.Unlock()
	return nil
}

// Start is used to begin the service
func (sap *Service) Start(*p2p.Server) error {
	log.Info("Starting statediff service")
	wg := new(sync.WaitGroup)
	payloadChan := make(chan IPLDPayload)
	quitChan := make(chan bool)
	sap.SyncAndPublish(wg, payloadChan, quitChan)
	sap.ScreenAndServe(wg, payloadChan, quitChan)
	return nil
}

// Stop is used to close down the service
func (sap *Service) Stop() error {
	log.Info("Stopping statediff service")
	close(sap.QuitChan)
	return nil
}

// serve is used to send screened payloads to their requesting sub
func (sap *Service) serve(id rpc.ID, payload ResponsePayload) {
	sap.Lock()
	sub, ok := sap.Subscriptions[id]
	if ok {
		select {
		case sub.PayloadChan <- payload:
			log.Infof("sending state diff payload to subscription %s", id)
		default:
			log.Infof("unable to send payload to subscription %s; channel has no receiver", id)
		}
	}
	sap.Unlock()
}

// close is used to close all listening subscriptions
func (sap *Service) close() {
	sap.Lock()
	for id, sub := range sap.Subscriptions {
		select {
		case sub.QuitChan <- true:
			delete(sap.Subscriptions, id)
			log.Infof("closing subscription %s", id)
		default:
			log.Infof("unable to close subscription %s; channel has no receiver", id)
		}
	}
	sap.Unlock()
}
