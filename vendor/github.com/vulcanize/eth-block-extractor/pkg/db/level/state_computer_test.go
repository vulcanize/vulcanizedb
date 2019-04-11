package level_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/eth-block-extractor/pkg/db/level"
	"github.com/vulcanize/eth-block-extractor/test_helpers"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/wrappers/core"
	state_wrapper "github.com/vulcanize/eth-block-extractor/test_helpers/mocks/wrappers/core/state"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/wrappers/trie"
)

var _ = Describe("", func() {
	It("initializes state trie at parent block's root", func() {
		chain, db, processor, trieFactory, validator := getMocks()
		fakeDB := db.CreateFakeUnderlyingDatabase()
		db.SetReturnDatabase(fakeDB)
		computer := level.NewStateComputer(chain, db, processor, trieFactory, validator)
		currentBlock, parentBlock := getFakeBlocks()

		_, err := computer.ComputeBlockStateTrie(currentBlock, parentBlock)

		Expect(err).NotTo(HaveOccurred())
		trieFactory.AssertNewStateTrieCalledWith(parentBlock.Root(), fakeDB)
	})

	It("returns error if state trie initialization fails", func() {
		chain, db, processor, trieFactory, validator := getMocks()
		trieFactory.SetReturnErr(test_helpers.FakeError)
		computer := level.NewStateComputer(chain, db, processor, trieFactory, validator)
		currentBlock, parentBlock := getFakeBlocks()

		_, err := computer.ComputeBlockStateTrie(currentBlock, parentBlock)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(test_helpers.FakeError))
	})

	It("processes the block to build the state trie", func() {
		chain, db, processor, trieFactory, validator := getMocks()
		computer := level.NewStateComputer(chain, db, processor, trieFactory, validator)
		stateTrie := state_wrapper.NewMockStateDB()
		fakeStateDB := &state.StateDB{}
		stateTrie.SetStateDB(fakeStateDB)
		trieFactory.SetStateDB(stateTrie)
		currentBlock, parentBlock := getFakeBlocks()

		_, err := computer.ComputeBlockStateTrie(currentBlock, parentBlock)

		Expect(err).NotTo(HaveOccurred())
		processor.AssertProcessCalledWith(currentBlock, fakeStateDB)
	})

	It("returns error if processing block fails", func() {
		chain, db, processor, trieFactory, validator := getMocks()
		processor.SetReturnErr(test_helpers.FakeError)
		computer := level.NewStateComputer(chain, db, processor, trieFactory, validator)
		currentBlock, parentBlock := getFakeBlocks()

		_, err := computer.ComputeBlockStateTrie(currentBlock, parentBlock)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(test_helpers.FakeError))
	})

	It("validates state computed by processing blocks", func() {
		chain, db, processor, trieFactory, validator := getMocks()
		fakeReceipts := types.Receipts{}
		processor.SetReturnReceipts(fakeReceipts)
		fakeUsedGas := uint64(1234)
		processor.SetReturnUsedGas(fakeUsedGas)
		computer := level.NewStateComputer(chain, db, processor, trieFactory, validator)
		stateTrie := state_wrapper.NewMockStateDB()
		fakeStateDB := &state.StateDB{}
		stateTrie.SetStateDB(fakeStateDB)
		trieFactory.SetStateDB(stateTrie)
		currentBlock, parentBlock := getFakeBlocks()

		_, err := computer.ComputeBlockStateTrie(currentBlock, parentBlock)

		Expect(err).NotTo(HaveOccurred())
		validator.AssertValidateStateCalledWith(currentBlock, parentBlock, fakeStateDB, fakeReceipts, fakeUsedGas)
	})

	It("returns error if validating state fails", func() {
		chain, db, processor, trieFactory, validator := getMocks()
		validator.SetReturnErr(test_helpers.FakeError)
		computer := level.NewStateComputer(chain, db, processor, trieFactory, validator)
		currentBlock, parentBlock := getFakeBlocks()

		_, err := computer.ComputeBlockStateTrie(currentBlock, parentBlock)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(test_helpers.FakeError))
	})

	It("commits validated state to memory database", func() {
		chain, db, processor, trieFactory, validator := getMocks()
		computer := level.NewStateComputer(chain, db, processor, trieFactory, validator)
		stateTrie := state_wrapper.NewMockStateDB()
		trieFactory.SetStateDB(stateTrie)
		currentBlock, parentBlock := getFakeBlocks()

		_, err := computer.ComputeBlockStateTrie(currentBlock, parentBlock)

		Expect(err).NotTo(HaveOccurred())
		stateTrie.AssertCommitCalled()
	})

	It("returns error if committing state fails", func() {
		chain, db, processor, trieFactory, validator := getMocks()
		computer := level.NewStateComputer(chain, db, processor, trieFactory, validator)
		stateTrie := state_wrapper.NewMockStateDB()
		stateTrie.SetReturnErr(test_helpers.FakeError)
		trieFactory.SetStateDB(stateTrie)
		currentBlock, parentBlock := getFakeBlocks()

		_, err := computer.ComputeBlockStateTrie(currentBlock, parentBlock)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(test_helpers.FakeError))
	})

	It("returns computed state trie root", func() {
		chain, db, processor, trieFactory, validator := getMocks()
		computer := level.NewStateComputer(chain, db, processor, trieFactory, validator)
		fakeIterator := trie.NewMockIterator(2)
		fakeIterator.SetReturnHash(test_helpers.FakeHash)
		fakeTrie := state_wrapper.NewMockTrie()
		fakeTrie.SetReturnIterator(fakeIterator)
		db.SetReturnTrie(fakeTrie)
		currentBlock, parentBlock := getFakeBlocks()

		stateRoot, err := computer.ComputeBlockStateTrie(currentBlock, parentBlock)

		Expect(err).NotTo(HaveOccurred())
		Expect(stateRoot).To(Equal(test_helpers.FakeHash))
	})
})

func getMocks() (*core.MockBlockChain, *state_wrapper.MockStateDatabase, *core.MockProcessor, *state_wrapper.MockStateDBFactory, *core.MockValidator) {
	chain := core.NewMockBlockChain()
	db := state_wrapper.NewMockStateDatabase()
	fakeDB := db.CreateFakeUnderlyingDatabase()
	db.SetReturnDatabase(fakeDB)
	fakeIterator := trie.NewMockIterator(1)
	fakeTrie := state_wrapper.NewMockTrie()
	fakeTrie.SetReturnIterator(fakeIterator)
	db.SetReturnTrie(fakeTrie)
	processor := core.NewMockProcessor()
	trieFactory := state_wrapper.NewMockStateDBFactory()
	stateTrie := state_wrapper.NewMockStateDB()
	trieFactory.SetStateDB(stateTrie)
	validator := core.NewMockValidator()
	return chain, db, processor, trieFactory, validator
}

func getFakeBlocks() (*types.Block, *types.Block) {
	currentBlockHeader := &types.Header{
		Root:   test_helpers.FakeHash,
		Number: big.NewInt(456),
	}
	currentBlock := types.NewBlockWithHeader(currentBlockHeader)
	parentBlockHeader := &types.Header{
		Root:   common.HexToHash("0x789"),
		Number: big.NewInt(457),
	}
	parentBlock := types.NewBlockWithHeader(parentBlockHeader)
	return currentBlock, parentBlock
}
