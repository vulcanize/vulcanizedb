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

package repository_test

import (
	"fmt"
	"github.com/vulcanize/vulcanizedb/pkg/core"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/omni/light/repository"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers/test_helpers/mocks"
)

var _ = Describe("Repository", func() {
	var db *postgres.DB
	var omniHeaderRepo repository.HeaderRepository   // omni/light header repository
	var coreHeaderRepo repositories.HeaderRepository // pkg/datastore header repository
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
		omniHeaderRepo = repository.NewHeaderRepository(db)
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

			err = omniHeaderRepo.AddCheckColumn(eventIDs[0])
			Expect(err).ToNot(HaveOccurred())

			_, err = db.Exec(query)
			Expect(err).ToNot(HaveOccurred())
		})

		It("Caches column it creates so that it does not need to repeatedly query the database to check for it's existence", func() {
			_, ok := omniHeaderRepo.CheckCache(eventIDs[0])
			Expect(ok).To(Equal(false))

			err := omniHeaderRepo.AddCheckColumn(eventIDs[0])
			Expect(err).ToNot(HaveOccurred())

			v, ok := omniHeaderRepo.CheckCache(eventIDs[0])
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

			err := omniHeaderRepo.AddCheckColumns(eventIDs)
			Expect(err).ToNot(HaveOccurred())

			for _, id := range eventIDs {
				_, err := db.Exec(fmt.Sprintf("SELECT %s FROM checked_headers", id))
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("Caches columns it creates so that it does not need to repeatedly query the database to check for it's existence", func() {
			for _, id := range eventIDs {
				_, ok := omniHeaderRepo.CheckCache(id)
				Expect(ok).To(Equal(false))
			}

			err := omniHeaderRepo.AddCheckColumns(eventIDs)
			Expect(err).ToNot(HaveOccurred())

			for _, id := range eventIDs {
				v, ok := omniHeaderRepo.CheckCache(id)
				Expect(ok).To(Equal(true))
				Expect(v).To(Equal(true))
			}
		})
	})

	Describe("MissingHeaders", func() {
		It("Returns all unchecked headers for the given eventID", func() {
			addHeaders(coreHeaderRepo)
			err := omniHeaderRepo.AddCheckColumn(eventIDs[0])
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err := omniHeaderRepo.MissingHeaders(6194630, 6194635, eventIDs[0])
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(3))
		})

		It("Returns unchecked headers in ascending order", func() {
			addHeaders(coreHeaderRepo)
			err := omniHeaderRepo.AddCheckColumn(eventIDs[0])
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err := omniHeaderRepo.MissingHeaders(6194630, 6194635, eventIDs[0])
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(3))

			h1 := missingHeaders[0]
			h2 := missingHeaders[1]
			h3 := missingHeaders[2]
			Expect(h1.BlockNumber).To(Equal(int64(6194632)))
			Expect(h2.BlockNumber).To(Equal(int64(6194633)))
			Expect(h3.BlockNumber).To(Equal(int64(6194634)))
		})

		It("Fails if eventID does not yet exist in check_headers table", func() {
			addHeaders(coreHeaderRepo)
			err := omniHeaderRepo.AddCheckColumn(eventIDs[0])
			Expect(err).ToNot(HaveOccurred())

			_, err = omniHeaderRepo.MissingHeaders(6194630, 6194635, "notEventId")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("MissingHeadersForAll", func() { // HERE
		It("Returns all headers that have not been checked for all of the ids provided", func() {
			addHeaders(coreHeaderRepo)
			err := omniHeaderRepo.AddCheckColumns(eventIDs)
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err := omniHeaderRepo.MissingHeadersForAll(6194630, 6194635, eventIDs)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(3))

			err = omniHeaderRepo.MarkHeaderChecked(missingHeaders[0].Id, eventIDs[0])
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err = omniHeaderRepo.MissingHeadersForAll(6194630, 6194635, eventIDs)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(3))

			err = omniHeaderRepo.MarkHeaderChecked(missingHeaders[0].Id, eventIDs[1])
			Expect(err).ToNot(HaveOccurred())
			err = omniHeaderRepo.MarkHeaderChecked(missingHeaders[0].Id, eventIDs[2])
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err = omniHeaderRepo.MissingHeadersForAll(6194630, 6194635, eventIDs)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(2))
		})

		It("Fails if one of the eventIDs does not yet exist in check_headers table", func() {
			addHeaders(coreHeaderRepo)
			err := omniHeaderRepo.AddCheckColumns(eventIDs)
			Expect(err).ToNot(HaveOccurred())
			badEventIDs := append(eventIDs, "notEventId")

			_, err = omniHeaderRepo.MissingHeadersForAll(6194630, 6194635, badEventIDs)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("MarkHeaderChecked", func() {
		It("Marks the header checked for the given eventID", func() {
			addHeaders(coreHeaderRepo)
			err := omniHeaderRepo.AddCheckColumn(eventIDs[0])
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err := omniHeaderRepo.MissingHeaders(6194630, 6194635, eventIDs[0])
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(3))

			headerID := missingHeaders[0].Id
			err = omniHeaderRepo.MarkHeaderChecked(headerID, eventIDs[0])
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err = omniHeaderRepo.MissingHeaders(6194630, 6194635, eventIDs[0])
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(2))
		})

		It("Fails if eventID does not yet exist in check_headers table", func() {
			addHeaders(coreHeaderRepo)
			err := omniHeaderRepo.AddCheckColumn(eventIDs[0])
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err := omniHeaderRepo.MissingHeaders(6194630, 6194635, eventIDs[0])
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(3))

			headerID := missingHeaders[0].Id
			err = omniHeaderRepo.MarkHeaderChecked(headerID, "notEventId")
			Expect(err).To(HaveOccurred())

			missingHeaders, err = omniHeaderRepo.MissingHeaders(6194630, 6194635, eventIDs[0])
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(3))
		})
	})

	Describe("MarkHeaderCheckedForAll", func() {
		It("Marks the header checked for all provided column ids", func() {
			addHeaders(coreHeaderRepo)
			err := omniHeaderRepo.AddCheckColumns(eventIDs)
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err := omniHeaderRepo.MissingHeadersForAll(6194630, 6194635, eventIDs)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(3))

			headerID := missingHeaders[0].Id
			err = omniHeaderRepo.MarkHeaderCheckedForAll(headerID, eventIDs)
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err = omniHeaderRepo.MissingHeaders(6194630, 6194635, eventIDs[0])
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
				err := omniHeaderRepo.AddCheckColumn(id)
				Expect(err).ToNot(HaveOccurred())
				missingHeaders, err = omniHeaderRepo.MissingHeaders(6194630, 6194635, id)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(missingHeaders)).To(Equal(3))
			}

			err := omniHeaderRepo.MarkHeadersCheckedForAll(missingHeaders, methodIDs)
			Expect(err).ToNot(HaveOccurred())
			for _, id := range methodIDs {
				missingHeaders, err = omniHeaderRepo.MissingHeaders(6194630, 6194635, id)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(missingHeaders)).To(Equal(0))
			}
		})
	})

	Describe("MissingMethodsCheckedEventsIntersection", func() {
		It("Returns headers that have been checked for all the provided events but have not been checked for all the provided methods", func() {
			addHeaders(coreHeaderRepo)
			for i, id := range eventIDs {
				err := omniHeaderRepo.AddCheckColumn(id)
				Expect(err).ToNot(HaveOccurred())
				err = omniHeaderRepo.AddCheckColumn(methodIDs[i])
				Expect(err).ToNot(HaveOccurred())
			}

			missingHeaders, err := omniHeaderRepo.MissingHeaders(6194630, 6194635, eventIDs[0])
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(3))

			headerID := missingHeaders[0].Id
			headerID2 := missingHeaders[1].Id
			for i, id := range eventIDs {
				err = omniHeaderRepo.MarkHeaderChecked(headerID, id)
				Expect(err).ToNot(HaveOccurred())
				err = omniHeaderRepo.MarkHeaderChecked(headerID2, id)
				Expect(err).ToNot(HaveOccurred())
				err = omniHeaderRepo.MarkHeaderChecked(headerID, methodIDs[i])
				Expect(err).ToNot(HaveOccurred())
			}

			intersectionHeaders, err := omniHeaderRepo.MissingMethodsCheckedEventsIntersection(6194630, 6194635, methodIDs, eventIDs)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(intersectionHeaders)).To(Equal(1))
			Expect(intersectionHeaders[0].Id).To(Equal(headerID2))

		})
	})
})

func addHeaders(coreHeaderRepo repositories.HeaderRepository) {
	coreHeaderRepo.CreateOrUpdateHeader(mocks.MockHeader1)
	coreHeaderRepo.CreateOrUpdateHeader(mocks.MockHeader2)
	coreHeaderRepo.CreateOrUpdateHeader(mocks.MockHeader3)
}
