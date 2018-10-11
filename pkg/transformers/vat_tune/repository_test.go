package vat_tune_test

import (
	"database/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"encoding/json"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_tune"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Vat tune repository", func() {
	var (
		db                *postgres.DB
		vatTuneRepository vat_tune.Repository
		err               error
		headerRepository  datastore.HeaderRepository
		rawHeader         []byte
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(core.Node{})
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		vatTuneRepository = vat_tune.NewVatTuneRepository(db)
		rawHeader, err = json.Marshal(types.Header{})
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Create", func() {
		var headerID int64

		BeforeEach(func() {
			headerID, err = headerRepository.CreateOrUpdateHeader(core.Header{Raw: rawHeader})
			Expect(err).NotTo(HaveOccurred())

			err = vatTuneRepository.Create(headerID, []vat_tune.VatTuneModel{test_data.VatTuneModel})
			Expect(err).NotTo(HaveOccurred())
		})

		It("adds a vat tune event", func() {
			var dbVatTune vat_tune.VatTuneModel
			err = db.Get(&dbVatTune, `SELECT ilk, urn, v, w, dink, dart, tx_idx, raw_log FROM maker.vat_tune WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbVatTune.Ilk).To(Equal(test_data.VatTuneModel.Ilk))
			Expect(dbVatTune.Urn).To(Equal(test_data.VatTuneModel.Urn))
			Expect(dbVatTune.V).To(Equal(test_data.VatTuneModel.V))
			Expect(dbVatTune.W).To(Equal(test_data.VatTuneModel.W))
			Expect(dbVatTune.Dink).To(Equal(test_data.VatTuneModel.Dink))
			Expect(dbVatTune.Dart).To(Equal(test_data.VatTuneModel.Dart))
			Expect(dbVatTune.TransactionIndex).To(Equal(test_data.VatTuneModel.TransactionIndex))
			Expect(dbVatTune.Raw).To(MatchJSON(test_data.VatTuneModel.Raw))
		})

		It("marks header as checked for logs", func() {
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT vat_tune_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("does not duplicate pit file vat_tune events", func() {
			err = vatTuneRepository.Create(headerID, []vat_tune.VatTuneModel{test_data.VatTuneModel})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("removes pit file vat_tune if corresponding header is deleted", func() {
			_, err = db.Exec(`DELETE FROM headers WHERE id = $1`, headerID)

			Expect(err).NotTo(HaveOccurred())
			var dbVatTune vat_tune.VatTuneModel
			err = db.Get(&dbVatTune, `SELECT ilk, urn, v, w, dink, dart, tx_idx, raw_log FROM maker.vat_tune WHERE header_id = $1`, headerID)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(sql.ErrNoRows))
		})
	})

	Describe("MarkHeaderChecked", func() {
		var headerID int64

		BeforeEach(func() {
			headerID, err = headerRepository.CreateOrUpdateHeader(core.Header{Raw: rawHeader})
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
			startingBlockNumber = GinkgoRandomSeed()
			vatTuneBlockNumber = startingBlockNumber + 1
			endingBlockNumber = startingBlockNumber + 2

			blockNumbers = []int64{startingBlockNumber, vatTuneBlockNumber, endingBlockNumber, endingBlockNumber + 1}

			headerIDs = []int64{}
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{BlockNumber: n, Raw: rawHeader})
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
				_, err = headerRepositoryTwo.CreateOrUpdateHeader(core.Header{BlockNumber: n, Raw: rawHeader})
				Expect(err).NotTo(HaveOccurred())
			}
			vatTuneRepository := vat_tune.NewVatTuneRepository(db)
			vatTuneRepositoryTwo := vat_tune.NewVatTuneRepository(dbTwo)
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
