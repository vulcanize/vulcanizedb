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
	watchedTransferEvents []*core.WatchedEvent
	watchedApprovalEvents []*core.WatchedEvent
	Names                 []string
}

func (mwer *MockWatchedEventsRepository) SetWatchedEvents(watchedEvents []*core.WatchedEvent) {
	for _, event := range watchedEvents {
		if event.Name == "Transfer" {
			mwer.watchedTransferEvents = append(mwer.watchedTransferEvents, event)
		}
		if event.Name == "Approval" {
			mwer.watchedApprovalEvents = append(mwer.watchedApprovalEvents, event)
		}
	}
}

func (mwer *MockWatchedEventsRepository) GetWatchedEvents(name string) ([]*core.WatchedEvent, error) {
	mwer.Names = append(mwer.Names, name)
	var result []*core.WatchedEvent
	if name == "Transfer" {
		result = mwer.watchedTransferEvents
		// clear watched events once returned so same events are returned for every filter while testing
		mwer.watchedTransferEvents = []*core.WatchedEvent{}
	}
	if name == "Approval" {
		result = mwer.watchedApprovalEvents
		// clear watched events once returned so same events are returned for every filter while testing
		mwer.watchedApprovalEvents = []*core.WatchedEvent{}
	}
	return result, nil
}

type MockEventRepo struct {
	TransferLogs    []event_triggered.TransferModel
	ApprovalLogs    []event_triggered.ApprovalModel
	VulcanizeLogIDs []int64
}

func (molr *MockEventRepo) CreateTransfer(transferModel event_triggered.TransferModel, vulcanizeLogId int64) error {
	molr.TransferLogs = append(molr.TransferLogs, transferModel)
	molr.VulcanizeLogIDs = append(molr.VulcanizeLogIDs, vulcanizeLogId)
	return nil
}

func (molk *MockEventRepo) CreateApproval(approvalModel event_triggered.ApprovalModel, vulcanizeLogID int64) error {
	molk.ApprovalLogs = append(molk.ApprovalLogs, approvalModel)
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
