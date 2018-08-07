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

package every_block

import (
	"errors"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type FlipKickTransformer struct {
	Fetcher    LogFetcher
	Converter  Converter
	Repository Repository
	Config     TransformerConfig
}

type FlipKickTransformerInitializer struct {
	Config TransformerConfig
}

func (i FlipKickTransformerInitializer) NewFlipKickTransformer(db *postgres.DB, blockChain core.BlockChain) shared.Transformer {
	fetcher := NewFetcher(blockChain)
	repository := NewFlipKickRepository(db)
	transformer := FlipKickTransformer{
		Fetcher:    fetcher,
		Repository: repository,
		Converter:  FlipKickConverter{},
		Config:     i.Config,
	}

	return transformer
}

func (fkt *FlipKickTransformer) SetConfig(config TransformerConfig) {
	fkt.Config = config
}

const (
	FetcherError     = "Error fetching FlipKick log events for block number %d: %s"
	LogToEntityError   = "Error converting eth log to FlipKick entity for block number %d: %s"
	EntityToModelError   = "Error converting eth log to FlipKick entity for block number %d: %s"
	RepositoryError  = "Error creating flip_kick record for block number %d: %s"
	TransformerError = "There has been %d error(s) transforming FlipKick event logs, see the logs for more details."
)

type transformerError struct {
	err         string
	blockNumber int64
	msg         string
}

func (te *transformerError) Error() string {
	return fmt.Sprintf(te.msg, te.blockNumber, te.err)
}

func newTransformerError(err error, blockNumber int64, msg string) error {
	e := transformerError{err.Error(), blockNumber, msg}
	log.Println(e.Error())
	return &e
}

func (fkt FlipKickTransformer) Execute() error {
	config := fkt.Config
	topics := [][]common.Hash{{common.HexToHash(FlipKickSignature)}}

	headers, err := fkt.Repository.MissingHeaders(config.StartingBlockNumber, config.EndingBlockNumber)
	if err != nil {
		log.Println("Error:", err)
		return err
	}

	log.Printf("Fetching event logs for %d headers \n", len(headers))
	var resultingErrors []error
	for _, header := range headers {
		ethLogs, err := fkt.Fetcher.FetchLogs(config.ContractAddress, topics, header.BlockNumber)
		if err != nil {
			resultingErrors = append(resultingErrors, newTransformerError(err, header.BlockNumber, FetcherError))
		}

		for _, ethLog := range ethLogs {
			entity, err := fkt.Converter.ToEntity(config.ContractAddress, config.ContractAbi, ethLog)
			if err != nil {
				resultingErrors = append(resultingErrors, newTransformerError(err, header.BlockNumber, LogToEntityError))
			}
			model, err := fkt.Converter.ToModel(*entity)
			if err != nil {
				resultingErrors = append(resultingErrors, newTransformerError(err, header.BlockNumber, EntityToModelError))
			}

			err = fkt.Repository.Create(header.Id, model)
			if err != nil {
				resultingErrors = append(resultingErrors, newTransformerError(err, header.BlockNumber, RepositoryError))
			}
		}
	}

	if len(resultingErrors) > 0 {
		for _, err := range resultingErrors {
			log.Println(err)
		}

		msg := fmt.Sprintf(TransformerError, len(resultingErrors))
		return errors.New(msg)
	}

	return nil
}
