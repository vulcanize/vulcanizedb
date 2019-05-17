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

package cold_import_test

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/geth/cold_import"
)

var _ = Describe("Cold importer node builder", func() {
	Describe("when level path is not valid", func() {
		It("returns error if no chaindata extension", func() {
			gethPath := "path/to/geth"
			mockReader := fakes.NewMockFsReader()
			mockParser := fakes.NewMockCryptoParser()
			nodeBuilder := cold_import.NewColdImportNodeBuilder(mockReader, mockParser)

			_, err := nodeBuilder.GetNode([]byte{1, 2, 3, 4, 5}, gethPath)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(cold_import.NoChainDataErr))
		})

		It("returns error if no root geth path", func() {
			chaindataPath := "chaindata"
			mockReader := fakes.NewMockFsReader()
			mockParser := fakes.NewMockCryptoParser()
			nodeBuilder := cold_import.NewColdImportNodeBuilder(mockReader, mockParser)

			_, err := nodeBuilder.GetNode([]byte{1, 2, 3, 4, 5}, chaindataPath)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(cold_import.NoGethRootErr))
		})
	})

	Describe("when reader fails", func() {
		It("returns err", func() {
			mockReader := fakes.NewMockFsReader()
			fakeError := errors.New("Failed")
			mockReader.SetReturnErr(fakeError)
			mockParser := fakes.NewMockCryptoParser()
			nodeBuilder := cold_import.NewColdImportNodeBuilder(mockReader, mockParser)

			_, err := nodeBuilder.GetNode([]byte{1, 2, 3, 4, 5}, "path/to/geth/chaindata")

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakeError))
		})
	})

	Describe("when parser fails", func() {
		It("returns err", func() {
			mockReader := fakes.NewMockFsReader()
			mockParser := fakes.NewMockCryptoParser()
			fakeErr := errors.New("Failed")
			mockParser.SetReturnErr(fakeErr)
			nodeBuilder := cold_import.NewColdImportNodeBuilder(mockReader, mockParser)

			_, err := nodeBuilder.GetNode([]byte{1, 2, 3, 4, 5}, "path/to/geth/chaindata")

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakeErr))
		})
	})

	Describe("when path is valid and reader and parser succeed", func() {
		It("builds a node", func() {
			fakeGenesisBlock := []byte{1, 2, 3, 4, 5}
			fakeRootGethPath := "root/path/to/geth/"
			fakeLevelPath := fakeRootGethPath + "chaindata"
			fakeNodeKeyPath := fakeRootGethPath + "nodekey"
			fakePublicKeyBytes := []byte{5, 4, 3, 2, 1}
			fakePublicKeyString := "public_key"
			mockReader := fakes.NewMockFsReader()
			mockReader.SetReturnBytes(fakePublicKeyBytes)
			mockParser := fakes.NewMockCryptoParser()
			mockParser.SetReturnVal(fakePublicKeyString)
			nodeBuilder := cold_import.NewColdImportNodeBuilder(mockReader, mockParser)

			result, err := nodeBuilder.GetNode(fakeGenesisBlock, fakeLevelPath)

			Expect(err).NotTo(HaveOccurred())
			mockReader.AssertReadCalledWith(fakeNodeKeyPath)
			mockParser.AssertParsePublicKeyCalledWith(string(fakePublicKeyBytes))
			Expect(result).NotTo(BeNil())
			Expect(result.ClientName).To(Equal(cold_import.ColdImportClientName))
			expectedGenesisBlock := common.BytesToHash(fakeGenesisBlock).String()
			Expect(result.GenesisBlock).To(Equal(expectedGenesisBlock))
			Expect(result.ID).To(Equal(fakePublicKeyString))
			Expect(result.NetworkID).To(Equal(cold_import.ColdImportNetworkId))
		})
	})

})
