package debt_ceiling_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/debt_ceiling"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

var _ = Describe("", func() {
	It("converts a log to an model", func() {
		converter := debt_ceiling.PitFileDebtCeilingConverter{}

		model, err := converter.ToModel(test_data.PitAddress, shared.PitABI, test_data.EthPitFileDebtCeilingLog)

		Expect(err).NotTo(HaveOccurred())
		Expect(model).To(Equal(test_data.PitFileDebtCeilingModel))
	})
})
