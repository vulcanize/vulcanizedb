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
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/sirupsen/logrus"
)

// Transformer implements the EventTransformer interface, to be run by the Watcher
type Transformer struct {
	Config     transformer.EventTransformerConfig
	Converter  Converter
	DB         *postgres.DB
}

// NewTransformer instantiates a new transformer by passing the DB connection to the converter
func (t Transformer) NewTransformer(db *postgres.DB) transformer.EventTransformer {
	t.Converter.SetDB(db)
	return t
}

// Execute runs a transformer on a set of logs, converting data into models and persisting to the DB
func (t Transformer) Execute(logs []core.HeaderSyncLog) error {
	transformerName := t.Config.TransformerName
	config := t.Config

	if len(logs) < 1 {
		return nil
	}

	models, err := t.Converter.ToModels(config.ContractAbi, logs)
	if err != nil {
		logrus.Errorf("error converting entities to models in %v: %v", transformerName, err)
		return err
	}

	err = PersistModels(models, t.DB)
	if err != nil {
		logrus.Errorf("error persisting %v record: %v", transformerName, err)
		return err
	}

	return nil
}

// GetConfig returns the config for a given transformer
func (t Transformer) GetConfig() transformer.EventTransformerConfig {
	return t.Config
}
