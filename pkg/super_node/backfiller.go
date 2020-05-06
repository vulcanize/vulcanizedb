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

	utils "github.com/vulcanize/vulcanizedb/libraries/shared/utilities"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

const (
	DefaultMaxBatchSize   uint64 = 100
	DefaultMaxBatchNumber int64  = 50
)

// BackFillInterface for filling in gaps in the super node
type BackFillInterface interface {
	// Method for the super node to periodically check for and fill in gaps in its data using an archival node
	BackFill(wg *sync.WaitGroup)
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
	// Number of goroutines
	BatchNumber int64
	// Channel for receiving quit signal
	QuitChan chan bool
	// Chain type
	Chain shared.ChainType
	// Headers with times_validated lower than this will be resynced
	validationLevel int
}

// NewBackFillService returns a new BackFillInterface
func NewBackFillService(settings *Config, screenAndServeChan chan shared.ConvertedData) (BackFillInterface, error) {
	publisher, err := NewIPLDPublisher(settings.Chain, settings.IPFSPath, settings.DB, settings.IPFSMode)
	if err != nil {
		return nil, err
	}
	indexer, err := NewCIDIndexer(settings.Chain, settings.DB, settings.IPFSMode)
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
	fetcher, err := NewPaylaodFetcher(settings.Chain, settings.HTTPClient, settings.Timeout)
	if err != nil {
		return nil, err
	}
	batchSize := settings.BatchSize
	if batchSize == 0 {
		batchSize = DefaultMaxBatchSize
	}
	batchNumber := int64(settings.BatchNumber)
	if batchNumber == 0 {
		batchNumber = DefaultMaxBatchNumber
	}
	return &BackFillService{
		Indexer:            indexer,
		Converter:          converter,
		Publisher:          publisher,
		Retriever:          retriever,
		Fetcher:            fetcher,
		GapCheckFrequency:  settings.Frequency,
		BatchSize:          batchSize,
		BatchNumber:        int64(batchNumber),
		ScreenAndServeChan: screenAndServeChan,
		QuitChan:           settings.Quit,
		Chain:              settings.Chain,
		validationLevel:    settings.ValidationLevel,
	}, nil
}

// BackFill periodically checks for and fills in gaps in the super node db
func (bfs *BackFillService) BackFill(wg *sync.WaitGroup) {
	ticker := time.NewTicker(bfs.GapCheckFrequency)
	wg.Add(1)

	go func() {
		for {
			select {
			case <-bfs.QuitChan:
				log.Infof("quiting %s FillGapsInSuperNode process", bfs.Chain.String())
				wg.Done()
				return
			case <-ticker.C:
				log.Infof("searching for gaps in the %s super node database", bfs.Chain.String())
				startingBlock, err := bfs.Retriever.RetrieveFirstBlockNumber()
				if err != nil {
					log.Errorf("super node db backfill RetrieveFirstBlockNumber error for chain %s: %v", bfs.Chain.String(), err)
					continue
				}
				if startingBlock != 0 && bfs.Chain == shared.Bitcoin || startingBlock != 1 && bfs.Chain == shared.Ethereum {
					log.Infof("found gap at the beginning of the %s sync", bfs.Chain.String())
					if err := bfs.backFill(0, uint64(startingBlock-1)); err != nil {
						log.Error(err)
					}
				}
				gaps, err := bfs.Retriever.RetrieveGapsInData(bfs.validationLevel)
				if err != nil {
					log.Errorf("super node db backfill RetrieveGapsInData error for chain %s: %v", bfs.Chain.String(), err)
					continue
				}
				for _, gap := range gaps {
					if err := bfs.backFill(gap.Start, gap.Stop); err != nil {
						log.Error(err)
					}
				}
			}
		}
	}()
	log.Infof("%s BackFill goroutine successfully spun up", bfs.Chain.String())
}

// backFill fetches, processes, and returns utils.StorageDiffs over a range of blocks
// It splits a large range up into smaller chunks, batch fetching and processing those chunks concurrently
func (bfs *BackFillService) backFill(startingBlock, endingBlock uint64) error {
	log.Infof("filling in %s gap from %d to %d", bfs.Chain.String(), startingBlock, endingBlock)
	if endingBlock < startingBlock {
		return fmt.Errorf("super node %s db backfill: ending block number needs to be greater than starting block number", bfs.Chain.String())
	}
	// break the range up into bins of smaller ranges
	blockRangeBins, err := utils.GetBlockHeightBins(startingBlock, endingBlock, bfs.BatchSize)
	if err != nil {
		return err
	}
	// int64 for atomic incrementing and decrementing to track the number of active processing goroutines we have
	var activeCount int64
	// channel for processing goroutines to signal when they are done
	processingDone := make(chan bool)
	forwardDone := make(chan bool)

	// for each block range bin spin up a goroutine to batch fetch and process data for that range
	go func() {
		for _, blockHeights := range blockRangeBins {
			// if we have reached our limit of active goroutines
			// wait for one to finish before starting the next
			if atomic.AddInt64(&activeCount, 1) > bfs.BatchNumber {
				// this blocks until a process signals it has finished
				<-forwardDone
			}
			go func(blockHeights []uint64) {
				payloads, err := bfs.Fetcher.FetchAt(blockHeights)
				if err != nil {
					log.Errorf("%s super node historical data fetcher error: %s", bfs.Chain.String(), err.Error())
				}
				for _, payload := range payloads {
					ipldPayload, err := bfs.Converter.Convert(payload)
					if err != nil {
						log.Errorf("%s super node historical data converter error: %s", bfs.Chain.String(), err.Error())
					}
					// If there is a ScreenAndServe process listening, forward payload to it
					select {
					case bfs.ScreenAndServeChan <- ipldPayload:
					default:
					}
					cidPayload, err := bfs.Publisher.Publish(ipldPayload)
					if err != nil {
						log.Errorf("%s super node historical data publisher error: %s", bfs.Chain.String(), err.Error())
					}
					if err := bfs.Indexer.Index(cidPayload); err != nil {
						log.Errorf("%s super node historical data indexer error: %s", bfs.Chain.String(), err.Error())
					}
				}
				// when this goroutine is done, send out a signal
				log.Infof("finished filling in %s gap from %d to %d", bfs.Chain.String(), blockHeights[0], blockHeights[len(blockHeights)-1])
				processingDone <- true
			}(blockHeights)
		}
	}()

	// listen on the processingDone chan
	// keeps track of the number of processing goroutines that have finished
	// when they have all finished, return
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
			if goroutinesFinished >= len(blockRangeBins) {
				return nil
			}
		}
	}
}
