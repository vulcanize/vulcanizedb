// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package event

import (
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type LogNoteTransformer struct {
	Config     transformer.EventTransformerConfig
	Converter  LogNoteConverter
	Repository Repository
}

func (tr LogNoteTransformer) NewLogNoteTransformer(db *postgres.DB) transformer.EventTransformer {
	tr.Repository.SetDB(db)
	return tr
}

func (tr LogNoteTransformer) Execute(logs []types.Log, header core.Header, recheckedHeader constants.TransformerExecution) error {
	transformerName := tr.Config.TransformerName

	// No matching logs, mark the header as checked for this type of logs
	if len(logs) < 1 {
		err := tr.Repository.MarkHeaderChecked(header.Id)
		if err != nil {
			log.Printf("Error marking header as checked in %v: %v", transformerName, err)
			return err
		}
		return nil
	}

	models, err := tr.Converter.ToModels(logs)
	if err != nil {
		log.Printf("Error converting logs in %v: %v", transformerName, err)
		return err
	}

	err = tr.Repository.Create(header.Id, models)
	if err != nil {
		log.Printf("Error persisting %v record: %v", transformerName, err)
		return err
	}
	return nil
}

func (tr LogNoteTransformer) GetName() string {
	return tr.Config.TransformerName
}

func (tr LogNoteTransformer) GetConfig() transformer.EventTransformerConfig {
	return tr.Config
}
