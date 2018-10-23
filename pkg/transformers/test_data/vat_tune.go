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
	Address: common.HexToAddress("0xCD726790550aFcd77e9a7a47e86A3F9010af126B"),
	Topics: []common.Hash{
		common.HexToHash("0x5dd6471a00000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x4554480000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000004f26ffbe5f04ed43630fdc30a87638d53d0b0876"),
		common.HexToHash("0x0000000000000000000000004f26ffbe5f04ed43630fdc30a87638d53d0b0876"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000c45dd6471a45544800000000000000000000000000000000000000000000000000000000000000000000000000000000004f26ffbe5f04ed43630fdc30a87638d53d0b08760000000000000000000000004f26ffbe5f04ed43630fdc30a87638d53d0b08760000000000000000000000004f26ffbe5f04ed43630fdc30a87638d53d0b08760000000000000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffe43e9298b1380000"),
	BlockNumber: 8761670,
	TxHash:      common.HexToHash("0x95eb3d6cbd83032efa29714d4a391ce163d7d215db668aadd7d33dd5c20b1ec7"),
	TxIndex:     0,
	BlockHash:   common.HexToHash("0xd35188df40a741f4bf4ccacc6be0154e338652660d0b56cc9058c1bca53b09ee"),
	Index:       6,
	Removed:     false,
}

var rawVatTuneLog, _ = json.Marshal(EthVatTuneLog)
var dartString = "115792089237316195423570985008687907853269984665640564039455584007913129639936"
var vatTuneDart, _ = new(big.Int).SetString(dartString, 10)
var urnAddress = "0x4F26FfBe5F04ED43630fdC30A87638d53D0b0876"
var VatTuneModel = vat_tune.VatTuneModel{
	Ilk:              "ETH",
	Urn:              urnAddress,
	V:                urnAddress,
	W:                urnAddress,
	Dink:             big.NewInt(0).String(),
	Dart:             vatTuneDart.String(),
	TransactionIndex: EthVatTuneLog.TxIndex,
	LogIndex:         EthVatTuneLog.Index,
	Raw:              rawVatTuneLog,
}
