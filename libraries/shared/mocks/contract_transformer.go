package mocks

import (
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type MockContractTransformer struct {
	ExecuteWasCalled bool
	ExecuteError     error
	Config           config.ContractConfig
}

func (mc *MockContractTransformer) Init() error {
	return nil
}

func (mc *MockContractTransformer) Execute() error {
	if mc.ExecuteError != nil {
		return mc.ExecuteError
	}
	mc.ExecuteWasCalled = true
	return nil
}

func (mc *MockContractTransformer) GetConfig() config.ContractConfig {
	return mc.Config
}

func (mc *MockContractTransformer) FakeTransformerInitializer(db *postgres.DB, bc core.BlockChain) transformer.ContractTransformer {
	return mc
}
