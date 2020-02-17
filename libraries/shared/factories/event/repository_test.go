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

package event_test

import (
	"fmt"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/test_data"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	"github.com/makerdao/vulcanizedb/test_config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"math/big"
)

var _ = Describe("Repository", func() {
	var db *postgres.DB

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
	})

	Describe("PersistModels", func() {
		const createTestEventTableQuery = `CREATE TABLE public.testEvent(
		id        SERIAL PRIMARY KEY,
		header_id INTEGER NOT NULL REFERENCES public.headers (id) ON DELETE CASCADE,
		log_id    BIGINT  NOT NULL REFERENCES public.event_logs (id) ON DELETE CASCADE,
		variable1 TEXT,
		UNIQUE (header_id, log_id)
		);`

		var (
			headerID, logID  int64
			headerRepository repositories.HeaderRepository
			testModel        event.InsertionModel
		)

		BeforeEach(func() {
			_, tableErr := db.Exec(createTestEventTableQuery)
			Expect(tableErr).NotTo(HaveOccurred())
			headerRepository = repositories.NewHeaderRepository(db)
			var insertHeaderErr error
			headerID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(insertHeaderErr).NotTo(HaveOccurred())
			eventLog := test_data.CreateTestLog(headerID, db)
			logID = eventLog.ID

			testModel = event.InsertionModel{
				SchemaName: "public",
				TableName:  "testEvent",
				OrderedColumns: []event.ColumnName{
					event.HeaderFK, event.LogFK, "variable1",
				},
				ColumnValues: event.ColumnValues{
					event.HeaderFK: headerID,
					event.LogFK:    logID,
					"variable1":    "value1",
				},
			}
		})

		AfterEach(func() {
			db.MustExec(`DROP TABLE public.testEvent;`)
		})

		// Needs to run before the other tests, since those insert keys in map
		It("memoizes queries", func() {
			Expect(len(event.ModelToQuery)).To(Equal(0))
			event.GetMemoizedQuery(testModel)
			Expect(len(event.ModelToQuery)).To(Equal(1))
			event.GetMemoizedQuery(testModel)
			Expect(len(event.ModelToQuery)).To(Equal(1))
		})

		It("persists a model to postgres", func() {
			createErr := event.PersistModels([]event.InsertionModel{testModel}, db)
			Expect(createErr).NotTo(HaveOccurred())

			var res TestEvent
			dbErr := db.Get(&res, `SELECT log_id, variable1 FROM public.testEvent;`)
			Expect(dbErr).NotTo(HaveOccurred())

			Expect(res.LogID).To(Equal(fmt.Sprint(testModel.ColumnValues[event.LogFK])))
			Expect(res.Variable1).To(Equal(testModel.ColumnValues["variable1"]))
		})

		Describe("returns errors", func() {
			It("for empty model slice", func() {
				err := event.PersistModels([]event.InsertionModel{}, db)
				Expect(err).To(MatchError("repository got empty model slice"))
			})

			It("for failed SQL inserts", func() {
				header := fakes.GetFakeHeader(1)
				headerID, headerErr := headerRepository.CreateOrUpdateHeader(header)
				Expect(headerErr).NotTo(HaveOccurred())

				brokenModel := event.InsertionModel{
					SchemaName: "public",
					TableName:  "testEvent",
					// Wrong name of last column compared to DB, will generate incorrect query
					OrderedColumns: []event.ColumnName{
						event.HeaderFK, event.LogFK, "variable2",
					},
					ColumnValues: event.ColumnValues{
						event.HeaderFK: headerID,
						event.LogFK:    logID,
						"variable1":    "value1",
					},
				}

				// Remove cached queries, or we won't generate a new (incorrect) one
				delete(event.ModelToQuery, "publictestEvent")

				createErr := event.PersistModels([]event.InsertionModel{brokenModel}, db)
				// Remove incorrect query, so other tests won't get it
				delete(event.ModelToQuery, "publictestEvent")

				Expect(createErr).To(HaveOccurred())
			})

			It("for unsupported types in ColumnValue", func() {
				unsupportedValue := big.NewInt(5)
				testModel = event.InsertionModel{
					SchemaName: "public",
					TableName:  "testEvent",
					OrderedColumns: []event.ColumnName{
						event.HeaderFK, event.LogFK, "variable1",
					},
					ColumnValues: event.ColumnValues{
						event.HeaderFK: headerID,
						event.LogFK:    logID,
						"variable1":    unsupportedValue,
					},
				}

				createErr := event.PersistModels([]event.InsertionModel{testModel}, db)
				Expect(createErr).To(MatchError(event.ErrUnsupportedValue(unsupportedValue)))
			})
		})

		It("upserts queries with conflicting source", func() {
			conflictingModel := event.InsertionModel{
				SchemaName: "public",
				TableName:  "testEvent",
				OrderedColumns: []event.ColumnName{
					event.HeaderFK, event.LogFK, "variable1",
				},
				ColumnValues: event.ColumnValues{
					event.HeaderFK: headerID,
					event.LogFK:    logID,
					"variable1":    "conflictingValue",
				},
			}

			createErr := event.PersistModels([]event.InsertionModel{testModel, conflictingModel}, db)
			Expect(createErr).NotTo(HaveOccurred())

			var res TestEvent
			dbErr := db.Get(&res, `SELECT log_id, variable1 FROM public.testEvent;`)
			Expect(dbErr).NotTo(HaveOccurred())
			Expect(res.Variable1).To(Equal(conflictingModel.ColumnValues["variable1"]))
		})

		It("generates correct queries", func() {
			actualQuery := event.GenerateInsertionQuery(testModel)
			expectedQuery := `INSERT INTO public.testEvent (header_id, log_id, variable1) VALUES($1, $2, $3)
		ON CONFLICT (header_id, log_id) DO UPDATE SET header_id = $1, log_id = $2, variable1 = $3;`
			Expect(actualQuery).To(Equal(expectedQuery))
		})

		It("marks log transformed", func() {
			createErr := event.PersistModels([]event.InsertionModel{testModel}, db)
			Expect(createErr).NotTo(HaveOccurred())

			var logTransformed bool
			getErr := db.Get(&logTransformed, `SELECT transformed FROM public.event_logs WHERE id = $1`, logID)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(logTransformed).To(BeTrue())
		})
	})
})

type TestEvent struct {
	LogID     string `db:"log_id"`
	Variable1 string
}
