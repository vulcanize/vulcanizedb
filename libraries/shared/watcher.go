package shared

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Watcher struct {
	Transformers []Transformer
	DB           postgres.DB
	Blockchain   core.BlockChain
}

func (watcher *Watcher) AddTransformers(us []TransformerInitializer, con ContractConfig) error {
	for _, transformerInitializer := range us {
		transformer, err := transformerInitializer(&watcher.DB, watcher.Blockchain, con)
		if err != nil {
			return err
		}
		watcher.Transformers = append(watcher.Transformers, transformer)
	}

	return nil
}

func (watcher *Watcher) AddTransformer(t Transformer) {
	watcher.Transformers = append(watcher.Transformers, t)
}

func (watcher *Watcher) Execute() error {
	var err error
	for _, transformer := range watcher.Transformers {
		err = transformer.Execute()
	}
	return err
}
