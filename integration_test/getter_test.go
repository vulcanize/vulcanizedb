package integration

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/constants"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/getter"
	"github.com/makerdao/vulcanizedb/pkg/eth"
	"github.com/makerdao/vulcanizedb/pkg/eth/client"
	"github.com/makerdao/vulcanizedb/pkg/eth/converters"
	"github.com/makerdao/vulcanizedb/pkg/eth/node"
	"github.com/makerdao/vulcanizedb/test_config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Interface Getter", func() {
	Describe("GetAbi", func() {
		It("Constructs and returns a custom abi based on results from supportsInterface calls", func() {
			expectedABI := `[` + constants.AddrChangeInterface + `,` + constants.NameChangeInterface + `,` + constants.ContentChangeInterface + `,` + constants.AbiChangeInterface + `,` + constants.PubkeyChangeInterface + `]`
			con := test_config.TestClient
			testIPC := con.IPCPath
			blockNumber := int64(6885696)
			rawRpcClient, err := rpc.Dial(testIPC)
			Expect(err).NotTo(HaveOccurred())
			rpcClient := client.NewRpcClient(rawRpcClient, testIPC)
			ethClient := ethclient.NewClient(rawRpcClient)
			blockChainClient := client.NewEthClient(ethClient)
			node := node.MakeNode(rpcClient)
			transactionConverter := converters.NewTransactionConverter(ethClient)
			blockChain := eth.NewBlockChain(blockChainClient, rpcClient, node, transactionConverter)
			interfaceGetter := getter.NewInterfaceGetter(blockChain)
			abi, err := interfaceGetter.GetABI(constants.PublicResolverAddress, blockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(abi).To(Equal(expectedABI))
			_, err = eth.ParseAbi(abi)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
