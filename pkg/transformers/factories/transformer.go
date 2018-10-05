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

package factories

import (
	"github.com/ethereum/go-ethereum/core/types"
	"log"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type TransformerInitializer struct {
	Config 		shared.TransformerConfig
	Converter	Converter
	Repository	Repository
}

type Model struct {
	TransactionIndex uint   `db:"tx_idx"`
	Raw              []byte `db:"raw_log"`
}

type Converter interface {
	ToModel(ethLog types.Log) (Model, error)
}

type Repository interface {
	Create(headerID int64, models []Model) error
	MarkHeaderChecked(headerID int64) error
	MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error)
}


func (initializer TransformerInitializer) NewTransformer(db *postgres.DB, blockChain core.BlockChain) shared.Transformer {
	fetcher := shared.NewFetcher(blockChain)
	return Transformer{
		Config:     initializer.Config,
		Converter:  initializer.Converter,
		Fetcher:    fetcher,
		Repository: initializer.Repository,
	}
}

type Transformer struct {
	Config     shared.TransformerConfig
	Converter  Converter
	Fetcher    shared.LogFetcher
	Repository Repository
}

func (transformer Transformer) Execute() error {
	missingHeaders, err := transformer.Repository.MissingHeaders(transformer.Config.StartingBlockNumber, transformer.Config.EndingBlockNumber)
	if err != nil {
		return err
	}

	log.Printf("Fetching vat move event logs for %d headers \n", len(missingHeaders))
	for _, header := range missingHeaders {
		topics := [][]common.Hash{{common.HexToHash(shared.VatMoveSignature)}}
		matchingLogs, err := transformer.Fetcher.FetchLogs(transformer.Config.ContractAddresses, topics, header.BlockNumber)
		if err != nil {
			return err
		}

		if len(matchingLogs) < 1 {
			err := transformer.Repository.MarkHeaderChecked(header.Id)
			if err != nil {
				return err
			}
		}

		for _, log := range matchingLogs {
			model, err := transformer.Converter.ToModel(log)
			if err != nil {
				return err
			}

			err = transformer.Repository.Create(header.Id, []Model{model})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
