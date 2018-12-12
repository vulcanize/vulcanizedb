package shared

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type WatcherRepository interface {
	GetCheckedColumnNames(db *postgres.DB) ([]string, error)
	CreateNotCheckedSQL(boolColumns []string) string
	MissingHeaders(startingBlockNumber int64, endingBlockNumber int64, db *postgres.DB, notCheckedSQL string) ([]core.Header, error)
}

type Watcher struct {
	Transformers []shared.Transformer
	DB           *postgres.DB
	Fetcher      shared.LogFetcher
	Chunker      shared.Chunker
	Addresses    []common.Address
	Topics       []common.Hash
	Repository   WatcherRepository
}

func NewWatcher(db *postgres.DB, fetcher shared.LogFetcher, repository WatcherRepository, chunker shared.Chunker) Watcher {
	return Watcher{
		DB:         db,
		Fetcher:    fetcher,
		Chunker:    chunker,
		Repository: repository,
	}
}

// Adds transformers to the watcher, each needs an initializer and the associated config.
// This also changes the configuration of the chunker, so that it will consider the new transformers.
func (watcher *Watcher) AddTransformers(initializers []shared.TransformerInitializer, configs []shared.TransformerConfig) {
	if len(initializers) != len(configs) {
		panic("Mismatch in number of transformers initializers and configs!")
	}

	for _, initializer := range initializers {
		transformer := initializer(watcher.DB)
		watcher.Transformers = append(watcher.Transformers, transformer)
	}

	var contractAddresses []common.Address
	var topic0s []common.Hash

	for _, config := range configs {
		for _, address := range config.ContractAddresses {
			contractAddresses = append(contractAddresses, common.HexToAddress(address))
		}
		topic0s = append(topic0s, common.HexToHash(config.Topic))
	}

	watcher.Addresses = append(watcher.Addresses, contractAddresses...)
	watcher.Topics = append(watcher.Topics, topic0s...)
	watcher.Chunker.AddConfigs(configs)
}

func (watcher *Watcher) Execute() error {
	checkedColumnNames, err := watcher.Repository.GetCheckedColumnNames(watcher.DB)
	if err != nil {
		return err
	}
	notCheckedSQL := watcher.Repository.CreateNotCheckedSQL(checkedColumnNames)

	// TODO Handle start and end numbers in transformers
	missingHeaders, err := watcher.Repository.MissingHeaders(0, -1, watcher.DB, notCheckedSQL)
	if err != nil {
		log.Error("Fetching of missing headers failed in watcher!")
		return err
	}

	for _, header := range missingHeaders {
		// TODO Extend FetchLogs for doing several blocks at a time
		logs, err := watcher.Fetcher.FetchLogs(watcher.Addresses, watcher.Topics, header)
		if err != nil {
			// TODO Handle fetch error in watcher
			log.Error("Error while fetching logs for header %v in watcher", header.Id)
			return err
		}

		chunkedLogs := watcher.Chunker.ChunkLogs(logs)

		for _, transformer := range watcher.Transformers {
			logChunk := chunkedLogs[transformer.Name()]
			err = transformer.Execute(logChunk, header)
			if err != nil {
				log.Error("%v transformer failed to execute in watcher: %v", transformer.Name(), err)
				return err
			}
		}
	}
	return err
}
