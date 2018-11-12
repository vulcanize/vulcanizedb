package vat_tune_test

import (
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_tune"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Vat tune repository", func() {
	var (
		db         *postgres.DB
		repository vat_tune.VatTuneRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repository = vat_tune.VatTuneRepository{}
		repository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.VatTuneModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.VatTuneChecked,
			LogEventTableName:        "maker.vat_heal",
			TestModel:                test_data.VatTuneModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &repository,
		}
		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a vat tune event", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			err = repository.Create(headerID, []interface{}{test_data.VatTuneModel})
			Expect(err).NotTo(HaveOccurred())

			var dbVatTune vat_tune.VatTuneModel
			err = db.Get(&dbVatTune, `SELECT ilk, urn, v, w, dink, dart, tx_idx, log_idx, raw_log FROM maker.vat_tune WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbVatTune.Ilk).To(Equal(test_data.VatTuneModel.Ilk))
			Expect(dbVatTune.Urn).To(Equal(test_data.VatTuneModel.Urn))
			Expect(dbVatTune.V).To(Equal(test_data.VatTuneModel.V))
			Expect(dbVatTune.W).To(Equal(test_data.VatTuneModel.W))
			Expect(dbVatTune.Dink).To(Equal(test_data.VatTuneModel.Dink))
			Expect(dbVatTune.Dart).To(Equal(test_data.VatTuneModel.Dart))
			Expect(dbVatTune.TransactionIndex).To(Equal(test_data.VatTuneModel.TransactionIndex))
			Expect(dbVatTune.LogIndex).To(Equal(test_data.VatTuneModel.LogIndex))
			Expect(dbVatTune.Raw).To(MatchJSON(test_data.VatTuneModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.VatTuneChecked,
			Repository:              &repository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})

	Describe("MissingHeaders", func() {
		inputs := shared_behaviors.MissingHeadersBehaviorInputs{
			Repository:    &repository,
			RepositoryTwo: &vat_tune.VatTuneRepository{},
		}

		shared_behaviors.SharedRepositoryMissingHeadersBehaviors(&inputs)
	})
})
