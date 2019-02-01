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

func (bv BlockValidator) ValidateBlocks() ValidationWindow {
	window := MakeValidationWindow(bv.blockchain, bv.windowSize)
	blockNumbers := MakeRange(window.LowerBound, window.UpperBound)
	RetrieveAndUpdateBlocks(bv.blockchain, bv.blockRepository, blockNumbers)
	lastBlock := bv.blockchain.LastBlock().Int64()
	bv.blockRepository.SetBlocksStatus(lastBlock)
	return window
}
