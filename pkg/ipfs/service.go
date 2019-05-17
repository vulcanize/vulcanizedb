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

// SyncAndPublish is an interface for streaming, converting to IPLDs, publishing, and indexing all Ethereum data
// This is the top-level interface used by the syncAndPublish command
type SyncPublishAndServe interface {
	// APIs(), Protocols(), Start() and Stop()
	node.Service
	// Main event loop for syncAndPublish processes
	SyncAndPublish(wg *sync.WaitGroup, forwardPayloadChan chan<- IPLDPayload, forwardQuitchan chan<- bool) error
	// Main event loop for handling client pub-sub
	Serve(wg *sync.WaitGroup, receivePayloadChan <-chan IPLDPayload, receiveQuitchan <-chan bool)
	// Method to subscribe to receive state diff processing output
	Subscribe(id rpc.ID, sub chan<- ResponsePayload, quitChan chan<- bool, params *Params)
	// Method to unsubscribe from state diff processing
	Unsubscribe(id rpc.ID) error
}

// Processor is the underlying struct for the SyncAndPublish interface
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
	// Chan the processor uses to subscribe to state diff payloads from the Streamer
	PayloadChan chan statediff.Payload
	// Used to signal shutdown of the service
	QuitChan chan bool
	// A mapping of rpc.IDs to their subscription channels
	Subscriptions map[rpc.ID]Subscription
}

// NewIPFSProcessor creates a new Processor interface using an underlying Processor struct
func NewIPFSProcessor(ipfsPath string, db *postgres.DB, ethClient core.EthClient, rpcClient core.RpcClient, qc chan bool) (SyncPublishAndServe, error) {
	publisher, err := NewIPLDPublisher(ipfsPath)
	if err != nil {
		return nil, err
	}
	return &Service{
		Streamer:    NewStateDiffStreamer(rpcClient),
		Repository:  NewCIDRepository(db),
		Converter:   NewPayloadConverter(ethClient),
		Publisher:   publisher,
		PayloadChan: make(chan statediff.Payload, payloadChanBufferSize),
		QuitChan:    qc,
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

func (sap *Service) Serve(wg *sync.WaitGroup, receivePayloadChan <-chan IPLDPayload, receiveQuitchan <-chan bool) {
	wg.Add(1)
	go func() {
		for {
			select {
			case payload := <-receivePayloadChan:
				println(payload.BlockNumber.Int64())
				// Method for using subscription parameters to filter payload and stream relevent info to sub channel
			case <-receiveQuitchan:
				log.Info("quiting Serve process")
				wg.Done()
				return
			}
		}
	}()
}

// Subscribe is used by the API to subscribe to the StateDiffingService loop
func (sap *Service) Subscribe(id rpc.ID, sub chan<- ResponsePayload, quitChan chan<- bool, params *Params) {
	log.Info("Subscribing to the statediff service")
	sap.Lock()
	sap.Subscriptions[id] = Subscription{
		PayloadChan: sub,
		QuitChan:    quitChan,
	}
	sap.Unlock()
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

// Start is used to begin the StateDiffingService
func (sap *Service) Start(*p2p.Server) error {
	log.Info("Starting statediff service")
	wg := new(sync.WaitGroup)
	payloadChan := make(chan IPLDPayload)
	quitChan := make(chan bool)
	go sap.SyncAndPublish(wg, payloadChan, quitChan)
	go sap.Serve(wg, payloadChan, quitChan)
	return nil
}

// Stop is used to close down the StateDiffingService
func (sap *Service) Stop() error {
	log.Info("Stopping statediff service")
	close(sap.QuitChan)
	return nil
}

// send is used to fan out and serve a payload to any subscriptions
func (sap *Service) send(payload ResponsePayload) {
	sap.Lock()
	for id, sub := range sap.Subscriptions {
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
