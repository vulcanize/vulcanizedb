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

package watcher

import (
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/super_node"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/watcher/shared"
	"github.com/vulcanize/vulcanizedb/pkg/wasm"
)

// SuperNodeWatcher is the top level interface for watching data from super node
type SuperNodeWatcher interface {
	Init() error
	Watch(wg *sync.WaitGroup) error
}

// Service is the underlying struct for the SuperNodeWatcher
type Service struct {
	// Config
	WatcherConfig Config
	// Interface for streaming data from super node
	SuperNodeStreamer shared.SuperNodeStreamer
	// Interface for db operations
	Repository shared.Repository
	// WASM instantiator
	WASMIniter *wasm.Instantiator

	// Channels for process communication/data relay
	PayloadChan chan super_node.SubscriptionPayload
	QuitChan    chan bool

	// Indexes
	// use atomic operations on these ONLY
	payloadIndex *int64
	endingIndex  *int64
	backFilling  *int32 // 0 => not backfilling; 1 => backfilling
}

// NewSuperNodeWatcher returns a new Service which satisfies the SuperNodeWatcher interface
func NewSuperNodeWatcher(c Config, quitChan chan bool) (SuperNodeWatcher, error) {
	repo, err := NewRepository(c.SubscriptionConfig.ChainType(), c.DB, c.TriggerFunctions)
	if err != nil {
		return nil, err
	}
	return &Service{
		WatcherConfig:     c,
		SuperNodeStreamer: NewSuperNodeStreamer(c.Client),
		Repository:        repo,
		WASMIniter:        wasm.NewWASMInstantiator(c.DB, c.WASMInstances),
		PayloadChan:       make(chan super_node.SubscriptionPayload, super_node.PayloadChanBufferSize),
		QuitChan:          quitChan,
	}, nil
}

// Init is used to initialize the Postgres WASM and trigger functions
func (s *Service) Init() error {
	// Instantiate the Postgres WASM functions
	if err := s.WASMIniter.Instantiate(); err != nil {
		return err
	}
	// Load the Postgres trigger functions that (can) use
	return s.Repository.LoadTriggers()
}

// Watch is the top level loop for watching super node
func (s *Service) Watch(wg *sync.WaitGroup) error {
	sub, err := s.SuperNodeStreamer.Stream(s.PayloadChan, s.WatcherConfig.SubscriptionConfig)
	if err != nil {
		return err
	}
	atomic.StoreInt64(s.payloadIndex, s.WatcherConfig.SubscriptionConfig.StartingBlock().Int64())
	atomic.StoreInt64(s.endingIndex, s.WatcherConfig.SubscriptionConfig.EndingBlock().Int64()) // less than 0 => never end
	backFilling := s.WatcherConfig.SubscriptionConfig.HistoricalData()
	if backFilling {
		atomic.StoreInt32(s.backFilling, 1)
	} else {
		atomic.StoreInt32(s.backFilling, 0)
	}
	backFillOnly := s.WatcherConfig.SubscriptionConfig.HistoricalDataOnly()
	if backFillOnly { // we are only processing historical data => handle single contiguous stream
		s.backFillOnlyQueuing(wg, sub)
	} else { // otherwise we need to be prepared to handle out-of-order data
		s.combinedQueuing(wg, sub)
	}
	return nil
}

// combinedQueuing assumes data is not necessarily going to come in linear order
// this is true when we are backfilling and streaming at the head or when we are
// only streaming at the head since reorgs can occur
func (s *Service) combinedQueuing(wg *sync.WaitGroup, sub *rpc.ClientSubscription) {
	wg.Add(1)
	go func() {
		for {
			select {
			case payload := <-s.PayloadChan:
				// If there is an error associated with the payload, log it and continue
				if payload.Error() != nil {
					logrus.Error(payload.Error())
					continue
				}
				if payload.Data.Height() == atomic.LoadInt64(s.payloadIndex) {
					// If the data is at our current index it is ready to be processed; add it to the ready data queue
					if err := s.Repository.ReadyData(payload); err != nil {
						logrus.Error(err)
					}
					atomic.AddInt64(s.payloadIndex, 1)
				} else { // otherwise add it to the wait queue
					if err := s.Repository.QueueData(payload); err != nil {
						logrus.Error(err)
					}
				}
			case err := <-sub.Err():
				logrus.Error(err)
			case <-s.QuitChan:
				logrus.Info("WatchContract shutting down")
				wg.Done()
				return
			}
		}
	}()
}

// backFillOnlyQueuing assumes the data is coming in contiguously and behind the head
// it puts all data on the ready queue
// it continues until the watcher is told to quit or we receive notification that the backfill is finished
func (s *Service) backFillOnlyQueuing(wg *sync.WaitGroup, sub *rpc.ClientSubscription) {
	wg.Add(1)
	go func() {
		for {
			select {
			case payload := <-s.PayloadChan:
				// If there is an error associated with the payload, log it and continue
				if payload.Error() != nil {
					logrus.Error(payload.Error())
					continue
				}
				// If the payload signals that backfilling has completed, shut down the process
				if payload.BackFillComplete() {
					logrus.Info("Backfill complete, WatchContract shutting down")
					wg.Done()
					return
				}
				// Add the payload the ready data queue
				if err := s.Repository.ReadyData(payload); err != nil {
					logrus.Error(err)
				}
			case err := <-sub.Err():
				logrus.Error(err)
			case <-s.QuitChan:
				logrus.Info("WatchContract shutting down")
				wg.Done()
				return
			}
		}
	}()
}
