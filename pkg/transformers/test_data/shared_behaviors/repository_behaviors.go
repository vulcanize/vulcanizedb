// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shared_behaviors

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
	"math/rand"
)

var (
	db               *postgres.DB
	headerRepository datastore.HeaderRepository
	err              error
	headerId         int64
)

type CreateBehaviorInputs struct {
	CheckedHeaderColumnName  string
	LogEventTableName        string
	TestModel                interface{}
	ModelWithDifferentLogIdx interface{}
	Repository               factories.Repository
}

type MarkedHeaderCheckedBehaviorInputs struct {
	CheckedHeaderColumnName string
	Repository              factories.Repository
}

type MissingHeadersBehaviorInputs struct {
	Repository    factories.Repository
	RepositoryTwo factories.Repository
}

func SharedRepositoryCreateBehaviors(inputs *CreateBehaviorInputs) {
	Describe("Create", func() {
		var headerID int64
		var repository = inputs.Repository
		var checkedHeaderColumn = inputs.CheckedHeaderColumnName
		var logEventModel = inputs.TestModel

		BeforeEach(func() {
			headerRepository = GetHeaderRepository()
			headerID, err = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
		})

		It("marks header as checked for logs", func() {
			err = repository.Create(headerID, []interface{}{logEventModel})

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			query := `SELECT ` + checkedHeaderColumn + ` FROM public.checked_headers WHERE header_id = $1`
			err = db.Get(&headerChecked, query, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("updates the header to checked if checked headers row already exists", func() {
			_, err := db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerID)
			Expect(err).NotTo(HaveOccurred())

			err = repository.Create(headerID, []interface{}{logEventModel})

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			query := `SELECT ` + checkedHeaderColumn + ` FROM public.checked_headers WHERE header_id = $1`
			err = db.Get(&headerChecked, query, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("returns an error if inserting duplicate log events", func() {
			err = repository.Create(headerID, []interface{}{logEventModel})
			Expect(err).NotTo(HaveOccurred())

			err = repository.Create(headerID, []interface{}{logEventModel})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("allows for multiple log events of the same type in one transaction if they have different log indexes", func() {
			err = repository.Create(headerID, []interface{}{logEventModel})
			Expect(err).NotTo(HaveOccurred())

			err = repository.Create(headerID, []interface{}{inputs.ModelWithDifferentLogIdx})
			Expect(err).NotTo(HaveOccurred())
		})

		It("removes the log event record if the corresponding header is deleted", func() {
			err = repository.Create(headerID, []interface{}{logEventModel})
			Expect(err).NotTo(HaveOccurred())

			_, err = db.Exec(`DELETE FROM headers WHERE id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())

			var count int
			query := `SELECT count(*) from ` + inputs.LogEventTableName
			err = db.QueryRow(query).Scan(&count)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(0))
		})

		It("returns an error if model is of wrong type", func() {
			err = repository.Create(headerID, []interface{}{test_data.WrongModel{}})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("model of type"))
		})
	})
}

func SharedRepositoryMissingHeadersBehaviors(inputs *MissingHeadersBehaviorInputs) {
	Describe("MissingHeaders", func() {
		var (
			repository               = inputs.Repository
			startingBlockNumber      int64
			endingBlockNumber        int64
			eventSpecificBlockNumber int64
			blockNumbers             []int64
			headerIDs                []int64
		)

		BeforeEach(func() {
			headerRepository = GetHeaderRepository()
			startingBlockNumber = rand.Int63()
			eventSpecificBlockNumber = startingBlockNumber + 1
			endingBlockNumber = startingBlockNumber + 2
			outOfRangeBlockNumber := endingBlockNumber + 1

			blockNumbers = []int64{startingBlockNumber, eventSpecificBlockNumber, endingBlockNumber, outOfRangeBlockNumber}

			headerIDs = []int64{}
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				headerIDs = append(headerIDs, headerID)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("returns headers that haven't been checked", func() {
			err := repository.MarkHeaderChecked(headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := repository.MissingHeaders(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(2))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber)))
		})

		It("only treats headers as checked if the event specific logs have been checked", func() {
			_, err := db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := repository.MissingHeaders(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(3))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(eventSpecificBlockNumber)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(eventSpecificBlockNumber)))
			Expect(headers[2].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(eventSpecificBlockNumber)))
		})

		It("only returns headers associated with the current node", func() {
			dbTwo := test_config.NewTestDB(core.Node{ID: "second"})
			headerRepositoryTwo := repositories.NewHeaderRepository(dbTwo)
			for _, n := range blockNumbers {
				_, err = headerRepositoryTwo.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				Expect(err).NotTo(HaveOccurred())
			}
			repositoryTwo := inputs.RepositoryTwo
			repositoryTwo.SetDB(dbTwo)

			err := repository.MarkHeaderChecked(headerIDs[0])
			Expect(err).NotTo(HaveOccurred())
			nodeOneMissingHeaders, err := repository.MissingHeaders(startingBlockNumber, endingBlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeOneMissingHeaders)).To(Equal(2))
			Expect(nodeOneMissingHeaders[0].BlockNumber).To(Or(Equal(eventSpecificBlockNumber), Equal(endingBlockNumber)))
			Expect(nodeOneMissingHeaders[1].BlockNumber).To(Or(Equal(eventSpecificBlockNumber), Equal(endingBlockNumber)))

			nodeTwoMissingHeaders, err := repositoryTwo.MissingHeaders(startingBlockNumber, endingBlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeTwoMissingHeaders)).To(Equal(3))
			Expect(nodeTwoMissingHeaders[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(eventSpecificBlockNumber), Equal(endingBlockNumber)))
			Expect(nodeTwoMissingHeaders[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(eventSpecificBlockNumber), Equal(endingBlockNumber)))
			Expect(nodeTwoMissingHeaders[2].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(eventSpecificBlockNumber), Equal(endingBlockNumber)))
		})
	})
}

func SharedRepositoryMarkHeaderCheckedBehaviors(inputs *MarkedHeaderCheckedBehaviorInputs) {
	var repository = inputs.Repository
	var checkedHeaderColumn = inputs.CheckedHeaderColumnName

	Describe("MarkHeaderChecked", func() {
		BeforeEach(func() {
			headerRepository = GetHeaderRepository()
			headerId, err = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
		})

		It("creates a row for a new headerId", func() {
			err = repository.MarkHeaderChecked(headerId)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			query := `SELECT ` + checkedHeaderColumn + ` FROM public.checked_headers WHERE header_id = $1`
			err = db.Get(&headerChecked, query, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("updates row when headerID already exists", func() {
			_, err = db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerId)

			err = repository.MarkHeaderChecked(headerId)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			query := `SELECT ` + checkedHeaderColumn + ` FROM public.checked_headers WHERE header_id = $1`
			err = db.Get(&headerChecked, query, headerId)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("returns an error if upserting a record fails", func() {
			err = repository.MarkHeaderChecked(0)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("violates foreign key constraint"))
		})
	})
}

func GetHeaderRepository() repositories.HeaderRepository {
	db = test_config.NewTestDB(test_config.NewTestNode())
	test_config.CleanTestDB(db)

	return repositories.NewHeaderRepository(db)
}
