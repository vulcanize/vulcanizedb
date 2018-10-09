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
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"log"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type Transformer struct {
	Config     shared.TransformerConfig
	Converter  Converter
	Repository Repository
	Fetcher    shared.SettableLogFetcher
}

type Converter interface {
	ToModels(ethLog []types.Log) ([]interface{}, error)
}

type Repository interface {
	Create(headerID int64, models []interface{}) error
	MarkHeaderChecked(headerID int64) error
	MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error)
	SetDB(db *postgres.DB)
}

func (transformer Transformer) NewTransformer(db *postgres.DB, bc core.BlockChain) shared.Transformer {
	transformer.Repository.SetDB(db)
	transformer.Fetcher.SetBC(bc)
	return transformer
}

func (transformer Transformer) Execute() error {
	missingHeaders, err := transformer.Repository.MissingHeaders(transformer.Config.StartingBlockNumber, transformer.Config.EndingBlockNumber)
	if err != nil {
		return err
	}

	log.Printf("Fetching vat move event logs for %d headers \n", len(missingHeaders))
	for _, header := range missingHeaders {
		// TODO Needs signature in config
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

		models, err := transformer.Converter.ToModels(matchingLogs)
		if err != nil {
			return err
		}

		// Can't assert the whole collection, need to wash types individually
		var typelessModels []interface{}
		for _, m := range models {
			typelessModels = append(typelessModels, m.(interface{}))
		}

		err = transformer.Repository.Create(header.Id, typelessModels)
		if err != nil {
			return err
		}
	}
	return nil
}
