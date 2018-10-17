package vat_tune

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"log"
)

type VatTuneTransformerInitializer struct {
	Config shared.TransformerConfig
}

func (initializer VatTuneTransformerInitializer) NewVatTuneTransformer(db *postgres.DB, blockChain core.BlockChain) shared.Transformer {
	converter := VatTuneConverter{}
	fetcher := shared.NewFetcher(blockChain)
	repository := NewVatTuneRepository(db)
	return VatTuneTransformer{
		Config:     initializer.Config,
		Converter:  converter,
		Fetcher:    fetcher,
		Repository: repository,
	}
}

type VatTuneTransformer struct {
	Config     shared.TransformerConfig
	Converter  Converter
	Fetcher    shared.LogFetcher
	Repository Repository
}

func (transformer VatTuneTransformer) Execute() error {
	missingHeaders, err := transformer.Repository.MissingHeaders(transformer.Config.StartingBlockNumber, transformer.Config.EndingBlockNumber)
	if err != nil {
		return err
	}
	log.Printf("Fetching vat tune event logs for %d headers \n", len(missingHeaders))
	for _, header := range missingHeaders {
		topics := [][]common.Hash{{common.HexToHash(shared.VatTuneSignature)}}
		matchingLogs, err := transformer.Fetcher.FetchLogs(VatTuneConfig.ContractAddresses, topics, header.BlockNumber)
		if err != nil {
			return err
		}
		if len(matchingLogs) < 1 {
			err = transformer.Repository.MarkHeaderChecked(header.Id)
			if err != nil {
				return err
			}
		}
		models, err := transformer.Converter.ToModels(matchingLogs)
		if err != nil {
			return err
		}
		err = transformer.Repository.Create(header.Id, models)
		if err != nil {
			return err
		}
	}
	return nil
}
