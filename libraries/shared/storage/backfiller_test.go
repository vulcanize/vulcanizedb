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

package storage_test

import (
	"errors"

	"github.com/ethereum/go-ethereum/statediff"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/shared/mocks"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/libraries/shared/test_data"
)

var _ = Describe("BackFiller", func() {
	Describe("BackFill", func() {
		var (
			mockFetcher *mocks.StateDiffFetcher
			backFiller  storage.BackFiller
		)
		BeforeEach(func() {
			mockFetcher = new(mocks.StateDiffFetcher)
			mockFetcher.PayloadsToReturn = map[uint64]statediff.Payload{
				test_data.BlockNumber.Uint64():  test_data.MockStatediffPayload,
				test_data.BlockNumber2.Uint64(): test_data.MockStatediffPayload2,
			}
		})

		It("batch calls statediff_stateDiffAt", func() {
			backFiller = storage.NewStorageBackFiller(mockFetcher, test_data.BlockNumber.Uint64(), 100)
			backFill := make(chan utils.StorageDiff)
			done := make(chan bool)
			errChan := make(chan error)
			backFillInitErr := backFiller.BackFill(
				test_data.BlockNumber2.Uint64(),
				backFill,
				errChan,
				done)
			Expect(backFillInitErr).ToNot(HaveOccurred())
			var diffs []utils.StorageDiff
			for {
				select {
				case diff := <-backFill:
					diffs = append(diffs, diff)
					continue
				case err := <-errChan:
					Expect(err).ToNot(HaveOccurred())
					continue
				case <-done:
					break
				}
				break
			}
			Expect(mockFetcher.CalledTimes).To(Equal(int64(1)))
			Expect(len(diffs)).To(Equal(4))
			Expect(containsDiff(diffs, test_data.CreatedExpectedStorageDiff)).To(BeTrue())
			Expect(containsDiff(diffs, test_data.UpdatedExpectedStorageDiff)).To(BeTrue())
			Expect(containsDiff(diffs, test_data.DeletedExpectedStorageDiff)).To(BeTrue())
			Expect(containsDiff(diffs, test_data.UpdatedExpectedStorageDiff2)).To(BeTrue())
		})

		It("has a configurable batch size", func() {
			backFiller = storage.NewStorageBackFiller(mockFetcher, test_data.BlockNumber.Uint64(), 1)
			backFill := make(chan utils.StorageDiff)
			done := make(chan bool)
			errChan := make(chan error)
			backFillInitErr := backFiller.BackFill(
				test_data.BlockNumber2.Uint64(),
				backFill,
				errChan,
				done)
			Expect(backFillInitErr).ToNot(HaveOccurred())
			var diffs []utils.StorageDiff
			for {
				select {
				case diff := <-backFill:
					diffs = append(diffs, diff)
					continue
				case err := <-errChan:
					Expect(err).ToNot(HaveOccurred())
					continue
				case <-done:
					break
				}
				break
			}
			Expect(mockFetcher.CalledTimes).To(Equal(int64(2)))
			Expect(len(diffs)).To(Equal(4))
			Expect(containsDiff(diffs, test_data.CreatedExpectedStorageDiff)).To(BeTrue())
			Expect(containsDiff(diffs, test_data.UpdatedExpectedStorageDiff)).To(BeTrue())
			Expect(containsDiff(diffs, test_data.DeletedExpectedStorageDiff)).To(BeTrue())
			Expect(containsDiff(diffs, test_data.UpdatedExpectedStorageDiff2)).To(BeTrue())
		})

		It("handles bin numbers in excess of the goroutine limit (100)", func() {
			payloadsToReturn := make(map[uint64]statediff.Payload, 1001)
			for i := test_data.BlockNumber.Uint64(); i <= test_data.BlockNumber.Uint64()+1000; i++ {
				payloadsToReturn[i] = test_data.MockStatediffPayload
			}
			mockFetcher.PayloadsToReturn = payloadsToReturn
			// batch size of 2 with 1001 block range => 501 bins
			backFiller = storage.NewStorageBackFiller(mockFetcher, test_data.BlockNumber.Uint64(), 2)
			backFill := make(chan utils.StorageDiff)
			done := make(chan bool)
			errChan := make(chan error)
			backFillInitErr := backFiller.BackFill(
				test_data.BlockNumber.Uint64()+1000,
				backFill,
				errChan,
				done)
			Expect(backFillInitErr).ToNot(HaveOccurred())
			var diffs []utils.StorageDiff
			for {
				select {
				case diff := <-backFill:
					diffs = append(diffs, diff)
					continue
				case err := <-errChan:
					Expect(err).ToNot(HaveOccurred())
					continue
				case <-done:
					break
				}
				break
			}
			Expect(mockFetcher.CalledTimes).To(Equal(int64(501)))
			Expect(len(diffs)).To(Equal(3003))
			Expect(containsDiff(diffs, test_data.CreatedExpectedStorageDiff)).To(BeTrue())
			Expect(containsDiff(diffs, test_data.UpdatedExpectedStorageDiff)).To(BeTrue())
			Expect(containsDiff(diffs, test_data.DeletedExpectedStorageDiff)).To(BeTrue())
		})

		It("passes fetcher errors forward", func() {
			mockFetcher.FetchErrs = map[uint64]error{
				test_data.BlockNumber.Uint64(): errors.New("mock fetcher error"),
			}
			backFiller = storage.NewStorageBackFiller(mockFetcher, test_data.BlockNumber.Uint64(), 1)
			backFill := make(chan utils.StorageDiff)
			done := make(chan bool)
			errChan := make(chan error)
			backFillInitErr := backFiller.BackFill(
				test_data.BlockNumber2.Uint64(),
				backFill,
				errChan,
				done)
			Expect(backFillInitErr).ToNot(HaveOccurred())
			var numOfErrs int
			var diffs []utils.StorageDiff
			for {
				select {
				case diff := <-backFill:
					diffs = append(diffs, diff)
					continue
				case err := <-errChan:
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("mock fetcher error"))
					numOfErrs++
					continue
				case <-done:
					break
				}
				break
			}
			Expect(mockFetcher.CalledTimes).To(Equal(int64(2)))
			Expect(numOfErrs).To(Equal(1))
			Expect(len(diffs)).To(Equal(1))
			Expect(containsDiff(diffs, test_data.UpdatedExpectedStorageDiff2)).To(BeTrue())

			mockFetcher.FetchErrs = map[uint64]error{
				test_data.BlockNumber.Uint64():  errors.New("mock fetcher error"),
				test_data.BlockNumber2.Uint64(): errors.New("mock fetcher error"),
			}
			mockFetcher.CalledTimes = 0
			backFiller = storage.NewStorageBackFiller(mockFetcher, test_data.BlockNumber.Uint64(), 1)
			backFill = make(chan utils.StorageDiff)
			done = make(chan bool)
			errChan = make(chan error)
			backFillInitErr = backFiller.BackFill(
				test_data.BlockNumber2.Uint64(),
				backFill,
				errChan,
				done)
			Expect(backFillInitErr).ToNot(HaveOccurred())
			numOfErrs = 0
			diffs = []utils.StorageDiff{}
			for {
				select {
				case diff := <-backFill:
					diffs = append(diffs, diff)
					continue
				case err := <-errChan:
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("mock fetcher error"))
					numOfErrs++
					continue
				case <-done:
					break
				}
				break
			}
			Expect(mockFetcher.CalledTimes).To(Equal(int64(2)))
			Expect(numOfErrs).To(Equal(2))
			Expect(len(diffs)).To(Equal(0))
		})
	})
})

func containsDiff(diffs []utils.StorageDiff, diff utils.StorageDiff) bool {
	for _, d := range diffs {
		if d == diff {
			return true
		}
	}
	return false
}
