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

package trie

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth"
)

// Backend is the struct for performing top-level trie processes
type Backend struct {
	Retriever *eth.CIDRetriever
	Fetcher   *eth.IPLDFetcher
	db        *postgres.DB
}

func NewEthBackend(db *postgres.DB, ipfsPath string) (*Backend, error) {
	r := eth.NewCIDRetriever(db)
	f, err := eth.NewIPLDFetcher(ipfsPath)
	if err != nil {
		return nil, err
	}
	return &Backend{
		Retriever: r,
		Fetcher:   f,
		db:        db,
	}, nil
}

type StateCID struct {
	Path string `db:"state_path"`
	CID  string `db:"cid"`
}

type StorageCID struct {
	Path string `db:"storage_path"`
	CID  string `db:"cid"`
}

// getStatePaths returns all of the distinct state paths up to the given height
// it also returns the CID for the most recent node corresponding to each path
func (b *Backend) getStatePaths(height int64) ([]StateCID, error) {
	paths := make([]StateCID, 0)
	pgStr := `SELECT DISTINCT ON (state_path) state_path, state_cids.cid
			FROM eth.state_cids, eth.header_cids
			WHERE state_cids.header_id = header_cids.id
			AND header.block_number <= $1
			ORDER BY state_path, header.block_number DESC`
	return paths, b.db.Select(&paths, pgStr, height)
}

// getStoragePaths returns all of the distinct storage paths up to the given height and for
// the account address corresponding to the provided state key
// it also returns the CID for the most recent node corresponding to each path
func (b *Backend) getStoragePaths(height int64, stateKey common.Hash) ([]StorageCID, error) {
	paths := make([]StorageCID, 0)
	pgStr := `SELECT DISTINCT ON (storage_path) storage_path, storage_cids.cid
			FROM eth.storage_cids, eth.state_cids, eth.header_cids
			WHERE storage_cids.state_id = state_cids.id
			AND state_cids.header_id = header_cids.id
			AND header.block_number <= $1
			AND state_cids.state_leaf_key = $2
			ORDER BY storage_path, header.block_number DESC`
	return paths, b.db.Select(&paths, pgStr, height, stateKey.String())
}
