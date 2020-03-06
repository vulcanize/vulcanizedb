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

package repositories_test

import (
	"math/rand"

	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	"github.com/makerdao/vulcanizedb/test_config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Checked Headers repository", func() {
	var (
		db   = test_config.NewTestDB(test_config.NewTestNode())
		repo datastore.CheckedHeadersRepository
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repo = repositories.NewCheckedHeadersRepository(db)
	})

	Describe("MarkHeaderChecked", func() {
		It("marks passed header as checked on insert", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, headerErr := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(headerErr).NotTo(HaveOccurred())

			err := repo.MarkHeaderChecked(headerID)

			Expect(err).NotTo(HaveOccurred())
			var checkedCount int
			fetchErr := db.Get(&checkedCount, `SELECT check_count FROM public.headers WHERE id = $1`, headerID)
			Expect(fetchErr).NotTo(HaveOccurred())
			Expect(checkedCount).To(Equal(1))
		})

		It("increments check count on update", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, headerErr := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(headerErr).NotTo(HaveOccurred())

			insertErr := repo.MarkHeaderChecked(headerID)
			Expect(insertErr).NotTo(HaveOccurred())

			updateErr := repo.MarkHeaderChecked(headerID)
			Expect(updateErr).NotTo(HaveOccurred())

			var checkedCount int
			fetchErr := db.Get(&checkedCount, `SELECT check_count FROM public.headers WHERE id = $1`, headerID)
			Expect(fetchErr).NotTo(HaveOccurred())
			Expect(checkedCount).To(Equal(2))
		})
	})

	Describe("MarkHeadersUncheckedSince", func() {
		It("marks headers with matching block number as unchecked", func() {
			blockNumberOne := rand.Int63()
			blockNumberTwo := blockNumberOne + 1
			blockNumberThree := blockNumberOne + 2
			fakeHeaderOne := fakes.GetFakeHeader(blockNumberOne)
			fakeHeaderTwo := fakes.GetFakeHeader(blockNumberTwo)
			fakeHeaderThree := fakes.GetFakeHeader(blockNumberThree)
			headerRepository := repositories.NewHeaderRepository(db)
			// insert three headers with incrementing block number
			headerIdOne, insertHeaderOneErr := headerRepository.CreateOrUpdateHeader(fakeHeaderOne)
			Expect(insertHeaderOneErr).NotTo(HaveOccurred())
			headerIdTwo, insertHeaderTwoErr := headerRepository.CreateOrUpdateHeader(fakeHeaderTwo)
			Expect(insertHeaderTwoErr).NotTo(HaveOccurred())
			headerIdThree, insertHeaderThreeErr := headerRepository.CreateOrUpdateHeader(fakeHeaderThree)
			Expect(insertHeaderThreeErr).NotTo(HaveOccurred())
			// mark all headers checked
			markHeaderOneCheckedErr := repo.MarkHeaderChecked(headerIdOne)
			Expect(markHeaderOneCheckedErr).NotTo(HaveOccurred())
			markHeaderTwoCheckedErr := repo.MarkHeaderChecked(headerIdTwo)
			Expect(markHeaderTwoCheckedErr).NotTo(HaveOccurred())
			markHeaderThreeCheckedErr := repo.MarkHeaderChecked(headerIdThree)
			Expect(markHeaderThreeCheckedErr).NotTo(HaveOccurred())

			// mark headers unchecked since blockNumberTwo
			err := repo.MarkHeadersUncheckedSince(blockNumberTwo)

			Expect(err).NotTo(HaveOccurred())
			var headerOneCheckCount, headerTwoCheckCount, headerThreeCheckCount int
			getHeaderOneErr := db.Get(&headerOneCheckCount, `SELECT check_count FROM public.headers WHERE id = $1`, headerIdOne)
			Expect(getHeaderOneErr).NotTo(HaveOccurred())
			Expect(headerOneCheckCount).To(Equal(1))
			getHeaderTwoErr := db.Get(&headerTwoCheckCount, `SELECT check_count FROM public.headers WHERE id = $1`, headerIdTwo)
			Expect(getHeaderTwoErr).NotTo(HaveOccurred())
			Expect(headerTwoCheckCount).To(BeZero())
			getHeaderThreeErr := db.Get(&headerThreeCheckCount, `SELECT check_count FROM public.headers WHERE id = $1`, headerIdThree)
			Expect(getHeaderThreeErr).NotTo(HaveOccurred())
			Expect(headerThreeCheckCount).To(BeZero())
		})
	})

	Describe("MarkSingleHeaderUnchecked", func() {
		It("marks headers with matching block number as unchecked", func() {
			blockNumberOne := rand.Int63()
			blockNumberTwo := blockNumberOne + 1
			blockNumberThree := blockNumberOne + 2
			fakeHeaderOne := fakes.GetFakeHeader(blockNumberOne)
			fakeHeaderTwo := fakes.GetFakeHeader(blockNumberTwo)
			fakeHeaderThree := fakes.GetFakeHeader(blockNumberThree)
			headerRepository := repositories.NewHeaderRepository(db)
			// insert three headers with incrementing block number
			headerIdOne, insertHeaderOneErr := headerRepository.CreateOrUpdateHeader(fakeHeaderOne)
			Expect(insertHeaderOneErr).NotTo(HaveOccurred())
			headerIdTwo, insertHeaderTwoErr := headerRepository.CreateOrUpdateHeader(fakeHeaderTwo)
			Expect(insertHeaderTwoErr).NotTo(HaveOccurred())
			headerIdThree, insertHeaderThreeErr := headerRepository.CreateOrUpdateHeader(fakeHeaderThree)
			Expect(insertHeaderThreeErr).NotTo(HaveOccurred())
			// mark all headers checked
			markHeaderOneCheckedErr := repo.MarkHeaderChecked(headerIdOne)
			Expect(markHeaderOneCheckedErr).NotTo(HaveOccurred())
			markHeaderTwoCheckedErr := repo.MarkHeaderChecked(headerIdTwo)
			Expect(markHeaderTwoCheckedErr).NotTo(HaveOccurred())
			markHeaderThreeCheckedErr := repo.MarkHeaderChecked(headerIdThree)
			Expect(markHeaderThreeCheckedErr).NotTo(HaveOccurred())

			// mark header from blockNumberTwo unchecked
			err := repo.MarkSingleHeaderUnchecked(blockNumberTwo)

			Expect(err).NotTo(HaveOccurred())
			var headerOneCheckCount, headerTwoCheckCount, headerThreeCheckCount int
			getHeaderOneErr := db.Get(&headerOneCheckCount, `SELECT check_count FROM public.headers WHERE id = $1`, headerIdOne)
			Expect(getHeaderOneErr).NotTo(HaveOccurred())
			Expect(headerOneCheckCount).To(Equal(1))
			getHeaderTwoErr := db.Get(&headerTwoCheckCount, `SELECT check_count FROM public.headers WHERE id = $1`, headerIdTwo)
			Expect(getHeaderTwoErr).NotTo(HaveOccurred())
			Expect(headerTwoCheckCount).To(BeZero())
			getHeaderThreeErr := db.Get(&headerThreeCheckCount, `SELECT check_count FROM public.headers WHERE id = $1`, headerIdThree)
			Expect(getHeaderThreeErr).NotTo(HaveOccurred())
			Expect(headerThreeCheckCount).To(Equal(1))
		})
	})

	Describe("UncheckedHeaders", func() {
		var (
			headerRepository datastore.HeaderRepository
			firstBlock,
			secondBlock,
			thirdBlock,
			lastBlock,
			secondHeaderID,
			thirdHeaderID int64
			blockNumbers        []int64
			headerIDs           []int64
			err                 error
			uncheckedCheckCount = int64(1)
			recheckCheckCount   = int64(2)
		)

		BeforeEach(func() {
			headerRepository = repositories.NewHeaderRepository(db)

			lastBlock = rand.Int63()
			thirdBlock = lastBlock - 15
			secondBlock = lastBlock - (15 + 30)
			firstBlock = lastBlock - (15 + 30 + 45)

			blockNumbers = []int64{firstBlock, secondBlock, thirdBlock, lastBlock}

			headerIDs = []int64{}
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				headerIDs = append(headerIDs, headerID)
				Expect(err).NotTo(HaveOccurred())
			}
			secondHeaderID = headerIDs[1]
			thirdHeaderID = headerIDs[2]
		})

		Describe("when ending block is specified", func() {
			It("excludes headers that are out of range", func() {
				headers, err := repo.UncheckedHeaders(firstBlock, thirdBlock, uncheckedCheckCount)
				Expect(err).NotTo(HaveOccurred())

				headerBlockNumbers := getBlockNumbers(headers)
				Expect(headerBlockNumbers).To(ConsistOf(firstBlock, secondBlock, thirdBlock))
				Expect(headerBlockNumbers).NotTo(ContainElement(lastBlock))
			})

			It("excludes headers that have been checked more than the check count", func() {
				_, err = db.Exec(`UPDATE public.headers SET check_count = 1 WHERE id = $1`, secondHeaderID)
				Expect(err).NotTo(HaveOccurred())

				headers, err := repo.UncheckedHeaders(firstBlock, thirdBlock, uncheckedCheckCount)
				Expect(err).NotTo(HaveOccurred())

				headerBlockNumbers := getBlockNumbers(headers)
				Expect(headerBlockNumbers).To(ConsistOf(firstBlock, thirdBlock))
				Expect(headerBlockNumbers).NotTo(ContainElement(secondBlock))
			})

			Describe("when header has already been checked", func() {
				It("includes header with block number > 15 back from latest with check count of 1", func() {
					_, err = db.Exec(`UPDATE public.headers SET check_count = 1 WHERE id = $1`, thirdHeaderID)
					Expect(err).NotTo(HaveOccurred())

					headers, err := repo.UncheckedHeaders(firstBlock, lastBlock, recheckCheckCount)
					Expect(err).NotTo(HaveOccurred())

					headerBlockNumbers := getBlockNumbers(headers)
					Expect(headerBlockNumbers).To(ContainElement(thirdBlock))
				})

				It("excludes header with block number < 15 back from latest with check count of 1", func() {
					excludedHeader := fakes.GetFakeHeader(thirdBlock + 1)
					excludedHeaderID, createHeaderErr := headerRepository.CreateOrUpdateHeader(excludedHeader)
					Expect(createHeaderErr).NotTo(HaveOccurred())
					_, updateHeaderErr := db.Exec(`UPDATE public.headers SET check_count = 1 WHERE id = $1`, excludedHeaderID)
					Expect(updateHeaderErr).NotTo(HaveOccurred())

					headers, err := repo.UncheckedHeaders(firstBlock, lastBlock, recheckCheckCount)
					Expect(err).NotTo(HaveOccurred())

					headerBlockNumbers := getBlockNumbers(headers)
					Expect(headerBlockNumbers).NotTo(ContainElement(excludedHeader.BlockNumber))
				})

				It("includes header with block number > 45 back from latest with check count of 2", func() {
					_, err = db.Exec(`UPDATE public.headers SET check_count = 1 WHERE id = $1`, secondHeaderID)
					Expect(err).NotTo(HaveOccurred())

					headers, err := repo.UncheckedHeaders(firstBlock, lastBlock, recheckCheckCount)
					Expect(err).NotTo(HaveOccurred())

					headerBlockNumbers := getBlockNumbers(headers)
					Expect(headerBlockNumbers).To(ContainElement(secondBlock))
				})

				It("excludes header with block number < 45 back from latest with check count of 2", func() {
					excludedHeader := fakes.GetFakeHeader(secondBlock + 1)
					excludedHeaderID, createHeaderErr := headerRepository.CreateOrUpdateHeader(excludedHeader)
					Expect(createHeaderErr).NotTo(HaveOccurred())
					_, updateHeaderErr := db.Exec(`UPDATE public.headers SET check_count = 2 WHERE id = $1`, excludedHeaderID)
					Expect(updateHeaderErr).NotTo(HaveOccurred())

					headers, err := repo.UncheckedHeaders(firstBlock, lastBlock, 3)
					Expect(err).NotTo(HaveOccurred())

					headerBlockNumbers := getBlockNumbers(headers)
					Expect(headerBlockNumbers).NotTo(ContainElement(excludedHeader.BlockNumber))
				})
			})

			It("only returns headers associated with any node", func() {
				dbTwo := test_config.NewTestDB(core.Node{ID: "second"})
				headerRepositoryTwo := repositories.NewHeaderRepository(dbTwo)
				repoTwo := repositories.NewCheckedHeadersRepository(dbTwo)
				for _, n := range blockNumbers {
					_, err = headerRepositoryTwo.CreateOrUpdateHeader(fakes.GetFakeHeader(n + 10))
					Expect(err).NotTo(HaveOccurred())
				}
				allHeaders := []int64{firstBlock, firstBlock + 10, secondBlock, secondBlock + 10, thirdBlock, thirdBlock + 10}

				nodeOneMissingHeaders, err := repo.UncheckedHeaders(firstBlock, thirdBlock+10, uncheckedCheckCount)
				Expect(err).NotTo(HaveOccurred())
				nodeOneHeaderBlockNumbers := getBlockNumbers(nodeOneMissingHeaders)
				Expect(nodeOneHeaderBlockNumbers).To(ConsistOf(allHeaders))

				nodeTwoMissingHeaders, err := repoTwo.UncheckedHeaders(firstBlock, thirdBlock+10, uncheckedCheckCount)
				Expect(err).NotTo(HaveOccurred())
				nodeTwoHeaderBlockNumbers := getBlockNumbers(nodeTwoMissingHeaders)
				Expect(nodeTwoHeaderBlockNumbers).To(ConsistOf(allHeaders))
			})
		})

		Describe("when ending block is -1", func() {
			It("includes all non-checked headers when ending block is -1 ", func() {
				headers, err := repo.UncheckedHeaders(firstBlock, -1, uncheckedCheckCount)
				Expect(err).NotTo(HaveOccurred())

				headerBlockNumbers := getBlockNumbers(headers)
				Expect(headerBlockNumbers).To(ConsistOf(firstBlock, secondBlock, thirdBlock, lastBlock))
			})

			It("excludes headers that have been checked more than the check count", func() {
				_, err = db.Exec(`UPDATE public.headers SET check_count = 1 WHERE id = $1`, headerIDs[1])
				Expect(err).NotTo(HaveOccurred())

				headers, err := repo.UncheckedHeaders(firstBlock, -1, uncheckedCheckCount)
				Expect(err).NotTo(HaveOccurred())

				headerBlockNumbers := getBlockNumbers(headers)
				Expect(headerBlockNumbers).To(ConsistOf(firstBlock, thirdBlock, lastBlock))
				Expect(headerBlockNumbers).NotTo(ContainElement(secondBlock))
			})

			Describe("when header has already been checked", func() {
				It("includes header with block number > 15 back from latest with check count of 1", func() {
					_, err = db.Exec(`UPDATE public.headers SET check_count = 1 WHERE id = $1`, thirdHeaderID)
					Expect(err).NotTo(HaveOccurred())

					headers, err := repo.UncheckedHeaders(firstBlock, -1, recheckCheckCount)
					Expect(err).NotTo(HaveOccurred())

					headerBlockNumbers := getBlockNumbers(headers)
					Expect(headerBlockNumbers).To(ContainElement(thirdBlock))
				})

				It("excludes header with block number < 15 back from latest with check count of 1", func() {
					excludedHeader := fakes.GetFakeHeader(thirdBlock + 1)
					excludedHeaderID, createHeaderErr := headerRepository.CreateOrUpdateHeader(excludedHeader)
					Expect(createHeaderErr).NotTo(HaveOccurred())
					_, updateHeaderErr := db.Exec(`UPDATE public.headers SET check_count = 1 WHERE id = $1`, excludedHeaderID)
					Expect(updateHeaderErr).NotTo(HaveOccurred())

					headers, err := repo.UncheckedHeaders(firstBlock, -1, recheckCheckCount)
					Expect(err).NotTo(HaveOccurred())

					headerBlockNumbers := getBlockNumbers(headers)
					Expect(headerBlockNumbers).NotTo(ContainElement(excludedHeader.BlockNumber))
				})

				It("includes header with block number > 45 back from latest with check count of 2", func() {
					_, err = db.Exec(`UPDATE public.headers SET check_count = 1 WHERE id = $1`, secondHeaderID)
					Expect(err).NotTo(HaveOccurred())

					headers, err := repo.UncheckedHeaders(firstBlock, -1, recheckCheckCount)
					Expect(err).NotTo(HaveOccurred())

					headerBlockNumbers := getBlockNumbers(headers)
					Expect(headerBlockNumbers).To(ContainElement(secondBlock))
				})

				It("excludes header with block number < 45 back from latest with check count of 2", func() {
					excludedHeader := fakes.GetFakeHeader(secondBlock + 1)
					excludedHeaderID, createHeaderErr := headerRepository.CreateOrUpdateHeader(excludedHeader)
					Expect(createHeaderErr).NotTo(HaveOccurred())
					_, updateHeaderErr := db.Exec(`UPDATE public.headers SET check_count = 2 WHERE id = $1`, excludedHeaderID)
					Expect(updateHeaderErr).NotTo(HaveOccurred())

					headers, err := repo.UncheckedHeaders(firstBlock, -1, 3)
					Expect(err).NotTo(HaveOccurred())

					headerBlockNumbers := getBlockNumbers(headers)
					Expect(headerBlockNumbers).NotTo(ContainElement(excludedHeader.BlockNumber))
				})
			})

			It("returns headers associated with any node", func() {
				dbTwo := test_config.NewTestDB(core.Node{ID: "second"})
				headerRepositoryTwo := repositories.NewHeaderRepository(dbTwo)
				repoTwo := repositories.NewCheckedHeadersRepository(dbTwo)
				for _, n := range blockNumbers {
					_, err = headerRepositoryTwo.CreateOrUpdateHeader(fakes.GetFakeHeader(n + 10))
					Expect(err).NotTo(HaveOccurred())
				}
				allHeaders := []int64{firstBlock, firstBlock + 10, secondBlock, secondBlock + 10, thirdBlock, thirdBlock + 10, lastBlock, lastBlock + 10}

				nodeOneMissingHeaders, err := repo.UncheckedHeaders(firstBlock, -1, uncheckedCheckCount)
				Expect(err).NotTo(HaveOccurred())
				nodeOneBlockNumbers := getBlockNumbers(nodeOneMissingHeaders)
				Expect(nodeOneBlockNumbers).To(ConsistOf(allHeaders))

				nodeTwoMissingHeaders, err := repoTwo.UncheckedHeaders(firstBlock, -1, uncheckedCheckCount)
				Expect(err).NotTo(HaveOccurred())
				nodeTwoBlockNumbers := getBlockNumbers(nodeTwoMissingHeaders)
				Expect(nodeTwoBlockNumbers).To(ConsistOf(allHeaders))
			})
		})
	})
})

func getBlockNumbers(headers []core.Header) []int64 {
	var headerBlockNumbers []int64
	for _, header := range headers {
		headerBlockNumbers = append(headerBlockNumbers, header.BlockNumber)
	}
	return headerBlockNumbers
}
