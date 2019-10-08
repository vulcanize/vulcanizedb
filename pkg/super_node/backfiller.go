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
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/params"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
)

// BackFillInterface for filling in gaps in the super node
type BackFillInterface interface {
	// Method for the super node to periodically check for and fill in gaps in its data using an archival node
	FillGaps(wg *sync.WaitGroup, quitChan <-chan bool)
}

// BackFillService for filling in gaps in the super node
type BackFillService struct {
	// Interface for converting statediff payloads into ETH-IPLD object payloads
	Converter ipfs.PayloadConverter
	// Interface for publishing the ETH-IPLD payloads to IPFS
	Publisher ipfs.IPLDPublisher
	// Interface for indexing the CIDs of the published ETH-IPLDs in Postgres
	Repository CIDRepository
	// Interface for searching and retrieving CIDs from Postgres index
	Retriever CIDRetriever
	// State-diff fetcher; needs to be configured with an archival core.RpcClient
	StateDiffFetcher fetcher.IStateDiffFetcher
	// Check frequency
	GapCheckFrequency time.Duration
}

// NewBackFillService returns a new BackFillInterface
func NewBackFillService(ipfsPath string, db *postgres.DB, archivalNodeRPCClient core.RpcClient, freq time.Duration) (BackFillInterface, error) {
	publisher, err := ipfs.NewIPLDPublisher(ipfsPath)
	if err != nil {
		return nil, err
	}
	return &BackFillService{
		Repository:        NewCIDRepository(db),
		Converter:         ipfs.NewPayloadConverter(params.MainnetChainConfig),
		Publisher:         publisher,
		Retriever:         NewCIDRetriever(db),
		StateDiffFetcher:  fetcher.NewStateDiffFetcher(archivalNodeRPCClient),
		GapCheckFrequency: freq,
	}, nil
}

// FillGaps periodically checks for and fills in gaps in the super node db
// this requires a core.RpcClient that is pointed at an archival node with the StateDiffAt method exposed
func (bfs *BackFillService) FillGaps(wg *sync.WaitGroup, quitChan <-chan bool) {
	ticker := time.NewTicker(bfs.GapCheckFrequency)
	wg.Add(1)

	go func() {
		for {
			select {
			case <-quitChan:
				log.Info("quiting FillGaps process")
				wg.Done()
				return
			case <-ticker.C:
				log.Info("searching for gaps in the super node database")
				startingBlock, firstBlockErr := bfs.Retriever.RetrieveFirstBlockNumber()
				if firstBlockErr != nil {
					log.Error(firstBlockErr)
					continue
				}
				if startingBlock != 1 {
					startingGap := [2]int64{
						1,
						startingBlock - 1,
					}
					log.Info("found gap at the beginning of the sync")
					bfs.fillGaps(startingGap)
				}

				gaps, gapErr := bfs.Retriever.RetrieveGapsInData()
				if gapErr != nil {
					log.Error(gapErr)
					continue
				}
				for _, gap := range gaps {
					bfs.fillGaps(gap)
				}
			}
		}
	}()
	log.Info("fillGaps goroutine successfully spun up")
}

func (bfs *BackFillService) fillGaps(gap [2]int64) {
	log.Infof("filling in gap from block %d to block %d", gap[0], gap[1])
	blockHeights := make([]uint64, 0, gap[1]-gap[0]+1)
	for i := gap[0]; i <= gap[1]; i++ {
		blockHeights = append(blockHeights, uint64(i))
	}
	payloads, fetchErr := bfs.StateDiffFetcher.FetchStateDiffsAt(blockHeights)
	if fetchErr != nil {
		log.Error(fetchErr)
		return
	}
	for _, payload := range payloads {
		ipldPayload, convertErr := bfs.Converter.Convert(*payload)
		if convertErr != nil {
			log.Error(convertErr)
			continue
		}
		cidPayload, publishErr := bfs.Publisher.Publish(ipldPayload)
		if publishErr != nil {
			log.Error(publishErr)
			continue
		}
		indexErr := bfs.Repository.Index(cidPayload)
		if indexErr != nil {
			log.Error(indexErr)
		}
	}
}
