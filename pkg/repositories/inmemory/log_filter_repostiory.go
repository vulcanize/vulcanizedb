package inmemory

import (
	"errors"

	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

func (repository *InMemory) GetFilter(name string) (filters.LogFilter, error) {
	panic("implement me")
}

func (repository *InMemory) CreateFilter(filter filters.LogFilter) error {
	key := filter.Name
	if _, ok := repository.logFilters[key]; ok || key == "" {
		return errors.New("filter name not unique")
	}
	repository.logFilters[key] = filter
	return nil
}
