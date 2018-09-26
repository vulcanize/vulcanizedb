// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package flop_kick_test

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flop_kick"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

var _ = Describe("Integration tests", func() {
	It("unpacks an flop kick event log", func() {
		address := common.HexToAddress(shared.FlopperContractAddress)
		abi, err := geth.ParseAbi(shared.FlopperABI)
		Expect(err).NotTo(HaveOccurred())

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)
		entity := &flop_kick.Entity{}

		var eventLog = test_data.FlopKickLog

		err = contract.UnpackLog(entity, "Kick", eventLog)
		Expect(err).NotTo(HaveOccurred())

		expectedEntity := test_data.FlopKickEntity
		Expect(entity.Id).To(Equal(expectedEntity.Id))
		Expect(entity.Lot).To(Equal(expectedEntity.Lot))
		Expect(entity.Bid).To(Equal(expectedEntity.Bid))
		Expect(entity.Gal).To(Equal(expectedEntity.Gal))
		Expect(entity.End).To(Equal(expectedEntity.End))
	})
})
