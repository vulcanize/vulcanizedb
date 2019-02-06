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

package shared_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hpcloud/tail"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/storage"
	shared2 "github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Storage Watcher", func() {
	It("adds transformers", func() {
		fakeAddress := common.HexToAddress("0x12345")
		fakeTransformer := &mocks.MockStorageTransformer{Address: fakeAddress}
		watcher := shared.NewStorageWatcher(&fakes.MockTailer{}, test_config.NewTestDB(core.Node{}))

		watcher.AddTransformers([]storage.TransformerInitializer{fakeTransformer.FakeTransformerInitializer})

		Expect(watcher.Transformers[fakeAddress]).To(Equal(fakeTransformer))
	})

	It("reads the tail of the storage diffs file", func() {
		mockTailer := fakes.NewMockTailer()
		watcher := shared.NewStorageWatcher(mockTailer, test_config.NewTestDB(core.Node{}))

		assert(func(err error) {
			Expect(err).To(BeNil())
			Expect(mockTailer.TailCalled).To(BeTrue())
		}, watcher, mockTailer, []*tail.Line{})
	})

	It("returns error if row parsing fails", func() {
		mockTailer := fakes.NewMockTailer()
		watcher := shared.NewStorageWatcher(mockTailer, test_config.NewTestDB(core.Node{}))
		line := &tail.Line{Text: "oops"}

		assert(func(err error) {
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(shared2.ErrRowMalformed{Length: 1}))
		}, watcher, mockTailer, []*tail.Line{line})
	})

	It("logs error if no transformer can parse storage row", func() {
		mockTailer := fakes.NewMockTailer()
		line := &tail.Line{
			Text: "12345,block_hash,123,storage_key,storage_value",
			Time: time.Time{},
			Err:  nil,
		}
		watcher := shared.NewStorageWatcher(mockTailer, test_config.NewTestDB(core.Node{}))
		tempFile, err := ioutil.TempFile("", "log")
		defer os.Remove(tempFile.Name())
		Expect(err).NotTo(HaveOccurred())
		logrus.SetOutput(tempFile)

		assert(func(err error) {
			Expect(err).NotTo(HaveOccurred())
			logContent, readErr := ioutil.ReadFile(tempFile.Name())
			Expect(readErr).NotTo(HaveOccurred())
			Expect(string(logContent)).To(ContainSubstring(shared2.ErrContractNotFound{Contract: common.HexToAddress("0x12345").Hex()}.Error()))
		}, watcher, mockTailer, []*tail.Line{line})
	})

	It("executes transformer with storage row", func() {
		address := []byte{1, 2, 3}
		blockHash := []byte{4, 5, 6}
		blockHeight := int64(789)
		storageKey := []byte{9, 8, 7}
		storageValue := []byte{6, 5, 4}
		mockTailer := fakes.NewMockTailer()
		line := &tail.Line{
			Text: fmt.Sprintf("%s,%s,%d,%s,%s", common.Bytes2Hex(address), common.Bytes2Hex(blockHash), blockHeight, common.Bytes2Hex(storageKey), common.Bytes2Hex(storageValue)),
			Time: time.Time{},
			Err:  nil,
		}
		watcher := shared.NewStorageWatcher(mockTailer, test_config.NewTestDB(core.Node{}))
		fakeTransformer := &mocks.MockStorageTransformer{Address: common.BytesToAddress(address)}
		watcher.AddTransformers([]storage.TransformerInitializer{fakeTransformer.FakeTransformerInitializer})

		assert(func(err error) {
			Expect(err).To(BeNil())
			expectedRow, err := shared2.FromStrings(strings.Split(line.Text, ","))
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeTransformer.PassedRow).To(Equal(expectedRow))
		}, watcher, mockTailer, []*tail.Line{line})
	})

	It("logs error if executing transformer fails", func() {
		address := []byte{1, 2, 3}
		blockHash := []byte{4, 5, 6}
		blockHeight := int64(789)
		storageKey := []byte{9, 8, 7}
		storageValue := []byte{6, 5, 4}
		mockTailer := fakes.NewMockTailer()
		line := &tail.Line{
			Text: fmt.Sprintf("%s,%s,%d,%s,%s", common.Bytes2Hex(address), common.Bytes2Hex(blockHash), blockHeight, common.Bytes2Hex(storageKey), common.Bytes2Hex(storageValue)),
			Time: time.Time{},
			Err:  nil,
		}
		watcher := shared.NewStorageWatcher(mockTailer, test_config.NewTestDB(core.Node{}))
		executionError := errors.New("storage watcher failed attempting to execute transformer")
		fakeTransformer := &mocks.MockStorageTransformer{Address: common.BytesToAddress(address), ExecuteErr: executionError}
		watcher.AddTransformers([]storage.TransformerInitializer{fakeTransformer.FakeTransformerInitializer})
		tempFile, err := ioutil.TempFile("", "log")
		defer os.Remove(tempFile.Name())
		Expect(err).NotTo(HaveOccurred())
		logrus.SetOutput(tempFile)

		assert(func(err error) {
			Expect(err).NotTo(HaveOccurred())
			logContent, readErr := ioutil.ReadFile(tempFile.Name())
			Expect(readErr).NotTo(HaveOccurred())
			Expect(string(logContent)).To(ContainSubstring(executionError.Error()))
		}, watcher, mockTailer, []*tail.Line{line})
	})
})

func assert(assertion func(err error), watcher shared.StorageWatcher, mockTailer *fakes.MockTailer, lines []*tail.Line) {
	errs := make(chan error, 1)
	done := make(chan bool, 1)
	go execute(watcher, mockTailer, errs, done)
	for _, line := range lines {
		mockTailer.Lines <- line
	}
	close(mockTailer.Lines)

	select {
	case err := <-errs:
		assertion(err)
		break
	case <-done:
		assertion(nil)
		break
	}
}

func execute(watcher shared.StorageWatcher, tailer *fakes.MockTailer, errs chan error, done chan bool) {
	err := watcher.Execute()
	if err != nil {
		errs <- err
	} else {
		done <- true
	}
}
