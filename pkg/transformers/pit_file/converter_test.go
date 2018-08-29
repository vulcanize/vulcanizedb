package pit_file_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

var _ = Describe("Pit file converter", func() {
	It("converts a log to an entity", func() {
		converter := pit_file.PitFileConverter{}

		model, err := converter.ToModel(test_data.PitAddress, shared.PitABI, test_data.EthPitFileLog)

		Expect(err).NotTo(HaveOccurred())
		Expect(model).To(Equal(test_data.PitFileModel))
	})
})
