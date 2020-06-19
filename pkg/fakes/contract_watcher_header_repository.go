package fakes

import "github.com/makerdao/vulcanizedb/pkg/core"

type MockContractWatcherHeaderRepository struct {
}

func (*MockContractWatcherHeaderRepository) AddCheckColumn(id string) error {
	return nil
}

func (*MockContractWatcherHeaderRepository) AddCheckColumns(ids []string) error {
	panic("implement me")
}

func (*MockContractWatcherHeaderRepository) MarkHeaderChecked(headerID int64, eventID string) error {
	panic("implement me")
}

func (*MockContractWatcherHeaderRepository) MarkHeaderCheckedForAll(headerID int64, ids []string) error {
	panic("implement me")
}

func (*MockContractWatcherHeaderRepository) MarkHeadersCheckedForAll(headers []core.Header, ids []string) error {
	panic("implement me")
}

func (*MockContractWatcherHeaderRepository) MissingHeaders(startingBlockNumber int64, endingBlockNumber int64, eventID string) ([]core.Header, error) {
	panic("implement me")
}

func (*MockContractWatcherHeaderRepository) MissingHeadersForAll(startingBlockNumber, endingBlockNumber int64, ids []string) ([]core.Header, error) {
	panic("implement me")
}

func (*MockContractWatcherHeaderRepository) CheckCache(key string) (interface{}, bool) {
	panic("implement me")
}
