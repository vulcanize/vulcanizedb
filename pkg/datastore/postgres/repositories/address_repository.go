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

package repositories

import (
	"database/sql"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type AddressRepository struct{}

func (repo AddressRepository) GetOrCreateAddress(db *postgres.DB, address string) (int, error) {
	stringAddressToCommonAddress := common.HexToAddress(address)
	hexAddress := stringAddressToCommonAddress.Hex()

	var addressId int
	getErr := db.Get(&addressId, `SELECT id FROM public.addresses WHERE address = $1`, hexAddress)
	if getErr == sql.ErrNoRows {
		insertErr := db.QueryRow(`INSERT INTO public.addresses (address) VALUES($1) RETURNING id`, hexAddress).Scan(&addressId)
		return addressId, insertErr
	}

	return addressId, getErr
}

func (repo AddressRepository) GetOrCreateAddressInTransaction(tx *sqlx.Tx, address string) (int, error) {
	stringAddressToCommonAddress := common.HexToAddress(address)
	hexAddress := stringAddressToCommonAddress.Hex()

	var addressId int
	getErr := tx.Get(&addressId, `SELECT id FROM public.addresses WHERE address = $1`, hexAddress)
	if getErr == sql.ErrNoRows {
		insertErr := tx.QueryRow(`INSERT INTO public.addresses (address) VALUES($1) RETURNING id`, hexAddress).Scan(&addressId)
		return addressId, insertErr
	}

	return addressId, getErr
}

func (repo AddressRepository) GetAddressById(db *postgres.DB, id int) (string, error){
	var address string
	getErr := db.Get(&address, `SELECT address FROM public.addresses WHERE id = $1`, id)
	if getErr != nil {
		return "", getErr
	}
	return address, nil
}
