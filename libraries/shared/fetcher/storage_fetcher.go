package fetcher

import (
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fs"
	"strings"
)

type IStorageFetcher interface {
	FetchStorageDiffs(chan<- utils.StorageDiffRow, chan<- error)
}

type CsvTailStorageFetcher struct {
	tailer fs.Tailer
}

func NewCsvTailStorageFetcher(tailer fs.Tailer) CsvTailStorageFetcher {
	return CsvTailStorageFetcher{tailer: tailer}
}

func (storageFetcher CsvTailStorageFetcher) FetchStorageDiffs(out chan<- utils.StorageDiffRow, errs chan<- error) {
	t, tailErr := storageFetcher.tailer.Tail()
	if tailErr != nil {
		errs <- tailErr
	}
	for line := range t.Lines {
		row, parseErr := utils.FromStrings(strings.Split(line.Text, ","))
		if parseErr != nil {
			errs <- parseErr
		} else {
			out <- row
		}
	}
}
