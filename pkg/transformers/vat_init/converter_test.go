package vat_init_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_init"
)

var _ = Describe("Vat init converter", func() {
	It("returns err if log missing topics", func() {
		converter := vat_init.VatInitConverter{}
		badLog := types.Log{}

		_, err := converter.ToModel(badLog)

		Expect(err).To(HaveOccurred())
	})

	It("converts a log to an model", func() {
		converter := vat_init.VatInitConverter{}

		model, err := converter.ToModel(test_data.EthVatInitLog)

		Expect(err).NotTo(HaveOccurred())
		Expect(model).To(Equal(test_data.VatInitModel))
	})
})
