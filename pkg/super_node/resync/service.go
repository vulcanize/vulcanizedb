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

package resync

import (
	"fmt"
	"github.com/sirupsen/logrus"

	utils "github.com/vulcanize/vulcanizedb/libraries/shared/utilities"
	"github.com/vulcanize/vulcanizedb/pkg/super_node"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

type Resync interface {
	Resync() error
}

type Service struct {
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
	// Interface for cleaning out data before resyncing (if clearOldCache is on)
	Cleaner shared.Cleaner
	// Size of batch fetches
	BatchSize uint64
	// Number of goroutines
	BatchNumber int64
	// Channel for receiving quit signal
	quitChan chan bool
	// Chain type
	chain shared.ChainType
	// Resync data type
	data shared.DataType
	// Resync ranges
	ranges [][2]uint64
	// Flag to turn on or off old cache destruction
	clearOldCache bool
	// Flag to turn on or off validation level reset
	resetValidation bool
}

// NewResyncService creates and returns a resync service from the provided settings
func NewResyncService(settings *Config) (Resync, error) {
	publisher, err := super_node.NewIPLDPublisher(settings.Chain, settings.IPFSPath, settings.DB, settings.IPFSMode)
	if err != nil {
		return nil, err
	}
	indexer, err := super_node.NewCIDIndexer(settings.Chain, settings.DB, settings.IPFSMode)
	if err != nil {
		return nil, err
	}
	converter, err := super_node.NewPayloadConverter(settings.Chain)
	if err != nil {
		return nil, err
	}
	retriever, err := super_node.NewCIDRetriever(settings.Chain, settings.DB)
	if err != nil {
		return nil, err
	}
	fetcher, err := super_node.NewPaylaodFetcher(settings.Chain, settings.HTTPClient, settings.Timeout)
	if err != nil {
		return nil, err
	}
	cleaner, err := super_node.NewCleaner(settings.Chain, settings.DB)
	if err != nil {
		return nil, err
	}
	batchSize := settings.BatchSize
	if batchSize == 0 {
		batchSize = super_node.DefaultMaxBatchSize
	}
	batchNumber := int64(settings.BatchNumber)
	if batchNumber == 0 {
		batchNumber = super_node.DefaultMaxBatchNumber
	}
	return &Service{
		Indexer:         indexer,
		Converter:       converter,
		Publisher:       publisher,
		Retriever:       retriever,
		Fetcher:         fetcher,
		Cleaner:         cleaner,
		BatchSize:       batchSize,
		BatchNumber:     int64(batchNumber),
		quitChan:        make(chan bool),
		chain:           settings.Chain,
		ranges:          settings.Ranges,
		data:            settings.ResyncType,
		clearOldCache:   settings.ClearOldCache,
		resetValidation: settings.ResetValidation,
	}, nil
}

func (rs *Service) Resync() error {
	if rs.resetValidation {
		logrus.Infof("resetting validation level")
		if err := rs.Cleaner.ResetValidation(rs.ranges); err != nil {
			return fmt.Errorf("validation reset failed: %v", err)
		}
	}
	if rs.clearOldCache {
		logrus.Infof("cleaning out old data from Postgres")
		if err := rs.Cleaner.Clean(rs.ranges, rs.data); err != nil {
			return fmt.Errorf("%s %s data resync cleaning error: %v", rs.chain.String(), rs.data.String(), err)
		}
	}
	// spin up worker goroutines
	heightsChan := make(chan []uint64)
	for i := 1; i <= int(rs.BatchNumber); i++ {
		go rs.resync(i, heightsChan)
	}
	for _, rng := range rs.ranges {
		if rng[1] < rng[0] {
			logrus.Errorf("%s resync range ending block number needs to be greater than the starting block number", rs.chain.String())
			continue
		}
		logrus.Infof("resyncing %s data from %d to %d", rs.chain.String(), rng[0], rng[1])
		// break the range up into bins of smaller ranges
		blockRangeBins, err := utils.GetBlockHeightBins(rng[0], rng[1], rs.BatchSize)
		if err != nil {
			return err
		}
		for _, heights := range blockRangeBins {
			heightsChan <- heights
		}
	}
	// send a quit signal to each worker
	// this blocks until each worker has finished its current task and can receive from the quit channel
	for i := 1; i <= int(rs.BatchNumber); i++ {
		rs.quitChan <- true
	}
	return nil
}

func (rs *Service) resync(id int, heightChan chan []uint64) {
	for {
		select {
		case heights := <-heightChan:
			logrus.Debugf("%s resync worker %d processing section from %d to %d", rs.chain.String(), id, heights[0], heights[len(heights)-1])
			payloads, err := rs.Fetcher.FetchAt(heights)
			if err != nil {
				logrus.Errorf("%s resync worker %d fetcher error: %s", rs.chain.String(), id, err.Error())
			}
			for _, payload := range payloads {
				ipldPayload, err := rs.Converter.Convert(payload)
				if err != nil {
					logrus.Errorf("%s resync worker %d converter error: %s", rs.chain.String(), id, err.Error())
				}
				cidPayload, err := rs.Publisher.Publish(ipldPayload)
				if err != nil {
					logrus.Errorf("%s resync worker %d publisher error: %s", rs.chain.String(), id, err.Error())
				}
				if err := rs.Indexer.Index(cidPayload); err != nil {
					logrus.Errorf("%s resync worker %d indexer error: %s", rs.chain.String(), id, err.Error())
				}
			}
			logrus.Infof("%s resync worker %d finished section from %d to %d", rs.chain.String(), id, heights[0], heights[len(heights)-1])
		case <-rs.quitChan:
			logrus.Infof("%s resync worker %d goroutine shutting down", rs.chain.String(), id)
			return
		}
	}
}
