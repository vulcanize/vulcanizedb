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

package btc

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

// IPLDPGFetcher satisfies the IPLDFetcher interface for ethereum
// it interfaces directly with PG-IPFS instead of going through a node-interface or remote node
type IPLDPGFetcher struct {
	db *postgres.DB
}

// NewIPLDPGFetcher creates a pointer to a new IPLDPGFetcher
func NewIPLDPGFetcher(db *postgres.DB) *IPLDPGFetcher {
	return &IPLDPGFetcher{
		db: db,
	}
}

// Fetch is the exported method for fetching and returning all the IPLDS specified in the CIDWrapper
func (f *IPLDPGFetcher) Fetch(cids shared.CIDsForFetching) (shared.IPLDs, error) {
	cidWrapper, ok := cids.(*CIDWrapper)
	if !ok {
		return nil, fmt.Errorf("btc fetcher: expected cids type %T got %T", &CIDWrapper{}, cids)
	}
	log.Debug("fetching iplds")
	iplds := IPLDs{}
	iplds.BlockNumber = cidWrapper.BlockNumber

	tx, err := f.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			shared.Rollback(tx)
			panic(p)
		} else if err != nil {
			shared.Rollback(tx)
		} else {
			err = tx.Commit()
		}
	}()

	iplds.Header, err = f.FetchHeader(tx, cidWrapper.Header)
	if err != nil {
		return nil, fmt.Errorf("btc pg fetcher: header fetching error: %s", err.Error())
	}
	iplds.Transactions, err = f.FetchTrxs(tx, cidWrapper.Transactions)
	if err != nil {
		return nil, fmt.Errorf("btc pg fetcher: transaction fetching error: %s", err.Error())
	}
	return iplds, err
}

// FetchHeaders fetches headers
func (f *IPLDPGFetcher) FetchHeader(tx *sqlx.Tx, c HeaderModel) (ipfs.BlockModel, error) {
	log.Debug("fetching header ipld")
	headerBytes, err := shared.FetchIPLD(tx, c.CID)
	if err != nil {
		return ipfs.BlockModel{}, err
	}
	return ipfs.BlockModel{
		Data: headerBytes,
		CID:  c.CID,
	}, nil
}

// FetchTrxs fetches transactions
func (f *IPLDPGFetcher) FetchTrxs(tx *sqlx.Tx, cids []TxModel) ([]ipfs.BlockModel, error) {
	log.Debug("fetching transaction iplds")
	trxIPLDs := make([]ipfs.BlockModel, len(cids))
	for i, c := range cids {
		trxBytes, err := shared.FetchIPLD(tx, c.CID)
		if err != nil {
			return nil, err
		}
		trxIPLDs[i] = ipfs.BlockModel{
			Data: trxBytes,
			CID:  c.CID,
		}
	}
	return trxIPLDs, nil
}
