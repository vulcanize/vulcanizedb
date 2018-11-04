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

package retriever

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/omni/contract"
)

// Address retriever is used to retrieve the addresses associated with a contract
// It requires a vDB synced database with blocks, transactions, receipts, logs,
// AND all of the targeted events persisted
type AddressRetriever interface {
	RetrieveTokenHolderAddresses(info contract.Contract) (map[common.Address]bool, error)
}

type addressRetriever struct {
	*postgres.DB
}

func NewAddressRetriever(db *postgres.DB) (r *addressRetriever) {

	return &addressRetriever{
		DB: db,
	}
}

// Method to retrieve list of token-holding/contract-related addresses by iterating over available events
// This generic method should work whether or not the argument/input names of the events meet the expected standard
// This could be generalized to iterate over ALL events and pull out any address arguments
func (r *addressRetriever) RetrieveTokenHolderAddresses(info contract.Contract) (map[common.Address]bool, error) {
	addrList := make([]string, 0)

	_, ok := info.Filters["Transfer"]
	if ok {
		addrs, err := r.retrieveTransferAddresses(info)
		if err != nil {
			return nil, err
		}
		addrList = append(addrList, addrs...)
	}

	_, ok = info.Filters["Mint"]
	if ok {
		addrs, err := r.retrieveTokenMintees(info)
		if err != nil {
			return nil, err
		}
		addrList = append(addrList, addrs...)
	}

	contractAddresses := make(map[common.Address]bool)
	for _, addr := range addrList {
		contractAddresses[common.HexToAddress(addr)] = true
	}

	return contractAddresses, nil
}

func (r *addressRetriever) retrieveTransferAddresses(contract contract.Contract) ([]string, error) {
	transferAddrs := make([]string, 0)
	event := contract.Events["Transfer"]

	for _, field := range event.Fields { // Iterate over event fields, finding the ones with address type

		if field.Type.T == abi.AddressTy { // If they have address type, retrieve those addresses
			addrs := make([]string, 0)
			pgStr := fmt.Sprintf("SELECT _%s FROM %s.%s", field.Name, contract.Name, event.Name)
			err := r.DB.Select(&addrs, pgStr)
			if err != nil {
				return []string{}, err
			}

			transferAddrs = append(transferAddrs, addrs...) // And append them to the growing list
		}
	}

	return transferAddrs, nil
}

func (r *addressRetriever) retrieveTokenMintees(contract contract.Contract) ([]string, error) {
	mintAddrs := make([]string, 0)
	event := contract.Events["Mint"]

	for _, field := range event.Fields { // Iterate over event fields, finding the ones with address type

		if field.Type.T == abi.AddressTy { // If they have address type, retrieve those addresses
			addrs := make([]string, 0)
			pgStr := fmt.Sprintf("SELECT _%s FROM %s.%s", field.Name, contract.Name, event.Name)
			err := r.DB.Select(&addrs, pgStr)
			if err != nil {
				return []string{}, err
			}

			mintAddrs = append(mintAddrs, addrs...) // And append them to the growing list
		}
	}

	return mintAddrs, nil
}
