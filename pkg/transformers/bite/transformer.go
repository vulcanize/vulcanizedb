/*
 *  Copyright 2018 Vulcanize
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package bite

import (
	"log"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type BiteTransformer struct {
	Repository Repository
	Fetcher    shared.LogFetcher
	Converter  Converter
	Config     shared.TransformerConfig
}

type BiteTransformerInitializer struct {
	Config shared.TransformerConfig
}

func (i BiteTransformerInitializer) NewBiteTransformer(db *postgres.DB, blockChain core.BlockChain) shared.Transformer {
	fetcher := shared.NewFetcher(blockChain)
	repository := NewBiteRepository(db)
	transformer := BiteTransformer{
		Fetcher:    fetcher,
		Repository: repository,
		Converter:  BiteConverter{},
		Config:     i.Config,
	}

	return transformer
}

func (b BiteTransformer) Execute() error {
	config := b.Config
	topics := [][]common.Hash{{common.HexToHash(shared.BiteSignature)}}

	missingHeaders, err := b.Repository.MissingHeaders(config.StartingBlockNumber, config.EndingBlockNumber)
	if err != nil {
		log.Println("Error fetching missing headers:", err)
		return err
	}

	log.Printf("Fetching bite event logs for %d headers \n", len(missingHeaders))
	for _, header := range missingHeaders {
		ethLogs, err := b.Fetcher.FetchLogs(config.ContractAddresses, topics, header.BlockNumber)
		if err != nil {
			log.Println("Error fetching matching logs:", err)
			return err
		}

		for _, ethLog := range ethLogs {
			entity, err := b.Converter.ToEntity(ethLog.Address.Hex(), config.ContractAbi, ethLog)
			model, err := b.Converter.ToModel(entity)
			if err != nil {
				log.Println("Error converting logs:", err)
				return err
			}

			err = b.Repository.Create(header.Id, model)
			if err != nil {
				log.Println("Error persisting bite record:", err)
				return err
			}
		}
	}

	return nil
}
func (b BiteTransformer) SetConfig(config shared.TransformerConfig) {
	b.Config = config
}
