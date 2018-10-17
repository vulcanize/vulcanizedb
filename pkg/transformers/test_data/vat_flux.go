/*
 *  Copyright 2018 Vulcanize
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package test_data

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_flux"
)

var VatFluxLog = types.Log{
	Address: common.HexToAddress("0x239e6f0ab02713f1f8aa90ebeded9fc66dc96cd6"),
	Topics: []common.Hash{
		common.HexToHash("0xa6e4182100000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000007340e006f4135ba6970d43bf43d88dcad4e7a8ca"),
		common.HexToHash("0x0000000000000000000000007fa9ef6609ca7921112231f8f195138ebba29770"),
		common.HexToHash("0x00000000000000000000000093086347c52a8878af71bb818509d484c6a2e1bf"),
	},
	Data:        hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000084a6e418217340e006f4135ba6970d43bf43d88dcad4e7a8ca00000000000000000000000007fa9ef6609ca7921112231f8f195138ebba297700000000000000000000000093086347c52a8878af71bb818509d484c6a2e1bf000000000000000000000000000000000000000000000000000000000000000000000000000000000000007b"),
	BlockNumber: 23,
	TxHash:      common.HexToHash("0xf98681bab9b8c75bd8aa4a7d0a8142ff527c5ea8fa54f3c2835d4533838b2e6f"),
	TxIndex:     0,
	BlockHash:   common.HexToHash("0xc3fe212ad4f81ade1265af6de2b4bb50d962b1a4db06aabc982e7f9cb0972c2d"),
	Index:       0,
	Removed:     false,
}

var rawFluxLog, _ = json.Marshal(VatFluxLog)
var VatFluxModel = vat_flux.VatFluxModel{
	Ilk:              "0x7340e006f4135BA6970D43bf43d88DCAD4e7a8CA",
	Src:              "0x07Fa9eF6609cA7921112231F8f195138ebbA2977",
	Dst:              "0x93086347c52a8878af71bB818509d484c6a2e1bF",
	Rad:              "123",
	TransactionIndex: 0,
	Raw:              rawFluxLog,
}
