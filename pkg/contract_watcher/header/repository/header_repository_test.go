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

package repository_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/header/repository"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/helpers/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/helpers/test_helpers/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
)

var _ = Describe("Repository", func() {
	var db *postgres.DB
	var contractHeaderRepo repository.HeaderRepository // contract_watcher headerSync header repository
	var coreHeaderRepo repositories.HeaderRepository   // pkg/datastore header repository
	var eventIDs = []string{
		"eventName_contractAddr",
		"eventName_contractAddr2",
		"eventName_contractAddr3",
	}
	var methodIDs = []string{
		"methodName_contractAddr",
		"methodName_contractAddr2",
		"methodName_contractAddr3",
	}

	BeforeEach(func() {
		db, _ = test_helpers.SetupDBandBC()
		contractHeaderRepo = repository.NewHeaderRepository(db)
		coreHeaderRepo = repositories.NewHeaderRepository(db)
	})

	AfterEach(func() {
		test_helpers.TearDown(db)
	})

	Describe("AddCheckColumn", func() {
		It("Creates a column for the given eventID to mark if the header has been checked for that event", func() {
			query := fmt.Sprintf("SELECT %s FROM checked_headers", eventIDs[0])
			_, err := db.Exec(query)
			Expect(err).To(HaveOccurred())

			err = contractHeaderRepo.AddCheckColumn(eventIDs[0])
			Expect(err).ToNot(HaveOccurred())

			_, err = db.Exec(query)
			Expect(err).ToNot(HaveOccurred())
		})

		It("Caches column it creates so that it does not need to repeatedly query the database to check for it's existence", func() {
			_, ok := contractHeaderRepo.CheckCache(eventIDs[0])
			Expect(ok).To(Equal(false))

			err := contractHeaderRepo.AddCheckColumn(eventIDs[0])
			Expect(err).ToNot(HaveOccurred())

			v, ok := contractHeaderRepo.CheckCache(eventIDs[0])
			Expect(ok).To(Equal(true))
			Expect(v).To(Equal(true))
		})
	})

	Describe("AddCheckColumns", func() {
		It("Creates a column for the given eventIDs to mark if the header has been checked for those events", func() {
			for _, id := range eventIDs {
				_, err := db.Exec(fmt.Sprintf("SELECT %s FROM checked_headers", id))
				Expect(err).To(HaveOccurred())
			}

			err := contractHeaderRepo.AddCheckColumns(eventIDs)
			Expect(err).ToNot(HaveOccurred())

			for _, id := range eventIDs {
				_, err := db.Exec(fmt.Sprintf("SELECT %s FROM checked_headers", id))
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("Caches columns it creates so that it does not need to repeatedly query the database to check for it's existence", func() {
			for _, id := range eventIDs {
				_, ok := contractHeaderRepo.CheckCache(id)
				Expect(ok).To(Equal(false))
			}

			err := contractHeaderRepo.AddCheckColumns(eventIDs)
			Expect(err).ToNot(HaveOccurred())

			for _, id := range eventIDs {
				v, ok := contractHeaderRepo.CheckCache(id)
				Expect(ok).To(Equal(true))
				Expect(v).To(Equal(true))
			}
		})
	})

	Describe("MissingHeaders", func() {
		It("Returns all unchecked headers for the given eventID", func() {
			addHeaders(coreHeaderRepo)
			err := contractHeaderRepo.AddCheckColumn(eventIDs[0])
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err := contractHeaderRepo.MissingHeaders(mocks.MockHeader1.BlockNumber, mocks.MockHeader4.BlockNumber, eventIDs[0])
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(3))
		})

		It("Returns unchecked headers in ascending order", func() {
			addHeaders(coreHeaderRepo)
			err := contractHeaderRepo.AddCheckColumn(eventIDs[0])
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err := contractHeaderRepo.MissingHeaders(mocks.MockHeader1.BlockNumber, mocks.MockHeader4.BlockNumber, eventIDs[0])
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(3))

			h1 := missingHeaders[0]
			h2 := missingHeaders[1]
			h3 := missingHeaders[2]
			Expect(h1.BlockNumber).To(Equal(mocks.MockHeader1.BlockNumber))
			Expect(h2.BlockNumber).To(Equal(mocks.MockHeader2.BlockNumber))
			Expect(h3.BlockNumber).To(Equal(mocks.MockHeader3.BlockNumber))
		})

		It("Returns only contiguous chunks of headers", func() {
			addDiscontinuousHeaders(coreHeaderRepo)
			err := contractHeaderRepo.AddCheckColumns(eventIDs)
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err := contractHeaderRepo.MissingHeaders(mocks.MockHeader1.BlockNumber, mocks.MockHeader4.BlockNumber, eventIDs[0])
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(2))
			Expect(missingHeaders[0].BlockNumber).To(Equal(mocks.MockHeader1.BlockNumber))
			Expect(missingHeaders[1].BlockNumber).To(Equal(mocks.MockHeader2.BlockNumber))
		})

		It("Fails if eventID does not yet exist in check_headers table", func() {
			addHeaders(coreHeaderRepo)
			err := contractHeaderRepo.AddCheckColumn(eventIDs[0])
			Expect(err).ToNot(HaveOccurred())

			_, err = contractHeaderRepo.MissingHeaders(mocks.MockHeader1.BlockNumber, mocks.MockHeader4.BlockNumber, "notEventId")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("MissingHeadersForAll", func() { // HERE
		It("Returns all headers that have not been checked for all of the ids provided", func() {
			addHeaders(coreHeaderRepo)
			err := contractHeaderRepo.AddCheckColumns(eventIDs)
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err := contractHeaderRepo.MissingHeadersForAll(mocks.MockHeader1.BlockNumber, mocks.MockHeader4.BlockNumber, eventIDs)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(3))

			err = contractHeaderRepo.MarkHeaderChecked(missingHeaders[0].ID, eventIDs[0])
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err = contractHeaderRepo.MissingHeadersForAll(mocks.MockHeader1.BlockNumber, mocks.MockHeader4.BlockNumber, eventIDs)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(3))

			err = contractHeaderRepo.MarkHeaderChecked(missingHeaders[0].ID, eventIDs[1])
			Expect(err).ToNot(HaveOccurred())
			err = contractHeaderRepo.MarkHeaderChecked(missingHeaders[0].ID, eventIDs[2])
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err = contractHeaderRepo.MissingHeadersForAll(mocks.MockHeader2.BlockNumber, mocks.MockHeader4.BlockNumber, eventIDs)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(2))
		})

		It("Returns only contiguous chunks of headers", func() {
			addDiscontinuousHeaders(coreHeaderRepo)
			err := contractHeaderRepo.AddCheckColumns(eventIDs)
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err := contractHeaderRepo.MissingHeadersForAll(mocks.MockHeader1.BlockNumber, mocks.MockHeader4.BlockNumber, eventIDs)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(2))
			Expect(missingHeaders[0].BlockNumber).To(Equal(mocks.MockHeader1.BlockNumber))
			Expect(missingHeaders[1].BlockNumber).To(Equal(mocks.MockHeader2.BlockNumber))
		})

		It("returns headers after starting header if starting header not missing", func() {
			addLaterHeaders(coreHeaderRepo)
			err := contractHeaderRepo.AddCheckColumns(eventIDs)
			Expect(err).NotTo(HaveOccurred())

			missingHeaders, err := contractHeaderRepo.MissingHeadersForAll(mocks.MockHeader1.BlockNumber, -1, eventIDs)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(2))
			Expect(missingHeaders[0].BlockNumber).To(Equal(mocks.MockHeader3.BlockNumber))
			Expect(missingHeaders[1].BlockNumber).To(Equal(mocks.MockHeader4.BlockNumber))
		})

		It("Fails if one of the eventIDs does not yet exist in check_headers table", func() {
			addHeaders(coreHeaderRepo)
			err := contractHeaderRepo.AddCheckColumns(eventIDs)
			Expect(err).ToNot(HaveOccurred())
			badEventIDs := append(eventIDs, "notEventId")

			_, err = contractHeaderRepo.MissingHeadersForAll(mocks.MockHeader1.BlockNumber, mocks.MockHeader4.BlockNumber, badEventIDs)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("MarkHeaderChecked", func() {
		It("Marks the header checked for the given eventID", func() {
			addHeaders(coreHeaderRepo)
			err := contractHeaderRepo.AddCheckColumn(eventIDs[0])
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err := contractHeaderRepo.MissingHeaders(mocks.MockHeader1.BlockNumber, mocks.MockHeader4.BlockNumber, eventIDs[0])
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(3))

			headerID := missingHeaders[0].ID
			err = contractHeaderRepo.MarkHeaderChecked(headerID, eventIDs[0])
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err = contractHeaderRepo.MissingHeaders(mocks.MockHeader2.BlockNumber, mocks.MockHeader4.BlockNumber, eventIDs[0])
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(2))
		})

		It("Fails if eventID does not yet exist in check_headers table", func() {
			addHeaders(coreHeaderRepo)
			err := contractHeaderRepo.AddCheckColumn(eventIDs[0])
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err := contractHeaderRepo.MissingHeaders(mocks.MockHeader1.BlockNumber, mocks.MockHeader4.BlockNumber, eventIDs[0])
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(3))

			headerID := missingHeaders[0].ID
			err = contractHeaderRepo.MarkHeaderChecked(headerID, "notEventId")
			Expect(err).To(HaveOccurred())

			missingHeaders, err = contractHeaderRepo.MissingHeaders(mocks.MockHeader1.BlockNumber, mocks.MockHeader4.BlockNumber, eventIDs[0])
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(3))
		})
	})

	Describe("MarkHeaderCheckedForAll", func() {
		It("Marks the header checked for all provided column ids", func() {
			addHeaders(coreHeaderRepo)
			err := contractHeaderRepo.AddCheckColumns(eventIDs)
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err := contractHeaderRepo.MissingHeadersForAll(mocks.MockHeader1.BlockNumber, mocks.MockHeader4.BlockNumber, eventIDs)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(3))

			headerID := missingHeaders[0].ID
			err = contractHeaderRepo.MarkHeaderCheckedForAll(headerID, eventIDs)
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err = contractHeaderRepo.MissingHeaders(mocks.MockHeader2.BlockNumber, mocks.MockHeader4.BlockNumber, eventIDs[0])
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(2))
		})
	})

	Describe("MarkHeadersCheckedForAll", func() {
		It("Marks the headers checked for all provided column ids", func() {
			addHeaders(coreHeaderRepo)
			methodIDs := []string{
				"methodName_contractAddr",
				"methodName_contractAddr2",
				"methodName_contractAddr3",
			}

			var missingHeaders []core.Header
			for _, id := range methodIDs {
				err := contractHeaderRepo.AddCheckColumn(id)
				Expect(err).ToNot(HaveOccurred())
				missingHeaders, err = contractHeaderRepo.MissingHeaders(mocks.MockHeader1.BlockNumber, mocks.MockHeader4.BlockNumber, id)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(missingHeaders)).To(Equal(3))
			}

			err := contractHeaderRepo.MarkHeadersCheckedForAll(missingHeaders, methodIDs)
			Expect(err).ToNot(HaveOccurred())
			for _, id := range methodIDs {
				missingHeaders, err = contractHeaderRepo.MissingHeaders(mocks.MockHeader1.BlockNumber, mocks.MockHeader4.BlockNumber, id)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(missingHeaders)).To(Equal(0))
			}
		})
	})

	Describe("MissingMethodsCheckedEventsIntersection", func() {
		It("Returns headers that have been checked for all the provided events but have not been checked for all the provided methods", func() {
			addHeaders(coreHeaderRepo)
			for i, id := range eventIDs {
				err := contractHeaderRepo.AddCheckColumn(id)
				Expect(err).ToNot(HaveOccurred())
				err = contractHeaderRepo.AddCheckColumn(methodIDs[i])
				Expect(err).ToNot(HaveOccurred())
			}

			missingHeaders, err := contractHeaderRepo.MissingHeaders(mocks.MockHeader1.BlockNumber, mocks.MockHeader4.BlockNumber, eventIDs[0])
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(3))

			headerID := missingHeaders[0].ID
			headerID2 := missingHeaders[1].ID
			for i, id := range eventIDs {
				err = contractHeaderRepo.MarkHeaderChecked(headerID, id)
				Expect(err).ToNot(HaveOccurred())
				err = contractHeaderRepo.MarkHeaderChecked(headerID2, id)
				Expect(err).ToNot(HaveOccurred())
				err = contractHeaderRepo.MarkHeaderChecked(headerID, methodIDs[i])
				Expect(err).ToNot(HaveOccurred())
			}

			intersectionHeaders, err := contractHeaderRepo.MissingMethodsCheckedEventsIntersection(mocks.MockHeader1.BlockNumber, mocks.MockHeader4.BlockNumber, methodIDs, eventIDs)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(intersectionHeaders)).To(Equal(1))
			Expect(intersectionHeaders[0].ID).To(Equal(headerID2))
		})
	})
})

