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

package deal

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"log"
)

type DealTransformer struct {
	Config     shared.TransformerConfig
	Converter  Converter
	Fetcher    shared.LogFetcher
	Repository Repository
}

type DealTransformerInitializer struct {
	Config shared.TransformerConfig
}

func (i DealTransformerInitializer) NewDealTransformer(db *postgres.DB, blockChain core.BlockChain) shared.Transformer {
	converter := NewDealConverter()
	fetcher := shared.NewFetcher(blockChain)
	repository := NewDealRepository(db)
	return DealTransformer{
		Config:     i.Config,
		Converter:  converter,
		Fetcher:    fetcher,
		Repository: repository,
	}
}

func (t DealTransformer) Execute() error {
	config := t.Config
	topics := [][]common.Hash{{common.HexToHash(shared.DealSignature)}}

	headers, err := t.Repository.MissingHeaders(config.StartingBlockNumber, config.EndingBlockNumber)
	if err != nil {
		return err
	}

	for _, header := range headers {
		ethLogs, err := t.Fetcher.FetchLogs(config.ContractAddress, topics, header.BlockNumber)
		if err != nil {
			log.Println("Error fetching deal logs:", err)
			return err
		}
		for _, ethLog := range ethLogs {
			model, err := t.Converter.ToModel(ethLog)
			if err != nil {
				log.Println("Error converting deal log", err)
				return err
			}
			err = t.Repository.Create(header.Id, model)
			if err != nil {
				log.Println("Error persisting deal record", err)
				return err
			}
		}
	}
	return err
}
