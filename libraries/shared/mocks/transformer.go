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
	"github.com/ethereum/go-ethereum/core/types"

	shared_t "github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type MockTransformer struct {
	ExecuteWasCalled bool
	ExecuteError     error
	PassedLogs       []types.Log
	PassedHeader     core.Header
	config           shared_t.EventTransformerConfig
}

func (mh *MockTransformer) Execute(logs []types.Log, header core.Header) error {
	if mh.ExecuteError != nil {
		return mh.ExecuteError
	}
	mh.ExecuteWasCalled = true
	mh.PassedLogs = logs
	mh.PassedHeader = header
	return nil
}

func (mh *MockTransformer) GetConfig() shared_t.EventTransformerConfig {
	return mh.config
}

func (mh *MockTransformer) SetTransformerConfig(config shared_t.EventTransformerConfig) {
	mh.config = config
}

func (mh *MockTransformer) FakeTransformerInitializer(db *postgres.DB) shared_t.EventTransformer {
	return mh
}

var FakeTransformerConfig = shared_t.EventTransformerConfig{
	TransformerName:   "FakeTransformer",
	ContractAddresses: []string{"FakeAddress"},
	Topic:             "FakeTopic",
}
