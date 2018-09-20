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

package tend

import (
	"log"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type TendTransformer struct {
	Repository Repository
	Fetcher    shared.LogFetcher
	Converter  Converter
	Config     shared.TransformerConfig
}

type TendTransformerInitializer struct {
	Config shared.TransformerConfig
}

func (i TendTransformerInitializer) NewTendTransformer(db *postgres.DB, blockChain core.BlockChain) shared.Transformer {
	converter := NewTendConverter()
	fetcher := shared.NewFetcher(blockChain)
	repository := NewTendRepository(db)
	return TendTransformer{
		Fetcher:    fetcher,
		Repository: repository,
		Converter:  converter,
		Config:     i.Config,
	}
}

func (t TendTransformer) Execute() error {
	config := t.Config
	topics := [][]common.Hash{{common.HexToHash(shared.TendFunctionSignature)}}

	missingHeaders, err := t.Repository.MissingHeaders(config.StartingBlockNumber, config.EndingBlockNumber)
	if err != nil {
		log.Println("Error fetching missing headers:", err)
		return err
	}

	log.Printf("Fetching tend event logs for %d headers \n", len(missingHeaders))
	for _, header := range missingHeaders {
		ethLogs, err := t.Fetcher.FetchLogs(config.ContractAddress, topics, header.BlockNumber)
		if err != nil {
			log.Println("Error fetching matching logs:", err)
			return err
		}
		if len(ethLogs) < 1 {
			err := t.Repository.MarkHeaderChecked(header.Id)
			if err != nil {
				return err
			}
		}

		models, err := t.Converter.Convert(ethLogs)
		if err != nil {
			log.Println("Error converting logs:", err)
			return err
		}

		err = t.Repository.Create(header.Id, models)
		if err != nil {
			log.Println("Error persisting tend record:", err)
			return err
		}
	}

	return nil
}
