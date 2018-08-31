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

package generic

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

// Retriever is used to iterate over addresses going into or out of a contract
// address in an attempt to generate a list of token holder addresses

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
	GetMinteesError   = "Error fetching token mintees from contract %s: %s"
	GetBurnersError   = "Error fetching token burners from contract %s: %s"
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

func (rt TokenHolderRetriever) retrieveTokenMintees() ([]string, error) {

	mintees := make([]string, 0)

	err := rt.Database.DB.Select(
		&mintees,
		`SELECT mintee FROM token_mints
               WHERE token_address = $1`,
		rt.ContractAddress,
	)
	if err != nil {
		return []string{}, newRetrieverError(err, GetMinteesError, rt.ContractAddress)
	}

	return mintees, nil
}

func (rt TokenHolderRetriever) retrieveTokenBurners() ([]string, error) {

	burners := make([]string, 0)

	err := rt.Database.DB.Select(
		&burners,
		`SELECT burner FROM token_burns
               WHERE token_address = $1`,
		rt.ContractAddress,
	)
	if err != nil {
		return []string{}, newRetrieverError(err, GetBurnersError, rt.ContractAddress)
	}

	return burners, nil
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

	mintees, err := rt.retrieveTokenMintees()
	if err != nil {
		return nil, err
	}

	burners, err := rt.retrieveTokenBurners()
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

	for _, addr := range mintees {
		contractAddresses[common.HexToAddress(addr)] = true
	}

	for _, addr := range burners {
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
