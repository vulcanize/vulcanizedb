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

package fetcher

import (
	"strings"

	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/fs"
)

type CsvTailStorageFetcher struct {
	tailer       fs.Tailer
	statusWriter fs.StatusWriter
}

func NewCsvTailStorageFetcher(tailer fs.Tailer, statusWriter fs.StatusWriter) CsvTailStorageFetcher {
	return CsvTailStorageFetcher{
		tailer:       tailer,
		statusWriter: statusWriter,
	}
}

func (storageFetcher CsvTailStorageFetcher) FetchStorageDiffs(out chan<- types.RawDiff, errs chan<- error) {
	t, tailErr := storageFetcher.tailer.Tail()
	if tailErr != nil {
		errs <- tailErr
	}
	writeErr := storageFetcher.statusWriter.Write()
	if writeErr != nil {
		errs <- writeErr
	}

	for line := range t.Lines {
		diff, parseErr := types.FromParityCsvRow(strings.Split(line.Text, ","))
		if parseErr != nil {
			errs <- parseErr
		} else {
			out <- diff
		}
	}
}
