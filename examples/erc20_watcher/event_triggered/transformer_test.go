// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package event_triggered_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/examples/constants"
	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/event_triggered"
	"github.com/vulcanize/vulcanizedb/examples/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockERC20Converter struct {
	watchedEvents      []*core.WatchedEvent
	transfersToConvert []event_triggered.TransferEntity
	approvalsToConvert []event_triggered.ApprovalEntity
	block              int64
}

func (mlkc *MockERC20Converter) ToTransferModel(entity event_triggered.TransferEntity) event_triggered.TransferModel {
	mlkc.transfersToConvert = append(mlkc.transfersToConvert, entity)
	return event_triggered.TransferModel{}
}

func (mlkc *MockERC20Converter) ToTransferEntity(watchedEvent core.WatchedEvent) (*event_triggered.TransferEntity, error) {
	mlkc.watchedEvents = append(mlkc.watchedEvents, &watchedEvent)
	e := &event_triggered.TransferEntity{Block: watchedEvent.BlockNumber}
	mlkc.block++
	return e, nil
}

func (mlkc *MockERC20Converter) ToApprovalModel(entity event_triggered.ApprovalEntity) event_triggered.ApprovalModel {
	mlkc.approvalsToConvert = append(mlkc.approvalsToConvert, entity)
	return event_triggered.ApprovalModel{}
}

func (mlkc *MockERC20Converter) ToApprovalEntity(watchedEvent core.WatchedEvent) (*event_triggered.ApprovalEntity, error) {
	mlkc.watchedEvents = append(mlkc.watchedEvents, &watchedEvent)
	e := &event_triggered.ApprovalEntity{Block: watchedEvent.BlockNumber}
	mlkc.block++
	return e, nil
}

var blockID1 = int64(5428074)
var logID1 = int64(113)
var blockID2 = int64(5428405)
var logID2 = int64(100)

var fakeWatchedEvents = []*core.WatchedEvent{
	{
		LogID:       logID1,
		Name:        "Transfer",
		BlockNumber: blockID1,
		Address:     constants.DaiContractAddress,
		TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
		Index:       110,
		Topic0:      constants.TransferEventSignature,
		Topic1:      "0x000000000000000000000000000000000000000000000000000000000000af21",
		Topic2:      "0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
		Topic3:      "",
		Data:        "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
	},
	{
		LogID:       logID2,
		Name:        "Approval",
		BlockNumber: blockID2,
		Address:     constants.DaiContractAddress,
		TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
		Index:       110,
		Topic0:      constants.ApprovalEventSignature,
		Topic1:      "0x000000000000000000000000000000000000000000000000000000000000af21",
		Topic2:      "0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
		Topic3:      "",
		Data:        "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
	},
}

var _ = Describe("Mock ERC20 transformer", func() {
	var mockERC20Converter MockERC20Converter
	var watchedEventsRepo mocks.MockWatchedEventsRepository
	var mockEventRepo mocks.MockEventRepo
	var filterRepo mocks.MockFilterRepository
	var transformer event_triggered.DaiTransformer

	BeforeEach(func() {
		mockERC20Converter = MockERC20Converter{}
		watchedEventsRepo = mocks.MockWatchedEventsRepository{}
		watchedEventsRepo.SetWatchedEvents(fakeWatchedEvents)
		mockEventRepo = mocks.MockEventRepo{}
		filterRepo = mocks.MockFilterRepository{}
		transformer = event_triggered.DaiTransformer{
			Converter:              &mockERC20Converter,
			WatchedEventRepository: &watchedEventsRepo,
			FilterRepository:       filterRepo,
			Repository:             &mockEventRepo,
		}
	})

	It("calls the watched events repo with correct filter", func() {
		transformer.Execute()
		Expect(len(watchedEventsRepo.Names)).To(Equal(2))
		Expect(watchedEventsRepo.Names).To(ConsistOf([]string{"Transfer", "Approval"}))
	})

	It("calls the mock ERC20 converter with the watched events", func() {
		transformer.Execute()
		Expect(len(mockERC20Converter.watchedEvents)).To(Equal(2))
		Expect(mockERC20Converter.watchedEvents).To(ConsistOf(fakeWatchedEvents))
	})

	It("converts a Transfer entity to a model", func() {
		transformer.Execute()
		Expect(len(mockERC20Converter.transfersToConvert)).To(Equal(1))
		Expect(mockERC20Converter.transfersToConvert[0].Block).To(Equal(blockID1))

		Expect(len(mockERC20Converter.approvalsToConvert)).To(Equal(1))
		Expect(mockERC20Converter.approvalsToConvert[0].Block).To(Equal(blockID2))
	})

	It("persists Transfer and Approval data for each watched Transfer or Approval event", func() {
		transformer.Execute()
		Expect(len(mockEventRepo.TransferLogs)).To(Equal(1))
		Expect(len(mockEventRepo.ApprovalLogs)).To(Equal(1))
		Expect(mockEventRepo.VulcanizeLogIDs).To(ConsistOf(logID1, logID2))
	})

})
