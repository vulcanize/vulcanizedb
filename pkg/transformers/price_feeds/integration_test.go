package price_feeds_test

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	rpc2 "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
	"github.com/vulcanize/vulcanizedb/test_config"
	"math/big"
)

var _ = Describe("Price feeds transformer", func() {
	var (
		db         *postgres.DB
		blockChain core.BlockChain
	)

	BeforeEach(func() {
		ipc := "https://kovan.infura.io/J5Vd2fRtGsw0zZ0Ov3BL"
		rpcClient, ethClient, err := getClients(ipc)
		Expect(err).NotTo(HaveOccurred())
		blockChain, err = getBlockChain(rpcClient, ethClient)
		Expect(err).NotTo(HaveOccurred())
		db = test_config.NewTestDB(blockChain.Node())
		test_config.CleanTestDB(db)

		for i := 8763054; i < 8763063; i++ {
			err = persistHeader(rpcClient, db, int64(i))
			Expect(err).NotTo(HaveOccurred())
		}
	})

	It("persists a ETH/USD price feed event", func() {
		config := price_feeds.IPriceFeedConfig{
			ContractAddresses:   []string{"0x9FfFE440258B79c5d6604001674A4722FfC0f7Bc"},
			StartingBlockNumber: 8763054,
			EndingBlockNumber:   8763054,
		}
		transformerInitializer := price_feeds.PriceFeedTransformerInitializer{Config: config}
		transformer := transformerInitializer.NewPriceFeedTransformer(db, blockChain)

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		var model price_feeds.PriceFeedModel
		err = db.Get(&model, `SELECT block_number, medianizer_address, usd_value, tx_idx, raw_log FROM maker.price_feeds WHERE block_number = $1`, config.StartingBlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.UsdValue).To(Equal("207.314891143"))
	})

	It("persists a MKR/USD price feed event", func() {
		config := price_feeds.IPriceFeedConfig{
			ContractAddresses:   []string{"0xB1997239Cfc3d15578A3a09730f7f84A90BB4975"},
			StartingBlockNumber: 8763059,
			EndingBlockNumber:   8763059,
		}
		transformerInitializer := price_feeds.PriceFeedTransformerInitializer{Config: config}
		transformer := transformerInitializer.NewPriceFeedTransformer(db, blockChain)

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		var model price_feeds.PriceFeedModel
		err = db.Get(&model, `SELECT block_number, medianizer_address, usd_value, tx_idx, raw_log FROM maker.price_feeds WHERE block_number = $1`, config.StartingBlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.UsdValue).To(Equal("391.803979212"))
	})

	It("persists a REP/USD price feed event", func() {
		config := price_feeds.IPriceFeedConfig{
			ContractAddresses:   []string{"0xf88bBDc1E2718F8857F30A180076ec38d53cf296"},
			StartingBlockNumber: 8763062,
			EndingBlockNumber:   8763062,
		}
		transformerInitializer := price_feeds.PriceFeedTransformerInitializer{Config: config}
		transformer := transformerInitializer.NewPriceFeedTransformer(db, blockChain)

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		var model price_feeds.PriceFeedModel
		err = db.Get(&model, `SELECT block_number, medianizer_address, usd_value, tx_idx, raw_log FROM maker.price_feeds WHERE block_number = $1`, config.StartingBlockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(model.UsdValue).To(Equal("12.8169284827"))
	})
})

type POAHeader struct {
	ParentHash  common.Hash    `json:"parentHash"       gencodec:"required"`
	UncleHash   common.Hash    `json:"sha3Uncles"       gencodec:"required"`
	Coinbase    common.Address `json:"miner"            gencodec:"required"`
	Root        common.Hash    `json:"stateRoot"        gencodec:"required"`
	TxHash      common.Hash    `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash common.Hash    `json:"receiptsRoot"     gencodec:"required"`
	Bloom       types.Bloom    `json:"logsBloom"        gencodec:"required"`
	Difficulty  *hexutil.Big   `json:"difficulty"       gencodec:"required"`
	Number      *hexutil.Big   `json:"number"           gencodec:"required"`
	GasLimit    hexutil.Uint64 `json:"gasLimit"         gencodec:"required"`
	GasUsed     hexutil.Uint64 `json:"gasUsed"          gencodec:"required"`
	Time        *hexutil.Big   `json:"timestamp"        gencodec:"required"`
	Extra       hexutil.Bytes  `json:"extraData"        gencodec:"required"`
	Hash        common.Hash    `json:"hash"`
}

func getClients(ipc string) (client.RpcClient, *ethclient.Client, error) {
	raw, err := rpc.Dial(ipc)
	if err != nil {
		return client.RpcClient{}, &ethclient.Client{}, err
	}
	return client.NewRpcClient(raw, ipc), ethclient.NewClient(raw), nil
}

func getBlockChain(rpcClient client.RpcClient, ethClient *ethclient.Client) (core.BlockChain, error) {
	client := client.NewEthClient(ethClient)
	node := node.MakeNode(rpcClient)
	transactionConverter := rpc2.NewRpcTransactionConverter(client)
	blockChain := geth.NewBlockChain(client, node, transactionConverter)
	return blockChain, nil
}

func persistHeader(rpcClient client.RpcClient, db *postgres.DB, blockNumber int64) error {
	var poaHeader POAHeader
	blockNumberArg := hexutil.EncodeBig(big.NewInt(int64(blockNumber)))
	err := rpcClient.CallContext(context.Background(), &poaHeader, "eth_getBlockByNumber", blockNumberArg, false)
	if err != nil {
		return err
	}
	headerRepository := repositories.NewHeaderRepository(db)
	_, err = headerRepository.CreateOrUpdateHeader(core.Header{
		BlockNumber: poaHeader.Number.ToInt().Int64(),
		Hash:        poaHeader.Hash.String(),
	})
	return err
}
