package integration

/* WIP
import (
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/ethclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/helpers/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/test_config"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	rpc2 "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
)

var _ = Describe("IPFS Processor", func() {
	var processor ipfs.SyncAndPublish
	var err error
	var db *postgres.DB
	var bc core.BlockChain
	var ec core.EthClient
	var rc core.RpcClient
	var quitChan chan bool

	AfterEach(func() {
		test_helpers.TearDown(db)
	})

	BeforeEach(func() {
		db, bc, ec, rc = setup()
		quitChan = make(chan bool)
		processor, err = ipfs.NewIPFSProcessor("~/.ipfs", db, ec, rc, quitChan)
	})

	Describe("Process", func() {
		It("Polls specified contract methods using contract's argument list", func() {

		})
	})
})

func setup() (*postgres.DB, core.BlockChain, core.EthClient, core.RpcClient) {
	con := test_config.InfuraClient
	infuraIPC := con.IPCPath
	rawRpcClient, err := rpc.Dial(infuraIPC)
	Expect(err).NotTo(HaveOccurred())
	rpcClient := client.NewRpcClient(rawRpcClient, infuraIPC)
	ethClient := ethclient.NewClient(rawRpcClient)
	blockChainClient := client.NewEthClient(ethClient)
	node := node.MakeNode(rpcClient)
	transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
	blockChain := geth.NewBlockChain(blockChainClient, rpcClient, node, transactionConverter)

	db, err := postgres.NewDB(config.Database{
		Hostname: "localhost",
		Name:     "vulcanize_private",
		Port:     5432,
	}, blockChain.Node())
	Expect(err).NotTo(HaveOccurred())

	return db, blockChain, ethClient, rpcClient
}
*/
