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
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/params"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
)

const (
	DefaultMaxBatchSize   uint64 = 100
	defaultMaxBatchNumber int64  = 10
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
	Fetcher fetcher.StateDiffFetcher
	// Check frequency
	GapCheckFrequency time.Duration
	// size of batch fetches
	BatchSize uint64
}

// NewBackFillService returns a new BackFillInterface
func NewBackFillService(ipfsPath string, db *postgres.DB, archivalNodeRPCClient core.RPCClient, freq time.Duration, batchSize uint64) (BackFillInterface, error) {
	publisher, err := ipfs.NewIPLDPublisher(ipfsPath)
	if err != nil {
		return nil, err
	}
	return &BackFillService{
		Repository:        NewCIDRepository(db),
		Converter:         ipfs.NewPayloadConverter(params.MainnetChainConfig),
		Publisher:         publisher,
		Retriever:         NewCIDRetriever(db),
		Fetcher:           fetcher.NewStateDiffFetcher(archivalNodeRPCClient),
		GapCheckFrequency: freq,
		BatchSize:         batchSize,
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
					log.Info("found gap at the beginning of the sync")
					bfs.fillGaps(1, uint64(startingBlock-1))
				}

				gaps, gapErr := bfs.Retriever.RetrieveGapsInData()
				if gapErr != nil {
					log.Error(gapErr)
					continue
				}
				for _, gap := range gaps {
					bfs.fillGaps(gap[0], gap[1])
				}
			}
		}
	}()
	log.Info("fillGaps goroutine successfully spun up")
}

func (bfs *BackFillService) fillGaps(startingBlock, endingBlock uint64) {
	log.Infof("going to fill in gap from %d to %d", startingBlock, endingBlock)
	errChan := make(chan error)
	done := make(chan bool)
	backFillInitErr := bfs.backFill(startingBlock, endingBlock, errChan, done)
	if backFillInitErr != nil {
		log.Error(backFillInitErr)
		return
	}
	for {
		select {
		case err := <-errChan:
			log.Error(err)
		case <-done:
			log.Infof("finished filling in gap from %d to %d", startingBlock, endingBlock)
			return
		}
	}
}

// backFill fetches, processes, and returns utils.StorageDiffs over a range of blocks
// It splits a large range up into smaller chunks, batch fetching and processing those chunks concurrently
func (bfs *BackFillService) backFill(startingBlock, endingBlock uint64, errChan chan error, done chan bool) error {
	if endingBlock < startingBlock {
		return errors.New("backfill: ending block number needs to be greater than starting block number")
	}
	//
	// break the range up into bins of smaller ranges
	blockRangeBins, err := utils.GetBlockHeightBins(startingBlock, endingBlock, bfs.BatchSize)
	if err != nil {
		return err
	}
	// int64 for atomic incrementing and decrementing to track the number of active processing goroutines we have
	var activeCount int64
	// channel for processing goroutines to signal when they are done
	processingDone := make(chan [2]uint64)
	forwardDone := make(chan bool)

	// for each block range bin spin up a goroutine to batch fetch and process state diffs for that range
	go func() {
		for _, blockHeights := range blockRangeBins {
			// if we have reached our limit of active goroutines
			// wait for one to finish before starting the next
			if atomic.AddInt64(&activeCount, 1) > defaultMaxBatchNumber {
				// this blocks until a process signals it has finished
				<-forwardDone
			}
			go func(blockHeights []uint64) {
				payloads, fetchErr := bfs.Fetcher.FetchStateDiffsAt(blockHeights)
				if fetchErr != nil {
					errChan <- fetchErr
				}
				for _, payload := range payloads {
					ipldPayload, convertErr := bfs.Converter.Convert(payload)
					if convertErr != nil {
						errChan <- convertErr
						continue
					}
					cidPayload, publishErr := bfs.Publisher.Publish(ipldPayload)
					if publishErr != nil {
						errChan <- publishErr
						continue
					}
					indexErr := bfs.Repository.Index(cidPayload)
					if indexErr != nil {
						errChan <- indexErr
					}
				}
				// when this goroutine is done, send out a signal
				processingDone <- [2]uint64{blockHeights[0], blockHeights[len(blockHeights)-1]}
			}(blockHeights)
		}
	}()

	// goroutine that listens on the processingDone chan
	// keeps track of the number of processing goroutines that have finished
	// when they have all finished, sends the final signal out
	go func() {
		goroutinesFinished := 0
		for {
			select {
			case doneWithHeights := <-processingDone:
				atomic.AddInt64(&activeCount, -1)
				select {
				// if we are waiting for a process to finish, signal that one has
				case forwardDone <- true:
				default:
				}
				log.Infof("finished filling in gap sub-bin from %d to %d", doneWithHeights[0], doneWithHeights[1])
				goroutinesFinished++
				if goroutinesFinished >= len(blockRangeBins) {
					done <- true
					return
				}
			}
		}
	}()

	return nil
}
