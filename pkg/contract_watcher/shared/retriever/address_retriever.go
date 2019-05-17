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

package retriever

import (
	"fmt"
	"strings"

	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/contract"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

// Address retriever is used to retrieve the addresses associated with a contract
type AddressRetriever interface {
	RetrieveTokenHolderAddresses(info contract.Contract) (map[common.Address]bool, error)
}

type addressRetriever struct {
	db   *postgres.DB
	mode types.Mode
}

func NewAddressRetriever(db *postgres.DB, mode types.Mode) (r *addressRetriever) {
	return &addressRetriever{
		db:   db,
		mode: mode,
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

func (r *addressRetriever) retrieveTransferAddresses(con contract.Contract) ([]string, error) {
	transferAddrs := make([]string, 0)
	event := con.Events["Transfer"]

	for _, field := range event.Fields { // Iterate over event fields, finding the ones with address type

		if field.Type.T == abi.AddressTy { // If they have address type, retrieve those addresses
			addrs := make([]string, 0)
			pgStr := fmt.Sprintf("SELECT %s_ FROM %s_%s.%s_event", strings.ToLower(field.Name), r.mode.String(), strings.ToLower(con.Address), strings.ToLower(event.Name))
			err := r.db.Select(&addrs, pgStr)
			if err != nil {
				return []string{}, err
			}

			transferAddrs = append(transferAddrs, addrs...) // And append them to the growing list
		}
	}

	return transferAddrs, nil
}

func (r *addressRetriever) retrieveTokenMintees(con contract.Contract) ([]string, error) {
	mintAddrs := make([]string, 0)
	event := con.Events["Mint"]

	for _, field := range event.Fields { // Iterate over event fields, finding the ones with address type

		if field.Type.T == abi.AddressTy { // If they have address type, retrieve those addresses
			addrs := make([]string, 0)
			pgStr := fmt.Sprintf("SELECT %s_ FROM %s_%s.%s_event", strings.ToLower(field.Name), r.mode.String(), strings.ToLower(con.Address), strings.ToLower(event.Name))
			err := r.db.Select(&addrs, pgStr)
			if err != nil {
				return []string{}, err
			}

			mintAddrs = append(mintAddrs, addrs...) // And append them to the growing list
		}
	}

	return mintAddrs, nil
}
