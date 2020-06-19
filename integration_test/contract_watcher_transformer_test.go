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
	"strings"

	"github.com/makerdao/vulcanizedb/pkg/config"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/constants"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/helpers/test_helpers"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/helpers/test_helpers/mocks"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/transformer"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("contractWatcher transformer", func() {
	var (
		db               *postgres.DB
		err              error
		blockChain       core.BlockChain
		headerRepository datastore.HeaderRepository
		headerID         int64
		ensAddr          = strings.ToLower(constants.EnsContractAddress)  // 0x314159265dd8dbb310642f98f50c066173c1259b
		tusdAddr         = strings.ToLower(constants.TusdContractAddress) // 0x8dd5fbce2f6a956c3022ba3663759011dd51e73e
	)

	BeforeEach(func() {
		db, blockChain = test_helpers.SetupDBandBC()
		headerRepository = repositories.NewHeaderRepository(db)
	})

	AfterEach(func() {
		test_helpers.TearDown(db)
	})

	Describe("Init", func() {
		It("Initializes transformer's contract objects", func() {
			_, insertErr := headerRepository.CreateOrUpdateHeader(mocks.MockHeader1)
			Expect(insertErr).NotTo(HaveOccurred())
			_, insertErrTwo := headerRepository.CreateOrUpdateHeader(mocks.MockHeader3)
			Expect(insertErrTwo).NotTo(HaveOccurred())
			t := transformer.NewTransformer(test_helpers.TusdConfig, blockChain, db)
			err = t.Init("")
			Expect(err).ToNot(HaveOccurred())

			c, ok := t.Contracts[tusdAddr]
			Expect(ok).To(Equal(true))

			// TODO: Fix this
			// This test sometimes randomly fails because
			// for some reason the starting block number is not updated from
			// its original value (5197514) to the block number (6194632)
			// of the earliest header (mocks.MockHeader1) in the repository
			// It is not clear how this happens without one of the above insertErrs
			// having been thrown and without any errors thrown during the Init() call
			Expect(c.StartingBlock).To(Equal(int64(6194632)))
			Expect(c.Abi).To(Equal(constants.TusdAbiString))
			Expect(c.Address).To(Equal(tusdAddr))
		})

		It("initializes when no headers available in db", func() {
			t := transformer.NewTransformer(test_helpers.TusdConfig, blockChain, db)
			err = t.Init("")
			Expect(err).ToNot(HaveOccurred())
		})

		It("Does nothing if nothing if no addresses are configured", func() {
			_, insertErr := headerRepository.CreateOrUpdateHeader(mocks.MockHeader1)
			Expect(insertErr).NotTo(HaveOccurred())
			_, insertErrTwo := headerRepository.CreateOrUpdateHeader(mocks.MockHeader3)
			Expect(insertErrTwo).NotTo(HaveOccurred())
			var testConf config.ContractConfig
			testConf = test_helpers.TusdConfig
			testConf.Addresses = nil
			t := transformer.NewTransformer(testConf, blockChain, db)
			err = t.Init("")
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
			_, err = headerRepository.CreateOrUpdateHeader(header1)
			Expect(err).NotTo(HaveOccurred())
			headerID, err = headerRepository.CreateOrUpdateHeader(header2)
			Expect(err).ToNot(HaveOccurred())
			_, err = headerRepository.CreateOrUpdateHeader(header3)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Transforms watched contract data into custom repositories", func() {
			t := transformer.NewTransformer(test_helpers.TusdConfig, blockChain, db)
			err = t.Init("")
			Expect(err).ToNot(HaveOccurred())
			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())

			log := test_helpers.TransferLog{}
			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM cw_%s.transfer_event", tusdAddr)).StructScan(&log)
			Expect(err).ToNot(HaveOccurred())
			// We don't know vulcID, so compare individual fields instead of complete structures
			Expect(log.HeaderID).To(Equal(headerID))
			Expect(log.From).To(Equal("0x1062a747393198f70F71ec65A582423Dba7E5Ab3"))
			Expect(log.To).To(Equal("0x2930096dB16b4A44Ecd4084EA4bd26F7EeF1AEf0"))
			Expect(log.Value).To(Equal("9998940000000000000000"))
		})

		It("Fails if initialization has not been done", func() {
			t := transformer.NewTransformer(test_helpers.TusdConfig, blockChain, db)
			err = t.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("transformer has no initialized contracts"))
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
			_, err = headerRepository.CreateOrUpdateHeader(header1)
			Expect(err).NotTo(HaveOccurred())
			headerID, err = headerRepository.CreateOrUpdateHeader(header2)
			Expect(err).ToNot(HaveOccurred())
			_, err = headerRepository.CreateOrUpdateHeader(header3)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Transforms watched contract data into custom repositories", func() {
			t := transformer.NewTransformer(test_helpers.ENSConfig, blockChain, db)
			err = t.Init("")
			Expect(err).ToNot(HaveOccurred())
			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(t.Start).To(Equal(int64(6885698)))

			log := test_helpers.NewOwnerLog{}
			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM cw_%s.newowner_event", ensAddr)).StructScan(&log)
			Expect(err).ToNot(HaveOccurred())
			// We don't know vulcID, so compare individual fields instead of complete structures
			Expect(log.HeaderID).To(Equal(headerID))
			Expect(log.Node).To(Equal("0x93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae"))
			Expect(log.Label).To(Equal("0x95832c7a47ff8a7840e28b78ce695797aaf402b1c186bad9eca28842625b5047"))
			Expect(log.Owner).To(Equal("0x6090A6e47849629b7245Dfa1Ca21D94cd15878Ef"))
		})

		It("It does not persist events if they do not pass the emitted arg filter", func() {
			var testConf config.ContractConfig
			testConf = test_helpers.ENSConfig
			testConf.EventArgs = map[string][]string{
				ensAddr: {"fake_filter_value"},
			}
			t := transformer.NewTransformer(testConf, blockChain, db)
			err = t.Init("")
			Expect(err).ToNot(HaveOccurred())
			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())

			log := test_helpers.NewOwnerLog{}
			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM cw_%s.newowner_event", ensAddr)).StructScan(&log)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("does not exist"))
		})
	})

	Describe("Execute- against both ENS and TrueUSD", func() {
		BeforeEach(func() {
			for i := 6885692; i <= 6885701; i++ {
				header, err := blockChain.GetHeaderByNumber(int64(i))
				Expect(err).ToNot(HaveOccurred())
				_, err = headerRepository.CreateOrUpdateHeader(header)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("Transforms watched contract data into custom repositories", func() {
			t := transformer.NewTransformer(test_helpers.ENSandTusdConfig, blockChain, db)
			err = t.Init("")
			Expect(err).ToNot(HaveOccurred())
			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(t.Start).To(Equal(int64(6885702)))

			newOwnerLog := test_helpers.NewOwnerLog{}
			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM cw_%s.newowner_event", ensAddr)).StructScan(&newOwnerLog)
			Expect(err).ToNot(HaveOccurred())
			// We don't know vulcID, so compare individual fields instead of complete structures
			Expect(newOwnerLog.Node).To(Equal("0x93cdeb708b7545dc668eb9280176169d1c33cfd8ed6f04690a0bcc88a93fc4ae"))
			Expect(newOwnerLog.Label).To(Equal("0x95832c7a47ff8a7840e28b78ce695797aaf402b1c186bad9eca28842625b5047"))
			Expect(newOwnerLog.Owner).To(Equal("0x6090A6e47849629b7245Dfa1Ca21D94cd15878Ef"))

			transferLog := test_helpers.TransferLog{}
			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM cw_%s.transfer_event", tusdAddr)).StructScan(&transferLog)
			Expect(err).ToNot(HaveOccurred())
			// We don't know vulcID, so compare individual fields instead of complete structures
			Expect(transferLog.From).To(Equal("0x8cA465764873E71CEa525F5EB6AE973d650c22C2"))
			Expect(transferLog.To).To(Equal("0xc338482360651E5D30BEd77b7c85358cbBFB2E0e"))
			Expect(transferLog.Value).To(Equal("2800000000000000000000"))
		})

		It("Marks header checked for a contract that has no logs at that header", func() {
			t := transformer.NewTransformer(test_helpers.ENSandTusdConfig, blockChain, db)
			err = t.Init("")
			Expect(err).ToNot(HaveOccurred())
			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(t.Start).To(Equal(int64(6885702)))

			newOwnerLog := test_helpers.NewOwnerLog{}
			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM cw_%s.newowner_event", ensAddr)).StructScan(&newOwnerLog)
			Expect(err).ToNot(HaveOccurred())
			transferLog := test_helpers.TransferLog{}
			err = db.QueryRowx(fmt.Sprintf("SELECT * FROM cw_%s.transfer_event", tusdAddr)).StructScan(&transferLog)
			Expect(err).ToNot(HaveOccurred())
			Expect(transferLog.HeaderID).ToNot(Equal(newOwnerLog.HeaderID))

			type checkedHeader struct {
				ID       int64 `db:"id"`
				HeaderID int64 `db:"header_id"`
				NewOwner int64 `db:"newowner_0x314159265dd8dbb310642f98f50c066173c1259b"`
				Transfer int64 `db:"transfer_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e"`
			}

			transferCheckedHeader := new(checkedHeader)
			err = db.QueryRowx("SELECT * FROM public.checked_headers WHERE header_id = $1", transferLog.HeaderID).StructScan(transferCheckedHeader)
			Expect(err).ToNot(HaveOccurred())
			Expect(transferCheckedHeader.Transfer).To(Equal(int64(1)))
			Expect(transferCheckedHeader.NewOwner).To(Equal(int64(1)))

			newOwnerCheckedHeader := new(checkedHeader)
			err = db.QueryRowx("SELECT * FROM public.checked_headers WHERE header_id = $1", newOwnerLog.HeaderID).StructScan(newOwnerCheckedHeader)
			Expect(err).ToNot(HaveOccurred())
			Expect(newOwnerCheckedHeader.NewOwner).To(Equal(int64(1)))
			Expect(newOwnerCheckedHeader.Transfer).To(Equal(int64(1)))
		})
	})
})
