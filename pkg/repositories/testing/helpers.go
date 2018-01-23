package testing

import (
	"sort"
	"strconv"

	"math/big"

	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/8thlight/vulcanizedb/pkg/filters"
	"github.com/8thlight/vulcanizedb/pkg/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func ClearData(postgres repositories.Postgres) {
	postgres.Db.MustExec("DELETE FROM watched_contracts")
	postgres.Db.MustExec("DELETE FROM transactions")
	postgres.Db.MustExec("DELETE FROM blocks")
	postgres.Db.MustExec("DELETE FROM logs")
	postgres.Db.MustExec("DELETE FROM receipts")
	postgres.Db.MustExec("DELETE FROM log_filters")
}

func AssertRepositoryBehavior(buildRepository func(node core.Node) repositories.Repository) {
	var repository repositories.Repository

	BeforeEach(func() {
		node := core.Node{
			GenesisBlock: "GENESIS",
			NetworkId:    1,
			Id:           "b6f90c0fdd8ec9607aed8ee45c69322e47b7063f0bfb7a29c8ecafab24d0a22d24dd2329b5ee6ed4125a03cb14e57fd584e67f9e53e6c631055cbbd82f080845",
			ClientName:   "Geth/v1.7.2-stable-1db4ecdc/darwin-amd64/go1.9",
		}
		repository = buildRepository(node)
	})

	Describe("Saving blocks", func() {
		It("starts with no blocks", func() {
			count := repository.BlockCount()
			Expect(count).Should(Equal(0))
		})

		It("increments the block count", func() {
			block := core.Block{Number: 123}

			repository.CreateOrUpdateBlock(block)

			Expect(repository.BlockCount()).To(Equal(1))
		})

		It("associates blocks to a node", func() {
			block := core.Block{
				Number: 123,
			}
			repository.CreateOrUpdateBlock(block)
			nodeTwo := core.Node{
				GenesisBlock: "0x456",
				NetworkId:    1,
				Id:           "x123456",
				ClientName:   "Geth",
			}
			repositoryTwo := buildRepository(nodeTwo)

			_, err := repositoryTwo.FindBlockByNumber(123)
			Expect(err).To(HaveOccurred())
		})

		It("saves the attributes of the block", func() {
			blockNumber := int64(123)
			gasLimit := int64(1000000)
			gasUsed := int64(10)
			blockHash := "x123"
			blockParentHash := "x456"
			blockNonce := "0x881db2ca900682e9a9"
			miner := "x123"
			extraData := "xextraData"
			blockTime := int64(1508981640)
			uncleHash := "x789"
			blockSize := int64(1000)
			difficulty := int64(10)
			blockReward := float64(5.132)
			unclesReward := float64(3.580)
			block := core.Block{
				Reward:       blockReward,
				Difficulty:   difficulty,
				GasLimit:     gasLimit,
				GasUsed:      gasUsed,
				Hash:         blockHash,
				ExtraData:    extraData,
				Nonce:        blockNonce,
				Miner:        miner,
				Number:       blockNumber,
				ParentHash:   blockParentHash,
				Size:         blockSize,
				Time:         blockTime,
				UncleHash:    uncleHash,
				UnclesReward: unclesReward,
			}

			repository.CreateOrUpdateBlock(block)

			savedBlock, err := repository.FindBlockByNumber(blockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(savedBlock.Reward).To(Equal(blockReward))
			Expect(savedBlock.Difficulty).To(Equal(difficulty))
			Expect(savedBlock.GasLimit).To(Equal(gasLimit))
			Expect(savedBlock.GasUsed).To(Equal(gasUsed))
			Expect(savedBlock.Hash).To(Equal(blockHash))
			Expect(savedBlock.Nonce).To(Equal(blockNonce))
			Expect(savedBlock.Miner).To(Equal(miner))
			Expect(savedBlock.ExtraData).To(Equal(extraData))
			Expect(savedBlock.Number).To(Equal(blockNumber))
			Expect(savedBlock.ParentHash).To(Equal(blockParentHash))
			Expect(savedBlock.Size).To(Equal(blockSize))
			Expect(savedBlock.Time).To(Equal(blockTime))
			Expect(savedBlock.UncleHash).To(Equal(uncleHash))
			Expect(savedBlock.UnclesReward).To(Equal(unclesReward))
		})

		It("does not find a block when searching for a number that does not exist", func() {
			_, err := repository.FindBlockByNumber(111)

			Expect(err).To(HaveOccurred())
		})

		It("saves one transaction associated to the block", func() {
			block := core.Block{
				Number:       123,
				Transactions: []core.Transaction{{}},
			}

			repository.CreateOrUpdateBlock(block)

			savedBlock, _ := repository.FindBlockByNumber(123)
			Expect(len(savedBlock.Transactions)).To(Equal(1))
		})

		It("saves two transactions associated to the block", func() {
			block := core.Block{
				Number:       123,
				Transactions: []core.Transaction{{}, {}},
			}

			repository.CreateOrUpdateBlock(block)

			savedBlock, _ := repository.FindBlockByNumber(123)
			Expect(len(savedBlock.Transactions)).To(Equal(2))
		})

		It(`replaces blocks and transactions associated to the block
			when a more new block is in conflict (same block number + nodeid)`, func() {
			blockOne := core.Block{
				Number:       123,
				Hash:         "xabc",
				Transactions: []core.Transaction{{Hash: "x123"}, {Hash: "x345"}},
			}
			blockTwo := core.Block{
				Number:       123,
				Hash:         "xdef",
				Transactions: []core.Transaction{{Hash: "x678"}, {Hash: "x9ab"}},
			}

			repository.CreateOrUpdateBlock(blockOne)
			repository.CreateOrUpdateBlock(blockTwo)

			savedBlock, _ := repository.FindBlockByNumber(123)
			Expect(len(savedBlock.Transactions)).To(Equal(2))
			Expect(savedBlock.Transactions[0].Hash).To(Equal("x678"))
			Expect(savedBlock.Transactions[1].Hash).To(Equal("x9ab"))
		})

		It(`does not replace blocks when block number is not unique
			     but block number + node id is`, func() {
			blockOne := core.Block{
				Number:       123,
				Transactions: []core.Transaction{{Hash: "x123"}, {Hash: "x345"}},
			}
			blockTwo := core.Block{
				Number:       123,
				Transactions: []core.Transaction{{Hash: "x678"}, {Hash: "x9ab"}},
			}
			repository.CreateOrUpdateBlock(blockOne)
			nodeTwo := core.Node{
				GenesisBlock: "0x456",
				NetworkId:    1,
			}
			repositoryTwo := buildRepository(nodeTwo)

			repository.CreateOrUpdateBlock(blockOne)
			repositoryTwo.CreateOrUpdateBlock(blockTwo)
			retrievedBlockOne, _ := repository.FindBlockByNumber(123)
			retrievedBlockTwo, _ := repositoryTwo.FindBlockByNumber(123)

			Expect(retrievedBlockOne.Transactions[0].Hash).To(Equal("x123"))
			Expect(retrievedBlockTwo.Transactions[0].Hash).To(Equal("x678"))
		})

		It("saves the attributes associated to a transaction", func() {
			gasLimit := int64(5000)
			gasPrice := int64(3)
			nonce := uint64(10000)
			to := "1234567890"
			from := "0987654321"
			var value = new(big.Int)
			value.SetString("34940183920000000000", 10)
			inputData := "0xf7d8c8830000000000000000000000000000000000000000000000000000000000037788000000000000000000000000000000000000000000000000000000000003bd14"
			transaction := core.Transaction{
				Hash:     "x1234",
				GasPrice: gasPrice,
				GasLimit: gasLimit,
				Nonce:    nonce,
				To:       to,
				From:     from,
				Value:    value.String(),
				Data:     inputData,
			}
			block := core.Block{
				Number:       123,
				Transactions: []core.Transaction{transaction},
			}

			repository.CreateOrUpdateBlock(block)

			savedBlock, _ := repository.FindBlockByNumber(123)
			Expect(len(savedBlock.Transactions)).To(Equal(1))
			savedTransaction := savedBlock.Transactions[0]
			Expect(savedTransaction.Data).To(Equal(transaction.Data))
			Expect(savedTransaction.Hash).To(Equal(transaction.Hash))
			Expect(savedTransaction.To).To(Equal(to))
			Expect(savedTransaction.From).To(Equal(from))
			Expect(savedTransaction.Nonce).To(Equal(nonce))
			Expect(savedTransaction.GasLimit).To(Equal(gasLimit))
			Expect(savedTransaction.GasPrice).To(Equal(gasPrice))
			Expect(savedTransaction.Value).To(Equal(value.String()))
		})

	})

	Describe("The missing block numbers", func() {
		It("is empty the starting block number is the highest known block number", func() {
			repository.CreateOrUpdateBlock(core.Block{Number: 1})

			Expect(len(repository.MissingBlockNumbers(1, 1))).To(Equal(0))
		})

		It("is the only missing block number", func() {
			repository.CreateOrUpdateBlock(core.Block{Number: 2})

			Expect(repository.MissingBlockNumbers(1, 2)).To(Equal([]int64{1}))
		})

		It("is both missing block numbers", func() {
			repository.CreateOrUpdateBlock(core.Block{Number: 3})

			Expect(repository.MissingBlockNumbers(1, 3)).To(Equal([]int64{1, 2}))
		})

		It("goes back to the starting block number", func() {
			repository.CreateOrUpdateBlock(core.Block{Number: 6})

			Expect(repository.MissingBlockNumbers(4, 6)).To(Equal([]int64{4, 5}))
		})

		It("only includes missing block numbers", func() {
			repository.CreateOrUpdateBlock(core.Block{Number: 4})
			repository.CreateOrUpdateBlock(core.Block{Number: 6})

			Expect(repository.MissingBlockNumbers(4, 6)).To(Equal([]int64{5}))
		})

		It("is a list with multiple gaps", func() {
			repository.CreateOrUpdateBlock(core.Block{Number: 4})
			repository.CreateOrUpdateBlock(core.Block{Number: 5})
			repository.CreateOrUpdateBlock(core.Block{Number: 8})
			repository.CreateOrUpdateBlock(core.Block{Number: 10})

			Expect(repository.MissingBlockNumbers(3, 10)).To(Equal([]int64{3, 6, 7, 9}))
		})

		It("returns empty array when lower bound exceeds upper bound", func() {
			Expect(repository.MissingBlockNumbers(10000, 1)).To(Equal([]int64{}))
		})

		It("only returns requested range even when other gaps exist", func() {
			repository.CreateOrUpdateBlock(core.Block{Number: 3})
			repository.CreateOrUpdateBlock(core.Block{Number: 8})

			Expect(repository.MissingBlockNumbers(1, 5)).To(Equal([]int64{1, 2, 4, 5}))
		})
	})

	Describe("The max block numbers", func() {
		It("returns the block number when a single block", func() {
			repository.CreateOrUpdateBlock(core.Block{Number: 1})

			Expect(repository.MaxBlockNumber()).To(Equal(int64(1)))
		})

		It("returns highest known block number when multiple blocks", func() {
			repository.CreateOrUpdateBlock(core.Block{Number: 1})
			repository.CreateOrUpdateBlock(core.Block{Number: 10})

			Expect(repository.MaxBlockNumber()).To(Equal(int64(10)))
		})
	})

	Describe("The block status", func() {
		It("sets the status of blocks within n-20 of chain HEAD as final", func() {
			blockNumberOfChainHead := 25
			for i := 0; i < blockNumberOfChainHead; i++ {
				repository.CreateOrUpdateBlock(core.Block{Number: int64(i), Hash: strconv.Itoa(i)})
			}

			repository.SetBlocksStatus(int64(blockNumberOfChainHead))

			blockOne, err := repository.FindBlockByNumber(1)
			Expect(err).ToNot(HaveOccurred())
			Expect(blockOne.IsFinal).To(Equal(true))
			blockTwo, err := repository.FindBlockByNumber(24)
			Expect(err).ToNot(HaveOccurred())
			Expect(blockTwo.IsFinal).To(BeFalse())
		})

	})

	Describe("Creating contracts", func() {
		It("returns the contract when it exists", func() {
			repository.CreateContract(core.Contract{Hash: "x123"})

			contract, err := repository.FindContract("x123")
			Expect(err).NotTo(HaveOccurred())
			Expect(contract.Hash).To(Equal("x123"))

			Expect(repository.ContractExists("x123")).To(BeTrue())
			Expect(repository.ContractExists("x456")).To(BeFalse())
		})

		It("returns err if contract does not exist", func() {
			_, err := repository.FindContract("x123")
			Expect(err).To(HaveOccurred())
		})

		It("returns empty array when no transactions 'To' a contract", func() {
			repository.CreateContract(core.Contract{Hash: "x123"})
			contract, err := repository.FindContract("x123")
			Expect(err).ToNot(HaveOccurred())
			Expect(contract.Transactions).To(BeEmpty())
		})

		It("returns transactions 'To' a contract", func() {
			block := core.Block{
				Number: 123,
				Transactions: []core.Transaction{
					{Hash: "TRANSACTION1", To: "x123"},
					{Hash: "TRANSACTION2", To: "x345"},
					{Hash: "TRANSACTION3", To: "x123"},
				},
			}
			repository.CreateOrUpdateBlock(block)

			repository.CreateContract(core.Contract{Hash: "x123"})
			contract, err := repository.FindContract("x123")
			Expect(err).ToNot(HaveOccurred())
			sort.Slice(contract.Transactions, func(i, j int) bool {
				return contract.Transactions[i].Hash < contract.Transactions[j].Hash
			})
			Expect(contract.Transactions).To(
				Equal([]core.Transaction{
					{Hash: "TRANSACTION1", To: "x123"},
					{Hash: "TRANSACTION3", To: "x123"},
				}))
		})

		It("stores the ABI of the contract", func() {
			repository.CreateContract(core.Contract{
				Abi:  "{\"some\": \"json\"}",
				Hash: "x123",
			})
			contract, err := repository.FindContract("x123")
			Expect(err).ToNot(HaveOccurred())
			Expect(contract.Abi).To(Equal("{\"some\": \"json\"}"))
		})

		It("updates the ABI of the contract if hash already present", func() {
			repository.CreateContract(core.Contract{
				Abi:  "{\"some\": \"json\"}",
				Hash: "x123",
			})
			repository.CreateContract(core.Contract{
				Abi:  "{\"some\": \"different json\"}",
				Hash: "x123",
			})
			contract, err := repository.FindContract("x123")
			Expect(err).ToNot(HaveOccurred())
			Expect(contract.Abi).To(Equal("{\"some\": \"different json\"}"))
		})
	})

	Describe("Saving logs", func() {
		It("returns the log when it exists", func() {
			repository.CreateLogs([]core.Log{{
				BlockNumber: 1,
				Index:       0,
				Address:     "x123",
				TxHash:      "x456",
				Topics:      core.Topics{0: "x777", 1: "x888", 2: "x999"},
				Data:        "xabc",
			}},
			)

			log := repository.FindLogs("x123", 1)

			Expect(log).NotTo(BeNil())
			Expect(log[0].BlockNumber).To(Equal(int64(1)))
			Expect(log[0].Address).To(Equal("x123"))
			Expect(log[0].Index).To(Equal(int64(0)))
			Expect(log[0].TxHash).To(Equal("x456"))
			Expect(log[0].Topics[0]).To(Equal("x777"))
			Expect(log[0].Topics[1]).To(Equal("x888"))
			Expect(log[0].Topics[2]).To(Equal("x999"))
			Expect(log[0].Data).To(Equal("xabc"))
		})

		It("returns nil if log does not exist", func() {
			log := repository.FindLogs("x123", 1)
			Expect(log).To(BeNil())
		})

		It("filters to the correct block number and address", func() {
			repository.CreateLogs([]core.Log{{
				BlockNumber: 1,
				Index:       0,
				Address:     "x123",
				TxHash:      "x456",
				Topics:      core.Topics{0: "x777", 1: "x888", 2: "x999"},
				Data:        "xabc",
			}},
			)
			repository.CreateLogs([]core.Log{{
				BlockNumber: 1,
				Index:       1,
				Address:     "x123",
				TxHash:      "x789",
				Topics:      core.Topics{0: "x111", 1: "x222", 2: "x333"},
				Data:        "xdef",
			}},
			)
			repository.CreateLogs([]core.Log{{
				BlockNumber: 2,
				Index:       0,
				Address:     "x123",
				TxHash:      "x456",
				Topics:      core.Topics{0: "x777", 1: "x888", 2: "x999"},
				Data:        "xabc",
			}},
			)

			log := repository.FindLogs("x123", 1)

			type logIndex struct {
				blockNumber int64
				Index       int64
			}
			var uniqueBlockNumbers []logIndex
			for _, log := range log {
				uniqueBlockNumbers = append(uniqueBlockNumbers,
					logIndex{log.BlockNumber, log.Index})
			}
			sort.Slice(uniqueBlockNumbers, func(i, j int) bool {
				if uniqueBlockNumbers[i].blockNumber < uniqueBlockNumbers[j].blockNumber {
					return true
				}
				if uniqueBlockNumbers[i].blockNumber > uniqueBlockNumbers[j].blockNumber {
					return false
				}
				return uniqueBlockNumbers[i].Index < uniqueBlockNumbers[j].Index
			})

			Expect(log).NotTo(BeNil())
			Expect(len(log)).To(Equal(2))
			Expect(uniqueBlockNumbers).To(Equal(
				[]logIndex{
					{blockNumber: 1, Index: 0},
					{blockNumber: 1, Index: 1}},
			))
		})

		It("saves the logs attached to a receipt", func() {
			logs := []core.Log{{
				Address:     "0x8a4774fe82c63484afef97ca8d89a6ea5e21f973",
				BlockNumber: 4745407,
				Data:        "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000645a68669900000000000000000000000000000000000000000000003397684ab5869b0000000000000000000000000000000000000000000000000000000000005a36053200000000000000000000000099041f808d598b782d5a3e498681c2452a31da08",
				Index:       86,
				Topics: core.Topics{
					0: "0x5a68669900000000000000000000000000000000000000000000000000000000",
					1: "0x000000000000000000000000d0148dad63f73ce6f1b6c607e3413dcf1ff5f030",
					2: "0x00000000000000000000000000000000000000000000003397684ab5869b0000",
					3: "0x000000000000000000000000000000000000000000000000000000005a360532",
				},
				TxHash: "0x002c4799161d809b23f67884eb6598c9df5894929fe1a9ead97ca175d360f547",
			}, {
				Address:     "0x99041f808d598b782d5a3e498681c2452a31da08",
				BlockNumber: 4745407,
				Data:        "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000418178358",
				Index:       87,
				Topics: core.Topics{
					0: "0x1817835800000000000000000000000000000000000000000000000000000000",
					1: "0x0000000000000000000000008a4774fe82c63484afef97ca8d89a6ea5e21f973",
					2: "0x0000000000000000000000000000000000000000000000000000000000000000",
					3: "0x0000000000000000000000000000000000000000000000000000000000000000",
				},
				TxHash: "0x002c4799161d809b23f67884eb6598c9df5894929fe1a9ead97ca175d360f547",
			}, {
				Address:     "0x99041f808d598b782d5a3e498681c2452a31da08",
				BlockNumber: 4745407,
				Data:        "0x00000000000000000000000000000000000000000000003338f64c8423af4000",
				Index:       88,
				Topics: core.Topics{
					0: "0x296ba4ca62c6c21c95e828080cb8aec7481b71390585605300a8a76f9e95b527",
				},
				TxHash: "0x002c4799161d809b23f67884eb6598c9df5894929fe1a9ead97ca175d360f547",
			},
			}
			receipt := core.Receipt{
				ContractAddress:   "",
				CumulativeGasUsed: 7481414,
				GasUsed:           60711,
				Logs:              logs,
				Bloom:             "0x00000800000000000000001000000000000000400000000080000000000000000000400000010000000000000000000000000000040000000000000002000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000800004000000000000001000000000000000000000000000002000000480000000000000002000000000000000020000000000000000000000000000000000000000080000000000180000c00000000000000002000002000000040000000000000000000000000000010000000000020000000000000000000002000000000000000000000000400800000000000000000",
				Status:            1,
				TxHash:            "0x002c4799161d809b23f67884eb6598c9df5894929fe1a9ead97ca175d360f547",
			}
			transaction :=
				core.Transaction{
					Hash:    receipt.TxHash,
					Receipt: receipt,
				}

			block := core.Block{Transactions: []core.Transaction{transaction}}
			err := repository.CreateOrUpdateBlock(block)
			Expect(err).To(Not(HaveOccurred()))
			retrievedLogs := repository.FindLogs("0x99041f808d598b782d5a3e498681c2452a31da08", 4745407)

			expected := logs[1:]
			Expect(retrievedLogs).To(Equal(expected))
		})

		It("still saves receipts without logs", func() {
			receipt := core.Receipt{
				TxHash: "0x002c4799161d809b23f67884eb6598c9df5894929fe1a9ead97ca175d360f547",
			}
			transaction := core.Transaction{
				Hash:    receipt.TxHash,
				Receipt: receipt,
			}

			block := core.Block{
				Transactions: []core.Transaction{transaction},
			}
			repository.CreateOrUpdateBlock(block)

			_, err := repository.FindReceipt(receipt.TxHash)

			Expect(err).To(Not(HaveOccurred()))
		})
	})

	Describe("Saving receipts", func() {
		It("returns the receipt when it exists", func() {
			expected := core.Receipt{
				ContractAddress:   "0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae",
				CumulativeGasUsed: 7996119,
				GasUsed:           21000,
				Logs:              []core.Log{},
				StateRoot:         "0x88abf7e73128227370aa7baa3dd4e18d0af70e92ef1f9ef426942fbe2dddb733",
				Status:            1,
				TxHash:            "0xe340558980f89d5f86045ac11e5cc34e4bcec20f9f1e2a427aa39d87114e8223",
			}

			transaction := core.Transaction{
				Hash:    expected.TxHash,
				Receipt: expected,
			}

			block := core.Block{Transactions: []core.Transaction{transaction}}
			repository.CreateOrUpdateBlock(block)
			receipt, err := repository.FindReceipt("0xe340558980f89d5f86045ac11e5cc34e4bcec20f9f1e2a427aa39d87114e8223")

			Expect(err).ToNot(HaveOccurred())
			//Not currently serializing bloom logs
			Expect(receipt.Bloom).To(Equal(core.Receipt{}.Bloom))
			Expect(receipt.TxHash).To(Equal(expected.TxHash))
			Expect(receipt.CumulativeGasUsed).To(Equal(expected.CumulativeGasUsed))
			Expect(receipt.GasUsed).To(Equal(expected.GasUsed))
			Expect(receipt.StateRoot).To(Equal(expected.StateRoot))
			Expect(receipt.Status).To(Equal(expected.Status))
		})

		It("returns ErrReceiptDoesNotExist when receipt does not exist", func() {
			receipt, err := repository.FindReceipt("DOES NOT EXIST")
			Expect(err).To(HaveOccurred())
			Expect(receipt).To(BeZero())
		})
	})

	Describe("LogFilter", func() {

		It("inserts filter into watched events", func() {

			logFilter := filters.LogFilter{
				Name:      "TestFilter",
				FromBlock: 1,
				ToBlock:   2,
				Address:   "0x8888f1f195afa192cfee860698584c030f4c9db1",
				Topics: core.Topics{
					"0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
					"",
					"0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
					"",
				},
			}
			err := repository.AddFilter(logFilter)
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns error if name is not provided", func() {

			logFilter := filters.LogFilter{
				FromBlock: 1,
				ToBlock:   2,
				Address:   "0x8888f1f195afa192cfee860698584c030f4c9db1",
				Topics: core.Topics{
					"0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
					"",
					"0x000000000000000000000000a94f5374fce5edbc8e2a8697c15331677e6ebf0b",
					"",
				},
			}
			err := repository.AddFilter(logFilter)
			Expect(err).To(HaveOccurred())
		})
	})
}
