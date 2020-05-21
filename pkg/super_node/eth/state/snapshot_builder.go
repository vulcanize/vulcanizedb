// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free sofssbare: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Sofssbare Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package state

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/statediff"
	"github.com/jmoiron/sqlx"
	"github.com/multiformats/go-multihash"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs/ipld"
	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

var (
	method   = "statediff_stateTrieAt"
	nullHash = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
)

type SnapsShotBuilder struct {
	db      *postgres.DB
	client  *rpc.Client
	params  statediff.Params
	timeout time.Duration
}

// NewSnapsShotBuilder returns a new SnapsShotBuilder
func NewSnapsShotBuilder(db *postgres.DB, httpClient *rpc.Client) *SnapsShotBuilder {
	return &SnapsShotBuilder{
		db:     db,
		client: httpClient,
		params: statediff.Params{
			IncludeTD:    true,
			IncludeBlock: true,
		},
		timeout: time.Minute * 10,
	}
}

// BuildSnapShotAt creates a state snapshot at the provided height
func (ssb *SnapsShotBuilder) BuildSnapShotAt(height uint64) error {
	payload, err := ssb.fetch(height)
	if err != nil {
		return err
	}
	convertedPayload, err := convert(payload)
	if err != nil {
		return err
	}
	return ssb.publish(convertedPayload)
}

// calls StateTrieAt(ctx context.Context, blockNumber uint64, params Params) (*Payload, error)
func (ssb *SnapsShotBuilder) fetch(height uint64) (*statediff.Payload, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ssb.timeout)
	defer cancel()
	res := new(statediff.Payload)
	if err := ssb.client.CallContext(ctx, res, method, height, ssb.params); err != nil {
		return nil, fmt.Errorf("ethereum StateTrieAt err for block %d: %s", height, err.Error())
	}
	return res, nil
}

