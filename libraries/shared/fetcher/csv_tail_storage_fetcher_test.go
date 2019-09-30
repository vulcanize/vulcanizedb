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

package fetcher_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hpcloud/tail"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Csv Tail Storage Fetcher", func() {
	var (
		errorsChannel  chan error
		mockTailer     *fakes.MockTailer
		diffsChannel   chan utils.StorageDiff
		storageFetcher fetcher.CsvTailStorageFetcher
	)

	BeforeEach(func() {
		errorsChannel = make(chan error)
		diffsChannel = make(chan utils.StorageDiff)
		mockTailer = fakes.NewMockTailer()
		storageFetcher = fetcher.NewCsvTailStorageFetcher(mockTailer)
	})

	It("adds error to errors channel if tailing file fails", func(done Done) {
		mockTailer.TailErr = fakes.FakeError

		go storageFetcher.FetchStorageDiffs(diffsChannel, errorsChannel)

		Expect(<-errorsChannel).To(MatchError(fakes.FakeError))
		close(done)
	})

	It("adds parsed csv row to rows channel for storage diff", func(done Done) {
		line := getFakeLine()

		go storageFetcher.FetchStorageDiffs(diffsChannel, errorsChannel)
		mockTailer.Lines <- line

		expectedRow, err := utils.FromParityCsvRow(strings.Split(line.Text, ","))
		Expect(err).NotTo(HaveOccurred())
		Expect(<-diffsChannel).To(Equal(expectedRow))
		close(done)
	})

	It("adds error to errors channel if parsing csv fails", func(done Done) {
		line := &tail.Line{Text: "invalid"}

		go storageFetcher.FetchStorageDiffs(diffsChannel, errorsChannel)
		mockTailer.Lines <- line

		Expect(<-errorsChannel).To(HaveOccurred())
		select {
		case <-diffsChannel:
			Fail("value passed to rows channel on error")
		default:
			Succeed()
		}
		close(done)
	})
})

func getFakeLine() *tail.Line {
	address := common.HexToAddress("0x1234567890abcdef")
	blockHash := []byte{4, 5, 6}
	blockHeight := int64(789)
	storageKey := []byte{9, 8, 7}
	storageValue := []byte{6, 5, 4}
	return &tail.Line{
		Text: fmt.Sprintf("%s,%s,%d,%s,%s", common.Bytes2Hex(address.Bytes()), common.Bytes2Hex(blockHash),
			blockHeight, common.Bytes2Hex(storageKey), common.Bytes2Hex(storageValue)),
		Time: time.Time{},
		Err:  nil,
	}
}
