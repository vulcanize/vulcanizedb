package mocks

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

type MockTransformer struct {
	ExecuteWasCalled bool
	ExecuteError     error
	PassedLogs       []types.Log
	PassedHeader     core.Header
	config           shared.TransformerConfig
}

func (mh *MockTransformer) Execute(logs []types.Log, header core.Header, recheckHeaders constants.TransformerExecution) error {
	if mh.ExecuteError != nil {
		return mh.ExecuteError
	}
	mh.ExecuteWasCalled = true
	mh.PassedLogs = logs
	mh.PassedHeader = header
	return nil
}

func (mh *MockTransformer) GetConfig() shared.TransformerConfig {
	return mh.config
}

func (mh *MockTransformer) SetTransformerConfig(config shared.TransformerConfig) {
	mh.config = config
}

func (mh *MockTransformer) FakeTransformerInitializer(db *postgres.DB) shared.Transformer {
	return mh
}

var FakeTransformerConfig = shared.TransformerConfig{
	TransformerName:   "FakeTransformer",
	ContractAddresses: []string{"FakeAddress"},
	Topic:             "FakeTopic",
}
