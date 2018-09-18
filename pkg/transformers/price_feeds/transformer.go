// Copyright © 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package price_feeds

import (
	"log"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type PriceFeedTransformerInitializer struct {
	Config IPriceFeedConfig
}

func (initializer PriceFeedTransformerInitializer) NewPriceFeedTransformer(db *postgres.DB, blockChain core.BlockChain) shared.Transformer {
	converter := PriceFeedConverter{}
	fetcher := NewPriceFeedFetcher(blockChain, initializer.Config.ContractAddresses)
	repository := NewPriceFeedRepository(db)
	return PriceFeedTransformer{
		Config:     initializer.Config,
		Converter:  converter,
		Fetcher:    fetcher,
		Repository: repository,
	}
}

type PriceFeedTransformer struct {
	Config     IPriceFeedConfig
	Converter  Converter
	Fetcher    IPriceFeedFetcher
	Repository IPriceFeedRepository
}

func (transformer PriceFeedTransformer) Execute() error {
	headers, err := transformer.Repository.MissingHeaders(transformer.Config.StartingBlockNumber, transformer.Config.EndingBlockNumber)
	if err != nil {
		return err
	}
	log.Printf("Fetching price feed event logs for %d headers \n", len(headers))
	for _, header := range headers {
		logs, err := transformer.Fetcher.FetchLogValues(header.BlockNumber)
		if err != nil {
			return err
		}
		if len(logs) < 1 {
			err := transformer.Repository.MarkHeaderChecked(header.Id)
			if err != nil {
				return err
			}
		}
		for _, log := range logs {
			model, err := transformer.Converter.ToModel(log, header.Id)
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
