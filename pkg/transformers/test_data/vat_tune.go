package test_data

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_tune"
)

var EthVatTuneLog = types.Log{
	Address: common.HexToAddress("0x239E6f0AB02713f1F8AA90ebeDeD9FC66Dc96CD6"),
	Topics: []common.Hash{
		common.HexToHash("0x5dd6471a00000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x7d7bee5fcfd8028cf7b00876c5b1421c800561a6000000000000000000000000"),
		common.HexToHash("0xa3e37186e017747dba34042e83e3f76ad3cce9b0000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000c45dd6471a66616b6520696c6b0000000000000000000000000000000000000000000000007d7bee5fcfd8028cf7b00876c5b1421c800561a6000000000000000000000000a3e37186e017747dba34042e83e3f76ad3cce9b00000000000000000000000000f243e26db94b5426032e6dfa6007802dea2a61400000000000000000000000000000000000000000000000000000000000000000000000000000000075bcd15000000000000000000000000000000000000000000000000000000003ade68b1"),
	BlockNumber: 22,
	TxHash:      common.HexToHash("0x2b19ef3965420d2d5eb8792eafbd5e365b2866200cfc4473835971f6965f831c"),
	TxIndex:     3,
	BlockHash:   common.HexToHash("0xf3a4a9ecb79a721421b347424dc9fa072ca762f3ba340557b54a205f058f59a6"),
	Index:       6,
	Removed:     false,
}

var rawVatTuneLog, _ = json.Marshal(EthVatTuneLog)
var VatTuneModel = vat_tune.VatTuneModel{
	Ilk:              "fake ilk",
	Urn:              "0x7d7bEe5fCfD8028cf7b00876C5b1421c800561A6",
	V:                "0xA3E37186E017747DbA34042e83e3F76Ad3CcE9b0",
	W:                "0x0F243E26db94B5426032E6DFA6007802Dea2a614",
	Dink:             big.NewInt(123456789).String(),
	Dart:             big.NewInt(987654321).String(),
	TransactionIndex: EthVatTuneLog.TxIndex,
	Raw:              rawVatTuneLog,
}
