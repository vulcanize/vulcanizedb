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

package mocks

import (
	"errors"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
)

// BackFiller mock for tests
type BackFiller struct {
	StorageDiffsToReturn []utils.StorageDiff
	BackFillErrs         []error
	PassedEndingBlock    uint64
}

// SetStorageDiffsToReturn for tests
func (backFiller *BackFiller) SetStorageDiffsToReturn(diffs []utils.StorageDiff) {
	backFiller.StorageDiffsToReturn = diffs
}

// BackFill mock method
func (backFiller *BackFiller) BackFill(startingBlock, endingBlock uint64, backFill chan utils.StorageDiff, errChan chan error, done chan bool) error {
	if endingBlock < startingBlock {
		return errors.New("backfill: ending block number needs to be greater than starting block number")
	}
	backFiller.PassedEndingBlock = endingBlock
	go func(backFill chan utils.StorageDiff, errChan chan error, done chan bool) {
		errLen := len(backFiller.BackFillErrs)
		for i, diff := range backFiller.StorageDiffsToReturn {
			if i < errLen {
				err := backFiller.BackFillErrs[i]
				if err != nil {
					errChan <- err
					continue
				}
			}
			backFill <- diff
		}
		done <- true
	}(backFill, errChan, done)
	return nil
}
