package fakes

import "github.com/vulcanize/vulcanizedb/pkg/core"

type MockLightHeaderRepository struct {
}

func (*MockLightHeaderRepository) AddCheckColumn(id string) error {
	return nil
}

func (*MockLightHeaderRepository) AddCheckColumns(ids []string) error {
	panic("implement me")
}

func (*MockLightHeaderRepository) MarkHeaderChecked(headerID int64, eventID string) error {
	panic("implement me")
}

func (*MockLightHeaderRepository) MarkHeaderCheckedForAll(headerID int64, ids []string) error {
	panic("implement me")
}

func (*MockLightHeaderRepository) MarkHeadersCheckedForAll(headers []core.Header, ids []string) error {
	panic("implement me")
}

func (*MockLightHeaderRepository) MissingHeaders(startingBlockNumber int64, endingBlockNumber int64, eventID string) ([]core.Header, error) {
	panic("implement me")
}

func (*MockLightHeaderRepository) MissingMethodsCheckedEventsIntersection(startingBlockNumber, endingBlockNumber int64, methodIds, eventIds []string) ([]core.Header, error) {
	panic("implement me")
}

func (*MockLightHeaderRepository) MissingHeadersForAll(startingBlockNumber, endingBlockNumber int64, ids []string) ([]core.Header, error) {
	panic("implement me")
}

func (*MockLightHeaderRepository) CheckCache(key string) (interface{}, bool) {
	panic("implement me")
}
