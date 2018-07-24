package pep_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pep"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Pep repository", func() {
	It("creates a pep", func() {
		db := test_config.NewTestDB(core.Node{})
		test_config.CleanTestDB(db)
		repository := pep.NewPepRepository(db)
		header := core.Header{BlockNumber: 12345}
		headerRepository := repositories.NewHeaderRepository(db)
		headerID, err := headerRepository.CreateOrUpdateHeader(header)
		Expect(err).NotTo(HaveOccurred())
		pepToAdd := pep.Pep{
			BlockNumber: header.BlockNumber,
			HeaderID:    headerID,
			UsdValue:    "123.456",
		}

		err = repository.CreatePep(pepToAdd)

		Expect(err).NotTo(HaveOccurred())
		var dbPep pep.Pep
		err = db.Get(&dbPep, `SELECT block_number, header_id, usd_value FROM maker.peps WHERE header_id = $1`, pepToAdd.HeaderID)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbPep).To(Equal(pepToAdd))
	})
})
