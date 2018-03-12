package shared

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Watcher struct {
	Handlers   []Handler
	DB         postgres.DB
	Blockchain core.Blockchain
}

func (watcher *Watcher) AddHandlers(us []HandlerInitializer) {
	for _, handlerInitializer := range us {
		handler := handlerInitializer(&watcher.DB, watcher.Blockchain)
		watcher.Handlers = append(watcher.Handlers, handler)
	}
}

func (watcher *Watcher) Execute() error {
	var err error
	for _, handler := range watcher.Handlers {
		err = handler.Execute()
	}
	return err
}
