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

import "github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"

// ClosingStorageFetcher is a mock fetcher for use in tests without backfilling
type ClosingStorageFetcher struct {
	DiffsToReturn []utils.StorageDiff
	ErrsToReturn  []error
}

// NewClosingStorageFetcher returns a new ClosingStorageFetcher
func NewClosingStorageFetcher() *ClosingStorageFetcher {
	return &ClosingStorageFetcher{}
}

// FetchStorageDiffs mock method
func (fetcher *ClosingStorageFetcher) FetchStorageDiffs(out chan<- utils.StorageDiff, errs chan<- error) {
	defer close(out)
	defer close(errs)
	for _, err := range fetcher.ErrsToReturn {
		errs <- err
	}
	for _, diff := range fetcher.DiffsToReturn {
		out <- diff
	}
}

// StorageFetcher is a mock fetcher for use in tests with backfilling
type StorageFetcher struct {
	DiffsToReturn []utils.StorageDiff
	ErrsToReturn  []error
}

// NewStorageFetcher returns a new StorageFetcher
func NewStorageFetcher() *StorageFetcher {
	return &StorageFetcher{}
}

// FetchStorageDiffs mock method
func (fetcher *StorageFetcher) FetchStorageDiffs(out chan<- utils.StorageDiff, errs chan<- error) {
	for _, err := range fetcher.ErrsToReturn {
		errs <- err
	}
	for _, diff := range fetcher.DiffsToReturn {
		out <- diff
	}
}
