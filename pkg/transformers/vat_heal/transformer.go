// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vat_heal

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type VatHealTransformer struct {
	Config     shared.TransformerConfig
	Converter  Converter
	Fetcher    shared.LogFetcher
	Repository Repository
}

type VatHealTransformerInitializer struct {
	Config shared.TransformerConfig
}

func (i VatHealTransformerInitializer) NewVatHealTransformer(db *postgres.DB, blockChain core.BlockChain) shared.Transformer {
	fetcher := shared.NewFetcher(blockChain)
	repository := NewVatHealRepository(db)
	transformer := VatHealTransformer{
		Fetcher:    fetcher,
		Repository: repository,
		Converter:  VatHealConverter{},
		Config:     i.Config,
	}

	return transformer
}

func (transformer VatHealTransformer) Execute() error {
	config := transformer.Config
	topics := [][]common.Hash{{common.HexToHash(config.Topics[0])}}
	headers, err := transformer.Repository.MissingHeaders(config.StartingBlockNumber, config.EndingBlockNumber)
	if err != nil {
		return err
	}

	for _, header := range headers {
		logs, err := transformer.Fetcher.FetchLogs(config.ContractAddresses, topics, header.BlockNumber)
		if err != nil {
			return err
		}

		if len(logs) < 1 {
			err = transformer.Repository.MarkCheckedHeader(header.Id)
		}
		if err != nil {
			return err
		}

		models, err := transformer.Converter.ToModels(logs)
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
