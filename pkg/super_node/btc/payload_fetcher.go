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
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

// PayloadFetcher satisfies the PayloadFetcher interface for bitcoin
type PayloadFetcher struct {
	// PayloadFetcher is thread-safe as long as the underlying client is thread-safe, since it has/modifies no other state
	// http.Client is thread-safe
	client *rpcclient.Client
}

// NewStateDiffFetcher returns a PayloadFetcher
func NewPayloadFetcher(c *rpcclient.ConnConfig) (*PayloadFetcher, error) {
	client, err := rpcclient.New(c, nil)
	if err != nil {
		return nil, err
	}
	return &PayloadFetcher{
		client: client,
	}, nil
}

// FetchAt fetches the block payloads at the given block heights
func (fetcher *PayloadFetcher) FetchAt(blockHeights []uint64) ([]shared.RawChainData, error) {
	blockPayloads := make([]shared.RawChainData, len(blockHeights))
	for i, height := range blockHeights {
		hash, err := fetcher.client.GetBlockHash(int64(height))
		if err != nil {
			return nil, err
		}
		block, err := fetcher.client.GetBlock(hash)
		if err != nil {
			return nil, err
		}
		blockPayloads[i] = BlockPayload{
			Height: int64(height),
			Header: &block.Header,
			Txs:    msgTxsToUtilTxs(block.Transactions),
		}
	}
	return blockPayloads, nil
}

func msgTxsToUtilTxs(msgs []*wire.MsgTx) []*btcutil.Tx {
	txs := make([]*btcutil.Tx, len(msgs))
	for i, msg := range msgs {
		tx := btcutil.NewTx(msg)
		tx.SetIndex(i)
		txs[i] = tx
	}
	return txs
}
