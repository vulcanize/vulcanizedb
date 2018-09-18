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

package vow

import (
	"log"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_file"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type DripFileVowTransformerInitializer struct {
	Config shared.TransformerConfig
}

func (initializer DripFileVowTransformerInitializer) NewDripFileVowTransformer(db *postgres.DB, blockChain core.BlockChain) shared.Transformer {
	converter := DripFileVowConverter{}
	fetcher := shared.NewFetcher(blockChain)
	repository := NewDripFileVowRepository(db)
	return DripFileVowTransformer{
		Config:     initializer.Config,
		Converter:  converter,
		Fetcher:    fetcher,
		Repository: repository,
	}
}

type DripFileVowTransformer struct {
	Config     shared.TransformerConfig
	Converter  Converter
	Fetcher    shared.LogFetcher
	Repository Repository
}

func (transformer DripFileVowTransformer) Execute() error {
	missingHeaders, err := transformer.Repository.MissingHeaders(transformer.Config.StartingBlockNumber, transformer.Config.EndingBlockNumber)
	if err != nil {
		return err
	}
	log.Printf("Fetching drip file vow event logs for %d headers \n", len(missingHeaders))
	for _, header := range missingHeaders {
		topics := [][]common.Hash{{common.HexToHash(shared.DripFileVowSignature)}}
		matchingLogs, err := transformer.Fetcher.FetchLogs(drip_file.DripFileConfig.ContractAddress, topics, header.BlockNumber)
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
