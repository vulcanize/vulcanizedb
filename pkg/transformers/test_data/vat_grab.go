package test_data

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_grab"
)

var EthVatGrabLog = types.Log{
	Address: common.HexToAddress(shared.VatContractAddress),
	Topics: []common.Hash{
		common.HexToHash("0x3690ae4c00000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x07fa9ef6609ca7921112231f8f195138ebba2977000000000000000000000000"),
		common.HexToHash("0x7340e006f4135ba6970d43bf43d88dcad4e7a8ca000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000c43690ae4c66616b6520696c6b00000000000000000000000000000000000000000000000007fa9ef6609ca7921112231f8f195138ebba29770000000000000000000000007340e006f4135ba6970d43bf43d88dcad4e7a8ca0000000000000000000000007526eb4f95e2a1394797cb38a921fb1eba09291b00000000000000000000000000000000000000000000000000000000000000000000000000000000069f6bc7000000000000000000000000000000000000000000000000000000000d3ed78e"),
	BlockNumber: 23,
	TxHash:      common.HexToHash("0x7cb84c750ce4985f7811abf641d52ffcb35306d943081475226484cf1470c6fa"),
	TxIndex:     4,
	BlockHash:   common.HexToHash("0xf5a367d560e14c4658ef85e4877e08b5560a4773b69b39f6b8025910b666fade"),
	Index:       0,
	Removed:     false,
}

var rawVatGrabLog, _ = json.Marshal(EthVatGrabLog)
var VatGrabModel = vat_grab.VatGrabModel{
	Ilk:              "fake ilk",
	Urn:              "0x07Fa9eF6609cA7921112231F8f195138ebbA2977",
	V:                "0x7340e006f4135BA6970D43bf43d88DCAD4e7a8CA",
	W:                "0x7526EB4f95e2a1394797Cb38a921Fb1EbA09291B",
	Dink:             "111111111",
	Dart:             "222222222",
	TransactionIndex: EthVatGrabLog.TxIndex,
	Raw:              rawVatGrabLog,
}
