package history_test

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/inmemory"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/history"
)

var _ = Describe("Header validator", func() {
	It("replaces headers in the validation window that have changed", func() {
		blockNumber := int64(10)
		oldHash := common.HexToHash("0x0987654321").Hex()
		oldHeader := core.Header{
			BlockNumber: blockNumber,
			Hash:        oldHash,
		}
		inMemory := inmemory.NewInMemory()
		headerRepository := inmemory.NewHeaderRepository(inMemory)
		headerRepository.CreateOrUpdateHeader(oldHeader)
		newHash := common.HexToHash("0x123456789").Hex()
		newHeader := core.Header{
			BlockNumber: blockNumber,
			Hash:        newHash,
		}
		headers := []core.Header{newHeader}
		blockChain := fakes.NewMockBlockChainWithHeaders(headers)
		validator := history.NewHeaderValidator(blockChain, headerRepository, 1)

		validator.ValidateHeaders()

		dbHeader, _ := headerRepository.GetHeader(blockNumber)
		Expect(dbHeader.Hash).NotTo(Equal(oldHash))
		Expect(dbHeader.Hash).To(Equal(newHash))
	})
})
