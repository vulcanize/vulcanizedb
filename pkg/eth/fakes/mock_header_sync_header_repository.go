package fakes

import "github.com/vulcanize/vulcanizedb/pkg/eth/core"

type MockHeaderSyncHeaderRepository struct {
}

func (*MockHeaderSyncHeaderRepository) AddCheckColumn(id string) error {
	return nil
}

func (*MockHeaderSyncHeaderRepository) AddCheckColumns(ids []string) error {
	panic("implement me")
}

func (*MockHeaderSyncHeaderRepository) MarkHeaderChecked(headerID int64, eventID string) error {
	panic("implement me")
}

func (*MockHeaderSyncHeaderRepository) MarkHeaderCheckedForAll(headerID int64, ids []string) error {
	panic("implement me")
}

func (*MockHeaderSyncHeaderRepository) MarkHeadersCheckedForAll(headers []core.Header, ids []string) error {
	panic("implement me")
}

func (*MockHeaderSyncHeaderRepository) MissingHeaders(startingBlockNumber int64, endingBlockNumber int64, eventID string) ([]core.Header, error) {
	panic("implement me")
}

func (*MockHeaderSyncHeaderRepository) MissingMethodsCheckedEventsIntersection(startingBlockNumber, endingBlockNumber int64, methodIds, eventIds []string) ([]core.Header, error) {
	panic("implement me")
}

func (*MockHeaderSyncHeaderRepository) MissingHeadersForAll(startingBlockNumber, endingBlockNumber int64, ids []string) ([]core.Header, error) {
	panic("implement me")
}

func (*MockHeaderSyncHeaderRepository) CheckCache(key string) (interface{}, bool) {
	panic("implement me")
}
