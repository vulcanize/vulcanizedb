package vat_flux

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"log"
)

type VatFluxTransformerInitializer struct {
	Config shared.TransformerConfig
}

func (initializer VatFluxTransformerInitializer) NewVatFluxTransformer(db *postgres.DB, blockChain core.BlockChain) shared.Transformer {
	converter := VatFluxConverter{}
	fetcher := shared.NewFetcher(blockChain)
	repository := NewVatFluxRepository(db)
	return VatFluxTransformer{
		Config:     initializer.Config,
		Converter:  converter,
		Fetcher:    fetcher,
		Repository: repository,
	}
}

type VatFluxTransformer struct {
	Config     shared.TransformerConfig
	Converter  Converter
	Fetcher    shared.LogFetcher
	Repository Repository
}

func (transformer VatFluxTransformer) Execute() error {
	missingHeaders, err := transformer.Repository.MissingHeaders(transformer.Config.StartingBlockNumber, transformer.Config.EndingBlockNumber)
	if err != nil {
		return err
	}
	log.Printf("Fetching vat flux event logs for %d headers \n", len(missingHeaders))
	for _, header := range missingHeaders {
		topics := [][]common.Hash{{common.HexToHash(shared.VatFluxSignature)}}
		matchingLogs, err := transformer.Fetcher.FetchLogs(VatFluxConfig.ContractAddresses, topics, header.BlockNumber)
		if err != nil {
			return err
		}
		if len(matchingLogs) < 1 {
			err = transformer.Repository.MarkCheckedHeader(header.Id)
			if err != nil {
				return err
			}
			continue
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
