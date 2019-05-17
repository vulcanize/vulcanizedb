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

package integration

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/vulcanize/vulcanizedb/pkg/config"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/full/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/helpers/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/helpers/test_helpers/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
)

var _ = Describe("contractWatcher full transformer", func() {
	var db *postgres.DB
	var err error
	var blockChain core.BlockChain
	var blockRepository repositories.BlockRepository
	var ensAddr = strings.ToLower(constants.EnsContractAddress)
	var tusdAddr = strings.ToLower(constants.TusdContractAddress)
	rand.Seed(time.Now().UnixNano())

	BeforeEach(func() {
		db, blockChain = test_helpers.SetupDBandBC()
		blockRepository = *repositories.NewBlockRepository(db)
	})

	AfterEach(func() {
		test_helpers.TearDown(db)
	})

	Describe("Init", func() {
		It("Initializes transformer's contract objects", func() {
			_, insertErr := blockRepository.CreateOrUpdateBlock(mocks.TransferBlock1)
			Expect(insertErr).NotTo(HaveOccurred())
			_, insertErrTwo := blockRepository.CreateOrUpdateBlock(mocks.TransferBlock2)
			Expect(insertErrTwo).NotTo(HaveOccurred())
			t := transformer.NewTransformer(test_helpers.TusdConfig, blockChain, db)
			err = t.Init()
			Expect(err).ToNot(HaveOccurred())

			c, ok := t.Contracts[tusdAddr]
			Expect(ok).To(Equal(true))

			Expect(c.StartingBlock).To(Equal(int64(6194633)))
			Expect(c.Abi).To(Equal(constants.TusdAbiString))
			Expect(c.Name).To(Equal("TrueUSD"))
			Expect(c.Address).To(Equal(tusdAddr))
		})

		It("Fails to initialize if first and most recent blocks cannot be fetched from vDB", func() {
			t := transformer.NewTransformer(test_helpers.TusdConfig, blockChain, db)
			err = t.Init()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no rows in result set"))
		})

		It("Does nothing if watched events are unset", func() {
			_, insertErr := blockRepository.CreateOrUpdateBlock(mocks.TransferBlock1)
			Expect(insertErr).NotTo(HaveOccurred())
			_, insertErrTwo := blockRepository.CreateOrUpdateBlock(mocks.TransferBlock2)
			Expect(insertErrTwo).NotTo(HaveOccurred())
			var testConf config.ContractConfig
			testConf = test_helpers.TusdConfig
			testConf.Events = nil
			t := transformer.NewTransformer(testConf, blockChain, db)
			err = t.Init()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no filters created"))

			_, ok := t.Contracts[tusdAddr]
			Expect(ok).To(Equal(false))
		})
	})

	Describe("Execute", func() {
		BeforeEach(func() {
			_, insertErr := blockRepository.CreateOrUpdateBlock(mocks.TransferBlock1)
			Expect(insertErr).NotTo(HaveOccurred())
			_, insertErrTwo := blockRepository.CreateOrUpdateBlock(mocks.TransferBlock2)
			Expect(insertErrTwo).NotTo(HaveOccurred())
		})

		It("Transforms watched contract data into custom repositories", func() {
			t := transformer.NewTransformer(test_helpers.TusdConfig, blockChain, db)
			err = t.Init()
			Expect(err).ToNot(HaveOccurred())

			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())

			log := test_helpers.TransferLog{}
			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM full_%s.transfer_event WHERE block = 6194634", tusdAddr)).StructScan(&log)

			// We don't know vulcID, so compare individual fields instead of complete structures
			Expect(log.Tx).To(Equal("0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad654eee"))
			Expect(log.Block).To(Equal(int64(6194634)))
			Expect(log.From).To(Equal("0x000000000000000000000000000000000000Af21"))
			Expect(log.To).To(Equal("0x09BbBBE21a5975cAc061D82f7b843bCE061BA391"))
			Expect(log.Value).To(Equal("1097077688018008265106216665536940668749033598146"))
		})

		It("Keeps track of contract-related addresses while transforming event data if they need to be used for later method polling", func() {
			var testConf config.ContractConfig
			testConf = test_helpers.TusdConfig
			testConf.Methods = map[string][]string{
				tusdAddr: {"balanceOf"},
			}
			t := transformer.NewTransformer(testConf, blockChain, db)
			err = t.Init()
			Expect(err).ToNot(HaveOccurred())

			c, ok := t.Contracts[tusdAddr]
			Expect(ok).To(Equal(true))

			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())

			b, ok := c.EmittedAddrs[common.HexToAddress("0x000000000000000000000000000000000000Af21")]
			Expect(ok).To(Equal(true))
			Expect(b).To(Equal(true))

			b, ok = c.EmittedAddrs[common.HexToAddress("0x09BbBBE21a5975cAc061D82f7b843bCE061BA391")]
			Expect(ok).To(Equal(true))
			Expect(b).To(Equal(true))

			_, ok = c.EmittedAddrs[common.HexToAddress("0x09BbBBE21a5975cAc061D82f7b843b1234567890")]
			Expect(ok).To(Equal(false))

			_, ok = c.EmittedAddrs[common.HexToAddress("0x")]
			Expect(ok).To(Equal(false))

			_, ok = c.EmittedAddrs[""]
			Expect(ok).To(Equal(false))

			_, ok = c.EmittedAddrs[common.HexToAddress("0x09THISE21a5IS5cFAKE1D82fAND43bCE06MADEUP")]
			Expect(ok).To(Equal(false))
		})

		It("Polls given methods using generated token holder address", func() {
			var testConf config.ContractConfig
			testConf = test_helpers.TusdConfig
			testConf.Methods = map[string][]string{
				tusdAddr: {"balanceOf"},
			}
			t := transformer.NewTransformer(testConf, blockChain, db)
			err = t.Init()
			Expect(err).ToNot(HaveOccurred())

			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())

			res := test_helpers.BalanceOf{}

			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM full_%s.balanceof_method WHERE who_ = '0x000000000000000000000000000000000000Af21' AND block = '6194634'", tusdAddr)).StructScan(&res)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Balance).To(Equal("0"))
			Expect(res.TokenName).To(Equal("TrueUSD"))

			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM full_%s.balanceof_method WHERE who_ = '0x09BbBBE21a5975cAc061D82f7b843bCE061BA391' AND block = '6194634'", tusdAddr)).StructScan(&res)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Balance).To(Equal("0"))
			Expect(res.TokenName).To(Equal("TrueUSD"))

			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM full_%s.balanceof_method WHERE who_ = '0xfE9e8709d3215310075d67E3ed32A380CCf451C8' AND block = '6194634'", tusdAddr)).StructScan(&res)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no rows in result set"))
		})

		It("Fails if initialization has not been done", func() {
			t := transformer.NewTransformer(test_helpers.TusdConfig, blockChain, db)

			err = t.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("transformer has no initialized contracts to work with"))
		})
	})

	Describe("Execute- against ENS registry contract", func() {
		BeforeEach(func() {
			_, insertErr := blockRepository.CreateOrUpdateBlock(mocks.NewOwnerBlock1)
			Expect(insertErr).NotTo(HaveOccurred())
			_, insertErrTwo := blockRepository.CreateOrUpdateBlock(mocks.NewOwnerBlock2)
			Expect(insertErrTwo).NotTo(HaveOccurred())
		})

		It("Transforms watched contract data into custom repositories", func() {
			t := transformer.NewTransformer(test_helpers.ENSConfig, blockChain, db)

			err = t.Init()
			Expect(err).ToNot(HaveOccurred())

			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())

			log := test_helpers.NewOwnerLog{}
			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM full_%s.newowner_event", ensAddr)).StructScan(&log)

			// We don't know vulcID, so compare individual fields instead of complete structures
			Expect(log.Tx).To(Equal("0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad654bbb"))
			Expect(log.Block).To(Equal(int64(6194635)))
			Expect(log.Node).To(Equal("0x0000000000000000000000000000000000000000000000000000c02aaa39b223"))
			Expect(log.Label).To(Equal("0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391"))
			Expect(log.Owner).To(Equal("0x000000000000000000000000000000000000Af21"))
		})

		It("Keeps track of contract-related hashes while transforming event data if they need to be used for later method polling", func() {
			var testConf config.ContractConfig
			testConf = test_helpers.ENSConfig
			testConf.Methods = map[string][]string{
				ensAddr: {"owner"},
			}
			t := transformer.NewTransformer(testConf, blockChain, db)
			err = t.Init()
			Expect(err).ToNot(HaveOccurred())

			c, ok := t.Contracts[ensAddr]
			Expect(ok).To(Equal(true))

			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(len(c.EmittedHashes)).To(Equal(3))

			b, ok := c.EmittedHashes[common.HexToHash("0x0000000000000000000000000000000000000000000000000000c02aaa39b223")]
			Expect(ok).To(Equal(true))
			Expect(b).To(Equal(true))

			b, ok = c.EmittedHashes[common.HexToHash("0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391")]
			Expect(ok).To(Equal(true))
			Expect(b).To(Equal(true))

			// Doesn't keep track of address since it wouldn't be used in calling the 'owner' method
			_, ok = c.EmittedAddrs[common.HexToAddress("0x000000000000000000000000000000000000Af21")]
			Expect(ok).To(Equal(false))
		})

		It("Polls given methods using generated token holder address", func() {
			var testConf config.ContractConfig
			testConf = test_helpers.ENSConfig
			testConf.Methods = map[string][]string{
				ensAddr: {"owner"},
			}
			t := transformer.NewTransformer(testConf, blockChain, db)
			err = t.Init()
			Expect(err).ToNot(HaveOccurred())

			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())

			res := test_helpers.Owner{}
			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM full_%s.owner_method WHERE node_ = '0x0000000000000000000000000000000000000000000000000000c02aaa39b223' AND block = '6194636'", ensAddr)).StructScan(&res)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Address).To(Equal("0x0000000000000000000000000000000000000000"))
			Expect(res.TokenName).To(Equal(""))

			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM full_%s.owner_method WHERE node_ = '0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391' AND block = '6194636'", ensAddr)).StructScan(&res)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Address).To(Equal("0x0000000000000000000000000000000000000000"))
			Expect(res.TokenName).To(Equal(""))

			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM full_%s.owner_method WHERE node_ = '0x9THIS110dcc444fIS242510c09bbAbe21aFAKEcacNODE82f7b843HASH61ba391' AND block = '6194636'", ensAddr)).StructScan(&res)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no rows in result set"))
		})

		It("It does not perist events if they do not pass the emitted arg filter", func() {
			var testConf config.ContractConfig
			testConf = test_helpers.ENSConfig
			testConf.EventArgs = map[string][]string{
				ensAddr: {"fake_filter_value"},
			}
			t := transformer.NewTransformer(testConf, blockChain, db)

			err = t.Init()
			Expect(err).ToNot(HaveOccurred())

			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())

			log := test_helpers.HeaderSyncNewOwnerLog{}
			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM full_%s.newowner_event", ensAddr)).StructScan(&log)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("does not exist"))
		})

		It("If a method arg filter is applied, only those arguments are used in polling", func() {
			var testConf config.ContractConfig
			testConf = test_helpers.ENSConfig
			testConf.MethodArgs = map[string][]string{
				ensAddr: {"0x0000000000000000000000000000000000000000000000000000c02aaa39b223"},
			}
			testConf.Methods = map[string][]string{
				ensAddr: {"owner"},
			}
			t := transformer.NewTransformer(testConf, blockChain, db)
			err = t.Init()
			Expect(err).ToNot(HaveOccurred())

			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())

			res := test_helpers.Owner{}
			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM full_%s.owner_method WHERE node_ = '0x0000000000000000000000000000000000000000000000000000c02aaa39b223' AND block = '6194636'", ensAddr)).StructScan(&res)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Address).To(Equal("0x0000000000000000000000000000000000000000"))
			Expect(res.TokenName).To(Equal(""))

			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM full_%s.owner_method WHERE node_ = '0x9dd48110dcc444fdc242510c09bbbbe21a5975cac061d82f7b843bce061ba391' AND block = '6194636'", ensAddr)).StructScan(&res)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no rows in result set"))
		})
	})
})
