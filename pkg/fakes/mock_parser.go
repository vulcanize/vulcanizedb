package fakes

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/types"
)

type MockParser struct {
	AbiToReturn string
	EventName   string
	Event       types.Event
}

func (*MockParser) Parse(contractAddr string) error {
	return nil
}

func (*MockParser) ParseAbiStr(abiStr string) error {
	panic("implement me")
}

func (parser *MockParser) Abi() string {
	return parser.AbiToReturn
}

func (*MockParser) ParsedAbi() abi.ABI {
	return abi.ABI{}
}

func (*MockParser) GetMethods(wanted []string) []types.Method {
	panic("implement me")
}

func (*MockParser) GetSelectMethods(wanted []string) []types.Method {
	return []types.Method{}
}

func (parser *MockParser) GetEvents(wanted []string) map[string]types.Event {
	return map[string]types.Event{parser.EventName: parser.Event}
}
