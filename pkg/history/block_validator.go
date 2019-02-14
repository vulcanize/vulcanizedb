// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package history

import (
	"github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
)

type BlockValidator struct {
	blockchain      core.BlockChain
	blockRepository datastore.BlockRepository
	windowSize      int
}

func NewBlockValidator(blockchain core.BlockChain, blockRepository datastore.BlockRepository, windowSize int) *BlockValidator {
	return &BlockValidator{
		blockchain:      blockchain,
		blockRepository: blockRepository,
		windowSize:      windowSize,
	}
}

func (bv BlockValidator) ValidateBlocks() (ValidationWindow, error) {
	window, err := MakeValidationWindow(bv.blockchain, bv.windowSize)
	if err != nil {
		logrus.Error("ValidateBlocks: error creating validation window: ", err)
		return ValidationWindow{}, err
	}

	blockNumbers := MakeRange(window.LowerBound, window.UpperBound)
	_, err = RetrieveAndUpdateBlocks(bv.blockchain, bv.blockRepository, blockNumbers)
	if err != nil {
		logrus.Error("ValidateBlocks: error getting and updating blocks: ", err)
		return ValidationWindow{}, err
	}

	lastBlock, err := bv.blockchain.LastBlock()
	if err != nil {
		logrus.Error("ValidateBlocks: error getting last block: ", err)
		return ValidationWindow{}, err
	}

	err = bv.blockRepository.SetBlocksStatus(lastBlock.Int64())
	if err != nil {
		logrus.Error("ValidateBlocks: error setting block status: ", err)
		return ValidationWindow{}, err
	}
	return window, nil
}
