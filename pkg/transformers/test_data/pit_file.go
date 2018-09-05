package test_data

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	ilk2 "github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/ilk"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/stability_fee"
	"math/big"
)

var (
	PitAddress = "0xff3f2400f1600f3f493a9a92704a29b96795af1a"
)

var EthPitFileIlkLog = types.Log{
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

var rawPitFileIlkLog, _ = json.Marshal(EthPitFileIlkLog)
var PitFileIlkModel = ilk2.PitFileIlkModel{
	Ilk:              "fake ilk",
	What:             "spot",
	Data:             big.NewInt(123).String(),
	TransactionIndex: EthPitFileIlkLog.TxIndex,
	Raw:              rawPitFileIlkLog,
}

var EthPitFileStabilityFeeLog = types.Log{
	Address: common.HexToAddress("0x6b59c42097e2Fff7cad96cb08cEeFd601081aD9c"),
	Topics: []common.Hash{
		common.HexToHash("0xd4e8be8300000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x6472697000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000044d4e8be8364726970000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 12,
	TxHash:      common.HexToHash("0x78cdc62316ccf8e31515d09745cc724f557569f01a557d0d09b1066bf7079fd2"),
	TxIndex:     0,
	BlockHash:   common.HexToHash("0xe3d8e458421533170871b4033f978a3793ef10b7e33a9328a13c09e2fd90208d"),
	Index:       0,
	Removed:     false,
}

var rawPitFileStabilityFeeLog, _ = json.Marshal(EthPitFileStabilityFeeLog)
var PitFileStabilityFeeModel = stability_fee.PitFileStabilityFeeModel{
	What:             "drip",
	Data:             "0x64d922894153BE9EEf7b7218dc565d1D0Ce2a092",
	TransactionIndex: 0,
	Raw:              rawPitFileStabilityFeeLog,
}
