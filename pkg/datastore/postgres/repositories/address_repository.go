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

package repositories

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type AddressRepository struct{}

const getOrCreateAddressQuery = `WITH addressId AS (
			INSERT INTO addresses (address) VALUES ($1) ON CONFLICT DO NOTHING RETURNING id
		)
		SELECT id FROM addresses WHERE address = $1
		UNION
		SELECT id FROM addressId`

func (AddressRepository) GetOrCreateAddress(db *postgres.DB, address string) (int64, error) {
	stringAddressToCommonAddress := common.HexToAddress(address)
	hexAddress := stringAddressToCommonAddress.Hex()

	var addressId int64
	getOrCreateErr := db.Get(&addressId, getOrCreateAddressQuery, checksumAddress)

	return addressId, getOrCreateErr
}

func (AddressRepository) GetOrCreateAddressInTransaction(tx *sqlx.Tx, address string) (int64, error) {
	stringAddressToCommonAddress := common.HexToAddress(address)
	hexAddress := stringAddressToCommonAddress.Hex()

	var addressId int64
	getOrCreateErr := tx.Get(&addressId, getOrCreateAddressQuery, checksumAddress)

	return addressId, getOrCreateErr
}

func (AddressRepository) GetAddressById(db *postgres.DB, id int64) (string, error) {
	var address string
	getErr := db.Get(&address, `SELECT address FROM public.addresses WHERE id = $1`, id)
	return address, getErr
}
