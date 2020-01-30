// VulcanizeDB
// Copyright © 2019 Vulcanize

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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
	rand2 "math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff"
	"github.com/ipfs/go-block-format"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth"
	eth2 "github.com/vulcanize/vulcanizedb/pkg/super_node/eth"
)

// Test variables
var (
	// block data
	BlockNumber = big.NewInt(1)
	MockHeader  = types.Header{
		Time:        0,
		Number:      BlockNumber,
		Root:        common.HexToHash("0x0"),
		TxHash:      common.HexToHash("0x0"),
		ReceiptHash: common.HexToHash("0x0"),
	}
	MockTransactions, MockReceipts, senderAddr = createTransactionsAndReceipts()
	ReceiptsRlp, _                             = rlp.EncodeToBytes(MockReceipts)
	MockBlock                                  = types.NewBlock(&MockHeader, MockTransactions, nil, MockReceipts)
	MockBlockRlp, _                            = rlp.EncodeToBytes(MockBlock)
	MockHeaderRlp, _                           = rlp.EncodeToBytes(MockBlock.Header())
	Address                                    = common.HexToAddress("0xaE9BEa628c4Ce503DcFD7E305CaB4e29E7476592")
	AnotherAddress                             = common.HexToAddress("0xaE9BEa628c4Ce503DcFD7E305CaB4e29E7476593")
	mockTopic11                                = common.HexToHash("0x04")
	mockTopic12                                = common.HexToHash("0x06")
	mockTopic21                                = common.HexToHash("0x05")
	mockTopic22                                = common.HexToHash("0x07")
	MockTrxMeta                                = []eth.TxModel{
		{
			CID:    "", // This is empty until we go to publish to ipfs
			Src:    senderAddr.Hex(),
			Dst:    Address.String(),
			Index:  0,
			TxHash: MockTransactions[0].Hash().String(),
		},
		{
			CID:    "",
			Src:    senderAddr.Hex(),
			Dst:    AnotherAddress.String(),
			Index:  1,
			TxHash: MockTransactions[1].Hash().String(),
		},
	}
	MockTrxMetaPostPublsh = []eth.TxModel{
		{
			CID:    "mockTrxCID1", // This is empty until we go to publish to ipfs
			Src:    senderAddr.Hex(),
			Dst:    Address.String(),
			Index:  0,
			TxHash: MockTransactions[0].Hash().String(),
		},
		{
			CID:    "mockTrxCID2",
			Src:    senderAddr.Hex(),
			Dst:    AnotherAddress.String(),
			Index:  1,
			TxHash: MockTransactions[1].Hash().String(),
		},
	}
	MockRctMeta = []eth.ReceiptModel{
		{
			CID: "",
			Topic0s: []string{
				mockTopic11.String(),
			},
			Topic1s: []string{
				mockTopic12.String(),
			},
			Contract: Address.String(),
		},
		{
			CID: "",
			Topic0s: []string{
				mockTopic21.String(),
			},
			Topic1s: []string{
				mockTopic22.String(),
			},
			Contract: AnotherAddress.String(),
		},
	}
	MockRctMetaPostPublish = []eth.ReceiptModel{
		{
			CID: "mockRctCID1",
			Topic0s: []string{
				mockTopic11.String(),
			},
			Topic1s: []string{
				mockTopic12.String(),
			},
			Contract: Address.String(),
		},
		{
			CID: "mockRctCID2",
			Topic0s: []string{
				mockTopic21.String(),
			},
			Topic1s: []string{
				mockTopic22.String(),
			},
			Contract: AnotherAddress.String(),
		},
	}

	// statediff data
	CodeHash            = common.Hex2Bytes("0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470")
	NonceValue          = rand2.Uint64()
	anotherNonceValue   = rand2.Uint64()
	BalanceValue        = rand2.Int63()
	anotherBalanceValue = rand2.Int63()
	ContractRoot        = common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
	StoragePath         = common.HexToHash("0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470").Bytes()
	StorageKey          = common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001").Bytes()
	StorageValue        = common.Hex2Bytes("0x03")
	storage             = []statediff.StorageDiff{{
		Key:   StorageKey,
		Value: StorageValue,
		Path:  StoragePath,
		Proof: [][]byte{},
		Leaf:  true,
	}}
	emptyStorage           = make([]statediff.StorageDiff, 0)
	ContractLeafKey        = crypto.Keccak256Hash(Address.Bytes())
	AnotherContractLeafKey = crypto.Keccak256Hash(AnotherAddress.Bytes())
	testAccount            = state.Account{
		Nonce:    NonceValue,
		Balance:  big.NewInt(BalanceValue),
		Root:     ContractRoot,
		CodeHash: CodeHash,
	}
	anotherTestAccount = state.Account{
		Nonce:    anotherNonceValue,
		Balance:  big.NewInt(anotherBalanceValue),
		Root:     common.HexToHash("0x"),
		CodeHash: nil,
	}
	ValueBytes, _        = rlp.EncodeToBytes(testAccount)
	AnotherValueBytes, _ = rlp.EncodeToBytes(anotherTestAccount)
	CreatedAccountDiffs  = []statediff.AccountDiff{
		{
			Key:     ContractLeafKey.Bytes(),
			Value:   ValueBytes,
			Storage: storage,
			Leaf:    true,
		},
		{
			Key:     AnotherContractLeafKey.Bytes(),
			Value:   AnotherValueBytes,
			Storage: emptyStorage,
			Leaf:    true,
		},
	}

	MockStateDiff = statediff.StateDiff{
		BlockNumber:     BlockNumber,
		BlockHash:       MockBlock.Hash(),
		CreatedAccounts: CreatedAccountDiffs,
	}
	MockStateDiffBytes, _ = rlp.EncodeToBytes(MockStateDiff)
	MockStateNodes        = []eth.TrieNode{
		{
			Key:   ContractLeafKey,
			Value: ValueBytes,
			Leaf:  true,
		},
		{
			Key:   AnotherContractLeafKey,
			Value: AnotherValueBytes,
			Leaf:  true,
		},
	}
	MockStateMetaPostPublish = []eth.StateNodeModel{
		{
			CID:      "mockStateCID1",
			Leaf:     true,
			StateKey: ContractLeafKey.String(),
		},
		{
			CID:      "mockStateCID2",
			Leaf:     true,
			StateKey: AnotherContractLeafKey.String(),
		},
	}
	MockStorageNodes = map[common.Hash][]eth.TrieNode{
		ContractLeafKey: {
			{
				Key:   common.BytesToHash(StorageKey),
				Value: StorageValue,
				Leaf:  true,
			},
		},
	}

	// aggregate payloads
	MockStateDiffPayload = statediff.Payload{
		BlockRlp:        MockBlockRlp,
		StateDiffRlp:    MockStateDiffBytes,
		ReceiptsRlp:     ReceiptsRlp,
		TotalDifficulty: big.NewInt(1337),
	}

	MockIPLDPayload = &eth.IPLDPayload{
		TotalDifficulty: big.NewInt(1337),
		Block:           MockBlock,
		Receipts:        MockReceipts,
		TxMetaData:      MockTrxMeta,
		ReceiptMetaData: MockRctMeta,
		StorageNodes:    MockStorageNodes,
		StateNodes:      MockStateNodes,
	}

	MockCIDPayload = &eth.CIDPayload{
		HeaderCID: eth2.HeaderModel{
			BlockHash:       MockBlock.Hash().String(),
			BlockNumber:     MockBlock.Number().String(),
			CID:             "mockHeaderCID",
			ParentHash:      MockBlock.ParentHash().String(),
			TotalDifficulty: "1337",
		},
		UncleCIDs:       []eth2.UncleModel{},
		TransactionCIDs: MockTrxMetaPostPublsh,
		ReceiptCIDs: map[common.Hash]eth.ReceiptModel{
			MockTransactions[0].Hash(): MockRctMetaPostPublish[0],
			MockTransactions[1].Hash(): MockRctMetaPostPublish[1],
		},
		StateNodeCIDs: MockStateMetaPostPublish,
		StorageNodeCIDs: map[common.Hash][]eth.StorageNodeModel{
			ContractLeafKey: {
				{
					CID:        "mockStorageCID",
					StorageKey: "0x0000000000000000000000000000000000000000000000000000000000000001",
					Leaf:       true,
				},
			},
		},
	}

	MockCIDWrapper = &eth.CIDWrapper{
		BlockNumber: big.NewInt(1),
		Headers: []eth2.HeaderModel{
			{
				BlockNumber:     "1",
				BlockHash:       MockBlock.Hash().String(),
				ParentHash:      "0x0000000000000000000000000000000000000000000000000000000000000000",
				CID:             "mockHeaderCID",
				TotalDifficulty: "1337",
			},
		},
		Transactions: MockTrxMetaPostPublsh,
		Receipts:     MockRctMetaPostPublish,
		Uncles:       []eth2.UncleModel{},
		StateNodes:   MockStateMetaPostPublish,
		StorageNodes: []eth.StorageNodeWithStateKeyModel{
			{
				CID:        "mockStorageCID",
				Leaf:       true,
				StateKey:   ContractLeafKey.Hex(),
				StorageKey: "0x0000000000000000000000000000000000000000000000000000000000000001",
			},
		},
	}

	MockIPLDWrapper = &eth.IPLDWrapper{
		BlockNumber: big.NewInt(1),
		Headers: []blocks.Block{
			blocks.NewBlock(MockHeaderRlp),
		},
		Transactions: []blocks.Block{
			blocks.NewBlock(MockTransactions.GetRlp(0)),
			blocks.NewBlock(MockTransactions.GetRlp(1)),
		},
		Receipts: []blocks.Block{
			blocks.NewBlock(MockReceipts.GetRlp(0)),
			blocks.NewBlock(MockReceipts.GetRlp(1)),
		},
		StateNodes: map[common.Hash]blocks.Block{
			ContractLeafKey:        blocks.NewBlock(ValueBytes),
			AnotherContractLeafKey: blocks.NewBlock(AnotherValueBytes),
		},
		StorageNodes: map[common.Hash]map[common.Hash]blocks.Block{
			ContractLeafKey: {
				common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"): blocks.NewBlock(StorageValue),
			},
		},
	}

	MockSeedNodePayload = eth2.StreamPayload{
		BlockNumber:     big.NewInt(1),
		HeadersRlp:      [][]byte{MockHeaderRlp},
		UnclesRlp:       [][]byte{},
		TransactionsRlp: [][]byte{MockTransactions.GetRlp(0), MockTransactions.GetRlp(1)},
		ReceiptsRlp:     [][]byte{MockTransactions.GetRlp(0), MockTransactions.GetRlp(1)},
		StateNodesRlp: map[common.Hash][]byte{
			ContractLeafKey:        ValueBytes,
			AnotherContractLeafKey: AnotherValueBytes,
		},
		StorageNodesRlp: map[common.Hash]map[common.Hash][]byte{
			ContractLeafKey: {
				common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"): StorageValue,
			},
		},
	}
)

