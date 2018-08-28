// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mocks

import (
	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/event_triggered"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

type MockWatchedEventsRepository struct {
	watchedEvents []*core.WatchedEvent
	Names         []string
}

func (mwer *MockWatchedEventsRepository) SetWatchedEvents(watchedEvents []*core.WatchedEvent) {
	mwer.watchedEvents = watchedEvents
}

func (mwer *MockWatchedEventsRepository) GetWatchedEvents(name string) ([]*core.WatchedEvent, error) {
	mwer.Names = append(mwer.Names, name)
	result := mwer.watchedEvents
	// clear watched events once returned so same events are returned for every filter while testing
	mwer.watchedEvents = []*core.WatchedEvent{}
	return result, nil
}

type MockTransferRepo struct {
	LogMakes        []event_triggered.TransferModel
	VulcanizeLogIDs []int64
}

func (molr *MockTransferRepo) Create(offerModel event_triggered.TransferModel, vulcanizeLogId int64) error {
	molr.LogMakes = append(molr.LogMakes, offerModel)
	molr.VulcanizeLogIDs = append(molr.VulcanizeLogIDs, vulcanizeLogId)
	return nil
}

type MockApprovalRepo struct {
	LogKills        []event_triggered.ApprovalModel
	VulcanizeLogIDs []int64
}

func (molk *MockApprovalRepo) Create(model event_triggered.ApprovalModel, vulcanizeLogID int64) error {
	molk.LogKills = append(molk.LogKills, model)
	molk.VulcanizeLogIDs = append(molk.VulcanizeLogIDs, vulcanizeLogID)
	return nil
}

type MockFilterRepository struct {
}

func (MockFilterRepository) CreateFilter(filter filters.LogFilter) error {
	panic("implement me")
}

func (MockFilterRepository) GetFilter(name string) (filters.LogFilter, error) {
	panic("implement me")
}
