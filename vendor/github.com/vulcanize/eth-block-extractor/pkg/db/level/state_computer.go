package level

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/eth-block-extractor/pkg/wrappers/core"
	"github.com/vulcanize/eth-block-extractor/pkg/wrappers/core/state"
)

type IStateComputer interface {
	ComputeBlockStateTrie(currentBlock *types.Block, parentBlock *types.Block) (root common.Hash, err error)
}

type StateComputer struct {
	blockChain     core.GethCoreBlockChain
	db             state.GethStateDatabase
	processor      core.GethStateProcessor
	stateDBFactory state.GethStateDBFactory
	validator      core.GethBlockValidator
}

func NewStateComputer(blockChain core.GethCoreBlockChain, db state.GethStateDatabase, processor core.GethStateProcessor, stateDBFactory state.GethStateDBFactory, validator core.GethBlockValidator) *StateComputer {
	return &StateComputer{
		blockChain:     blockChain,
		db:             db,
		processor:      processor,
		stateDBFactory: stateDBFactory,
		validator:      validator,
	}
}

func (sc *StateComputer) ComputeBlockStateTrie(block, parent *types.Block) (root common.Hash, err error) {
	stateTrie, err := sc.stateDBFactory.NewStateDB(parent.Root(), sc.db.Database())
	if err != nil {
		return root, err
	}
	return sc.createStateTrieForBlock(block, parent, stateTrie)
}

func (sc *StateComputer) createStateTrieForBlock(block, parent *types.Block, stateTrie state.GethStateDB) (root common.Hash, err error) {
	receipts, _, usedGas, err := sc.processor.Process(block, stateTrie.StateDB())
	if err != nil {
		return root, err
	}
	err = sc.validator.ValidateState(block, parent, stateTrie.StateDB(), receipts, usedGas)
	if err != nil {
		return root, err
	}
	return stateTrie.Commit(sc.blockChain.Config().IsEIP158(block.Number()))
}
