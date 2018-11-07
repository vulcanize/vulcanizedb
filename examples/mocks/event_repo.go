// VulcanizeDB
// Copyright © 2018 Vulcanize

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

package mocks

import (
	et1 "github.com/vulcanize/vulcanizedb/examples/erc20_watcher/event_triggered"
	et2 "github.com/vulcanize/vulcanizedb/examples/generic/event_triggered"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

type MockWatchedEventsRepository struct {
	watchedTransferEvents []*core.WatchedEvent
	watchedApprovalEvents []*core.WatchedEvent
	watchedBurnEvents     []*core.WatchedEvent
	watchedMintEvents     []*core.WatchedEvent
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
		if event.Name == "Burn" {
			mwer.watchedBurnEvents = append(mwer.watchedBurnEvents, event)
		}
		if event.Name == "Mint" {
			mwer.watchedMintEvents = append(mwer.watchedMintEvents, event)
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
	if name == "Burn" {
		result = mwer.watchedBurnEvents
		// clear watched events once returned so same events are returned for every filter while testing
		mwer.watchedBurnEvents = []*core.WatchedEvent{}
	}
	if name == "Mint" {
		result = mwer.watchedMintEvents
		// clear watched events once returned so same events are returned for every filter while testing
		mwer.watchedMintEvents = []*core.WatchedEvent{}
	}
	return result, nil
}

type MockEventRepo struct {
	TransferLogs    []et1.TransferModel
	ApprovalLogs    []et1.ApprovalModel
	BurnLogs        []et2.BurnModel
	MintLogs        []et2.MintModel
	VulcanizeLogIDs []int64
}

func (molr *MockEventRepo) CreateTransfer(transferModel *et1.TransferModel, vulcanizeLogId int64) error {
	molr.TransferLogs = append(molr.TransferLogs, *transferModel)
	molr.VulcanizeLogIDs = append(molr.VulcanizeLogIDs, vulcanizeLogId)
	return nil
}

func (molk *MockEventRepo) CreateApproval(approvalModel *et1.ApprovalModel, vulcanizeLogID int64) error {
	molk.ApprovalLogs = append(molk.ApprovalLogs, *approvalModel)
	molk.VulcanizeLogIDs = append(molk.VulcanizeLogIDs, vulcanizeLogID)
	return nil
}

func (molr *MockEventRepo) CreateBurn(burnModel *et2.BurnModel, vulcanizeLogId int64) error {
	molr.BurnLogs = append(molr.BurnLogs, *burnModel)
	molr.VulcanizeLogIDs = append(molr.VulcanizeLogIDs, vulcanizeLogId)
	return nil
}

func (molk *MockEventRepo) CreateMint(mintModel *et2.MintModel, vulcanizeLogID int64) error {
	molk.MintLogs = append(molk.MintLogs, *mintModel)
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
