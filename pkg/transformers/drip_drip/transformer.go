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

package drip_drip

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type DripDripTransformerInitializer struct {
	Config shared.TransformerConfig
}

func (initializer DripDripTransformerInitializer) NewDripDripTransformer(db *postgres.DB, blockChain core.BlockChain) shared.Transformer {
	converter := DripDripConverter{}
	fetcher := shared.NewFetcher(blockChain)
	repository := NewDripDripRepository(db)
	return DripDripTransformer{
		Config:     initializer.Config,
		Converter:  converter,
		Fetcher:    fetcher,
		Repository: repository,
	}
}

type DripDripTransformer struct {
	Config     shared.TransformerConfig
	Converter  Converter
	Fetcher    shared.LogFetcher
	Repository Repository
}

func (transformer DripDripTransformer) Execute() error {
	missingHeaders, err := transformer.Repository.MissingHeaders(transformer.Config.StartingBlockNumber, transformer.Config.EndingBlockNumber)
	if err != nil {
		return err
	}
	for _, header := range missingHeaders {
		topics := [][]common.Hash{{common.HexToHash(shared.DripDripSignature)}}
		matchingLogs, err := transformer.Fetcher.FetchLogs(transformer.Config.ContractAddress, topics, header.BlockNumber)
		if err != nil {
			return err
		}
		for _, log := range matchingLogs {
			model, err := transformer.Converter.ToModel(log)
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
