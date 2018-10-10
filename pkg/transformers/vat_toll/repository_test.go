package vat_toll_test

import (
	"database/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_toll"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Vat toll repository", func() {
	Describe("Create", func() {
		var (
			db                *postgres.DB
			vatTollRepository vat_toll.Repository
			err               error
			headerID          int64
		)

		BeforeEach(func() {
			db = test_config.NewTestDB(core.Node{})
			test_config.CleanTestDB(db)
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err = headerRepository.CreateOrUpdateHeader(core.Header{})
			Expect(err).NotTo(HaveOccurred())
			vatTollRepository = vat_toll.NewVatTollRepository(db)
		})

		It("adds a vat toll event", func() {
			err = vatTollRepository.Create(headerID, []vat_toll.VatTollModel{test_data.VatTollModel})

			Expect(err).NotTo(HaveOccurred())
			var dbVatToll vat_toll.VatTollModel
			err = db.Get(&dbVatToll, `SELECT ilk, urn, take, tx_idx, raw_log FROM maker.vat_toll WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbVatToll.Ilk).To(Equal(test_data.VatTollModel.Ilk))
			Expect(dbVatToll.Urn).To(Equal(test_data.VatTollModel.Urn))
			Expect(dbVatToll.Take).To(Equal(test_data.VatTollModel.Take))
			Expect(dbVatToll.TransactionIndex).To(Equal(test_data.VatTollModel.TransactionIndex))
			Expect(dbVatToll.Raw).To(MatchJSON(test_data.VatTollModel.Raw))
		})

		It("marks header as checked for logs", func() {
			err = vatTollRepository.Create(headerID, []vat_toll.VatTollModel{test_data.VatTollModel})

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT vat_toll_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("does not duplicate vat toll events", func() {
			err = vatTollRepository.Create(headerID, []vat_toll.VatTollModel{test_data.VatTollModel})
			Expect(err).NotTo(HaveOccurred())

			err = vatTollRepository.Create(headerID, []vat_toll.VatTollModel{test_data.VatTollModel})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("removes vat toll if corresponding header is deleted", func() {
			err = vatTollRepository.Create(headerID, []vat_toll.VatTollModel{test_data.VatTollModel})
			Expect(err).NotTo(HaveOccurred())

			_, err = db.Exec(`DELETE FROM headers WHERE id = $1`, headerID)

			Expect(err).NotTo(HaveOccurred())
			var dbVatToll vat_toll.VatTollModel
			err = db.Get(&dbVatToll, `SELECT ilk, urn, take, tx_idx, raw_log FROM maker.vat_toll WHERE header_id = $1`, headerID)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(sql.ErrNoRows))
		})
	})

	Describe("MarkHeaderChecked", func() {
		var (
			db                *postgres.DB
			vatTollRepository vat_toll.Repository
			err               error
			headerID          int64
		)

		BeforeEach(func() {
			db = test_config.NewTestDB(core.Node{})
			test_config.CleanTestDB(db)
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err = headerRepository.CreateOrUpdateHeader(core.Header{})
			Expect(err).NotTo(HaveOccurred())
			vatTollRepository = vat_toll.NewVatTollRepository(db)
		})

		It("creates a row for a new headerID", func() {
			err = vatTollRepository.MarkHeaderChecked(headerID)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT vat_toll_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("updates row when headerID already exists", func() {
			_, err = db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerID)

			err = vatTollRepository.MarkHeaderChecked(headerID)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT vat_toll_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})
	})

	Describe("MissingHeaders", func() {
		It("returns headers that haven't been checked", func() {
			db := test_config.NewTestDB(core.Node{})
			test_config.CleanTestDB(db)
			headerRepository := repositories.NewHeaderRepository(db)
			startingBlockNumber := int64(1)
			pitFileBlockNumber := int64(2)
			endingBlockNumber := int64(3)
			blockNumbers := []int64{startingBlockNumber, pitFileBlockNumber, endingBlockNumber, endingBlockNumber + 1}
			var headerIDs []int64
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{BlockNumber: n})
				headerIDs = append(headerIDs, headerID)
				Expect(err).NotTo(HaveOccurred())
			}
			pitFileRepository := vat_toll.NewVatTollRepository(db)
			err := pitFileRepository.MarkHeaderChecked(headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := pitFileRepository.MissingHeaders(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(2))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber)))
		})

		It("only treats headers as checked if vat toll logs have been checked", func() {
			db := test_config.NewTestDB(core.Node{})
			test_config.CleanTestDB(db)
			headerRepository := repositories.NewHeaderRepository(db)
			startingBlockNumber := int64(1)
			vatTolldBlockNumber := int64(2)
			endingBlockNumber := int64(3)
			blockNumbers := []int64{startingBlockNumber, vatTolldBlockNumber, endingBlockNumber, endingBlockNumber + 1}
			var headerIDs []int64
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{BlockNumber: n})
				headerIDs = append(headerIDs, headerID)
				Expect(err).NotTo(HaveOccurred())
			}
			vatTollRepository := vat_toll.NewVatTollRepository(db)
			_, err := db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := vatTollRepository.MissingHeaders(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(3))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(vatTolldBlockNumber)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(vatTolldBlockNumber)))
			Expect(headers[2].BlockNumber).To(Or(Equal(startingBlockNumber), Equal(endingBlockNumber), Equal(vatTolldBlockNumber)))
		})

		It("only returns headers associated with the current node", func() {
			db := test_config.NewTestDB(core.Node{})
			test_config.CleanTestDB(db)
			blockNumbers := []int64{1, 2, 3}
			headerRepository := repositories.NewHeaderRepository(db)
			dbTwo := test_config.NewTestDB(core.Node{ID: "second"})
			headerRepositoryTwo := repositories.NewHeaderRepository(dbTwo)
			var headerIDs []int64
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{BlockNumber: n})
				Expect(err).NotTo(HaveOccurred())
				headerIDs = append(headerIDs, headerID)
				_, err = headerRepositoryTwo.CreateOrUpdateHeader(core.Header{BlockNumber: n})
				Expect(err).NotTo(HaveOccurred())
			}
			pitFileRepository := vat_toll.NewVatTollRepository(db)
			pitFileRepositoryTwo := vat_toll.NewVatTollRepository(dbTwo)
			err := pitFileRepository.MarkHeaderChecked(headerIDs[0])
			Expect(err).NotTo(HaveOccurred())

			nodeOneMissingHeaders, err := pitFileRepository.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeOneMissingHeaders)).To(Equal(len(blockNumbers) - 1))

			nodeTwoMissingHeaders, err := pitFileRepositoryTwo.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeTwoMissingHeaders)).To(Equal(len(blockNumbers)))
		})
	})
})
