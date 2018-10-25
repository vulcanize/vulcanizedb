package vat_slip_test

import (
	"database/sql"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_slip"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Vat slip repository", func() {
	var (
		db                *postgres.DB
		vatSlipRepository vat_slip.VatSlipRepository
		err               error
		headerRepository  datastore.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(core.Node{})
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		vatSlipRepository = vat_slip.VatSlipRepository{}
		vatSlipRepository.SetDB(db)
	})

	Describe("Create", func() {
		var headerID int64

		BeforeEach(func() {
			headerID, err = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
		})

		It("adds a vat slip event", func() {
			err = vatSlipRepository.Create(headerID, []interface{}{test_data.VatSlipModel})
			Expect(err).NotTo(HaveOccurred())

			var dbVatSlip vat_slip.VatSlipModel
			err = db.Get(&dbVatSlip, `SELECT ilk, guy, rad, tx_idx, log_idx, raw_log FROM maker.vat_slip WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbVatSlip.Ilk).To(Equal(test_data.VatSlipModel.Ilk))
			Expect(dbVatSlip.Guy).To(Equal(test_data.VatSlipModel.Guy))
			Expect(dbVatSlip.Rad).To(Equal(test_data.VatSlipModel.Rad))
			Expect(dbVatSlip.TransactionIndex).To(Equal(test_data.VatSlipModel.TransactionIndex))
			Expect(dbVatSlip.LogIndex).To(Equal(test_data.VatSlipModel.LogIndex))
			Expect(dbVatSlip.Raw).To(MatchJSON(test_data.VatSlipModel.Raw))
		})

		It("marks header as checked for logs", func() {
			err = vatSlipRepository.Create(headerID, []interface{}{test_data.VatSlipModel})
			Expect(err).NotTo(HaveOccurred())

			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT vat_slip_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("updates the header to checked if checked headers row already exists", func() {
			_, err := db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerID)
			Expect(err).NotTo(HaveOccurred())
			err = vatSlipRepository.Create(headerID, []interface{}{test_data.VatSlipModel})
			Expect(err).NotTo(HaveOccurred())

			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT vat_slip_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("does not duplicate vat slip events", func() {
			err = vatSlipRepository.Create(headerID, []interface{}{test_data.VatSlipModel})
			Expect(err).NotTo(HaveOccurred())

			err = vatSlipRepository.Create(headerID, []interface{}{test_data.VatSlipModel})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("allows for multiple vat slip events in one transaction if they have different log indexes", func() {
			err = vatSlipRepository.Create(headerID, []interface{}{test_data.VatSlipModel})
			Expect(err).NotTo(HaveOccurred())

			newVatSlip := test_data.VatSlipModel
			newVatSlip.LogIndex = newVatSlip.LogIndex + 1
			err := vatSlipRepository.Create(headerID, []interface{}{newVatSlip})

			Expect(err).NotTo(HaveOccurred())
		})

		It("removes vat slip if corresponding header is deleted", func() {
			err = vatSlipRepository.Create(headerID, []interface{}{test_data.VatSlipModel})
			Expect(err).NotTo(HaveOccurred())

			_, err = db.Exec(`DELETE FROM headers WHERE id = $1`, headerID)

			Expect(err).NotTo(HaveOccurred())
			var dbVatSlip vat_slip.VatSlipModel
			err = db.Get(&dbVatSlip, `SELECT ilk, guy, rad, tx_idx, raw_log FROM maker.vat_slip WHERE header_id = $1`, headerID)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(sql.ErrNoRows))
		})

		It("returns an error if model is of wrong type", func() {
			err = vatSlipRepository.Create(headerID, []interface{}{test_data.WrongModel{}})
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
			err = vatSlipRepository.MarkHeaderChecked(headerID)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT vat_slip_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("updates row when headerID already exists", func() {
			_, err = db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerID)

			err = vatSlipRepository.MarkHeaderChecked(headerID)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT vat_slip_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})
	})

	Describe("MissingHeaders", func() {
		var (
			startingBlock, endingBlock, vatSlipBlock int64
			blockNumbers, headerIDs                  []int64
		)

		BeforeEach(func() {
			startingBlock = rand.Int63()
			vatSlipBlock = startingBlock + 1
			endingBlock = startingBlock + 2

			blockNumbers = []int64{startingBlock, vatSlipBlock, endingBlock, endingBlock + 1}

			headerIDs = []int64{}
			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				Expect(err).NotTo(HaveOccurred())
				headerIDs = append(headerIDs, headerID)
			}
		})

		It("returns headers that haven't been checked", func() {
			err := vatSlipRepository.MarkHeaderChecked(headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := vatSlipRepository.MissingHeaders(startingBlock, endingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(2))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock)))
		})

		It("only treats headers as checked if vat slip logs have been checked", func() {
			_, err := db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := vatSlipRepository.MissingHeaders(startingBlock, endingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(3))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(vatSlipBlock)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(vatSlipBlock)))
			Expect(headers[2].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(vatSlipBlock)))
		})

		It("only returns headers associated with the current node", func() {
			dbTwo := test_config.NewTestDB(core.Node{ID: "second"})
			headerRepositoryTwo := repositories.NewHeaderRepository(dbTwo)
			for _, n := range blockNumbers {
				_, err = headerRepositoryTwo.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				Expect(err).NotTo(HaveOccurred())
			}
			vatSlipRepositoryTwo := vat_slip.VatSlipRepository{}
			vatSlipRepositoryTwo.SetDB(dbTwo)
			err := vatSlipRepository.MarkHeaderChecked(headerIDs[0])
			Expect(err).NotTo(HaveOccurred())

			nodeOneMissingHeaders, err := vatSlipRepository.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeOneMissingHeaders)).To(Equal(len(blockNumbers) - 1))

			nodeTwoMissingHeaders, err := vatSlipRepositoryTwo.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeTwoMissingHeaders)).To(Equal(len(blockNumbers)))
		})
	})
})
