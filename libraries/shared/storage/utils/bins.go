// VulcanizeDB
// Copyright © 2019 Vulcanize

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

package utils

import (
	"errors"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

func GetBlockHeightBins(startingBlock, endingBlock, batchSize uint64) ([][]uint64, error) {
	if endingBlock < startingBlock {
		return nil, errors.New("backfill: ending block number needs to be greater than starting block number")
	}
	if batchSize == 0 {
		return nil, errors.New("backfill: batchsize needs to be greater than zero")
	}
	length := endingBlock - startingBlock + 1
	numberOfBins := length / batchSize
	remainder := length % batchSize
	if remainder != 0 {
		numberOfBins++
	}
	blockRangeBins := make([][]uint64, numberOfBins)
	for i := range blockRangeBins {
		nextBinStart := startingBlock + batchSize
		blockRange := make([]uint64, 0, nextBinStart-startingBlock+1)
		for j := startingBlock; j < nextBinStart && j <= endingBlock; j++ {
			blockRange = append(blockRange, j)
		}
		startingBlock = nextBinStart
		blockRangeBins[i] = blockRange
	}
	return blockRangeBins, nil
}

func MissingHeightsToGaps(heights []uint64) []shared.Gap {
	validationGaps := make([]shared.Gap, 0)
	start := heights[0]
	lastHeight := start
	for i, height := range heights[1:] {
		if height != lastHeight+1 {
			validationGaps = append(validationGaps, shared.Gap{
				Start: start,
				Stop:  lastHeight,
			})
			start = height
		}
		if i+2 == len(heights) {
			validationGaps = append(validationGaps, shared.Gap{
				Start: start,
				Stop:  height,
			})
		}
		lastHeight = height
	}
	return validationGaps
}
