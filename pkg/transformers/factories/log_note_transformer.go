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
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

type LogNoteTransformer struct {
	Config     shared.TransformerConfig
	Converter  LogNoteConverter
	Repository Repository
}

func (transformer LogNoteTransformer) NewLogNoteTransformer(db *postgres.DB) shared.Transformer {
	transformer.Repository.SetDB(db)
	return transformer
}

func (transformer LogNoteTransformer) Execute(logs []types.Log, missingHeaders []core.Header) error {
	transformerName := transformer.Config.TransformerName

	for _, header := range missingHeaders {
		// No matching logs, mark the header as checked for this type of logs
		if len(logs) < 1 {
			err := transformer.Repository.MarkHeaderChecked(header.Id)
			if err != nil {
				log.Printf("Error marking header as checked in %v: %v", transformerName, err)
				return err
			}
			// Continue with the next header; nothing to persist
			continue
		}

		models, err := transformer.Converter.ToModels(logs)
		if err != nil {
			log.Printf("Error converting logs in %v: %v", transformerName, err)
			return err
		}

		err = transformer.Repository.Create(header.Id, models)
		if err != nil {
			log.Printf("Error persisting %v record: %v", transformerName, err)
			return err
		}
	}
	return nil
}
