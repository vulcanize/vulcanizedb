package shared

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type Watcher struct {
	Transformers []shared.Transformer
	DB           postgres.DB
	Blockchain   core.BlockChain
	Fetcher      shared.LogFetcher
	Chunker      shared.LogChunker
	Addresses    []common.Address
	Topics       []common.Hash
}

func NewWatcher(db postgres.DB, bc core.BlockChain) Watcher {
	transformerConfigs := transformers.TransformerConfigs()
	var contractAddresses []common.Address
	var topic0s []common.Hash

	for _, config := range transformerConfigs {
		for _, address := range config.ContractAddresses {
			contractAddresses = append(contractAddresses, common.HexToAddress(address))
		}
		topic0s = append(topic0s, common.HexToHash(config.Topic))
	}

	chunker := shared.NewLogChunker(transformerConfigs)
	fetcher := shared.NewFetcher(bc)

	return Watcher{
		DB: db,
		Blockchain: bc,
		Fetcher: fetcher,
		Chunker: chunker,
		Addresses: contractAddresses,
		Topics: topic0s,
	}
}

func (watcher *Watcher) AddTransformers(us []shared.TransformerInitializer) {
	for _, transformerInitializer := range us {
		transformer := transformerInitializer(&watcher.DB)
		watcher.Transformers = append(watcher.Transformers, transformer)
	}
}

func (watcher *Watcher) Execute() error {
	checkedColumnNames, err := shared.GetCheckedColumnNames(&watcher.DB)
	if err != nil {
		return err
	}
	notCheckedSQL := shared.CreateNotCheckedSQL(checkedColumnNames)

	// TODO Handle start and end numbers in transformers?
	missingHeaders, err := shared.MissingHeaders(0, -1, &watcher.DB, notCheckedSQL)

	for _, header := range missingHeaders {
		// TODO Extend FetchLogs for doing several blocks at a time
		logs, err := watcher.Fetcher.FetchLogs(watcher.Addresses, watcher.Topics, header)
		if err != nil {
			// TODO Handle fetch error in watcher
			return err
		}

		chunkedLogs := watcher.Chunker.ChunkLogs(logs)

		for _, transformer := range watcher.Transformers {
			logChunk := chunkedLogs[transformer.Name()]
			err = transformer.Execute(logChunk, header)
		}
	}
	return err
}
