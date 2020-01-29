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

package mocks

import (
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/eth/core"
	"github.com/vulcanize/vulcanizedb/pkg/eth/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/eth/fakes"
)

type MockEventTransformer struct {
	ExecuteWasCalled bool
	ExecuteError     error
	PassedLogs       []core.HeaderSyncLog
	config           transformer.EventTransformerConfig
}

func (t *MockEventTransformer) Execute(logs []core.HeaderSyncLog) error {
	if t.ExecuteError != nil {
		return t.ExecuteError
	}
	t.ExecuteWasCalled = true
	t.PassedLogs = logs
	return nil
}

func (t *MockEventTransformer) GetConfig() transformer.EventTransformerConfig {
	return t.config
}

func (t *MockEventTransformer) SetTransformerConfig(config transformer.EventTransformerConfig) {
	t.config = config
}

func (t *MockEventTransformer) FakeTransformerInitializer(db *postgres.DB) transformer.EventTransformer {
	return t
}

var FakeTransformerConfig = transformer.EventTransformerConfig{
	TransformerName:   "FakeTransformer",
	ContractAddresses: []string{fakes.FakeAddress.Hex()},
	Topic:             fakes.FakeHash.Hex(),
}
