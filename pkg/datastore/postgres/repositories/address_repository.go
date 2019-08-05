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
	"database/sql"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type AddressRepository struct{}

func (AddressRepository) GetOrCreateAddress(db *postgres.DB, address string) (int, error) {
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

func (AddressRepository) GetOrCreateAddressInTransaction(tx *sqlx.Tx, address string) (int, error) {
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

func (AddressRepository) GetAddressById(db *postgres.DB, id int) (string, error) {
	var address string
	getErr := db.Get(&address, `SELECT address FROM public.addresses WHERE id = $1`, id)
	if getErr != nil {
		return "", getErr
	}
	return address, nil
}
