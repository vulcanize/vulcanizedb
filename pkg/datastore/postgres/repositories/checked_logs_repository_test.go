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

package repositories_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	"github.com/makerdao/vulcanizedb/test_config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Checked logs repository", func() {
	var (
		db            *postgres.DB
		fakeAddress   = fakes.FakeAddress.Hex()
		fakeAddresses = []string{fakeAddress}
		fakeTopicZero = fakes.FakeHash.Hex()
		repository    datastore.CheckedLogsRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repository = repositories.NewCheckedLogsRepository(db)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("AlreadyWatchingLog", func() {
		It("returns true if all addresses and the topic0 are already present in the db", func() {
			_, insertErr := db.Exec(`INSERT INTO public.watched_logs (contract_address, topic_zero) VALUES ($1, $2)`, fakeAddress, fakeTopicZero)
			Expect(insertErr).NotTo(HaveOccurred())

			hasBeenChecked, err := repository.AlreadyWatchingLog(fakeAddresses, fakeTopicZero)

			Expect(err).NotTo(HaveOccurred())
			Expect(hasBeenChecked).To(BeTrue())
		})

		It("returns true if addresses and topic0 were fetched because of a combination of other transformers", func() {
			anotherFakeAddress := common.HexToAddress("0x" + fakes.RandomString(40)).Hex()
			anotherFakeTopicZero := common.HexToHash("0x" + fakes.RandomString(64)).Hex()
			// insert row with matching address but different topic0
			_, insertOneErr := db.Exec(`INSERT INTO public.watched_logs (contract_address, topic_zero) VALUES ($1, $2)`, fakeAddress, anotherFakeTopicZero)
			Expect(insertOneErr).NotTo(HaveOccurred())
			// insert row with matching topic0 but different address
			_, insertTwoErr := db.Exec(`INSERT INTO public.watched_logs (contract_address, topic_zero) VALUES ($1, $2)`, anotherFakeAddress, fakeTopicZero)
			Expect(insertTwoErr).NotTo(HaveOccurred())

			hasBeenChecked, err := repository.AlreadyWatchingLog(fakeAddresses, fakeTopicZero)

			Expect(err).NotTo(HaveOccurred())
			Expect(hasBeenChecked).To(BeTrue())
		})

		It("returns false if any address has not been checked", func() {
			anotherFakeAddress := common.HexToAddress("0x" + fakes.RandomString(40)).Hex()
			_, insertErr := db.Exec(`INSERT INTO public.watched_logs (contract_address, topic_zero) VALUES ($1, $2)`, fakeAddress, fakeTopicZero)
			Expect(insertErr).NotTo(HaveOccurred())

			hasBeenChecked, err := repository.AlreadyWatchingLog(append(fakeAddresses, anotherFakeAddress), fakeTopicZero)

			Expect(err).NotTo(HaveOccurred())
			Expect(hasBeenChecked).To(BeFalse())
		})

		It("returns false if topic0 has not been checked", func() {
			anotherFakeTopicZero := common.HexToHash("0x" + fakes.RandomString(64)).Hex()
			_, insertErr := db.Exec(`INSERT INTO public.watched_logs (contract_address, topic_zero) VALUES ($1, $2)`, fakeAddress, anotherFakeTopicZero)
			Expect(insertErr).NotTo(HaveOccurred())

			hasBeenChecked, err := repository.AlreadyWatchingLog(fakeAddresses, fakeTopicZero)

			Expect(err).NotTo(HaveOccurred())
			Expect(hasBeenChecked).To(BeFalse())
		})
	})

	Describe("MarkLogWatched", func() {
		It("adds a row for all of transformer's addresses + topic0", func() {
			anotherFakeAddress := common.HexToAddress("0x" + fakes.RandomString(40)).Hex()
			err := repository.MarkLogWatched(append(fakeAddresses, anotherFakeAddress), fakeTopicZero)

			Expect(err).NotTo(HaveOccurred())
			var comboOneExists, comboTwoExists bool
			getComboOneErr := db.Get(&comboOneExists, `SELECT EXISTS(SELECT 1 FROM public.watched_logs WHERE contract_address = $1 AND topic_zero = $2)`, fakeAddress, fakeTopicZero)
			Expect(getComboOneErr).NotTo(HaveOccurred())
			Expect(comboOneExists).To(BeTrue())
			getComboTwoErr := db.Get(&comboTwoExists, `SELECT EXISTS(SELECT 1 FROM public.watched_logs WHERE contract_address = $1 AND topic_zero = $2)`, anotherFakeAddress, fakeTopicZero)
			Expect(getComboTwoErr).NotTo(HaveOccurred())
			Expect(comboTwoExists).To(BeTrue())
		})
	})
})
