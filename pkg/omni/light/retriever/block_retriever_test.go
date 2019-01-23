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

package retriever_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/omni/light/retriever"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/helpers/test_helpers/mocks"
)

var _ = Describe("Block Retriever", func() {
	var db *postgres.DB
	var r retriever.BlockRetriever
	var headerRepository repositories.HeaderRepository

	BeforeEach(func() {
		db, _ = test_helpers.SetupDBandBC()
		headerRepository = repositories.NewHeaderRepository(db)
		r = retriever.NewBlockRetriever(db)
	})

	AfterEach(func() {
		test_helpers.TearDown(db)
	})

	Describe("RetrieveFirstBlock", func() {
		It("Retrieves block number of earliest header in the database", func() {
			headerRepository.CreateOrUpdateHeader(mocks.MockHeader1)
			headerRepository.CreateOrUpdateHeader(mocks.MockHeader2)
			headerRepository.CreateOrUpdateHeader(mocks.MockHeader3)

			i, err := r.RetrieveFirstBlock()
			Expect(err).NotTo(HaveOccurred())
			Expect(i).To(Equal(int64(6194632)))
		})

		It("Fails if no headers can be found in the database", func() {
			_, err := r.RetrieveFirstBlock()
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("RetrieveMostRecentBlock", func() {
		It("Retrieves the latest header's block number", func() {
			headerRepository.CreateOrUpdateHeader(mocks.MockHeader1)
			headerRepository.CreateOrUpdateHeader(mocks.MockHeader2)
			headerRepository.CreateOrUpdateHeader(mocks.MockHeader3)

			i, err := r.RetrieveMostRecentBlock()
			Expect(err).ToNot(HaveOccurred())
			Expect(i).To(Equal(int64(6194634)))
		})

		It("Fails if no headers can be found in the database", func() {
			i, err := r.RetrieveMostRecentBlock()
			Expect(err).To(HaveOccurred())
			Expect(i).To(Equal(int64(0)))
		})
	})
})
