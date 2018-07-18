package integration

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	rpc2 "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Rewards calculations", func() {

	It("calculates a block reward for a real block", func() {
		rpcClient, err := rpc.Dial(test_config.InfuraClient.IPCPath)
		Expect(err).NotTo(HaveOccurred())
		ethClient := ethclient.NewClient(rpcClient)
		blockChainClient := client.NewClient(ethClient)
		clientWrapper := node.ClientWrapper{
			ContextCaller: rpcClient,
			IPCPath:       test_config.InfuraClient.IPCPath,
		}
		node := node.MakeNode(clientWrapper)
		transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
		blockChain := geth.NewBlockChain(blockChainClient, node, transactionConverter)
		block, err := blockChain.GetBlockByNumber(1071819)
		Expect(err).ToNot(HaveOccurred())
		Expect(block.Reward).To(Equal(5.31355))
	})

	It("calculates an uncle reward for a real block", func() {
		rpcClient, err := rpc.Dial(test_config.InfuraClient.IPCPath)
		Expect(err).NotTo(HaveOccurred())
		ethClient := ethclient.NewClient(rpcClient)
		blockChainClient := client.NewClient(ethClient)
		clientWrapper := node.ClientWrapper{
			ContextCaller: rpcClient,
			IPCPath:       test_config.InfuraClient.IPCPath,
		}
		node := node.MakeNode(clientWrapper)
		transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
		blockChain := geth.NewBlockChain(blockChainClient, node, transactionConverter)
		block, err := blockChain.GetBlockByNumber(1071819)
		Expect(err).ToNot(HaveOccurred())
		Expect(block.UnclesReward).To(Equal(6.875))
	})

})