func (ssb *SnapsShotBuilder) publish(payload *eth.ConvertedPayload) error {
	headerNode, err := ipld.NewEthHeader(payload.Block.Header())
	if err != nil {
		return err
	}

	tx, err := ssb.db.Beginx()
	if err != nil {
		return err
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

	if err := shared.PublishIPLD(tx, headerNode); err != nil {
		return err
	}
	header := eth.HeaderModel{
		CID:             headerNode.Cid().String(),
		ParentHash:      payload.Block.ParentHash().String(),
		BlockNumber:     payload.Block.Number().String(),
		BlockHash:       payload.Block.Hash().String(),
		TotalDifficulty: payload.TotalDifficulty.String(),
		Reward:          "0",
		Bloom:           payload.Block.Bloom().Bytes(),
		StateRoot:       payload.Block.Root().String(),
		RctRoot:         payload.Block.ReceiptHash().String(),
		TxRoot:          payload.Block.TxHash().String(),
		UncleRoot:       payload.Block.UncleHash().String(),
		Timestamp:       payload.Block.Time(),
	}
	headerID, err := ssb.indexHeaderCID(tx, header)
	if err != nil {
		return err
	}

	err = ssb.publishAndIndexStateAndStorage(tx, payload, headerID)
	return err // return err variable explicitly so that we return the err = tx.Commit() assignment in the defer
}

func (ssb *SnapsShotBuilder) publishAndIndexStateAndStorage(tx *sqlx.Tx, ipldPayload *eth.ConvertedPayload, headerID int64) error {
	// Publish and index state and storage
	for _, stateNode := range ipldPayload.StateNodes {
		stateCIDStr, err := shared.PublishRaw(tx, ipld.MEthStateTrie, multihash.KECCAK_256, stateNode.Value)
		if err != nil {
			return err
		}
		stateModel := eth.StateNodeModel{
			Path:     stateNode.Path,
			StateKey: stateNode.LeafKey.String(),
			CID:      stateCIDStr,
			NodeType: eth.ResolveFromNodeType(stateNode.Type),
		}
		stateID, err := ssb.indexStateTrieCID(tx, stateModel, headerID)
		if err != nil {
			return err
		}
		// If we have a leaf, decode and index the account data and any associated storage diffs
		if stateNode.Type == statediff.Leaf {
			var i []interface{}
			if err := rlp.DecodeBytes(stateNode.Value, &i); err != nil {
				return err
			}
			if len(i) != 2 {
				return fmt.Errorf("eth IPLDPublisherAndIndexer expected state leaf node rlp to decode into ssbo elements")
			}
			for _, storageNode := range ipldPayload.StorageNodes[common.Bytes2Hex(stateNode.Path)] {
				storageCIDStr, err := shared.PublishRaw(tx, ipld.MEthStorageTrie, multihash.KECCAK_256, storageNode.Value)
				if err != nil {
					return err
				}
				storageModel := eth.StorageNodeModel{
					Path:       storageNode.Path,
					StorageKey: storageNode.LeafKey.Hex(),
					CID:        storageCIDStr,
					NodeType:   eth.ResolveFromNodeType(storageNode.Type),
				}
				if err := ssb.indexStorageTrieCID(tx, storageModel, stateID); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func convert(payload *statediff.Payload) (*eth.ConvertedPayload, error) {
	block := new(types.Block)
	if err := rlp.DecodeBytes(payload.BlockRlp, block); err != nil {
		return nil, err
	}
	convertedPayload := &eth.ConvertedPayload{
		TotalDifficulty: payload.TotalDifficulty,
		Block:           block,
		StateNodes:      make([]eth.TrieNode, 0),
		StorageNodes:    make(map[string][]eth.TrieNode),
	}

	stateDiff := new(statediff.StateObject)
	if err := rlp.DecodeBytes(payload.StateObjectRlp, stateDiff); err != nil {
		return nil, err
	}
	for _, stateNode := range stateDiff.Nodes {
		statePath := common.Bytes2Hex(stateNode.Path)
		convertedPayload.StateNodes = append(convertedPayload.StateNodes, eth.TrieNode{
			Path:    stateNode.Path,
			Value:   stateNode.NodeValue,
			Type:    stateNode.NodeType,
			LeafKey: common.BytesToHash(stateNode.LeafKey),
		})
		for _, storageNode := range stateNode.StorageNodes {
			convertedPayload.StorageNodes[statePath] = append(convertedPayload.StorageNodes[statePath], eth.TrieNode{
				Path:    storageNode.Path,
				Value:   storageNode.NodeValue,
				Type:    storageNode.NodeType,
				LeafKey: common.BytesToHash(storageNode.LeafKey),
			})
		}
	}

	return convertedPayload, nil
}

// header shares same table as regularly indexed headers
// but we are careful not to overwrite anything and leave the validation level at 0
// so that the regular processes will still collect diff data at that height
// and fill in things we miss by this process such as tx, uncles, receipts, and miner rewards
func (ssb *SnapsShotBuilder) indexHeaderCID(tx *sqlx.Tx, header eth.HeaderModel) (int64, error) {
	var headerID int64
	err := tx.QueryRowx(`INSERT INTO eth.header_cids (block_number, block_hash, parent_hash, cid, td, node_id, reward, state_root, tx_root, receipt_root, uncle_root, bloom, timestamp, times_validated)
								VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
								ON CONFLICT (block_number, block_hash) DO UPDATE SET block_number = header_cids.block_number
								RETURNING id`,
		header.BlockNumber, header.BlockHash, header.ParentHash, header.CID, header.TotalDifficulty, ssb.db.NodeID, header.Reward, header.StateRoot, header.TxRoot,
		header.RctRoot, header.UncleRoot, header.Bloom, header.Timestamp, 0).Scan(&headerID)
	return headerID, err
}

// we write state trie nodes collected in this fashion to a different table since they represent something different
func (ssb *SnapsShotBuilder) indexStateTrieCID(tx *sqlx.Tx, stateNode eth.StateNodeModel, headerID int64) (int64, error) {
	var stateID int64
	var stateKey string
	if stateNode.StateKey != nullHash.String() {
		stateKey = stateNode.StateKey
	}
	err := tx.QueryRowx(`INSERT INTO eth.state_trie_cids (header_id, state_leaf_key, cid, state_path, node_type) VALUES ($1, $2, $3, $4, $5)
									ON CONFLICT (header_id, state_path) DO UPDATE SET (state_leaf_key, cid, node_type) = ($2, $3, $5)
									RETURNING id`,
		headerID, stateKey, stateNode.CID, stateNode.Path, stateNode.NodeType).Scan(&stateID)
	return stateID, err
}

// we write storage trie nodes collected in this fashion to a different table since they represent something different
func (ssb *SnapsShotBuilder) indexStorageTrieCID(tx *sqlx.Tx, storageCID eth.StorageNodeModel, stateID int64) error {
	var storageKey string
	if storageCID.StorageKey != nullHash.String() {
		storageKey = storageCID.StorageKey
	}
	_, err := tx.Exec(`INSERT INTO eth.storage_trie_cids (state_id, storage_leaf_key, cid, storage_path, node_type) VALUES ($1, $2, $3, $4, $5) 
							  ON CONFLICT (state_id, storage_path) DO UPDATE SET (storage_leaf_key, cid, node_type) = ($2, $3, $5)`,
		stateID, storageKey, storageCID.CID, storageCID.Path, storageCID.NodeType)
	return err
}
