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
	Stop() error
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
	chain shared.ChainType
	// Headers with times_validated lower than this will be resynced
	validationLevel int
}

// NewBackFillService returns a new BackFillInterface
func NewBackFillService(settings *Config, screenAndServeChan chan shared.ConvertedData) (BackFillInterface, error) {
	publisher, err := NewIPLDPublisher(settings.Chain, settings.IPFSPath, settings.BackFillDBConn, settings.IPFSMode)
	if err != nil {
		return nil, err
	}
	indexer, err := NewCIDIndexer(settings.Chain, settings.BackFillDBConn, settings.IPFSMode)
	if err != nil {
		return nil, err
	}
	converter, err := NewPayloadConverter(settings.Chain)
	if err != nil {
		return nil, err
	}
	retriever, err := NewCIDRetriever(settings.Chain, settings.BackFillDBConn)
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
		QuitChan:           make(chan bool),
		chain:              settings.Chain,
		validationLevel:    settings.ValidationLevel,
	}, nil
}

// BackFill periodically checks for and fills in gaps in the super node db
func (bfs *BackFillService) BackFill(wg *sync.WaitGroup) {
	ticker := time.NewTicker(bfs.GapCheckFrequency)
	go func() {
		wg.Add(1)
		defer wg.Done()
		for {
			select {
			case <-bfs.QuitChan:
				log.Infof("quiting %s FillGapsInSuperNode process", bfs.chain.String())
				return
			case <-ticker.C:
				gaps, err := bfs.Retriever.RetrieveGapsInData(bfs.validationLevel)
				if err != nil {
					log.Errorf("%s super node db backFill RetrieveGapsInData error: %v", bfs.chain.String(), err)
					continue
				}
				// spin up worker goroutines for this search pass
				// we start and kill a new batch of workers for each pass
				// so that we know each of the previous workers is done before we search for new gaps
				heightsChan := make(chan []uint64)
				for i := 1; i <= int(bfs.BatchNumber); i++ {
					go bfs.backFill(wg, i, heightsChan)
				}
				for _, gap := range gaps {
					log.Infof("backFilling %s data from %d to %d", bfs.chain.String(), gap.Start, gap.Stop)
					blockRangeBins, err := utils.GetBlockHeightBins(gap.Start, gap.Stop, bfs.BatchSize)
					if err != nil {
						log.Errorf("%s super node db backFill GetBlockHeightBins error: %v", bfs.chain.String(), err)
						continue
					}
					for _, heights := range blockRangeBins {
						select {
						case <-bfs.QuitChan:
							log.Infof("quiting %s BackFill process", bfs.chain.String())
							return
						default:
							heightsChan <- heights
						}
					}
				}
				// send a quit signal to each worker
				// this blocks until each worker has finished its current task and is free to receive from the quit channel
				for i := 1; i <= int(bfs.BatchNumber); i++ {
					bfs.QuitChan <- true
				}
			}
		}
	}()
	log.Infof("%s BackFill goroutine successfully spun up", bfs.chain.String())
}

func (bfs *BackFillService) backFill(wg *sync.WaitGroup, id int, heightChan chan []uint64) {
	wg.Add(1)
	defer wg.Done()
	for {
		select {
		case heights := <-heightChan:
			log.Debugf("%s backFill worker %d processing section from %d to %d", bfs.chain.String(), id, heights[0], heights[len(heights)-1])
			payloads, err := bfs.Fetcher.FetchAt(heights)
			if err != nil {
				log.Errorf("%s backFill worker %d fetcher error: %s", bfs.chain.String(), id, err.Error())
			}
			for _, payload := range payloads {
				ipldPayload, err := bfs.Converter.Convert(payload)
				if err != nil {
					log.Errorf("%s backFill worker %d converter error: %s", bfs.chain.String(), id, err.Error())
				}
				// If there is a ScreenAndServe process listening, forward converted payload to it
				select {
				case bfs.ScreenAndServeChan <- ipldPayload:
					log.Debugf("%s backFill worker %d forwarded converted payload to server", bfs.chain.String(), id)
				default:
					log.Debugf("%s backFill worker %d unable to forward converted payload to server; no channel ready to receive", bfs.chain.String(), id)
				}
				cidPayload, err := bfs.Publisher.Publish(ipldPayload)
				if err != nil {
					log.Errorf("%s backFill worker %d publisher error: %s", bfs.chain.String(), id, err.Error())
					continue
				}
				if err := bfs.Indexer.Index(cidPayload); err != nil {
					log.Errorf("%s backFill worker %d indexer error: %s", bfs.chain.String(), id, err.Error())
				}
			}
			log.Infof("%s backFill worker %d finished section from %d to %d", bfs.chain.String(), id, heights[0], heights[len(heights)-1])
		case <-bfs.QuitChan:
			log.Infof("%s backFill worker %d shutting down", bfs.chain.String(), id)
			return
		}
	}
}

func (bfs *BackFillService) Stop() error {
	log.Infof("Stopping %s backFill service", bfs.chain.String())
	close(bfs.QuitChan)
	return nil
}
