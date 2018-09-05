package stability_fee_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/stability_fee"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

var _ = Describe("Pit file stability fee converter", func() {
	It("converts a log to an model", func() {
		converter := stability_fee.PitFileStabilityFeeConverter{}

		model, err := converter.ToModel(test_data.PitAddress, shared.PitABI, test_data.EthPitFileStabilityFeeLog)

		Expect(err).NotTo(HaveOccurred())
		Expect(model).To(Equal(test_data.PitFileStabilityFeeModel))
	})
})
