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
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
)

const (
	DefaultMaxBatchSize   uint64 = 5000
	defaultMaxBatchNumber int64 = 100
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
	batchSize uint64
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
		Fetcher:  fetcher.NewStateDiffFetcher(archivalNodeRPCClient),
		GapCheckFrequency: freq,
		batchSize: DefaultMaxBatchSize,
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
	errChan := make(chan error)
	done := make(chan bool)
	backFillInitErr := bfs.BackFill(startingBlock, endingBlock, errChan, done)
	if backFillInitErr != nil {
		log.Error(backFillInitErr)
		return
	}
	for {
		select {
		case err := <- errChan:
			log.Error(err)
		case <- done:
			return
		}
	}
}


// BackFill fetches, processes, and returns utils.StorageDiffs over a range of blocks
// It splits a large range up into smaller chunks, batch fetching and processing those chunks concurrently
func (bfs *BackFillService) BackFill(startingBlock, endingBlock uint64, errChan chan error, done chan bool) error {
	if endingBlock < startingBlock {
		return errors.New("backfill: ending block number needs to be greater than starting block number")
	}
	// break the range up into bins of smaller ranges
	length := endingBlock - startingBlock + 1
	numberOfBins := length / bfs.batchSize
	remainder := length % bfs.batchSize
	if remainder != 0 {
		numberOfBins++
	}
	blockRangeBins := make([][]uint64, numberOfBins)
	for i := range blockRangeBins {
		nextBinStart := startingBlock + uint64(bfs.batchSize)
		if nextBinStart > endingBlock {
			nextBinStart = endingBlock + 1
		}
		blockRange := make([]uint64, 0, nextBinStart-startingBlock+1)
		for j := startingBlock; j < nextBinStart; j++ {
			blockRange = append(blockRange, j)
		}
		startingBlock = nextBinStart
		blockRangeBins[i] = blockRange
	}

	// int64 for atomic incrementing and decrementing to track the number of active processing goroutines we have
	var activeCount int64
	// channel for processing goroutines to signal when they are done
	processingDone := make(chan bool)
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
					stateDiff := new(statediff.StateDiff)
					stateDiffDecodeErr := rlp.DecodeBytes(payload.StateDiffRlp, stateDiff)
					if stateDiffDecodeErr != nil {
						errChan <- stateDiffDecodeErr
						continue
					}
					ipldPayload, convertErr := bfs.Converter.Convert(payload)
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
				// when this goroutine is done, send out a signal
				processingDone <- true
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
			case <-processingDone:
				atomic.AddInt64(&activeCount, -1)
				select {
				// if we are waiting for a process to finish, signal that one has
				case forwardDone <- true:
				default:
				}
				goroutinesFinished++
				if goroutinesFinished == int(numberOfBins) {
					done <- true
					return
				}
			}
		}
	}()

	return nil
}