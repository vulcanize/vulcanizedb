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
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
)

// Test variables
var (
	// block data
	BlockNumber = big.NewInt(rand2.Int63())
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
	MockHeaderRlp, err                         = rlp.EncodeToBytes(MockBlock.Header())
	MockTrxMeta                                = []*ipfs.TrxMetaData{
		{
			CID: "", // This is empty until we go to publish to ipfs
			Src: senderAddr.Hex(),
			Dst: "0x0000000000000000000000000000000000000000",
		},
		{
			CID: "",
			Src: senderAddr.Hex(),
			Dst: "0x0000000000000000000000000000000000000001",
		},
	}
	MockRctMeta = []*ipfs.ReceiptMetaData{
		{
			CID: "",
			Topic0s: []string{
				"0x0000000000000000000000000000000000000000000000000000000000000004",
			},
			ContractAddress: "0x0000000000000000000000000000000000000000",
		},
		{
			CID: "",
			Topic0s: []string{
				"0x0000000000000000000000000000000000000000000000000000000000000005",
			},
			ContractAddress: "0x0000000000000000000000000000000000000001",
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
	Address                = common.HexToAddress("0xaE9BEa628c4Ce503DcFD7E305CaB4e29E7476592")
	ContractLeafKey        = ipfs.AddressToKey(Address)
	AnotherAddress         = common.HexToAddress("0xaE9BEa628c4Ce503DcFD7E305CaB4e29E7476593")
	AnotherContractLeafKey = ipfs.AddressToKey(AnotherAddress)
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
	valueBytes, _        = rlp.EncodeToBytes(testAccount)
	anotherValueBytes, _ = rlp.EncodeToBytes(anotherTestAccount)
	CreatedAccountDiffs  = []statediff.AccountDiff{
		{
			Key:     ContractLeafKey.Bytes(),
			Value:   valueBytes,
			Storage: storage,
			Leaf:    true,
		},
		{
			Key:     AnotherContractLeafKey.Bytes(),
			Value:   anotherValueBytes,
			Storage: emptyStorage,
			Leaf:    true,
		},
	}

	UpdatedAccountDiffs = []statediff.AccountDiff{{
		Key:     ContractLeafKey.Bytes(),
		Value:   valueBytes,
		Storage: storage,
		Leaf:    true,
	}}

	DeletedAccountDiffs = []statediff.AccountDiff{{
		Key:     ContractLeafKey.Bytes(),
		Value:   valueBytes,
		Storage: storage,
		Leaf:    true,
	}}

	MockStateDiff = statediff.StateDiff{
		BlockNumber:     BlockNumber,
		BlockHash:       MockBlock.Hash(),
		CreatedAccounts: CreatedAccountDiffs,
	}
	MockStateDiffBytes, _ = rlp.EncodeToBytes(MockStateDiff)
	MockStateNodes        = map[common.Hash]ipfs.StateNode{
		ContractLeafKey: {
			Value: valueBytes,
			Leaf:  true,
		},
		AnotherContractLeafKey: {
			Value: anotherValueBytes,
			Leaf:  true,
		},
	}
	MockStorageNodes = map[common.Hash][]ipfs.StorageNode{
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
		BlockRlp:     MockBlockRlp,
		StateDiffRlp: MockStateDiffBytes,
		ReceiptsRlp:  ReceiptsRlp,
	}

	EmptyStateDiffPayload = statediff.Payload{
		BlockRlp:     []byte{},
		StateDiffRlp: []byte{},
		ReceiptsRlp:  []byte{},
	}

	MockIPLDPayload = &ipfs.IPLDPayload{
		BlockNumber: big.NewInt(1),
		BlockHash:   MockBlock.Hash(),
		Receipts:    MockReceipts,
		HeaderRLP:   MockHeaderRlp,
		BlockBody:   MockBlock.Body(),
		TrxMetaData: []*ipfs.TrxMetaData{
			{
				CID: "",
				Src: senderAddr.Hex(),
				Dst: "0x0000000000000000000000000000000000000000",
			},
			{
				CID: "",
				Src: senderAddr.Hex(),
				Dst: "0x0000000000000000000000000000000000000001",
			},
		},
		ReceiptMetaData: []*ipfs.ReceiptMetaData{
			{
				CID: "",
				Topic0s: []string{
					"0x0000000000000000000000000000000000000000000000000000000000000004",
				},
				ContractAddress: "0x0000000000000000000000000000000000000000",
			},
			{
				CID: "",
				Topic0s: []string{
					"0x0000000000000000000000000000000000000000000000000000000000000005",
				},
				ContractAddress: "0x0000000000000000000000000000000000000001",
			},
		},
		StorageNodes: MockStorageNodes,
		StateNodes:   MockStateNodes,
	}

	MockCIDPayload = &ipfs.CIDPayload{
		BlockNumber: "1",
		BlockHash:   MockBlock.Hash(),
		HeaderCID:   "mockHeaderCID",
		UncleCIDS:   make(map[common.Hash]string),
		TransactionCIDs: map[common.Hash]*ipfs.TrxMetaData{
			MockTransactions[0].Hash(): {
				CID: "mockTrxCID1",
				Dst: "0x0000000000000000000000000000000000000000",
				Src: senderAddr.Hex(),
			},
			MockTransactions[1].Hash(): {
				CID: "mockTrxCID2",
				Dst: "0x0000000000000000000000000000000000000001",
				Src: senderAddr.Hex(),
			},
		},
		ReceiptCIDs: map[common.Hash]*ipfs.ReceiptMetaData{
			MockTransactions[0].Hash(): {
				CID:             "mockRctCID1",
				Topic0s:         []string{"0x0000000000000000000000000000000000000000000000000000000000000004"},
				ContractAddress: "0x0000000000000000000000000000000000000000",
			},
			MockTransactions[1].Hash(): {
				CID:             "mockRctCID2",
				Topic0s:         []string{"0x0000000000000000000000000000000000000000000000000000000000000005"},
				ContractAddress: "0x0000000000000000000000000000000000000001",
			},
		},
		StateNodeCIDs: map[common.Hash]ipfs.StateNodeCID{
			ContractLeafKey: {
				CID:  "mockStateCID1",
				Leaf: true,
				Key:  "",
			},
			AnotherContractLeafKey: {
				CID:  "mockStateCID2",
				Leaf: true,
				Key:  "",
			},
		},
		StorageNodeCIDs: map[common.Hash][]ipfs.StorageNodeCID{
			ContractLeafKey: {
				{
					CID:      "mockStorageCID",
					Key:      "0x0000000000000000000000000000000000000000000000000000000000000001",
					Leaf:     true,
					StateKey: "",
				},
			},
		},
	}

	MockCIDWrapper = &ipfs.CIDWrapper{
		BlockNumber:  big.NewInt(1),
		Headers:      []string{"mockHeaderCID"},
		Transactions: []string{"mockTrxCID1", "mockTrxCID2"},
		Receipts:     []string{"mockRctCID1", "mockRctCID2"},
		Uncles:       []string{},
		StateNodes: []ipfs.StateNodeCID{
			{
				CID:  "mockStateCID1",
				Leaf: true,
				Key:  ContractLeafKey.Hex(),
			},
			{
				CID:  "mockStateCID2",
				Leaf: true,
				Key:  AnotherContractLeafKey.Hex(),
			},
		},
		StorageNodes: []ipfs.StorageNodeCID{
			{
				CID:      "mockStorageCID",
				Leaf:     true,
				StateKey: ContractLeafKey.Hex(),
				Key:      "0x0000000000000000000000000000000000000000000000000000000000000001",
			},
		},
	}
)

// createTransactionsAndReceipts is a helper function to generate signed mock transactions and mock receipts with mock logs
func createTransactionsAndReceipts() (types.Transactions, types.Receipts, common.Address) {
	// make transactions
	trx1 := types.NewTransaction(0, common.HexToAddress("0x0"), big.NewInt(1000), 50, big.NewInt(100), nil)
	trx2 := types.NewTransaction(1, common.HexToAddress("0x1"), big.NewInt(2000), 100, big.NewInt(200), nil)
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
	mockTopic1 := common.HexToHash("0x04")
	mockReceipt1 := types.NewReceipt(common.HexToHash("0x0").Bytes(), false, 50)
	mockReceipt1.ContractAddress = common.HexToAddress("0xaE9BEa628c4Ce503DcFD7E305CaB4e29E7476592")
	mockLog1 := &types.Log{
		Topics: []common.Hash{mockTopic1},
	}
	mockReceipt1.Logs = []*types.Log{mockLog1}
	mockReceipt1.TxHash = signedTrx1.Hash()
	mockTopic2 := common.HexToHash("0x05")
	mockReceipt2 := types.NewReceipt(common.HexToHash("0x1").Bytes(), false, 100)
	mockReceipt2.ContractAddress = common.HexToAddress("0xaE9BEa628c4Ce503DcFD7E305CaB4e29E7476593")
	mockLog2 := &types.Log{
		Topics: []common.Hash{mockTopic2},
	}
	mockReceipt2.Logs = []*types.Log{mockLog2}
	mockReceipt2.TxHash = signedTrx2.Hash()
	return types.Transactions{signedTrx1, signedTrx2}, types.Receipts{mockReceipt1, mockReceipt2}, senderAddr
}
