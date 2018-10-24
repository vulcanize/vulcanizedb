package vat_tune_test

import (
	"database/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_tune"
	"github.com/vulcanize/vulcanizedb/test_config"
	"math/rand"
)

var _ = Describe("Vat tune repository", func() {
	var (
		db                *postgres.DB
		vatTuneRepository vat_tune.VatTuneRepository
		err               error
		headerRepository  datastore.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(core.Node{})
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		vatTuneRepository = vat_tune.VatTuneRepository{}
		vatTuneRepository.SetDB(db)
	})

	Describe("Create", func() {
		var headerID int64

		BeforeEach(func() {
			headerID, err = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
		})

		It("adds a vat tune event", func() {
			err = vatTuneRepository.Create(headerID, []interface{}{test_data.VatTuneModel})
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

		It("marks header as checked for logs", func() {
			err = vatTuneRepository.Create(headerID, []interface{}{test_data.VatTuneModel})
			Expect(err).NotTo(HaveOccurred())

			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT vat_tune_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("updates the header to checked if checked headers row already exists", func() {
			_, err := db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerID)
			Expect(err).NotTo(HaveOccurred())
			err = vatTuneRepository.Create(headerID, []interface{}{test_data.VatTuneModel})
			Expect(err).NotTo(HaveOccurred())

			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT vat_tune_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("does not duplicate vat tune events", func() {
			err = vatTuneRepository.Create(headerID, []interface{}{test_data.VatTuneModel})
			Expect(err).NotTo(HaveOccurred())

			err = vatTuneRepository.Create(headerID, []interface{}{test_data.VatTuneModel})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("allows for multiple vat tune events in one transaction if they have different log indexes", func() {
			err = vatTuneRepository.Create(headerID, []interface{}{test_data.VatTuneModel})
			Expect(err).NotTo(HaveOccurred())

			newVatTune := test_data.VatTuneModel
			newVatTune.LogIndex = newVatTune.LogIndex + 1
			err := vatTuneRepository.Create(headerID, []interface{}{newVatTune})

			Expect(err).NotTo(HaveOccurred())
		})

		It("removes vat tune if corresponding header is deleted", func() {
			err = vatTuneRepository.Create(headerID, []interface{}{test_data.VatTuneModel})
			Expect(err).NotTo(HaveOccurred())

			_, err = db.Exec(`DELETE FROM headers WHERE id = $1`, headerID)

			Expect(err).NotTo(HaveOccurred())
			var dbVatTune vat_tune.VatTuneModel
			err = db.Get(&dbVatTune, `SELECT ilk, urn, v, w, dink, dart, tx_idx, raw_log FROM maker.vat_tune WHERE header_id = $1`, headerID)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(sql.ErrNoRows))
		})

		It("returns an error if model is of wrong type", func() {
			err = vatTuneRepository.Create(headerID, []interface{}{test_data.WrongModel{}})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("model of type"))
		})
	})

	Describe("MarkHeaderChecked", func() {
		var headerID int64

		BeforeEach(func() {
			headerID, err = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
		})

		It("creates a row for a new headerID", func() {
			err = vatTuneRepository.MarkHeaderChecked(headerID)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT vat_tune_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("updates row when headerID already exists", func() {
			_, err = db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerID)

			err = vatTuneRepository.MarkHeaderChecked(headerID)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT vat_tune_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})
	})

	Describe("MissingHeaders", func() {
		var (
			startingBlockNumber, vatTuneBlockNumber, endingBlockNumber int64
			blockNumbers, headerIDs                                    []int64
		)

		BeforeEach(func() {
			startingBlockNumber = rand.Int63()
			vatTuneBlockNumber = startingBlockNumber + 1
			endingBlockNumber = startingBlockNumber + 2

			blockNumbers = []int64{startingBlockNumber, vatTuneBlockNumber, endingBlockNumber, endingBlockNumber + 1}

			headerIDs = []int64{}
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				Expect(err).NotTo(HaveOccurred())
				headerIDs = append(headerIDs, headerID)
			}
		})

		It("returns headers that haven't been checked", func() {
			err := vatTuneRepository.MarkHeaderChecked(headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := vatTuneRepository.MissingHeaders(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(2))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber)))
		})

		It("only treats headers as checked if vat tune logs have been checked", func() {
			_, err := db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := vatTuneRepository.MissingHeaders(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(3))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(vatTuneBlockNumber)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(vatTuneBlockNumber)))
			Expect(headers[2].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(vatTuneBlockNumber)))
		})

		It("only returns headers associated with the current node", func() {
			dbTwo := test_config.NewTestDB(core.Node{ID: "second"})
			headerRepositoryTwo := repositories.NewHeaderRepository(dbTwo)
			for _, n := range blockNumbers {
				_, err = headerRepositoryTwo.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				Expect(err).NotTo(HaveOccurred())
			}
			vatTuneRepositoryTwo := vat_tune.VatTuneRepository{}
			vatTuneRepositoryTwo.SetDB(dbTwo)
			err := vatTuneRepository.MarkHeaderChecked(headerIDs[0])
			Expect(err).NotTo(HaveOccurred())

			nodeOneMissingHeaders, err := vatTuneRepository.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeOneMissingHeaders)).To(Equal(len(blockNumbers) - 1))

			nodeTwoMissingHeaders, err := vatTuneRepositoryTwo.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeTwoMissingHeaders)).To(Equal(len(blockNumbers)))
		})
	})
})
