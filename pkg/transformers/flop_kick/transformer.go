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

package flop_kick

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type Transformer struct {
	Config     shared.TransformerConfig
	Converter  Converter
	Fetcher    shared.LogFetcher
	Repository Repository
}

type FlopKickTransformerInitializer struct {
	Config shared.TransformerConfig
}

func (i FlopKickTransformerInitializer) NewFlopKickTransformer(db *postgres.DB, blockChain core.BlockChain) shared.Transformer {
	return Transformer{
		Config:     i.Config,
		Converter:  FlopKickConverter{},
		Fetcher:    shared.NewFetcher(blockChain),
		Repository: NewFlopKickRepository(db),
	}
}

func (t Transformer) Execute() error {
	config := t.Config
	headers, err := t.Repository.MissingHeaders(config.StartingBlockNumber, config.EndingBlockNumber)
	if err != nil {
		return err
	}

	for _, header := range headers {
		topics := [][]common.Hash{{common.HexToHash(shared.FlopKickSignature)}}
		matchingLogs, err := t.Fetcher.FetchLogs(config.ContractAddress, topics, header.BlockNumber)
		if err != nil {
			return err
		}
		if len(matchingLogs) < 1 {
			err := t.Repository.MarkHeaderChecked(header.Id)
			if err != nil {
				return err
			}
		}

		entities, err := t.Converter.ToEntities(config.ContractAddress, config.ContractAbi, matchingLogs)
		if err != nil {
			return err
		}

		models, err := t.Converter.ToModels(entities)
		if err != nil {
			return err
		}

		err = t.Repository.Create(header.Id, models)
		if err != nil {
			return err
		}
	}
	return nil
}
