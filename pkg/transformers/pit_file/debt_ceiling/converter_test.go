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
		Expect(model.What).To(Equal(test_data.PitFileDebtCeilingModel.What))
		Expect(model.Data).To(Equal(test_data.PitFileDebtCeilingModel.Data))
		Expect(model.TransactionIndex).To(Equal(test_data.PitFileDebtCeilingModel.TransactionIndex))
		Expect(model.Raw).To(Equal(test_data.PitFileDebtCeilingModel.Raw))
		Expect(model).To(Equal(test_data.PitFileDebtCeilingModel))
	})
})
