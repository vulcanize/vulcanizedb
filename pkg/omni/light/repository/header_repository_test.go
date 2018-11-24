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
	var r repository.HeaderRepository
	var headerRepository repositories.HeaderRepository
	var eventID, query string

	BeforeEach(func() {
		db, _ = test_helpers.SetupDBandBC()
		r = repository.NewHeaderRepository(db)
		headerRepository = repositories.NewHeaderRepository(db)
		eventID = "eventName_contractAddr"
	})

	AfterEach(func() {
		test_helpers.TearDown(db)
	})

	Describe("AddCheckColumn", func() {
		It("Creates a column for the given eventID to mark if the header has been checked for that event", func() {
			query = fmt.Sprintf("SELECT %s FROM checked_headers", eventID)
			_, err := db.Exec(query)
			Expect(err).To(HaveOccurred())

			err = r.AddCheckColumn(eventID)
			Expect(err).ToNot(HaveOccurred())

			_, err = db.Exec(query)
			Expect(err).ToNot(HaveOccurred())
		})

		It("Caches column it creates so that it does not need to repeatedly query the database to check for it's existence", func() {
			_, ok := r.CheckCache(eventID)
			Expect(ok).To(Equal(false))

			err := r.AddCheckColumn(eventID)
			Expect(err).ToNot(HaveOccurred())

			v, ok := r.CheckCache(eventID)
			Expect(ok).To(Equal(true))
			Expect(v).To(Equal(true))
		})
	})

	Describe("MissingHeaders", func() {
		It("Returns all unchecked headers for the given eventID", func() {
			headerRepository.CreateOrUpdateHeader(mocks.MockHeader1)
			headerRepository.CreateOrUpdateHeader(mocks.MockHeader2)
			headerRepository.CreateOrUpdateHeader(mocks.MockHeader3)
			err := r.AddCheckColumn(eventID)
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err := r.MissingHeaders(6194630, 6194635, eventID)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(3))
		})

		It("Fails if eventID does not yet exist in check_headers table", func() {
			headerRepository.CreateOrUpdateHeader(mocks.MockHeader1)
			headerRepository.CreateOrUpdateHeader(mocks.MockHeader2)
			headerRepository.CreateOrUpdateHeader(mocks.MockHeader3)
			err := r.AddCheckColumn(eventID)
			Expect(err).ToNot(HaveOccurred())

			_, err = r.MissingHeaders(6194630, 6194635, "notEventId")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("MarkHeaderChecked", func() {
		It("Marks the header checked for the given eventID", func() {
			headerRepository.CreateOrUpdateHeader(mocks.MockHeader1)
			headerRepository.CreateOrUpdateHeader(mocks.MockHeader2)
			headerRepository.CreateOrUpdateHeader(mocks.MockHeader3)
			err := r.AddCheckColumn(eventID)
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err := r.MissingHeaders(6194630, 6194635, eventID)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(3))

			headerID := missingHeaders[0].Id
			err = r.MarkHeaderChecked(headerID, eventID)
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err = r.MissingHeaders(6194630, 6194635, eventID)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(2))
		})

		It("Fails if eventID does not yet exist in check_headers table", func() {
			headerRepository.CreateOrUpdateHeader(mocks.MockHeader1)
			headerRepository.CreateOrUpdateHeader(mocks.MockHeader2)
			headerRepository.CreateOrUpdateHeader(mocks.MockHeader3)
			err := r.AddCheckColumn(eventID)
			Expect(err).ToNot(HaveOccurred())

			missingHeaders, err := r.MissingHeaders(6194630, 6194635, eventID)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(3))

			headerID := missingHeaders[0].Id
			err = r.MarkHeaderChecked(headerID, "notEventId")
			Expect(err).To(HaveOccurred())

			missingHeaders, err = r.MissingHeaders(6194630, 6194635, eventID)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(missingHeaders)).To(Equal(3))
		})
	})
})
