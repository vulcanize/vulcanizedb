package fakes

import "github.com/makerdao/vulcanizedb/pkg/core"

type MockCheckedHeaderRepository struct {
}

func (*MockCheckedHeaderRepository) AddCheckColumn(id string) error {
	return nil
}

func (*MockCheckedHeaderRepository) AddCheckColumns(ids []string) error {
	panic("implement me")
}

func (*MockCheckedHeaderRepository) MarkHeaderChecked(headerID int64, eventID string) error {
	panic("implement me")
}

func (*MockCheckedHeaderRepository) MarkHeaderCheckedForAll(headerID int64, ids []string) error {
	panic("implement me")
}

func (*MockCheckedHeaderRepository) MarkHeadersCheckedForAll(headers []core.Header, ids []string) error {
	panic("implement me")
}

func (*MockCheckedHeaderRepository) MissingHeaders(startingBlockNumber int64, endingBlockNumber int64, eventID string) ([]core.Header, error) {
	panic("implement me")
}

func (*MockCheckedHeaderRepository) MissingMethodsCheckedEventsIntersection(startingBlockNumber, endingBlockNumber int64, methodIds, eventIds []string) ([]core.Header, error) {
	panic("implement me")
}

func (*MockCheckedHeaderRepository) MissingHeadersForAll(startingBlockNumber, endingBlockNumber int64, ids []string) ([]core.Header, error) {
	panic("implement me")
}

func (*MockCheckedHeaderRepository) CheckCache(key string) (interface{}, bool) {
	panic("implement me")
}
