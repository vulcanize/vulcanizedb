package ilk

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type PitFileIlkTransformer struct {
	Config     shared.TransformerConfig
	Converter  Converter
	Fetcher    shared.LogFetcher
	Repository Repository
}

type PitFileIlkTransformerInitializer struct {
	Config shared.TransformerConfig
}

func (initializer PitFileIlkTransformerInitializer) NewPitFileIlkTransformer(db *postgres.DB, blockChain core.BlockChain) shared.Transformer {
	converter := PitFileIlkConverter{}
	fetcher := shared.NewFetcher(blockChain)
	repository := NewPitFileIlkRepository(db)
	return PitFileIlkTransformer{
		Config:     initializer.Config,
		Converter:  converter,
		Fetcher:    fetcher,
		Repository: repository,
	}
}

func (transformer PitFileIlkTransformer) Execute() error {
	missingHeaders, err := transformer.Repository.MissingHeaders(transformer.Config.StartingBlockNumber, transformer.Config.EndingBlockNumber)
	if err != nil {
		return err
	}
	for _, header := range missingHeaders {
		topics := [][]common.Hash{{common.HexToHash(shared.PitFileIlkSignature)}}
		matchingLogs, err := transformer.Fetcher.FetchLogs(pit_file.PitFileConfig.ContractAddress, topics, header.BlockNumber)
		if err != nil {
			return err
		}
		for _, log := range matchingLogs {
			model, err := transformer.Converter.ToModel(pit_file.PitFileConfig.ContractAddress, shared.PitABI, log)
			if err != nil {
				return err
			}
			err = transformer.Repository.Create(header.Id, model)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
