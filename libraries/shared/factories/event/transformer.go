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

type Transformer struct {
	Config     transformer.EventTransformerConfig
	Converter  Converter
	Repository Repository
}

func (transformer Transformer) NewTransformer(db *postgres.DB) transformer.EventTransformer {
	transformer.Converter.SetDB(db)
	transformer.Repository.SetDB(db)
	return transformer
}

func (transformer Transformer) Execute(logs []core.HeaderSyncLog) error {
	transformerName := transformer.Config.TransformerName
	config := transformer.Config

	if len(logs) < 1 {
		return nil
	}

	models, err := transformer.Converter.ToModels(config.ContractAbi, logs)
	if err != nil {
		logrus.Errorf("error converting entities to models in %v: %v", transformerName, err)
		return err
	}

	err = transformer.Repository.Create(models)
	if err != nil {
		logrus.Errorf("error persisting %v record: %v", transformerName, err)
		return err
	}

	return nil
}

func (transformer Transformer) GetName() string {
	return transformer.Config.TransformerName
}

func (transformer Transformer) GetConfig() transformer.EventTransformerConfig {
	return transformer.Config
}
