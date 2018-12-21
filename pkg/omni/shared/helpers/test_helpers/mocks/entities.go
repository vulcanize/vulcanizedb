// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package mocks

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/constants"
)

var TransferBlock1 = core.Block{
	Hash:   "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad123ert",
	Number: 6194633,
	Transactions: []core.Transaction{{
		Hash: "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad654aaa",
		Receipt: core.Receipt{
			TxHash:          "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad654aaa",
			ContractAddress: "",
			Logs: []core.Log{{
				BlockNumber: 6194633,
				TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad654aaa",
				Address:     constants.TusdContractAddress,
				Topics: core.Topics{
					constants.TransferEvent.Signature(),
					"0x000000000000000000000000000000000000000000000000000000000000af21",
					"0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
					"",
				},
				Index: 1,
				Data:  "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
			}},
		},
	}},
}

var TransferBlock2 = core.Block{
	Hash:   "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad123ooo",
	Number: 6194634,
	Transactions: []core.Transaction{{
		Hash: "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad654eee",
		Receipt: core.Receipt{
			TxHash:          "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad654eee",
			ContractAddress: "",
			Logs: []core.Log{{
				BlockNumber: 6194634,
				TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad654eee",
				Address:     constants.TusdContractAddress,
				Topics: core.Topics{
					constants.TransferEvent.Signature(),
					"0x000000000000000000000000000000000000000000000000000000000000af21",
					"0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
					"",
				},
				Index: 1,
				Data:  "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
			}},
		},
	}},
}

var NewOwnerBlock1 = core.Block{
	Hash:   "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad123ppp",
	Number: 6194635,
	Transactions: []core.Transaction{{
		Hash: "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad654bbb",
		Receipt: core.Receipt{
			TxHash:          "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad654bbb",
			ContractAddress: "",
			Logs: []core.Log{{
				BlockNumber: 6194635,
				TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad654bbb",
				Address:     constants.EnsContractAddress,
				Topics: core.Topics{
					constants.NewOwnerEvent.Signature(),
					"0x0000000000000000000000000000000000000000000000000000c02aaa39b223",
					"0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
					"",
				},
				Index: 1,
				Data:  "0x000000000000000000000000000000000000000000000000000000000000af21",
			}},
		},
	}},
}

var NewOwnerBlock2 = core.Block{
	Hash:   "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad123ggg",
	Number: 6194636,
	Transactions: []core.Transaction{{
		Hash: "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad654lll",
		Receipt: core.Receipt{
			TxHash:          "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad654lll",
			ContractAddress: "",
			Logs: []core.Log{{
				BlockNumber: 6194636,
				TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad654lll",
				Address:     constants.EnsContractAddress,
				Topics: core.Topics{
					constants.NewOwnerEvent.Signature(),
					"0x0000000000000000000000000000000000000000000000000000c02aaa39b223",
					"0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba400",
					"",
				},
				Index: 1,
				Data:  "0x000000000000000000000000000000000000000000000000000000000000af21",
			}},
		},
	}},
}

var ExpectedTransferFilter = filters.LogFilter{
	Name:      "Transfer",
	Address:   constants.TusdContractAddress,
	ToBlock:   -1,
	FromBlock: 6194634,
	Topics:    core.Topics{constants.TransferEvent.Signature()},
}

var ExpectedApprovalFilter = filters.LogFilter{
	Name:      "Approval",
	Address:   constants.TusdContractAddress,
	ToBlock:   -1,
	FromBlock: 6194634,
	Topics:    core.Topics{constants.ApprovalEvent.Signature()},
}

var MockTranferEvent = core.WatchedEvent{
	LogID:       1,
	Name:        constants.TransferEvent.String(),
	BlockNumber: 5488076,
	Address:     constants.TusdContractAddress,
	TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
	Index:       110,
	Topic0:      constants.TransferEvent.Signature(),
	Topic1:      "0x000000000000000000000000000000000000000000000000000000000000af21",
	Topic2:      "0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
	Topic3:      "",
	Data:        "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
}

var rawFakeHeader, _ = json.Marshal(core.Header{})

var MockHeader1 = core.Header{
	Hash:        "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad123ert",
	BlockNumber: 6194632,
	Raw:         rawFakeHeader,
	Timestamp:   "50000000",
}

var MockHeader2 = core.Header{
	Hash:        "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad456yui",
	BlockNumber: 6194633,
	Raw:         rawFakeHeader,
	Timestamp:   "50000015",
}

var MockHeader3 = core.Header{
	Hash:        "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad234hfs",
	BlockNumber: 6194634,
	Raw:         rawFakeHeader,
	Timestamp:   "50000030",
}

var MockTransferLog1 = types.Log{
	Index:       1,
	Address:     common.HexToAddress(constants.TusdContractAddress),
	BlockNumber: 5488076,
	TxIndex:     110,
	TxHash:      common.HexToHash("0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae"),
	Topics: []common.Hash{
		common.HexToHash(constants.TransferEvent.Signature()),
		common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000af21"),
		common.HexToHash("0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391"),
	},
	Data: hexutil.MustDecode("0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe"),
}

var MockTransferLog2 = types.Log{
	Index:       3,
	Address:     common.HexToAddress(constants.TusdContractAddress),
	BlockNumber: 5488077,
	TxIndex:     2,
	TxHash:      common.HexToHash("0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546df"),
	Topics: []common.Hash{
		common.HexToHash(constants.TransferEvent.Signature()),
		common.HexToHash("0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391"),
		common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000af21"),
	},
	Data: hexutil.MustDecode("0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe"),
}

var MockMintLog = types.Log{
	Index:       10,
	Address:     common.HexToAddress(constants.TusdContractAddress),
	BlockNumber: 5488080,
	TxIndex:     50,
	TxHash:      common.HexToHash("0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6minty"),
	Topics: []common.Hash{
		common.HexToHash(constants.MintEvent.Signature()),
		common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000af21"),
	},
	Data: hexutil.MustDecode("0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe"),
}

var MockNewOwnerLog1 = types.Log{
	Index:       1,
	Address:     common.HexToAddress(constants.EnsContractAddress),
	BlockNumber: 5488076,
	TxIndex:     110,
	TxHash:      common.HexToHash("0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae"),
	Topics: []common.Hash{
		common.HexToHash(constants.NewOwnerEvent.Signature()),
		common.HexToHash("0x000000000000000000000000c02aaa39b223helloa0e5c4f27ead9083c752553"),
		common.HexToHash("0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391"),
	},
	Data: hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000af21"),
}

var MockNewOwnerLog2 = types.Log{
	Index:       3,
	Address:     common.HexToAddress(constants.EnsContractAddress),
	BlockNumber: 5488077,
	TxIndex:     2,
	TxHash:      common.HexToHash("0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546df"),
	Topics: []common.Hash{
		common.HexToHash(constants.NewOwnerEvent.Signature()),
		common.HexToHash("0x000000000000000000000000c02aaa39b223helloa0e5c4f27ead9083c752553"),
		common.HexToHash("0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba400"),
	},
	Data: hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000af21"),
}
