// VulcanizeDB
// Copyright © 2018 Vulcanize

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

package transformer_test

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/omni/full/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers/test_helpers/mocks"
)

var _ = Describe("Transformer", func() {
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

	Describe("SetEvents", func() {
		It("Sets which events to watch from the given contract address", func() {
			watchedEvents := []string{"Transfer", "Mint"}
			t := transformer.NewTransformer("", blockChain, db)
			t.SetEvents(constants.TusdContractAddress, watchedEvents)
			Expect(t.WatchedEvents[tusdAddr]).To(Equal(watchedEvents))
		})
	})

	Describe("SetEventAddrs", func() {
		It("Sets which account addresses to watch events for", func() {
			eventAddrs := []string{"test1", "test2"}
			t := transformer.NewTransformer("", blockChain, db)
			t.SetEventArgs(constants.TusdContractAddress, eventAddrs)
			Expect(t.EventArgs[tusdAddr]).To(Equal(eventAddrs))
		})
	})

	Describe("SetMethods", func() {
		It("Sets which methods to poll at the given contract address", func() {
			watchedMethods := []string{"balanceOf", "totalSupply"}
			t := transformer.NewTransformer("", blockChain, db)
			t.SetMethods(constants.TusdContractAddress, watchedMethods)
			Expect(t.WantedMethods[tusdAddr]).To(Equal(watchedMethods))
		})
	})

	Describe("SetMethodAddrs", func() {
		It("Sets which account addresses to poll methods against", func() {
			methodAddrs := []string{"test1", "test2"}
			t := transformer.NewTransformer("", blockChain, db)
			t.SetMethodArgs(constants.TusdContractAddress, methodAddrs)
			Expect(t.MethodArgs[tusdAddr]).To(Equal(methodAddrs))
		})
	})

	Describe("SetStartingBlock", func() {
		It("Sets the block range that the contract should be watched within", func() {
			t := transformer.NewTransformer("", blockChain, db)
			t.SetStartingBlock(constants.TusdContractAddress, 11)
			Expect(t.ContractStart[tusdAddr]).To(Equal(int64(11)))
		})
	})

	Describe("SetCreateAddrList", func() {
		It("Sets the block range that the contract should be watched within", func() {
			t := transformer.NewTransformer("", blockChain, db)
			t.SetCreateAddrList(constants.TusdContractAddress, true)
			Expect(t.CreateAddrList[tusdAddr]).To(Equal(true))
		})
	})

	Describe("SetCreateHashList", func() {
		It("Sets the block range that the contract should be watched within", func() {
			t := transformer.NewTransformer("", blockChain, db)
			t.SetCreateHashList(constants.TusdContractAddress, true)
			Expect(t.CreateHashList[tusdAddr]).To(Equal(true))
		})
	})

	Describe("Init", func() {
		It("Initializes transformer's contract objects", func() {
			blockRepository.CreateOrUpdateBlock(mocks.TransferBlock1)
			blockRepository.CreateOrUpdateBlock(mocks.TransferBlock2)
			t := transformer.NewTransformer("", blockChain, db)
			t.SetEvents(constants.TusdContractAddress, []string{"Transfer"})
			err = t.Init()
			Expect(err).ToNot(HaveOccurred())

			c, ok := t.Contracts[tusdAddr]
			Expect(ok).To(Equal(true))

			Expect(c.StartingBlock).To(Equal(int64(6194633)))
			Expect(c.LastBlock).To(Equal(int64(6194634)))
			Expect(c.Abi).To(Equal(constants.TusdAbiString))
			Expect(c.Name).To(Equal("TrueUSD"))
			Expect(c.Address).To(Equal(tusdAddr))
		})

		It("Fails to initialize if first and most recent blocks cannot be fetched from vDB", func() {
			t := transformer.NewTransformer("", blockChain, db)
			t.SetEvents(constants.TusdContractAddress, []string{"Transfer"})
			err = t.Init()
			Expect(err).To(HaveOccurred())
		})

		It("Does nothing if watched events are unset", func() {
			blockRepository.CreateOrUpdateBlock(mocks.TransferBlock1)
			blockRepository.CreateOrUpdateBlock(mocks.TransferBlock2)
			t := transformer.NewTransformer("", blockChain, db)
			err = t.Init()
			Expect(err).ToNot(HaveOccurred())

			_, ok := t.Contracts[tusdAddr]
			Expect(ok).To(Equal(false))
		})
	})

	Describe("Execute", func() {
		BeforeEach(func() {
			blockRepository.CreateOrUpdateBlock(mocks.TransferBlock1)
			blockRepository.CreateOrUpdateBlock(mocks.TransferBlock2)
		})

		It("Transforms watched contract data into custom repositories", func() {
			t := transformer.NewTransformer("", blockChain, db)
			t.SetEvents(constants.TusdContractAddress, []string{"Transfer"})
			t.SetMethods(constants.TusdContractAddress, nil)
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
			t := transformer.NewTransformer("", blockChain, db)
			t.SetEvents(constants.TusdContractAddress, []string{"Transfer"})
			t.SetMethods(constants.TusdContractAddress, []string{"balanceOf"})
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
			t := transformer.NewTransformer("", blockChain, db)
			t.SetEvents(constants.TusdContractAddress, []string{"Transfer"})
			t.SetMethods(constants.TusdContractAddress, []string{"balanceOf"})
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
		})

		It("Fails if initialization has not been done", func() {
			t := transformer.NewTransformer("", blockChain, db)
			t.SetEvents(constants.TusdContractAddress, []string{"Transfer"})
			t.SetMethods(constants.TusdContractAddress, nil)

			err = t.Execute()
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Execute- against ENS registry contract", func() {
		BeforeEach(func() {
			blockRepository.CreateOrUpdateBlock(mocks.NewOwnerBlock1)
			blockRepository.CreateOrUpdateBlock(mocks.NewOwnerBlock2)
		})

		It("Transforms watched contract data into custom repositories", func() {
			t := transformer.NewTransformer("", blockChain, db)
			t.SetEvents(constants.EnsContractAddress, []string{"NewOwner"})
			t.SetMethods(constants.EnsContractAddress, nil)

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
			t := transformer.NewTransformer("", blockChain, db)
			t.SetEvents(constants.EnsContractAddress, []string{"NewOwner"})
			t.SetMethods(constants.EnsContractAddress, []string{"owner"})
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
			t := transformer.NewTransformer("", blockChain, db)
			t.SetEvents(constants.EnsContractAddress, []string{"NewOwner"})
			t.SetMethods(constants.EnsContractAddress, []string{"owner"})
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

			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM full_%s.owner_method WHERE node_ = '0x95832c7a47ff8a7840e28b78ceMADEUPaaf4HASHc186badTHIS288IS625bFAKE' AND block = '6194636'", ensAddr)).StructScan(&res)
			Expect(err).To(HaveOccurred())
		})

		It("It does not perist events if they do not pass the emitted arg filter", func() {
			t := transformer.NewTransformer("", blockChain, db)
			t.SetEvents(constants.EnsContractAddress, []string{"NewOwner"})
			t.SetMethods(constants.EnsContractAddress, nil)
			t.SetEventArgs(constants.EnsContractAddress, []string{"fake_filter_value"})

			err = t.Init()
			Expect(err).ToNot(HaveOccurred())

			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())

			log := test_helpers.LightNewOwnerLog{}
			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM full_%s.newowner_event", ensAddr)).StructScan(&log)
			Expect(err).To(HaveOccurred())
		})

		It("If a method arg filter is applied, only those arguments are used in polling", func() {
			t := transformer.NewTransformer("", blockChain, db)
			t.SetEvents(constants.EnsContractAddress, []string{"NewOwner"})
			t.SetMethods(constants.EnsContractAddress, []string{"owner"})
			t.SetMethodArgs(constants.EnsContractAddress, []string{"0x0000000000000000000000000000000000000000000000000000c02aaa39b223"})
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
		})
	})
})
