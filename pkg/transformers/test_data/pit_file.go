package test_data

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file"
	"math/big"
)

var (
	PitAddress = "0xff3f2400f1600f3f493a9a92704a29b96795af1a"
)

var EthPitFileLog = types.Log{
	Address: common.HexToAddress(PitAddress),
	Topics: []common.Hash{
		common.HexToHash("0x1a0b287e00000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000000f243e26db94b5426032e6dfa6007802dea2a614"),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x73706f7400000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000641a0b287e66616b6520696c6b00000000000000000000000000000000000000000000000073706f7400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007b"),
	BlockNumber: 11,
	TxHash:      common.HexToHash("0x1ba8125f60fa045c85b35df3983bee37db8627fbc32e3442a5cf17c85bb83f09"),
	TxIndex:     111,
	BlockHash:   common.HexToHash("0x6dc284247c524b22b10a75ef1c9d1709a509208d04c15fa2b675a293db637d21"),
	Index:       0,
	Removed:     false,
}

var raw, _ = json.Marshal(EthPitFileLog)

var PitFileModel = pit_file.PitFileModel{
	Ilk:              "fake ilk",
	What:             "spot",
	Risk:             big.NewInt(123).String(),
	TransactionIndex: EthPitFileLog.TxIndex,
	Raw:              raw,
}
