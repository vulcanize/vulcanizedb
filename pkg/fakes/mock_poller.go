package fakes

import (
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/contract"
)

type MockPoller struct {
	ContractName string
}

func (*MockPoller) PollContract(con contract.Contract, lastBlock int64) error {
	panic("implement me")
}

func (*MockPoller) PollContractAt(con contract.Contract, blockNumber int64) error {
	panic("implement me")
}

func (poller *MockPoller) FetchContractData(contractAbi, contractAddress, method string, methodArgs []interface{}, result interface{}, blockNumber int64) error {
	if p, ok := result.(*string); ok {
		*p = poller.ContractName
	}
	return nil
}
