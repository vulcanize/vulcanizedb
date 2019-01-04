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
	"strings"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/omni/light/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers/test_helpers/mocks"
)

var _ = Describe("Transformer", func() {
	var db *postgres.DB
	var err error
	var blockChain core.BlockChain
	var headerRepository repositories.HeaderRepository
	var headerID, headerID2 int64
	var ensAddr = strings.ToLower(constants.EnsContractAddress)
	var tusdAddr = strings.ToLower(constants.TusdContractAddress)

	BeforeEach(func() {
		db, blockChain = test_helpers.SetupDBandBC()
		headerRepository = repositories.NewHeaderRepository(db)
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
			headerRepository.CreateOrUpdateHeader(mocks.MockHeader1)
			headerRepository.CreateOrUpdateHeader(mocks.MockHeader3)
			t := transformer.NewTransformer("", blockChain, db)
			t.SetEvents(constants.TusdContractAddress, []string{"Transfer"})
			err = t.Init()
			Expect(err).ToNot(HaveOccurred())

			c, ok := t.Contracts[tusdAddr]
			Expect(ok).To(Equal(true))

			Expect(c.StartingBlock).To(Equal(int64(6194632)))
			Expect(c.LastBlock).To(Equal(int64(-1)))
			Expect(c.Abi).To(Equal(constants.TusdAbiString))
			Expect(c.Name).To(Equal("TrueUSD"))
			Expect(c.Address).To(Equal(tusdAddr))
		})

		It("Fails to initialize if first and most recent block numbers cannot be fetched from vDB headers table", func() {
			t := transformer.NewTransformer("", blockChain, db)
			t.SetEvents(constants.TusdContractAddress, []string{"Transfer"})
			err = t.Init()
			Expect(err).To(HaveOccurred())
		})

		It("Does nothing if watched events are unset", func() {
			headerRepository.CreateOrUpdateHeader(mocks.MockHeader1)
			headerRepository.CreateOrUpdateHeader(mocks.MockHeader3)
			t := transformer.NewTransformer("", blockChain, db)
			err = t.Init()
			Expect(err).ToNot(HaveOccurred())

			_, ok := t.Contracts[tusdAddr]
			Expect(ok).To(Equal(false))
		})
	})

	Describe("Execute- against TrueUSD contract", func() {
		BeforeEach(func() {
			header1, err := blockChain.GetHeaderByNumber(6791668)
			Expect(err).ToNot(HaveOccurred())
			header2, err := blockChain.GetHeaderByNumber(6791669)
			Expect(err).ToNot(HaveOccurred())
			header3, err := blockChain.GetHeaderByNumber(6791670)
			Expect(err).ToNot(HaveOccurred())
			headerRepository.CreateOrUpdateHeader(header1)
			headerID, err = headerRepository.CreateOrUpdateHeader(header2)
			Expect(err).ToNot(HaveOccurred())
			headerRepository.CreateOrUpdateHeader(header3)
		})

		It("Transforms watched contract data into custom repositories", func() {
			t := transformer.NewTransformer("", blockChain, db)
			t.SetEvents(constants.TusdContractAddress, []string{"Transfer"})
			t.SetMethods(constants.TusdContractAddress, nil)
			err = t.Init()
			Expect(err).ToNot(HaveOccurred())
			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())

			log := test_helpers.LightTransferLog{}
			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.transfer_event", tusdAddr)).StructScan(&log)
			Expect(err).ToNot(HaveOccurred())
			// We don't know vulcID, so compare individual fields instead of complete structures
			Expect(log.HeaderID).To(Equal(headerID))
			Expect(log.From).To(Equal("0x1062a747393198f70F71ec65A582423Dba7E5Ab3"))
			Expect(log.To).To(Equal("0x2930096dB16b4A44Ecd4084EA4bd26F7EeF1AEf0"))
			Expect(log.Value).To(Equal("9998940000000000000000"))
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
			Expect(len(c.EmittedAddrs)).To(Equal(4))
			Expect(len(c.EmittedHashes)).To(Equal(0))

			b, ok := c.EmittedAddrs[common.HexToAddress("0x1062a747393198f70F71ec65A582423Dba7E5Ab3")]
			Expect(ok).To(Equal(true))
			Expect(b).To(Equal(true))

			b, ok = c.EmittedAddrs[common.HexToAddress("0x2930096dB16b4A44Ecd4084EA4bd26F7EeF1AEf0")]
			Expect(ok).To(Equal(true))
			Expect(b).To(Equal(true))

			b, ok = c.EmittedAddrs[common.HexToAddress("0x571A326f5B15E16917dC17761c340c1ec5d06f6d")]
			Expect(ok).To(Equal(true))
			Expect(b).To(Equal(true))

			b, ok = c.EmittedAddrs[common.HexToAddress("0xFBb1b73C4f0BDa4f67dcA266ce6Ef42f520fBB98")]
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
			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.balanceof_method WHERE who_ = '0x1062a747393198f70F71ec65A582423Dba7E5Ab3' AND block = '6791669'", tusdAddr)).StructScan(&res)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Balance).To(Equal("55849938025000000000000"))
			Expect(res.TokenName).To(Equal("TrueUSD"))

			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.balanceof_method WHERE who_ = '0x09BbBBE21a5975cAc061D82f7b843b1234567890' AND block = '6791669'", tusdAddr)).StructScan(&res)
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
			header1, err := blockChain.GetHeaderByNumber(6885695)
			Expect(err).ToNot(HaveOccurred())
			header2, err := blockChain.GetHeaderByNumber(6885696)
			Expect(err).ToNot(HaveOccurred())
			header3, err := blockChain.GetHeaderByNumber(6885697)
			Expect(err).ToNot(HaveOccurred())
			headerRepository.CreateOrUpdateHeader(header1)
			headerID, err = headerRepository.CreateOrUpdateHeader(header2)
			Expect(err).ToNot(HaveOccurred())
			headerRepository.CreateOrUpdateHeader(header3)
		})

		It("Transforms watched contract data into custom repositories", func() {
			t := transformer.NewTransformer("", blockChain, db)
			t.SetEvents(constants.EnsContractAddress, []string{"NewOwner"})
			t.SetMethods(constants.EnsContractAddress, nil)
			err = t.Init()
			Expect(err).ToNot(HaveOccurred())
			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())

			log := test_helpers.LightNewOwnerLog{}
			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.newowner_event", ensAddr)).StructScan(&log)
			Expect(err).ToNot(HaveOccurred())
			// We don't know vulcID, so compare individual fields instead of complete structures
			Expect(log.HeaderID).To(Equal(headerID))
			Expect(log.Node).To(Equal("0x93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae"))
			Expect(log.Label).To(Equal("0x95832c7a47ff8a7840e28b78ce695797aaf402b1c186bad9eca28842625b5047"))
			Expect(log.Owner).To(Equal("0x6090A6e47849629b7245Dfa1Ca21D94cd15878Ef"))
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
			Expect(len(c.EmittedHashes)).To(Equal(2))
			Expect(len(c.EmittedAddrs)).To(Equal(0))

			b, ok := c.EmittedHashes[common.HexToHash("0x93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae")]
			Expect(ok).To(Equal(true))
			Expect(b).To(Equal(true))

			b, ok = c.EmittedHashes[common.HexToHash("0x95832c7a47ff8a7840e28b78ce695797aaf402b1c186bad9eca28842625b5047")]
			Expect(ok).To(Equal(true))
			Expect(b).To(Equal(true))

			// Doesn't keep track of address since it wouldn't be used in calling the 'owner' method
			_, ok = c.EmittedAddrs[common.HexToAddress("0x6090A6e47849629b7245Dfa1Ca21D94cd15878Ef")]
			Expect(ok).To(Equal(false))
		})

		It("Polls given method using list of collected hashes", func() {
			t := transformer.NewTransformer("", blockChain, db)
			t.SetEvents(constants.EnsContractAddress, []string{"NewOwner"})
			t.SetMethods(constants.EnsContractAddress, []string{"owner"})
			err = t.Init()
			Expect(err).ToNot(HaveOccurred())
			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())

			res := test_helpers.Owner{}
			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.owner_method WHERE node_ = '0x93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae' AND block = '6885696'", ensAddr)).StructScan(&res)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Address).To(Equal("0x6090A6e47849629b7245Dfa1Ca21D94cd15878Ef"))
			Expect(res.TokenName).To(Equal(""))

			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.owner_method WHERE node_ = '0x95832c7a47ff8a7840e28b78ce695797aaf402b1c186bad9eca28842625b5047' AND block = '6885696'", ensAddr)).StructScan(&res)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Address).To(Equal("0x0000000000000000000000000000000000000000"))
			Expect(res.TokenName).To(Equal(""))

			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.owner_method WHERE node_ = '0x95832c7a47ff8a7840e28b78ceMADEUPaaf4HASHc186badTHIS288IS625bFAKE' AND block = '6885696'", ensAddr)).StructScan(&res)
			Expect(err).To(HaveOccurred())
		})

		It("It does not persist events if they do not pass the emitted arg filter", func() {
			t := transformer.NewTransformer("", blockChain, db)
			t.SetEvents(constants.EnsContractAddress, []string{"NewOwner"})
			t.SetMethods(constants.EnsContractAddress, nil)
			t.SetEventArgs(constants.EnsContractAddress, []string{"fake_filter_value"})
			err = t.Init()
			Expect(err).ToNot(HaveOccurred())
			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())

			log := test_helpers.LightNewOwnerLog{}
			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.newowner_event", ensAddr)).StructScan(&log)
			Expect(err).To(HaveOccurred())
		})

		It("If a method arg filter is applied, only those arguments are used in polling", func() {
			t := transformer.NewTransformer("", blockChain, db)
			t.SetEvents(constants.EnsContractAddress, []string{"NewOwner"})
			t.SetMethods(constants.EnsContractAddress, []string{"owner"})
			t.SetMethodArgs(constants.EnsContractAddress, []string{"0x93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae"})
			err = t.Init()
			Expect(err).ToNot(HaveOccurred())
			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())

			res := test_helpers.Owner{}
			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.owner_method WHERE node_ = '0x93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae' AND block = '6885696'", ensAddr)).StructScan(&res)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Address).To(Equal("0x6090A6e47849629b7245Dfa1Ca21D94cd15878Ef"))
			Expect(res.TokenName).To(Equal(""))

			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.owner_method WHERE node_ = '0x95832c7a47ff8a7840e28b78ce695797aaf402b1c186bad9eca28842625b5047' AND block = '6885696'", ensAddr)).StructScan(&res)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Execute- against both ENS and TrueUSD", func() {
		BeforeEach(func() {
			header1, err := blockChain.GetHeaderByNumber(6791668)
			Expect(err).ToNot(HaveOccurred())
			header2, err := blockChain.GetHeaderByNumber(6791669)
			Expect(err).ToNot(HaveOccurred())
			header3, err := blockChain.GetHeaderByNumber(6791670)
			Expect(err).ToNot(HaveOccurred())
			header4, err := blockChain.GetHeaderByNumber(6885695)
			Expect(err).ToNot(HaveOccurred())
			header5, err := blockChain.GetHeaderByNumber(6885696)
			Expect(err).ToNot(HaveOccurred())
			header6, err := blockChain.GetHeaderByNumber(6885697)
			Expect(err).ToNot(HaveOccurred())
			headerRepository.CreateOrUpdateHeader(header1)
			headerID, err = headerRepository.CreateOrUpdateHeader(header2)
			Expect(err).ToNot(HaveOccurred())
			headerRepository.CreateOrUpdateHeader(header3)
			headerRepository.CreateOrUpdateHeader(header4)
			headerID2, err = headerRepository.CreateOrUpdateHeader(header5)
			Expect(err).ToNot(HaveOccurred())
			headerRepository.CreateOrUpdateHeader(header6)
		})

		It("Transforms watched contract data into custom repositories", func() {
			t := transformer.NewTransformer("", blockChain, db)
			t.SetEvents(constants.EnsContractAddress, []string{"NewOwner"})
			t.SetMethods(constants.EnsContractAddress, nil)
			t.SetEvents(constants.TusdContractAddress, []string{"Transfer"})
			t.SetMethods(constants.TusdContractAddress, nil)
			err = t.Init()
			Expect(err).ToNot(HaveOccurred())
			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())

			newOwnerLog := test_helpers.LightNewOwnerLog{}
			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.newowner_event", ensAddr)).StructScan(&newOwnerLog)
			Expect(err).ToNot(HaveOccurred())
			// We don't know vulcID, so compare individual fields instead of complete structures
			Expect(newOwnerLog.HeaderID).To(Equal(headerID2))
			Expect(newOwnerLog.Node).To(Equal("0x93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae"))
			Expect(newOwnerLog.Label).To(Equal("0x95832c7a47ff8a7840e28b78ce695797aaf402b1c186bad9eca28842625b5047"))
			Expect(newOwnerLog.Owner).To(Equal("0x6090A6e47849629b7245Dfa1Ca21D94cd15878Ef"))

			transferLog := test_helpers.LightTransferLog{}
			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.transfer_event", tusdAddr)).StructScan(&transferLog)
			Expect(err).ToNot(HaveOccurred())
			// We don't know vulcID, so compare individual fields instead of complete structures
			Expect(transferLog.HeaderID).To(Equal(headerID))
			Expect(transferLog.From).To(Equal("0x1062a747393198f70F71ec65A582423Dba7E5Ab3"))
			Expect(transferLog.To).To(Equal("0x2930096dB16b4A44Ecd4084EA4bd26F7EeF1AEf0"))
			Expect(transferLog.Value).To(Equal("9998940000000000000000"))
		})

		It("Keeps track of contract-related hashes and addresses while transforming event data if they need to be used for later method polling", func() {
			t := transformer.NewTransformer("", blockChain, db)
			t.SetEvents(constants.EnsContractAddress, []string{"NewOwner"})
			t.SetMethods(constants.EnsContractAddress, []string{"owner"})
			t.SetEvents(constants.TusdContractAddress, []string{"Transfer"})
			t.SetMethods(constants.TusdContractAddress, []string{"balanceOf"})
			err = t.Init()
			Expect(err).ToNot(HaveOccurred())
			ens, ok := t.Contracts[ensAddr]
			Expect(ok).To(Equal(true))
			tusd, ok := t.Contracts[tusdAddr]
			Expect(ok).To(Equal(true))
			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(len(ens.EmittedHashes)).To(Equal(2))
			Expect(len(ens.EmittedAddrs)).To(Equal(0))
			Expect(len(tusd.EmittedAddrs)).To(Equal(4))
			Expect(len(tusd.EmittedHashes)).To(Equal(0))

			b, ok := ens.EmittedHashes[common.HexToHash("0x93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae")]
			Expect(ok).To(Equal(true))
			Expect(b).To(Equal(true))

			b, ok = ens.EmittedHashes[common.HexToHash("0x95832c7a47ff8a7840e28b78ce695797aaf402b1c186bad9eca28842625b5047")]
			Expect(ok).To(Equal(true))
			Expect(b).To(Equal(true))

			b, ok = tusd.EmittedAddrs[common.HexToAddress("0x1062a747393198f70F71ec65A582423Dba7E5Ab3")]
			Expect(ok).To(Equal(true))
			Expect(b).To(Equal(true))

			b, ok = tusd.EmittedAddrs[common.HexToAddress("0x2930096dB16b4A44Ecd4084EA4bd26F7EeF1AEf0")]
			Expect(ok).To(Equal(true))
			Expect(b).To(Equal(true))

			b, ok = tusd.EmittedAddrs[common.HexToAddress("0x571A326f5B15E16917dC17761c340c1ec5d06f6d")]
			Expect(ok).To(Equal(true))
			Expect(b).To(Equal(true))

			b, ok = tusd.EmittedAddrs[common.HexToAddress("0xFBb1b73C4f0BDa4f67dcA266ce6Ef42f520fBB98")]
			Expect(ok).To(Equal(true))
			Expect(b).To(Equal(true))
		})

		It("Polls given methods for each contract, using list of collected values", func() {
			t := transformer.NewTransformer("", blockChain, db)
			t.SetEvents(constants.EnsContractAddress, []string{"NewOwner"})
			t.SetMethods(constants.EnsContractAddress, []string{"owner"})
			t.SetEvents(constants.TusdContractAddress, []string{"Transfer"})
			t.SetMethods(constants.TusdContractAddress, []string{"balanceOf"})
			err = t.Init()
			Expect(err).ToNot(HaveOccurred())
			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())

			owner := test_helpers.Owner{}
			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.owner_method WHERE node_ = '0x93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae' AND block = '6885696'", ensAddr)).StructScan(&owner)
			Expect(err).ToNot(HaveOccurred())
			Expect(owner.Address).To(Equal("0x6090A6e47849629b7245Dfa1Ca21D94cd15878Ef"))
			Expect(owner.TokenName).To(Equal(""))

			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.owner_method WHERE node_ = '0x95832c7a47ff8a7840e28b78ce695797aaf402b1c186bad9eca28842625b5047' AND block = '6885696'", ensAddr)).StructScan(&owner)
			Expect(err).ToNot(HaveOccurred())
			Expect(owner.Address).To(Equal("0x0000000000000000000000000000000000000000"))
			Expect(owner.TokenName).To(Equal(""))

			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.owner_method WHERE node_ = '0x95832c7a47ff8a7840e28b78ceMADEUPaaf4HASHc186badTHIS288IS625bFAKE' AND block = '6885696'", ensAddr)).StructScan(&owner)
			Expect(err).To(HaveOccurred())

			bal := test_helpers.BalanceOf{}
			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.balanceof_method WHERE who_ = '0x1062a747393198f70F71ec65A582423Dba7E5Ab3' AND block = '6791669'", tusdAddr)).StructScan(&bal)
			Expect(err).ToNot(HaveOccurred())
			Expect(bal.Balance).To(Equal("55849938025000000000000"))
			Expect(bal.TokenName).To(Equal("TrueUSD"))

			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM light_%s.balanceof_method WHERE who_ = '0x09BbBBE21a5975cAc061D82f7b843b1234567890' AND block = '6791669'", tusdAddr)).StructScan(&bal)
			Expect(err).To(HaveOccurred())
		})
	})
})
