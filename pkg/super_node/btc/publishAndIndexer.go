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
	"strconv"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs/ipld"
	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

// IPLDPublisherAndIndexer satisfies the IPLDPublisher interface for bitcoin
// It interfaces directly with the public.blocks table of PG-IPFS rather than going through an ipfs intermediary
// It publishes and indexes IPLDs together in a single sqlx.Tx
type IPLDPublisherAndIndexer struct {
	indexer *CIDIndexer
}

// NewIPLDPublisherAndIndexer creates a pointer to a new IPLDPublisherAndIndexer which satisfies the IPLDPublisher interface
func NewIPLDPublisherAndIndexer(db *postgres.DB) *IPLDPublisherAndIndexer {
	return &IPLDPublisherAndIndexer{
		indexer: NewCIDIndexer(db),
	}
}

// Publish publishes an IPLDPayload to IPFS and returns the corresponding CIDPayload
func (pub *IPLDPublisherAndIndexer) Publish(payload shared.ConvertedData) (shared.CIDsForIndexing, error) {
	ipldPayload, ok := payload.(ConvertedPayload)
	if !ok {
		return nil, fmt.Errorf("btc publisher expected payload type %T got %T", ConvertedPayload{}, payload)
	}
	// Generate the iplds
	headerNode, txNodes, txTrieNodes, err := ipld.FromHeaderAndTxs(ipldPayload.Header, ipldPayload.Txs)
	if err != nil {
		return nil, err
	}

	// Begin new db tx
	tx, err := pub.indexer.db.Beginx()
	if err != nil {
		return nil, err
	}

	// Publish trie nodes
	for _, node := range txTrieNodes {
		if err := shared.PublishIPLD(tx, node); err != nil {
			shared.Rollback(tx)
			return nil, err
		}
	}

	// Publish and index header
	if err := shared.PublishIPLD(tx, headerNode); err != nil {
		shared.Rollback(tx)
		return nil, err
	}
	header := HeaderModel{
		CID:         headerNode.Cid().String(),
		ParentHash:  ipldPayload.Header.PrevBlock.String(),
		BlockNumber: strconv.Itoa(int(ipldPayload.BlockPayload.BlockHeight)),
		BlockHash:   ipldPayload.Header.BlockHash().String(),
		Timestamp:   ipldPayload.Header.Timestamp.UnixNano(),
		Bits:        ipldPayload.Header.Bits,
	}
	headerID, err := pub.indexer.indexHeaderCID(tx, header)
	if err != nil {
		shared.Rollback(tx)
		return nil, err
	}

	// Publish and index txs
	for i, txNode := range txNodes {
		if err := shared.PublishIPLD(tx, txNode); err != nil {
			shared.Rollback(tx)
			return nil, err
		}
		txModel := ipldPayload.TxMetaData[i]
		txModel.CID = txNode.Cid().String()
		txID, err := pub.indexer.indexTransactionCID(tx, txModel, headerID)
		if err != nil {
			shared.Rollback(tx)
			return nil, err
		}
		for _, input := range txModel.TxInputs {
			if err := pub.indexer.indexTxInput(tx, input, txID); err != nil {
				shared.Rollback(tx)
				return nil, err
			}
		}
		for _, output := range txModel.TxOutputs {
			if err := pub.indexer.indexTxOutput(tx, output, txID); err != nil {
				shared.Rollback(tx)
				return nil, err
			}
		}
	}

	// This IPLDPublisher does both publishing and indexing, we do not need to pass anything forward to the indexer
	return nil, tx.Commit()
}

// Index satisfies the shared.CIDIndexer interface
func (pub *IPLDPublisherAndIndexer) Index(cids shared.CIDsForIndexing) error {
	return nil
}
