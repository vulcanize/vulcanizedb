// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package transformer_test

import (
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/examples/constants"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	rpc2 "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
	"github.com/vulcanize/vulcanizedb/pkg/omni/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/omni/types"
)

var mockEvent = core.WatchedEvent{
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

var _ = Describe("Repository Test", func() {
	var db *postgres.DB
	var err error
	var con types.Config
	var blockRepository repositories.BlockRepository
	rand.Seed(time.Now().UnixNano())

	BeforeEach(func() {
		infuraIPC := "https://mainnet.infura.io/v3/b09888c1113640cc9ab42750ce750c05"
		rawRpcClient, err := rpc.Dial(infuraIPC)
		Expect(err).NotTo(HaveOccurred())
		rpcClient := client.NewRpcClient(rawRpcClient, infuraIPC)
		ethClient := ethclient.NewClient(rawRpcClient)
		blockChainClient := client.NewEthClient(ethClient)
		node := node.MakeNode(rpcClient)
		transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
		blockChain := geth.NewBlockChain(blockChainClient, node, transactionConverter)

		db, err = postgres.NewDB(config.Database{
			Hostname: "localhost",
			Name:     "vulcanize_private",
			Port:     5432,
		}, blockChain.Node())
		Expect(err).NotTo(HaveOccurred())

		con = types.Config{
			DB:      db,
			BC:      blockChain,
			Network: "",
		}

		blockRepository = *repositories.NewBlockRepository(db)
	})

	AfterEach(func() {
		db.Query(`DELETE FROM blocks`)
		db.Query(`DELETE FROM logs`)
		db.Query(`DELETE FROM transactions`)
		db.Query(`DELETE FROM receipts`)
		db.Query(`DROP SCHEMA IF EXISTS trueusd CASCADE`)
	})

	It("Fails to initialize if first and most recent blocks cannot be fetched from vDB", func() {
		t := transformer.NewTransformer(&con)
		t.SetEvents(constants.TusdContractAddress, []string{"Transfer"})
		err = t.Init()
		Expect(err).To(HaveOccurred())
	})

	It("Initializes and executes successfully if first and most recent blocks can be fetched from vDB", func() {
		log := core.Log{
			BlockNumber: 6194634,
			TxHash:      "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad234hfs",
			Address:     constants.TusdContractAddress,
			Topics: core.Topics{
				constants.TransferEvent.Signature(),
				"0x000000000000000000000000000000000000000000000000000000000000af21",
				"0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391",
				"",
			},
			Index: 1,
			Data:  "0x000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359000000000000000000000000000000000000000000000000392d2e2bda9c00000000000000000000000000000000000000000000000000927f41fa0a4a418000000000000000000000000000000000000000000000000000000000005adcfebe",
		}

		receipt := core.Receipt{
			TxHash:          "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
			ContractAddress: constants.TusdContractAddress,
			Logs:            []core.Log{log},
		}

		transaction := core.Transaction{
			Hash:    "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad6546ae",
			Receipt: receipt,
		}

		block := core.Block{
			Hash:         "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad123ert",
			Number:       6194634,
			Transactions: []core.Transaction{transaction},
		}

		blockRepository.CreateOrUpdateBlock(block)

		t := transformer.NewTransformer(&con)
		t.SetEvents(constants.TusdContractAddress, []string{"Transfer"})
		err = t.Init()
		Expect(err).ToNot(HaveOccurred())

		c, ok := t.Contracts[constants.TusdContractAddress]
		Expect(ok).To(Equal(true))

		err = t.Execute()
		Expect(err).ToNot(HaveOccurred())

		b, ok := c.Addresses["0x000000000000000000000000000000000000Af21"]
		Expect(ok).To(Equal(true))
		Expect(b).To(Equal(true))

		b, ok = c.Addresses["0x09BbBBE21a5975cAc061D82f7b843bCE061BA391"]
		Expect(ok).To(Equal(true))
		Expect(b).To(Equal(true))
	})
})
