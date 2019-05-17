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
	"sync"

	"github.com/ethereum/go-ethereum/statediff"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const payloadChanBufferSize = 800 // 1/10th max eth sub buffer size

// SyncAndPublish is an interface for streaming, converting to IPLDs, publishing, and indexing all Ethereum data
// This is the top-level interface used by the syncAndPublish command
type SyncAndPublish interface {
	Process(wg *sync.WaitGroup) error
}

// Processor is the underlying struct for the SyncAndPublish interface
type Processor struct {
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
	// Chan used to shut down the Processor
	QuitChan chan bool
}

// NewIPFSProcessor creates a new Processor interface using an underlying Processor struct
func NewIPFSProcessor(ipfsPath string, db *postgres.DB, ethClient core.EthClient, rpcClient core.RpcClient, qc chan bool) (SyncAndPublish, error) {
	publisher, err := NewIPLDPublisher(ipfsPath)
	if err != nil {
		return nil, err
	}
	return &Processor{
		Streamer:    NewStateDiffStreamer(rpcClient),
		Repository:  NewCIDRepository(db),
		Converter:   NewPayloadConverter(ethClient),
		Publisher:   publisher,
		PayloadChan: make(chan statediff.Payload, payloadChanBufferSize),
		QuitChan:    qc,
	}, nil
}

// Process is the main processing loop
func (i *Processor) Process(wg *sync.WaitGroup) error {
	sub, err := i.Streamer.Stream(i.PayloadChan)
	if err != nil {
		return err
	}
	wg.Add(1)
	go func() {
		for {
			select {
			case payload := <-i.PayloadChan:
				if payload.Err != nil {
					log.Error(err)
					continue
				}
				ipldPayload, err := i.Converter.Convert(payload)
				if err != nil {
					log.Error(err)
					continue
				}
				cidPayload, err := i.Publisher.Publish(ipldPayload)
				if err != nil {
					log.Error(err)
					continue
				}
				err = i.Repository.Index(cidPayload)
				if err != nil {
					log.Error(err)
				}
			case err = <-sub.Err():
				log.Error(err)
			case <-i.QuitChan:
				log.Info("quiting IPFSProcessor")
				wg.Done()
				return
			}
		}
	}()

	return nil
}