func addHeaders(coreHeaderRepo repositories.HeaderRepository) {
	_, err := coreHeaderRepo.CreateOrUpdateHeader(mocks.MockHeader1)
	Expect(err).NotTo(HaveOccurred())
	_, err = coreHeaderRepo.CreateOrUpdateHeader(mocks.MockHeader2)
	Expect(err).NotTo(HaveOccurred())
	_, err = coreHeaderRepo.CreateOrUpdateHeader(mocks.MockHeader3)
	Expect(err).NotTo(HaveOccurred())
}

func addDiscontinuousHeaders(coreHeaderRepo repositories.HeaderRepository) {
	_, err := coreHeaderRepo.CreateOrUpdateHeader(mocks.MockHeader1)
	Expect(err).NotTo(HaveOccurred())
	_, err = coreHeaderRepo.CreateOrUpdateHeader(mocks.MockHeader2)
	Expect(err).NotTo(HaveOccurred())
	_, err = coreHeaderRepo.CreateOrUpdateHeader(mocks.MockHeader4)
	Expect(err).NotTo(HaveOccurred())
}

func addLaterHeaders(coreHeaderRepo repositories.HeaderRepository) {
	_, err := coreHeaderRepo.CreateOrUpdateHeader(mocks.MockHeader3)
	Expect(err).NotTo(HaveOccurred())
	_, err = coreHeaderRepo.CreateOrUpdateHeader(mocks.MockHeader4)
	Expect(err).NotTo(HaveOccurred())
}
