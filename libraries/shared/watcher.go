package shared

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type Watcher struct {
	Transformers []shared.Transformer
	DB           *postgres.DB
	Fetcher      shared.LogFetcher
	Chunker      shared.Chunker
	Addresses    []common.Address
	Topics       []common.Hash
}

func NewWatcher(db *postgres.DB, bc core.BlockChain) Watcher {
	chunker := shared.NewLogChunker()
	fetcher := shared.NewFetcher(bc)
	return Watcher{
		DB:      db,
		Fetcher: fetcher,
		Chunker: chunker,
	}
}

// Adds transformers to the watcher and updates the chunker, so that it will consider the new transformers.
func (watcher *Watcher) AddTransformers(initializers []shared.TransformerInitializer) {
	var contractAddresses []common.Address
	var topic0s []common.Hash
	var configs []shared.TransformerConfig

	for _, initializer := range initializers {
		transformer := initializer(watcher.DB)
		watcher.Transformers = append(watcher.Transformers, transformer)

		config := transformer.GetConfig()
		configs = append(configs, config)

		addresses := shared.HexStringsToAddresses(config.ContractAddresses)
		contractAddresses = append(contractAddresses, addresses...)
		topic0s = append(topic0s, common.HexToHash(config.Topic))
	}

	watcher.Addresses = append(watcher.Addresses, contractAddresses...)
	watcher.Topics = append(watcher.Topics, topic0s...)
	watcher.Chunker.AddConfigs(configs)
}

func (watcher *Watcher) Execute() error {
	checkedColumnNames, err := shared.GetCheckedColumnNames(watcher.DB)
	if err != nil {
		return err
	}
	notCheckedSQL := shared.CreateNotCheckedSQL(checkedColumnNames)

	// TODO Handle start and end numbers in transformers
	missingHeaders, err := shared.MissingHeaders(0, -1, watcher.DB, notCheckedSQL)
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

		// Can't quit early and mark as checked if there are no logs. If we are running continuousLogSync,
		// not all logs we're interested in might have been fetched.
		for _, transformer := range watcher.Transformers {
			transformerName := transformer.GetConfig().TransformerName
			logChunk := chunkedLogs[transformerName]
			err = transformer.Execute(logChunk, header)
			if err != nil {
				log.Error("%v transformer failed to execute in watcher: %v", transformerName, err)
				return err
			}
		}
	}
	return err
}
