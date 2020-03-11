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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/statediff"
	"github.com/ethereum/go-ethereum/statediff/testhelpers"
	"github.com/ipfs/go-block-format"
	"github.com/multiformats/go-multihash"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs/ipld"
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
		Difficulty:  big.NewInt(5000000),
		Extra:       []byte{},
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
	MockLog1                                   = &types.Log{
		Topics: []common.Hash{mockTopic11, mockTopic12},
		Data:   []byte{},
	}
	MockLog2 = &types.Log{
		Topics: []common.Hash{mockTopic21, mockTopic22},
		Data:   []byte{},
	}
	HeaderCID, _  = ipld.RawdataToCid(ipld.MEthHeader, MockHeaderRlp, multihash.KECCAK_256)
	Trx1CID, _    = ipld.RawdataToCid(ipld.MEthTx, MockTransactions.GetRlp(0), multihash.KECCAK_256)
	Trx2CID, _    = ipld.RawdataToCid(ipld.MEthTx, MockTransactions.GetRlp(1), multihash.KECCAK_256)
	Rct1CID, _    = ipld.RawdataToCid(ipld.MEthTxReceipt, MockReceipts.GetRlp(0), multihash.KECCAK_256)
	Rct2CID, _    = ipld.RawdataToCid(ipld.MEthTxReceipt, MockReceipts.GetRlp(1), multihash.KECCAK_256)
	State1CID, _  = ipld.RawdataToCid(ipld.MEthStateTrie, ContractLeafNode, multihash.KECCAK_256)
	State2CID, _  = ipld.RawdataToCid(ipld.MEthStateTrie, AccountLeafNode, multihash.KECCAK_256)
	StorageCID, _ = ipld.RawdataToCid(ipld.MEthStorageTrie, StorageLeafNode, multihash.KECCAK_256)
	MockTrxMeta   = []eth.TxModel{
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
			CID:    Trx1CID.String(), // This is empty until we go to publish to ipfs
			Src:    senderAddr.Hex(),
			Dst:    Address.String(),
			Index:  0,
			TxHash: MockTransactions[0].Hash().String(),
		},
		{
			CID:    Trx2CID.String(),
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
			CID: Rct1CID.String(),
			Topic0s: []string{
				mockTopic11.String(),
			},
			Topic1s: []string{
				mockTopic12.String(),
			},
			Contract: Address.String(),
		},
		{
			CID: Rct2CID.String(),
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
	storageLocation    = common.HexToHash("0")
	StorageLeafKey     = crypto.Keccak256Hash(storageLocation[:]).Bytes()
	StorageValue       = common.Hex2Bytes("01")
	StoragePartialPath = common.Hex2Bytes("20290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563")
	StorageLeafNode, _ = rlp.EncodeToBytes([]interface{}{
		StoragePartialPath,
		StorageValue,
	})

	nonce1             = uint64(1)
	contractRoot       = "0x821e2556a290c86405f8160a2d662042a431ba456b9db265c79bb837c04be5f0"
	contractCodeHash   = common.HexToHash("0x753f98a8d4328b15636e46f66f2cb4bc860100aa17967cc145fcd17d1d4710ea")
	contractPathHash   = crypto.Keccak256Hash([]byte{'\x06'})
	ContractAddress    = common.HexToAddress("0x703c4b2bD70c169f5717101CaeE543299Fc946C7")
	ContractLeafKey    = testhelpers.AddressToLeafKey(ContractAddress)
	ContractAccount, _ = rlp.EncodeToBytes(state.Account{
		Nonce:    nonce1,
		Balance:  big.NewInt(0),
		CodeHash: contractCodeHash.Bytes(),
		Root:     common.HexToHash(contractRoot),
	})
	ContractPartialPath = common.Hex2Bytes("3114658a74d9cc9f7acf2c5cd696c3494d7c344d78bfec3add0d91ec4e8d1c45")
	ContractLeafNode, _ = rlp.EncodeToBytes([]interface{}{
		ContractPartialPath,
		ContractAccount,
	})

	nonce0          = uint64(0)
	accountRoot     = "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"
	accountCodeHash = common.HexToHash("0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470")
	AccountAddresss = common.HexToAddress("0x0D3ab14BBaD3D99F4203bd7a11aCB94882050E7e")
	AccountLeafKey  = testhelpers.Account2LeafKey
	Account, _      = rlp.EncodeToBytes(state.Account{
		Nonce:    nonce0,
		Balance:  big.NewInt(1000),
		CodeHash: accountCodeHash.Bytes(),
		Root:     common.HexToHash(accountRoot),
	})
	AccountPartialPath = common.Hex2Bytes("3957f3e2f04a0764c3a0491b175f69926da61efbcc8f61fa1455fd2d2b4cdd45")
	AccountLeafNode, _ = rlp.EncodeToBytes([]interface{}{
		AccountPartialPath,
		Account,
	})

	CreatedAccountDiffs = []statediff.AccountDiff{
		{
			Path:      []byte{'\x06'},
			NodeType:  statediff.Leaf,
			LeafKey:   ContractLeafKey,
			NodeValue: ContractLeafNode,
			Storage: []statediff.StorageDiff{
				{
					Path:      []byte{},
					NodeType:  statediff.Leaf,
					LeafKey:   StorageLeafKey,
					NodeValue: StorageLeafNode,
				},
			},
		},
		{
			Path:      []byte{'\x0c'},
			NodeType:  statediff.Leaf,
			LeafKey:   AccountLeafKey,
			NodeValue: AccountLeafNode,
			Storage:   []statediff.StorageDiff{},
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
			LeafKey: common.BytesToHash(ContractLeafKey),
			Path:    []byte{'\x06'},
			Value:   ContractLeafNode,
			Type:    statediff.Leaf,
		},
		{
			LeafKey: common.BytesToHash(AccountLeafKey),
			Path:    []byte{'\x0c'},
			Value:   AccountLeafNode,
			Type:    statediff.Leaf,
		},
	}
	MockStateMetaPostPublish = []eth.StateNodeModel{
		{
			CID:      State1CID.String(),
			Path:     []byte{'\x06'},
			NodeType: 2,
			StateKey: common.BytesToHash(ContractLeafKey).Hex(),
		},
		{
			CID:      State2CID.String(),
			Path:     []byte{'\x0c'},
			NodeType: 2,
			StateKey: common.BytesToHash(AccountLeafKey).Hex(),
		},
	}
	MockStorageNodes = map[common.Hash][]eth.TrieNode{
		contractPathHash: {
			{
				LeafKey: common.BytesToHash(StorageLeafKey),
				Value:   StorageLeafNode,
				Type:    statediff.Leaf,
				Path:    []byte{},
			},
		},
	}

	// aggregate payloads
	MockStateDiffPayload = statediff.Payload{
		BlockRlp:        MockBlockRlp,
		StateDiffRlp:    MockStateDiffBytes,
		ReceiptsRlp:     ReceiptsRlp,
		TotalDifficulty: MockBlock.Difficulty(),
	}

	MockConvertedPayload = eth.ConvertedPayload{
		TotalDifficulty: MockBlock.Difficulty(),
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
			CID:             HeaderCID.String(),
			ParentHash:      MockBlock.ParentHash().String(),
			TotalDifficulty: MockBlock.Difficulty().String(),
			Reward:          "5000000000000000000",
		},
		UncleCIDs:       []eth2.UncleModel{},
		TransactionCIDs: MockTrxMetaPostPublsh,
		ReceiptCIDs: map[common.Hash]eth.ReceiptModel{
			MockTransactions[0].Hash(): MockRctMetaPostPublish[0],
			MockTransactions[1].Hash(): MockRctMetaPostPublish[1],
		},
		StateNodeCIDs: MockStateMetaPostPublish,
		StorageNodeCIDs: map[common.Hash][]eth.StorageNodeModel{
			contractPathHash: {
				{
					CID:        StorageCID.String(),
					Path:       []byte{},
					StorageKey: common.BytesToHash(StorageLeafKey).Hex(),
					NodeType:   2,
				},
			},
		},
	}

	MockCIDWrapper = &eth.CIDWrapper{
		BlockNumber: big.NewInt(1),
		Header: eth2.HeaderModel{
			BlockNumber:     "1",
			BlockHash:       MockBlock.Hash().String(),
			ParentHash:      "0x0000000000000000000000000000000000000000000000000000000000000000",
			CID:             HeaderCID.String(),
			TotalDifficulty: MockBlock.Difficulty().String(),
			Reward:          "5000000000000000000",
		},
		Transactions: MockTrxMetaPostPublsh,
		Receipts:     MockRctMetaPostPublish,
		Uncles:       []eth2.UncleModel{},
		StateNodes:   MockStateMetaPostPublish,
		StorageNodes: []eth.StorageNodeWithStateKeyModel{
			{
				Path:       []byte{},
				CID:        StorageCID.String(),
				NodeType:   2,
				StateKey:   common.BytesToHash(ContractLeafKey).Hex(),
				StorageKey: common.BytesToHash(StorageLeafKey).Hex(),
			},
		},
	}

	HeaderIPLD, _  = blocks.NewBlockWithCid(MockHeaderRlp, HeaderCID)
	Trx1IPLD, _    = blocks.NewBlockWithCid(MockTransactions.GetRlp(0), Trx1CID)
	Trx2IPLD, _    = blocks.NewBlockWithCid(MockTransactions.GetRlp(1), Trx2CID)
	Rct1IPLD, _    = blocks.NewBlockWithCid(MockReceipts.GetRlp(0), Rct1CID)
	Rct2IPLD, _    = blocks.NewBlockWithCid(MockReceipts.GetRlp(1), Rct2CID)
	State1IPLD, _  = blocks.NewBlockWithCid(ContractLeafNode, State1CID)
	State2IPLD, _  = blocks.NewBlockWithCid(AccountLeafNode, State2CID)
	StorageIPLD, _ = blocks.NewBlockWithCid(StorageLeafNode, StorageCID)

	MockIPLDs = eth.IPLDs{
		BlockNumber: big.NewInt(1),
		Header: ipfs.BlockModel{
			Data: HeaderIPLD.RawData(),
			CID:  HeaderIPLD.Cid().String(),
		},
		Transactions: []ipfs.BlockModel{
			{
				Data: Trx1IPLD.RawData(),
				CID:  Trx1IPLD.Cid().String(),
			},
			{
				Data: Trx2IPLD.RawData(),
				CID:  Trx2IPLD.Cid().String(),
			},
		},
		Receipts: []ipfs.BlockModel{
			{
				Data: Rct1IPLD.RawData(),
				CID:  Rct1IPLD.Cid().String(),
			},
			{
				Data: Rct2IPLD.RawData(),
				CID:  Rct2IPLD.Cid().String(),
			},
		},
		StateNodes: []eth2.StateNode{
			{
				StateLeafKey: common.BytesToHash(ContractLeafKey),
				Type:         statediff.Leaf,
				IPLD: ipfs.BlockModel{
					Data: State1IPLD.RawData(),
					CID:  State1IPLD.Cid().String(),
				},
				Path: []byte{'\x06'},
			},
			{
				StateLeafKey: common.BytesToHash(AccountLeafKey),
				Type:         statediff.Leaf,
				IPLD: ipfs.BlockModel{
					Data: State2IPLD.RawData(),
					CID:  State2IPLD.Cid().String(),
				},
				Path: []byte{'\x0c'},
			},
		},
		StorageNodes: []eth2.StorageNode{
			{
				StateLeafKey:   common.BytesToHash(ContractLeafKey),
				StorageLeafKey: common.BytesToHash(StorageLeafKey),
				Type:           statediff.Leaf,
				IPLD: ipfs.BlockModel{
					Data: StorageIPLD.RawData(),
					CID:  StorageIPLD.Cid().String(),
				},
				Path: []byte{},
			},
		},
	}
)

// createTransactionsAndReceipts is a helper function to generate signed mock transactions and mock receipts with mock logs
func createTransactionsAndReceipts() (types.Transactions, types.Receipts, common.Address) {
	// make transactions
	trx1 := types.NewTransaction(0, Address, big.NewInt(1000), 50, big.NewInt(100), []byte{})
	trx2 := types.NewTransaction(1, AnotherAddress, big.NewInt(2000), 100, big.NewInt(200), []byte{})
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
	mockReceipt1.Logs = []*types.Log{MockLog1}
	mockReceipt1.TxHash = signedTrx1.Hash()
	mockReceipt2 := types.NewReceipt(common.HexToHash("0x1").Bytes(), false, 100)
	mockReceipt2.Logs = []*types.Log{MockLog2}
	mockReceipt2.TxHash = signedTrx2.Hash()
	return types.Transactions{signedTrx1, signedTrx2}, types.Receipts{mockReceipt1, mockReceipt2}, senderAddr
}
