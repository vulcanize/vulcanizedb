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

// Indexer is an interface for streaming, converting to IPLDs, publishing, and indexing all Ethereum data
// This is the top-level interface used by the syncAndPublish command
type Indexer interface {
	Index(wg *sync.WaitGroup) error
}

// ipfsIndexer is the underlying struct for the Indexer interface
// we are not exporting this to enforce proper initialization through the NewIPFSIndexer function
type ipfsIndexer struct {
	Streamer    Streamer
	Converter   Converter
	Publisher   Publisher
	Repository  Repository
	PayloadChan chan statediff.Payload
	QuitChan    chan bool
}

// NewIPFSIndexer creates a new Indexer interface using an underlying ipfsIndexer struct
func NewIPFSIndexer(ipfsPath string, db *postgres.DB, ethClient core.EthClient, rpcClient core.RpcClient, qc chan bool) (Indexer, error) {
	publisher, err := NewIPLDPublisher(ipfsPath)
	if err != nil {
		return nil, err
	}
	return &ipfsIndexer{
		Streamer:    NewStateDiffSyncer(rpcClient),
		Repository:  NewCIDRepository(db),
		Converter:   NewPayloadConverter(ethClient),
		Publisher:   publisher,
		PayloadChan: make(chan statediff.Payload, payloadChanBufferSize),
		QuitChan:    qc,
	}, nil
}

// Index is the main processing loop
func (i *ipfsIndexer) Index(wg *sync.WaitGroup) error {
	sub, err := i.Streamer.Stream(i.PayloadChan)
	if err != nil {
		return err
	}
	go func() {
		wg.Add(1)
		defer wg.Done()
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
				}
				cidPayload, err := i.Publisher.Publish(ipldPayload)
				if err != nil {
					log.Error(err)
				}
				err = i.Repository.Index(cidPayload)
				if err != nil {
					log.Error(err)
				}
			case err = <-sub.Err():
				log.Error(err)
			case <-i.QuitChan:
				log.Info("quiting IPFSIndexer")
				return
			}
		}
	}()

	return nil
}