// createTransactionsAndReceipts is a helper function to generate signed mock transactions and mock receipts with mock logs
func createTransactionsAndReceipts() (types.Transactions, types.Receipts, common.Address) {
	// make transactions
	trx1 := types.NewTransaction(0, Address, big.NewInt(1000), 50, big.NewInt(100), nil)
	trx2 := types.NewTransaction(1, AnotherAddress, big.NewInt(2000), 100, big.NewInt(200), nil)
	transactionSigner := types.MakeSigner(params.MainnetChainConfig, BlockNumber)
	mockCurve := elliptic.P256()
	mockPrvKey, err := ecdsa.GenerateKey(mockCurve, rand.Reader)
	if err != nil {
		log.Fatal(err)
	}
	signedTrx1, err := types.SignTx(trx1, transactionSigner, mockPrvKey)
	if err != nil {
		log.Fatal(err)
	}
	signedTrx2, err := types.SignTx(trx2, transactionSigner, mockPrvKey)
	if err != nil {
		log.Fatal(err)
	}
	senderAddr, err := types.Sender(transactionSigner, signedTrx1) // same for both trx
	if err != nil {
		log.Fatal(err)
	}
	// make receipts
	mockReceipt1 := types.NewReceipt(common.HexToHash("0x0").Bytes(), false, 50)
	mockLog1 := &types.Log{
		Topics: []common.Hash{mockTopic11, mockTopic12},
	}
	mockReceipt1.Logs = []*types.Log{mockLog1}
	mockReceipt1.TxHash = signedTrx1.Hash()
	mockReceipt2 := types.NewReceipt(common.HexToHash("0x1").Bytes(), false, 100)
	mockLog2 := &types.Log{
		Topics: []common.Hash{mockTopic21, mockTopic22},
	}
	mockReceipt2.Logs = []*types.Log{mockLog2}
	mockReceipt2.TxHash = signedTrx2.Hash()
	return types.Transactions{signedTrx1, signedTrx2}, types.Receipts{mockReceipt1, mockReceipt2}, senderAddr
}
