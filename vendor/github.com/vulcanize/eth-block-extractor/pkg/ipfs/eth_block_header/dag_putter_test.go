package eth_block_header_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_block_header"
	"github.com/vulcanize/eth-block-extractor/test_helpers"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/ipfs"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/wrappers/rlp"
)

var _ = Describe("Creating an IPLD for a block header", func() {
	It("decodes passed bytes into ethereum block header", func() {
		mockDecoder := rlp.NewMockDecoder()
		mockDecoder.SetReturnOut(&types.Header{})
		dagPutter := eth_block_header.NewBlockHeaderDagPutter(ipfs.NewMockAdder(), mockDecoder)
		fakeBytes := []byte{1, 2, 3, 4, 5}

		_, err := dagPutter.DagPut(fakeBytes)

		Expect(err).NotTo(HaveOccurred())
		mockDecoder.AssertDecodeCalledWith(fakeBytes, &types.Header{})
	})

	It("returns error if decoding fails", func() {
		mockDecoder := rlp.NewMockDecoder()
		mockDecoder.SetReturnOut(&types.Header{})
		mockDecoder.SetError(test_helpers.FakeError)
		dagPutter := eth_block_header.NewBlockHeaderDagPutter(ipfs.NewMockAdder(), mockDecoder)

		_, err := dagPutter.DagPut([]byte{1, 2, 3, 4, 5})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(test_helpers.FakeError))
	})

	It("adds ethereum block header to ipfs", func() {
		mockAdder := ipfs.NewMockAdder()
		mockDecoder := rlp.NewMockDecoder()
		mockDecoder.SetReturnOut(&types.Header{})
		dagPutter := eth_block_header.NewBlockHeaderDagPutter(mockAdder, mockDecoder)
		fakeBytes := []byte{1, 2, 3, 4, 5}

		_, err := dagPutter.DagPut(fakeBytes)

		Expect(err).NotTo(HaveOccurred())
		mockAdder.AssertAddCalled(1, &eth_block_header.EthBlockHeaderNode{})
	})

	It("returns error if adding to ipfs fails", func() {
		mockAdder := ipfs.NewMockAdder()
		mockAdder.SetError(test_helpers.FakeError)
		mockDecoder := rlp.NewMockDecoder()
		mockDecoder.SetReturnOut(&types.Header{})
		dagPutter := eth_block_header.NewBlockHeaderDagPutter(mockAdder, mockDecoder)

		_, err := dagPutter.DagPut([]byte{1, 2, 3, 4, 5})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(test_helpers.FakeError))
	})
})
