package vat_grab_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_grab"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Vat grab repository", func() {
	var (
		db                *postgres.DB
		vatGrabRepository vat_grab.VatGrabRepository
		headerRepository  datastore.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		vatGrabRepository = vat_grab.VatGrabRepository{}
		vatGrabRepository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.VatGrabModel
		modelWithDifferentLogIdx.LogIndex++

		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.VatGrabChecked,
			LogEventTableName:        "maker.vat_grab",
			TestModel:                test_data.VatGrabModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &vatGrabRepository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a vat grab event", func() {
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

			err = vatGrabRepository.Create(headerID, []interface{}{test_data.VatGrabModel})
			Expect(err).NotTo(HaveOccurred())
			var dbVatGrab vat_grab.VatGrabModel
			err = db.Get(&dbVatGrab, `SELECT ilk, urn, v, w, dink, dart, log_idx, tx_idx, raw_log FROM maker.vat_grab WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbVatGrab.Ilk).To(Equal(test_data.VatGrabModel.Ilk))
			Expect(dbVatGrab.Urn).To(Equal(test_data.VatGrabModel.Urn))
			Expect(dbVatGrab.V).To(Equal(test_data.VatGrabModel.V))
			Expect(dbVatGrab.W).To(Equal(test_data.VatGrabModel.W))
			Expect(dbVatGrab.Dink).To(Equal(test_data.VatGrabModel.Dink))
			Expect(dbVatGrab.Dart).To(Equal(test_data.VatGrabModel.Dart))
			Expect(dbVatGrab.LogIndex).To(Equal(test_data.VatGrabModel.LogIndex))
			Expect(dbVatGrab.TransactionIndex).To(Equal(test_data.VatGrabModel.TransactionIndex))
			Expect(dbVatGrab.Raw).To(MatchJSON(test_data.VatGrabModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.VatGrabChecked,
			Repository:              &vatGrabRepository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})
