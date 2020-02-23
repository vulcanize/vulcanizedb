// VulcanizeDB
// Copyright © 2019 Vulcanize
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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

package repository

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/postgres"
)

const getOrCreateAddressQuery = `WITH addressId AS (
			INSERT INTO addresses (address, hashed_address) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING id
		)
		SELECT id FROM addresses WHERE address = $1
		UNION
		SELECT id FROM addressId`

func GetOrCreateAddress(db *postgres.DB, address string) (int64, error) {
	checksumAddress := getChecksumAddress(address)
	hashedAddress := utils.HexToKeccak256Hash(checksumAddress).Hex()

	var addressID int64
	getOrCreateErr := db.Get(&addressID, getOrCreateAddressQuery, checksumAddress, hashedAddress)

	return addressID, getOrCreateErr
}

func GetOrCreateAddressInTransaction(tx *sqlx.Tx, address string) (int64, error) {
	checksumAddress := getChecksumAddress(address)
	hashedAddress := utils.HexToKeccak256Hash(checksumAddress).Hex()

	var addressID int64
	getOrCreateErr := tx.Get(&addressID, getOrCreateAddressQuery, checksumAddress, hashedAddress)

	return addressID, getOrCreateErr
}

func GetAddressByID(db *postgres.DB, id int64) (string, error) {
	var address string
	getErr := db.Get(&address, `SELECT address FROM public.addresses WHERE id = $1`, id)
	return address, getErr
}

func getChecksumAddress(address string) string {
	stringAddressToCommonAddress := common.HexToAddress(address)
	return stringAddressToCommonAddress.Hex()
}
