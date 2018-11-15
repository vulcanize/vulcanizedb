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

package erc20_watcher

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

// ERC20 retriever is used to iterate over addresses involved in token
// transfers and approvals to generate a list of token holder addresses

type TokenHolderRetrieverInterface interface {
	RetrieveTokenHolderAddresses() (map[common.Address]bool, error)
	retrieveTokenSenders() ([]string, error)
	retrieveTokenReceivers() ([]string, error)
	retrieveTokenOwners() ([]string, error)
	retrieveTokenSpenders() ([]string, error)
}

type TokenHolderRetriever struct {
	Database        *postgres.DB
	ContractAddress string
}

type retrieverError struct {
	err     string
	msg     string
	address string
}

// Retriever error method
func (re *retrieverError) Error() string {
	return fmt.Sprintf(re.msg, re.address, re.err)
}

// Used to create a new retriever error for a given error and fetch method
func newRetrieverError(err error, msg string, address string) error {
	e := retrieverError{err.Error(), msg, address}
	log.Println(e.Error())
	return &e
}

// Constant error definitions
const (
	GetSendersError   = "Error fetching token senders from contract %s: %s"
	GetReceiversError = "Error fetching token receivers from contract %s: %s"
	GetOwnersError    = "Error fetching token owners from contract %s: %s"
	GetSpendersError  = "Error fetching token spenders from contract %s: %s"
)

func NewTokenHolderRetriever(db *postgres.DB, address string) TokenHolderRetriever {
	return TokenHolderRetriever{
		Database:        db,
		ContractAddress: address,
	}
}

func (rt TokenHolderRetriever) retrieveTokenSenders() ([]string, error) {

	senders := make([]string, 0)

	err := rt.Database.DB.Select(
		&senders,
		`SELECT from_address FROM token_transfers
               WHERE token_address = $1`,
		rt.ContractAddress,
	)
	if err != nil {
		return []string{}, newRetrieverError(err, GetSendersError, rt.ContractAddress)
	}

	return senders, nil
}

func (rt TokenHolderRetriever) retrieveTokenReceivers() ([]string, error) {

	receivers := make([]string, 0)

	err := rt.Database.DB.Select(
		&receivers,
		`SELECT to_address FROM token_transfers
               WHERE token_address = $1`,
		rt.ContractAddress,
	)
	if err != nil {
		return []string{}, newRetrieverError(err, GetReceiversError, rt.ContractAddress)
	}
	return receivers, err
}

func (rt TokenHolderRetriever) retrieveTokenOwners() ([]string, error) {

	owners := make([]string, 0)

	err := rt.Database.DB.Select(
		&owners,
		`SELECT owner FROM token_approvals
               WHERE token_address = $1`,
		rt.ContractAddress,
	)
	if err != nil {
		return []string{}, newRetrieverError(err, GetOwnersError, rt.ContractAddress)
	}

	return owners, nil
}

func (rt TokenHolderRetriever) retrieveTokenSpenders() ([]string, error) {

	spenders := make([]string, 0)

	err := rt.Database.DB.Select(
		&spenders,
		`SELECT spender FROM token_approvals
               WHERE token_address = $1`,
		rt.ContractAddress,
	)
	if err != nil {
		return []string{}, newRetrieverError(err, GetSpendersError, rt.ContractAddress)
	}

	return spenders, nil
}

func (rt TokenHolderRetriever) RetrieveTokenHolderAddresses() (map[common.Address]bool, error) {

	senders, err := rt.retrieveTokenSenders()
	if err != nil {
		return nil, err
	}

	receivers, err := rt.retrieveTokenReceivers()
	if err != nil {
		return nil, err
	}

	owners, err := rt.retrieveTokenOwners()
	if err != nil {
		return nil, err
	}

	spenders, err := rt.retrieveTokenSenders()
	if err != nil {
		return nil, err
	}

	contractAddresses := make(map[common.Address]bool)

	for _, addr := range senders {
		contractAddresses[common.HexToAddress(addr)] = true
	}

	for _, addr := range receivers {
		contractAddresses[common.HexToAddress(addr)] = true
	}

	for _, addr := range owners {
		contractAddresses[common.HexToAddress(addr)] = true
	}

	for _, addr := range spenders {
		contractAddresses[common.HexToAddress(addr)] = true
	}

	return contractAddresses, nil
}
