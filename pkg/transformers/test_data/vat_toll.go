package test_data

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_toll"
	"math/big"
)

var EthVatTollLog = types.Log{
	Address: common.HexToAddress(KovanVatContractAddress),
	Topics: []common.Hash{
		common.HexToHash("0x09b7a0b500000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000000a3e37186e017747dba34042e83e3f76ad3cce9b"),
		common.HexToHash("0x00000000000000000000000000000000000000000000000000000000075bcd15"),
	},
	Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000006409b7a0b566616b6520696c6b000000000000000000000000000000000000000000000000a3e37186e017747dba34042e83e3f76ad3cce9b000000000000000000000000000000000000000000000000000000000000000000000000000000000075bcd15"),
	BlockNumber: 21,
	TxHash:      common.HexToHash("0x0d59cb158b033ffdfb9a021d1e80bfbbcd99594c62c501897ccee446bcd33828"),
	TxIndex:     2,
	BlockHash:   fakes.FakeHash,
	Index:       4,
	Removed:     false,
}

var rawVatTollLog, _ = json.Marshal(EthVatTollLog)
var VatTollModel = vat_toll.VatTollModel{
	Ilk:              "66616b6520696c6b000000000000000000000000000000000000000000000000",
	Urn:              "0000000000000000000000000a3e37186e017747dba34042e83e3f76ad3cce9b",
	Take:             big.NewInt(123456789).String(),
	TransactionIndex: EthVatTollLog.TxIndex,
	LogIndex:         EthVatTollLog.Index,
	Raw:              rawVatTollLog,
}
