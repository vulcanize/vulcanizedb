package vat_toll_test

import (
	"database/sql"
	"encoding/json"

	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_toll"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Vat toll repository", func() {
	var (
		db                *postgres.DB
		vatTollRepository vat_toll.Repository
		err               error
		headerRepository  datastore.HeaderRepository
		rawHeader         []byte
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(core.Node{})
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		vatTollRepository = vat_toll.NewVatTollRepository(db)
		rawHeader, err = json.Marshal(types.Header{})
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Create", func() {
		var headerID int64

		BeforeEach(func() {
			headerID, err = headerRepository.CreateOrUpdateHeader(core.Header{Raw: rawHeader})
			Expect(err).NotTo(HaveOccurred())

			err = vatTollRepository.Create(headerID, []vat_toll.VatTollModel{test_data.VatTollModel})
			Expect(err).NotTo(HaveOccurred())
		})

		It("adds a vat toll event", func() {
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
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT vat_toll_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("does not duplicate vat toll events", func() {
			err = vatTollRepository.Create(headerID, []vat_toll.VatTollModel{test_data.VatTollModel})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("removes vat toll if corresponding header is deleted", func() {
			_, err = db.Exec(`DELETE FROM headers WHERE id = $1`, headerID)

			Expect(err).NotTo(HaveOccurred())
			var dbVatToll vat_toll.VatTollModel
			err = db.Get(&dbVatToll, `SELECT ilk, urn, take, tx_idx, raw_log FROM maker.vat_toll WHERE header_id = $1`, headerID)
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
		var (
			startingBlock, endingBlock, vatTollBlock int64
			blockNumbers, headerIDs                  []int64
		)

		BeforeEach(func() {
			startingBlock = GinkgoRandomSeed()
			vatTollBlock = startingBlock + 1
			endingBlock = startingBlock + 2

			blockNumbers = []int64{startingBlock, vatTollBlock, endingBlock, endingBlock + 1}

			headerIDs = []int64{}
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(core.Header{BlockNumber: n, Raw: rawHeader})
				Expect(err).NotTo(HaveOccurred())
				headerIDs = append(headerIDs, headerID)
			}
		})

		It("returns headers that haven't been checked", func() {
			err := vatTollRepository.MarkHeaderChecked(headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := vatTollRepository.MissingHeaders(startingBlock, endingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(2))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock)))
		})

		It("only treats headers as checked if vat toll logs have been checked", func() {
			_, err := db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := vatTollRepository.MissingHeaders(startingBlock, endingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(3))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(vatTollBlock)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(vatTollBlock)))
			Expect(headers[2].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(vatTollBlock)))
		})

		It("only returns headers associated with the current node", func() {
			dbTwo := test_config.NewTestDB(core.Node{ID: "second"})
			headerRepositoryTwo := repositories.NewHeaderRepository(dbTwo)
			for _, n := range blockNumbers {
				_, err = headerRepositoryTwo.CreateOrUpdateHeader(core.Header{BlockNumber: n, Raw: rawHeader})
				Expect(err).NotTo(HaveOccurred())
			}
			vatTollRepositoryTwo := vat_toll.NewVatTollRepository(dbTwo)
			err := vatTollRepository.MarkHeaderChecked(headerIDs[0])
			Expect(err).NotTo(HaveOccurred())

			nodeOneMissingHeaders, err := vatTollRepository.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeOneMissingHeaders)).To(Equal(len(blockNumbers) - 1))

			nodeTwoMissingHeaders, err := vatTollRepositoryTwo.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeTwoMissingHeaders)).To(Equal(len(blockNumbers)))
		})
	})
})
