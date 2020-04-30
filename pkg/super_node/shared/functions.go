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

package shared

import (
	"bytes"

	"github.com/ipfs/go-cid"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipfs/go-ipfs-ds-help"
	node "github.com/ipfs/go-ipld-format"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
)

// ListContainsString used to check if a list of strings contains a particular string
func ListContainsString(sss []string, s string) bool {
	for _, str := range sss {
		if s == str {
			return true
		}
	}
	return false
}

// IPLDsContainBytes used to check if a list of strings contains a particular string
func IPLDsContainBytes(iplds []ipfs.BlockModel, b []byte) bool {
	for _, ipld := range iplds {
		if bytes.Equal(ipld.Data, b) {
			return true
		}
	}
	return false
}

// ListContainsGap used to check if a list of Gaps contains a particular Gap
func ListContainsGap(gapList []Gap, gap Gap) bool {
	for _, listGap := range gapList {
		if listGap == gap {
			return true
		}
	}
	return false
}

// HandleNullAddrPointer will return an emtpy string for a nil address pointer
func HandleNullAddrPointer(to *common.Address) string {
	if to == nil {
		return ""
	}
	return to.Hex()
}

// HandleNullAddr will return an empty string for a a null address
func HandleNullAddr(to common.Address) string {
	if to.Hex() == "0x0000000000000000000000000000000000000000" {
		return ""
	}
	return to.Hex()
}

// Rollback sql transaction and log any error
func Rollback(tx *sqlx.Tx) {
	if err := tx.Rollback(); err != nil {
		logrus.Error(err)
	}
}

// PublishIPLD is used to insert an ipld into Postgres blockstore with the provided tx
func PublishIPLD(tx *sqlx.Tx, i node.Node) error {
	dbKey := dshelp.CidToDsKey(i.Cid())
	prefixedKey := blockstore.BlockPrefix.String() + dbKey.String()
	raw := i.RawData()
	_, err := tx.Exec(`INSERT INTO public.blocks (key, data) VALUES ($1, $2) ON CONFLICT (key) DO NOTHING`, prefixedKey, raw)
	return err
}

// FetchIPLD is used to retrieve an ipld from Postgres blockstore with the provided tx
func FetchIPLD(tx *sqlx.Tx, cid string) ([]byte, error) {
	mhKey, err := MultihashKeyFromCIDString(cid)
	if err != nil {
		return nil, err
	}
	pgStr := `SELECT data FROM public.blocks WHERE key = $1`
	var block []byte
	return block, tx.Get(&block, pgStr, mhKey)
}

// MultihashKeyFromCIDString converts a cid string into a blockstore-prefixed multihash db key string
func MultihashKeyFromCIDString(c string) (string, error) {
	dc, err := cid.Decode(c)
	if err != nil {
		return "", err
	}
	dbKey := dshelp.CidToDsKey(dc)
	return blockstore.BlockPrefix.String() + dbKey.String(), nil
}
