package fakes

import "github.com/vulcanize/vulcanizedb/pkg/filters"

type MockFilterRepository struct {
}

func (*MockFilterRepository) CreateFilter(filter filters.LogFilter) error {
	return nil
}

func (*MockFilterRepository) GetFilter(name string) (filters.LogFilter, error) {
	panic("implement me")
}
