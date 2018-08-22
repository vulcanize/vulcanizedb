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
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"log"
)

// Retriever is used to iterate over addresses going into or out of a contract
// address in an attempt to generate a list of token holder addresses

type RetrieverInterface interface {
	RetrieveSendingAddresses() ([]string, error)
	RetrieveReceivingAddresses() ([]string, error)
	RetrieveContractAssociatedAddresses() (map[common.Address]bool, error)
}

type Retriever struct {
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
	GetSenderError   = "Error fetching addresses receiving from contract %s: %s"
	GetReceiverError = "Error fetching addresses sending to contract %s: %s"
)

func NewRetriever(db *postgres.DB, address string) Retriever {
	return Retriever{
		Database:        db,
		ContractAddress: address,
	}
}

func (rt Retriever) RetrieveReceivingAddresses() ([]string, error) {

	receiversFromContract := make([]string, 0)

	err := rt.Database.DB.Select(
		&receiversFromContract,
		`SELECT tx_to FROM TRANSACTIONS
               WHERE tx_from = $1
			   LIMIT 20`,
		rt.ContractAddress,
	)
	if err != nil {
		return []string{}, newRetrieverError(err, GetReceiverError, rt.ContractAddress)
	}
	return receiversFromContract, err
}

func (rt Retriever) RetrieveSendingAddresses() ([]string, error) {

	sendersToContract := make([]string, 0)

	err := rt.Database.DB.Select(
		&sendersToContract,
		`SELECT tx_from FROM TRANSACTIONS
			WHERE tx_to = $1
		    LIMIT 20`,
		rt.ContractAddress,
	)
	if err != nil {
		return []string{}, newRetrieverError(err, GetSenderError, rt.ContractAddress)
	}
	return sendersToContract, err
}

func (rt Retriever) RetrieveContractAssociatedAddresses() (map[common.Address]bool, error) {

	sending, err := rt.RetrieveSendingAddresses()
	if err != nil {
		return nil, err
	}

	receiving, err := rt.RetrieveReceivingAddresses()
	if err != nil {
		return nil, err
	}

	contractAddresses := make(map[common.Address]bool)

	for _, addr := range sending {
		contractAddresses[common.HexToAddress(addr)] = true
	}

	for _, addr := range receiving {
		contractAddresses[common.HexToAddress(addr)] = true
	}

	return contractAddresses, nil
}
