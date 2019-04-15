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
	log "github.com/sirupsen/logrus"
	"sync"

	"github.com/ethereum/go-ethereum/statediff"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const payloadChanBufferSize = 800 // 1/10th max size

type Indexer interface {
	Index() error
}

type IPFSIndexer struct {
	Syncer    Syncer
	Converter Converter
	Publisher Publisher
	Repository Repository
	PayloadChan chan statediff.Payload
	QuitChan    chan bool
}

func NewIPFSIndexer(ipfsPath string, db *postgres.DB, ethClient core.EthClient, rpcClient core.RpcClient, qc chan bool) (*IPFSIndexer, error) {
	publisher, err := NewIPLDPublisher(ipfsPath)
	if err != nil {
		return nil, err
	}
	return &IPFSIndexer{
		Syncer: NewStateDiffSyncer(rpcClient),
		Repository: NewCIDRepository(db),
		Converter: NewPayloadConverter(ethClient),
		Publisher: publisher,
		PayloadChan: make(chan statediff.Payload, payloadChanBufferSize),
		QuitChan: qc,
	}, nil
}

// The main processing loop for the syncAndPublish
func (i *IPFSIndexer) Index(wg sync.WaitGroup) error {
	sub, err := i.Syncer.Sync(i.PayloadChan)
	if err != nil {
		return err
	}
	go func() {
		wg.Add(1)
		defer wg.Done()
		for {
			select {
			case payload := <-i.PayloadChan:
				ipldPayload, err := i.Converter.Convert(payload)
				if err != nil {
					log.Error(err)
				}
				cidPayload, err := i.Publisher.Publish(ipldPayload)
				if err != nil {
					log.Error(err)
				}
				err = i.Repository.IndexCIDs(cidPayload)
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
