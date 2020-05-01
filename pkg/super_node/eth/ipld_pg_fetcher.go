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

package eth

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

// IPLDPGFetcher satisfies the IPLDFetcher interface for ethereum
// It interfaces directly with PG-IPFS
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
		return nil, fmt.Errorf("eth fetcher: expected cids type %T got %T", &CIDWrapper{}, cids)
	}
	log.Debug("fetching iplds")
	iplds := IPLDs{}
	iplds.TotalDifficulty, ok = new(big.Int).SetString(cidWrapper.Header.TotalDifficulty, 10)
	if !ok {
		return nil, errors.New("eth fetcher: unable to set total difficulty")
	}
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
		return nil, fmt.Errorf("eth pg fetcher: header fetching error: %s", err.Error())
	}
	iplds.Uncles, err = f.FetchUncles(tx, cidWrapper.Uncles)
	if err != nil {
		return nil, fmt.Errorf("eth pg fetcher: uncle fetching error: %s", err.Error())
	}
	iplds.Transactions, err = f.FetchTrxs(tx, cidWrapper.Transactions)
	if err != nil {
		return nil, fmt.Errorf("eth pg fetcher: transaction fetching error: %s", err.Error())
	}
	iplds.Receipts, err = f.FetchRcts(tx, cidWrapper.Receipts)
	if err != nil {
		return nil, fmt.Errorf("eth pg fetcher: receipt fetching error: %s", err.Error())
	}
	iplds.StateNodes, err = f.FetchState(tx, cidWrapper.StateNodes)
	if err != nil {
		return nil, fmt.Errorf("eth pg fetcher: state fetching error: %s", err.Error())
	}
	iplds.StorageNodes, err = f.FetchStorage(tx, cidWrapper.StorageNodes)
	if err != nil {
		return nil, fmt.Errorf("eth pg fetcher: storage fetching error: %s", err.Error())
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

// FetchUncles fetches uncles
func (f *IPLDPGFetcher) FetchUncles(tx *sqlx.Tx, cids []UncleModel) ([]ipfs.BlockModel, error) {
	log.Debug("fetching uncle iplds")
	uncleIPLDs := make([]ipfs.BlockModel, len(cids))
	for i, c := range cids {
		uncleBytes, err := shared.FetchIPLD(tx, c.CID)
		if err != nil {
			return nil, err
		}
		uncleIPLDs[i] = ipfs.BlockModel{
			Data: uncleBytes,
			CID:  c.CID,
		}
	}
	return uncleIPLDs, nil
}

// FetchTrxs fetches transactions
func (f *IPLDPGFetcher) FetchTrxs(tx *sqlx.Tx, cids []TxModel) ([]ipfs.BlockModel, error) {
	log.Debug("fetching transaction iplds")
	trxIPLDs := make([]ipfs.BlockModel, len(cids))
	for i, c := range cids {
		txBytes, err := shared.FetchIPLD(tx, c.CID)
		if err != nil {
			return nil, err
		}
		trxIPLDs[i] = ipfs.BlockModel{
			Data: txBytes,
			CID:  c.CID,
		}
	}
	return trxIPLDs, nil
}

// FetchRcts fetches receipts
func (f *IPLDPGFetcher) FetchRcts(tx *sqlx.Tx, cids []ReceiptModel) ([]ipfs.BlockModel, error) {
	log.Debug("fetching receipt iplds")
	rctIPLDs := make([]ipfs.BlockModel, len(cids))
	for i, c := range cids {
		rctBytes, err := shared.FetchIPLD(tx, c.CID)
		if err != nil {
			return nil, err
		}
		rctIPLDs[i] = ipfs.BlockModel{
			Data: rctBytes,
			CID:  c.CID,
		}
	}
	return rctIPLDs, nil
}

// FetchState fetches state nodes
func (f *IPLDPGFetcher) FetchState(tx *sqlx.Tx, cids []StateNodeModel) ([]StateNode, error) {
	log.Debug("fetching state iplds")
	stateNodes := make([]StateNode, 0, len(cids))
	for _, stateNode := range cids {
		if stateNode.CID == "" {
			continue
		}
		stateBytes, err := shared.FetchIPLD(tx, stateNode.CID)
		if err != nil {
			return nil, err
		}
		stateNodes = append(stateNodes, StateNode{
			IPLD: ipfs.BlockModel{
				Data: stateBytes,
				CID:  stateNode.CID,
			},
			StateLeafKey: common.HexToHash(stateNode.StateKey),
			Type:         ResolveToNodeType(stateNode.NodeType),
			Path:         stateNode.Path,
		})
	}
	return stateNodes, nil
}

// FetchStorage fetches storage nodes
func (f *IPLDPGFetcher) FetchStorage(tx *sqlx.Tx, cids []StorageNodeWithStateKeyModel) ([]StorageNode, error) {
	log.Debug("fetching storage iplds")
	storageNodes := make([]StorageNode, 0, len(cids))
	for _, storageNode := range cids {
		if storageNode.CID == "" || storageNode.StateKey == "" {
			continue
		}
		storageBytes, err := shared.FetchIPLD(tx, storageNode.CID)
		if err != nil {
			return nil, err
		}
		storageNodes = append(storageNodes, StorageNode{
			IPLD: ipfs.BlockModel{
				Data: storageBytes,
				CID:  storageNode.CID,
			},
			StateLeafKey:   common.HexToHash(storageNode.StateKey),
			StorageLeafKey: common.HexToHash(storageNode.StorageKey),
			Type:           ResolveToNodeType(storageNode.NodeType),
			Path:           storageNode.Path,
		})
	}
	return storageNodes, nil
}
