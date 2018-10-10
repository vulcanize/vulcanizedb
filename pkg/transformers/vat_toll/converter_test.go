package vat_toll_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_toll"
)

var _ = Describe("Vat toll converter", func() {
	It("returns err if log is missing topics", func() {
		converter := vat_toll.VatTollConverter{}
		badLog := types.Log{
			Data: []byte{1, 1, 1, 1, 1},
		}

		_, err := converter.ToModels([]types.Log{badLog})

		Expect(err).To(HaveOccurred())
	})

	It("converts a log to an model", func() {
		converter := vat_toll.VatTollConverter{}

		models, err := converter.ToModels([]types.Log{test_data.EthVatTollLog})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(models)).To(Equal(1))
		Expect(models[0]).To(Equal(test_data.VatTollModel))
	})
})
