package vat_slip_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"strconv"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_slip"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Vat slip repository", func() {
	var (
		db         *postgres.DB
		repository vat_slip.VatSlipRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repository = vat_slip.VatSlipRepository{}
		repository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.VatSlipModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.VatSlipChecked,
			LogEventTableName:        "maker.vat_slip",
			TestModel:                test_data.VatSlipModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &repository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a vat slip event", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			err = repository.Create(headerID, []interface{}{test_data.VatSlipModel})
			Expect(err).NotTo(HaveOccurred())

			var dbVatSlip vat_slip.VatSlipModel
			err = db.Get(&dbVatSlip, `SELECT ilk, guy, rad, tx_idx, log_idx, raw_log FROM maker.vat_slip WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_data.VatSlipModel.Ilk, db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbVatSlip.Ilk).To(Equal(strconv.Itoa(ilkID)))
			Expect(dbVatSlip.Guy).To(Equal(test_data.VatSlipModel.Guy))
			Expect(dbVatSlip.Rad).To(Equal(test_data.VatSlipModel.Rad))
			Expect(dbVatSlip.TransactionIndex).To(Equal(test_data.VatSlipModel.TransactionIndex))
			Expect(dbVatSlip.LogIndex).To(Equal(test_data.VatSlipModel.LogIndex))
			Expect(dbVatSlip.Raw).To(MatchJSON(test_data.VatSlipModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.VatSlipChecked,
			Repository:              &repository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
