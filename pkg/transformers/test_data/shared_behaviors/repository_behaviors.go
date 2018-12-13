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
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/factories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
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
			headerRepository = getHeaderRepository()
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

func SharedRepositoryMarkHeaderCheckedBehaviors(inputs *MarkedHeaderCheckedBehaviorInputs) {
	var repository = inputs.Repository
	var checkedHeaderColumn = inputs.CheckedHeaderColumnName

	Describe("MarkHeaderChecked", func() {
		BeforeEach(func() {
			headerRepository = getHeaderRepository()
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

func getHeaderRepository() repositories.HeaderRepository {
	db = test_config.NewTestDB(test_config.NewTestNode())
	test_config.CleanTestDB(db)

	return repositories.NewHeaderRepository(db)
}
