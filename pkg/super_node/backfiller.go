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
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

const (
	DefaultMaxBatchSize   uint64 = 100
	defaultMaxBatchNumber int64  = 10
)

// BackFillInterface for filling in gaps in the super node
type BackFillInterface interface {
	// Method for the super node to periodically check for and fill in gaps in its data using an archival node
	FillGapsInSuperNode(wg *sync.WaitGroup)
}

// BackFillService for filling in gaps in the super node
type BackFillService struct {
	// Interface for converting payloads into IPLD object payloads
	Converter shared.PayloadConverter
	// Interface for publishing the IPLD payloads to IPFS
	Publisher shared.IPLDPublisher
	// Interface for indexing the CIDs of the published IPLDs in Postgres
	Indexer shared.CIDIndexer
	// Interface for searching and retrieving CIDs from Postgres index
	Retriever shared.CIDRetriever
	// Interface for fetching payloads over at historical blocks; over http
	Fetcher shared.PayloadFetcher
	// Channel for forwarding backfill payloads to the ScreenAndServe process
	ScreenAndServeChan chan shared.ConvertedData
	// Check frequency
	GapCheckFrequency time.Duration
	// Size of batch fetches
	BatchSize uint64
	// Channel for receiving quit signal
	QuitChan chan bool
	// Chain type
	chain shared.ChainType
}

// NewBackFillService returns a new BackFillInterface
func NewBackFillService(settings *Config, screenAndServeChan chan shared.ConvertedData) (BackFillInterface, error) {
	publisher, err := NewIPLDPublisher(settings.Chain, settings.IPFSPath)
	if err != nil {
		return nil, err
	}
	indexer, err := NewCIDIndexer(settings.Chain, settings.DB)
	if err != nil {
		return nil, err
	}
	converter, err := NewPayloadConverter(settings.Chain)
	if err != nil {
		return nil, err
	}
	retriever, err := NewCIDRetriever(settings.Chain, settings.DB)
	if err != nil {
		return nil, err
	}
	fetcher, err := NewPaylaodFetcher(settings.Chain, settings.HTTPClient)
	if err != nil {
		return nil, err
	}
	batchSize := settings.BatchSize
	if batchSize == 0 {
		batchSize = DefaultMaxBatchSize
	}
	return &BackFillService{
		Indexer:            indexer,
		Converter:          converter,
		Publisher:          publisher,
		Retriever:          retriever,
		Fetcher:            fetcher,
		GapCheckFrequency:  settings.Frequency,
		BatchSize:          batchSize,
		ScreenAndServeChan: screenAndServeChan,
		QuitChan:           settings.Quit,
		chain:              settings.Chain,
	}, nil
}

// FillGapsInSuperNode periodically checks for and fills in gaps in the super node db
func (bfs *BackFillService) FillGapsInSuperNode(wg *sync.WaitGroup) {
	ticker := time.NewTicker(bfs.GapCheckFrequency)
	wg.Add(1)

	go func() {
		for {
			select {
			case <-bfs.QuitChan:
				log.Infof("quiting %s FillGaps process", bfs.chain.String())
				wg.Done()
				return
			case <-ticker.C:
				log.Infof("searching for gaps in the %s super node database", bfs.chain.String())
				startingBlock, err := bfs.Retriever.RetrieveFirstBlockNumber()
				if err != nil {
					log.Errorf("super node db backfill RetrieveFirstBlockNumber error for chain %s: %v", bfs.chain.String(), err)
					continue
				}
				if startingBlock != 0 {
					log.Infof("found gap at the beginning of the %s sync", bfs.chain.String())
					bfs.fillGaps(0, uint64(startingBlock-1))
				}
				gaps, err := bfs.Retriever.RetrieveGapsInData()
				if err != nil {
					log.Error("super node db backfill RetrieveGapsInData error for chain %s: %v", bfs.chain.String(), err)
					continue
				}
				for _, gap := range gaps {
					if err := bfs.fillGaps(gap.Start, gap.Stop); err != nil {
						log.Error(err)
					}
				}
			}
		}
	}()
	log.Infof("%s fillGaps goroutine successfully spun up", bfs.chain.String())
}

func (bfs *BackFillService) fillGaps(startingBlock, endingBlock uint64) error {
	log.Infof("going to fill in %s gap from %d to %d", bfs.chain.String(), startingBlock, endingBlock)
	errChan := make(chan error)
	done := make(chan bool)
	err := bfs.backFill(startingBlock, endingBlock, errChan, done)
	if err != nil {
		return err
	}
	for {
		select {
		case err := <-errChan:
			log.Errorf("super node db backfill error for chain %s: %v", bfs.chain.String(), err)
		case <-done:
			log.Infof("finished filling in %s gap from %d to %d", bfs.chain.String(), startingBlock, endingBlock)
			return nil
		}
	}
}

// backFill fetches, processes, and returns utils.StorageDiffs over a range of blocks
// It splits a large range up into smaller chunks, batch fetching and processing those chunks concurrently
func (bfs *BackFillService) backFill(startingBlock, endingBlock uint64, errChan chan error, done chan bool) error {
	if endingBlock < startingBlock {
		return fmt.Errorf("super node %s db backfill: ending block number needs to be greater than starting block number", bfs.chain.String())
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
				payloads, err := bfs.Fetcher.FetchAt(blockHeights)
				if err != nil {
					errChan <- err
				}
				for _, payload := range payloads {
					ipldPayload, err := bfs.Converter.Convert(payload)
					if err != nil {
						errChan <- err
						continue
					}
					// If there is a ScreenAndServe process listening, forward payload to it
					select {
					case bfs.ScreenAndServeChan <- ipldPayload:
					default:
					}
					cidPayload, err := bfs.Publisher.Publish(ipldPayload)
					if err != nil {
						errChan <- err
						continue
					}
					if err := bfs.Indexer.Index(cidPayload); err != nil {
						errChan <- err
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
				log.Infof("finished filling in %s gap sub-bin from %d to %d", bfs.chain.String(), doneWithHeights[0], doneWithHeights[1])
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
