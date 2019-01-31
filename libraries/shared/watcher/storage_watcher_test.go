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

package watcher_test

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

	"github.com/vulcanize/vulcanizedb/libraries/shared/mocks"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/libraries/shared/utils"
	"github.com/vulcanize/vulcanizedb/libraries/shared/watcher"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Storage Watcher", func() {
	It("adds transformers", func() {
		fakeAddress := common.HexToAddress("0x12345")
		fakeTransformer := &mocks.MockStorageTransformer{Address: fakeAddress}
		w := watcher.NewStorageWatcher(&fakes.MockTailer{}, test_config.NewTestDB(core.Node{}))

		w.AddTransformers([]transformer.StorageTransformerInitializer{fakeTransformer.FakeTransformerInitializer})

		Expect(w.Transformers[fakeAddress]).To(Equal(fakeTransformer))
	})

	It("reads the tail of the storage diffs file", func() {
		mockTailer := fakes.NewMockTailer()
		w := watcher.NewStorageWatcher(mockTailer, test_config.NewTestDB(core.Node{}))

		assert(func(err error) {
			Expect(err).To(BeNil())
			Expect(mockTailer.TailCalled).To(BeTrue())
		}, w, mockTailer, []*tail.Line{})
	})

	It("returns error if row parsing fails", func() {
		mockTailer := fakes.NewMockTailer()
		w := watcher.NewStorageWatcher(mockTailer, test_config.NewTestDB(core.Node{}))
		line := &tail.Line{Text: "oops"}

		assert(func(err error) {
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrRowMalformed{Length: 1}))
		}, w, mockTailer, []*tail.Line{line})
	})

	It("logs error if no transformer can parse storage row", func() {
		mockTailer := fakes.NewMockTailer()
		address := common.HexToAddress("0x12345")
		line := getFakeLine(address.Bytes())
		w := watcher.NewStorageWatcher(mockTailer, test_config.NewTestDB(core.Node{}))
		tempFile, err := ioutil.TempFile("", "log")
		defer os.Remove(tempFile.Name())
		Expect(err).NotTo(HaveOccurred())
		logrus.SetOutput(tempFile)

		assert(func(err error) {
			Expect(err).NotTo(HaveOccurred())
			logContent, readErr := ioutil.ReadFile(tempFile.Name())
			Expect(readErr).NotTo(HaveOccurred())
			Expect(string(logContent)).To(ContainSubstring(utils.ErrContractNotFound{Contract: address.Hex()}.Error()))
		}, w, mockTailer, []*tail.Line{line})
	})

	It("executes transformer with storage row", func() {
		address := []byte{1, 2, 3}
		line := getFakeLine(address)
		mockTailer := fakes.NewMockTailer()
		w := watcher.NewStorageWatcher(mockTailer, test_config.NewTestDB(core.Node{}))
		fakeTransformer := &mocks.MockStorageTransformer{Address: common.BytesToAddress(address)}
		w.AddTransformers([]transformer.StorageTransformerInitializer{fakeTransformer.FakeTransformerInitializer})

		assert(func(err error) {
			Expect(err).To(BeNil())
			expectedRow, err := utils.FromStrings(strings.Split(line.Text, ","))
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeTransformer.PassedRow).To(Equal(expectedRow))
		}, w, mockTailer, []*tail.Line{line})
	})

	Describe("when executing transformer fails", func() {
		It("queues row when error is storage key not found", func() {
			address := []byte{1, 2, 3}
			line := getFakeLine(address)
			mockTailer := fakes.NewMockTailer()
			w := watcher.NewStorageWatcher(mockTailer, test_config.NewTestDB(core.Node{}))
			mockQueue := &mocks.MockStorageQueue{}
			w.Queue = mockQueue
			keyNotFoundError := utils.ErrStorageKeyNotFound{Key: "unknown_storage_key"}
			fakeTransformer := &mocks.MockStorageTransformer{Address: common.BytesToAddress(address), ExecuteErr: keyNotFoundError}
			w.AddTransformers([]transformer.StorageTransformerInitializer{fakeTransformer.FakeTransformerInitializer})

			assert(func(err error) {
				Expect(err).NotTo(HaveOccurred())
				Expect(mockQueue.AddCalled).To(BeTrue())
			}, w, mockTailer, []*tail.Line{line})
		})

		It("logs error if queuing row fails", func() {
			address := []byte{1, 2, 3}
			line := getFakeLine(address)
			mockTailer := fakes.NewMockTailer()
			w := watcher.NewStorageWatcher(mockTailer, test_config.NewTestDB(core.Node{}))
			mockQueue := &mocks.MockStorageQueue{}
			mockQueue.AddError = fakes.FakeError
			w.Queue = mockQueue
			keyNotFoundError := utils.ErrStorageKeyNotFound{Key: "unknown_storage_key"}
			fakeTransformer := &mocks.MockStorageTransformer{Address: common.BytesToAddress(address), ExecuteErr: keyNotFoundError}
			w.AddTransformers([]transformer.StorageTransformerInitializer{fakeTransformer.FakeTransformerInitializer})
			tempFile, err := ioutil.TempFile("", "log")
			defer os.Remove(tempFile.Name())
			Expect(err).NotTo(HaveOccurred())
			logrus.SetOutput(tempFile)

			assert(func(err error) {
				Expect(err).NotTo(HaveOccurred())
				Expect(mockQueue.AddCalled).To(BeTrue())
				logContent, readErr := ioutil.ReadFile(tempFile.Name())
				Expect(readErr).NotTo(HaveOccurred())
				Expect(string(logContent)).To(ContainSubstring(fakes.FakeError.Error()))
			}, w, mockTailer, []*tail.Line{line})
		})

		It("logs any other error", func() {
			address := []byte{1, 2, 3}
			line := getFakeLine(address)
			mockTailer := fakes.NewMockTailer()
			w := watcher.NewStorageWatcher(mockTailer, test_config.NewTestDB(core.Node{}))
			executionError := errors.New("storage watcher failed attempting to execute transformer")
			fakeTransformer := &mocks.MockStorageTransformer{Address: common.BytesToAddress(address), ExecuteErr: executionError}
			w.AddTransformers([]transformer.StorageTransformerInitializer{fakeTransformer.FakeTransformerInitializer})
			tempFile, err := ioutil.TempFile("", "log")
			defer os.Remove(tempFile.Name())
			Expect(err).NotTo(HaveOccurred())
			logrus.SetOutput(tempFile)

			assert(func(err error) {
				Expect(err).NotTo(HaveOccurred())
				logContent, readErr := ioutil.ReadFile(tempFile.Name())
				Expect(readErr).NotTo(HaveOccurred())
				Expect(string(logContent)).To(ContainSubstring(executionError.Error()))
			}, w, mockTailer, []*tail.Line{line})
		})
	})
})

func assert(assertion func(err error), watcher watcher.StorageWatcher, mockTailer *fakes.MockTailer, lines []*tail.Line) {
	errs := make(chan error, 1)
	done := make(chan bool, 1)
	go execute(watcher, errs, done)
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

func execute(w watcher.StorageWatcher, errs chan error, done chan bool) {
	err := w.Execute()
	if err != nil {
		errs <- err
	} else {
		done <- true
	}
}

func getFakeLine(address []byte) *tail.Line {
	blockHash := []byte{4, 5, 6}
	blockHeight := int64(789)
	storageKey := []byte{9, 8, 7}
	storageValue := []byte{6, 5, 4}
	return &tail.Line{
		Text: fmt.Sprintf("%s,%s,%d,%s,%s", common.Bytes2Hex(address), common.Bytes2Hex(blockHash), blockHeight, common.Bytes2Hex(storageKey), common.Bytes2Hex(storageValue)),
		Time: time.Time{},
		Err:  nil,
	}
}
