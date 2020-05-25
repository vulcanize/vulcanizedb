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

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/statediff"
	"github.com/jmoiron/sqlx"
	"github.com/multiformats/go-multihash"
	"github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs/ipld"
	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

var (
	nullHash    = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	txSizeLimit = 1024 * 500
	queuePutStr = "INSERT INTO eth.queued_nodes (header_id, rlp) VALUES ($1, $2)"
)

type SnapshotBuilder struct {
	db       *postgres.DB
	http, ws *rpc.Client
	txSize   int
}

// NewSnapshotBuilder returns a new SnapshotBuilder
func NewSnapshotBuilder(db *postgres.DB, http, ws *rpc.Client) *SnapshotBuilder {
	return &SnapshotBuilder{
		db:     db,
		http:   http,
		ws:     ws,
		txSize: 0,
	}
}

// BuildSnapShotAt creates a state snapshot at the provided height
func (ssb *SnapshotBuilder) BuildSnapshotAt(height uint64) error {
	header, err := ssb.getHeader(height)
	if err != nil {
		return err
	}
	errChan := make(chan error)
	doneChan := make(chan bool)
	go ssb.queueStateNodes(header, errChan, doneChan)
	<-doneChan
	return nil
}

func (ssb *SnapshotBuilder) queueStateNodes(header eth.HeaderModel, errChan chan<- error, done chan bool) {
	stateNodes := make(chan []byte, 20000)
	ctx, cancel := context.WithCancel(context.Background())
	sub, err := ssb.stream(ctx, stateNodes, common.HexToHash(header.BlockHash))
	if err != nil {
		logrus.Fatal(err)
	}
	ticker := time.NewTicker(time.Second * 30)
	tx, err := ssb.db.Beginx()
	if err != nil {
		logrus.Fatal(err)
	}
	defer cancel()
	defer close(done)
	for {
		select {
		case <-ticker.C:
			if ssb.txSize > txSizeLimit {
				logrus.Infof("SnapshotBuilder committing %d bytes to db", ssb.txSize)
				if err := tx.Commit(); err != nil {
					errChan <- err
				}
				logrus.Info("SnapshotBuilder beginning new tx")
				tx, err = ssb.db.Beginx()
				if err != nil {
					errChan <- err
				}
			}
			continue
		case stateNode := <-stateNodes:
			if err := ssb.queueStateNode(tx, stateNode, header.ID); err != nil {
				errChan <- err
			}
			continue
		default:
		}
		select {
		case stateNode := <-stateNodes:
			if err := ssb.queueStateNode(tx, stateNode, header.ID); err != nil {
				errChan <- err
			}
		case err := <-sub.Err():
			if err != nil {
				logrus.Errorf("SnapshotBuilder subscription err: %s", err.Error())
				shared.Rollback(tx)
				return
			}
			logrus.Infof("SnapshotBuilder committing %d nodes to db", ssb.txSize)
			logrus.Info("SnapshotBuilder process has completed")
			tx.Commit()
			return
		}
	}
}

func (ssb *SnapshotBuilder) queueStateNode(tx *sqlx.Tx, stateNodeRLP []byte, headerID int64) error {
	ssb.txSize += len(stateNodeRLP)
	_, err := tx.Exec(queuePutStr, headerID, stateNodeRLP)
	return err
}

func (ssb *SnapshotBuilder) writeStateNode(tx *sqlx.Tx, stateNode statediff.StateNode, headerID int64) error {
	stateCIDStr, err := shared.PublishRaw(tx, ipld.MEthStateTrie, multihash.KECCAK_256, stateNode.NodeValue)
	if err != nil {
		return err
	}
	stateModel := eth.StateNodeModel{
		Path:     stateNode.Path,
		StateKey: common.Bytes2Hex(stateNode.LeafKey),
		CID:      stateCIDStr,
		NodeType: eth.ResolveFromNodeType(stateNode.NodeType),
	}
	stateID, err := ssb.indexStateTrieCID(tx, stateModel, headerID)
	if err != nil {
		return err
	}
	ssb.txSize += 1
	for _, storageNode := range stateNode.StorageNodes {
		storageCIDStr, err := shared.PublishRaw(tx, ipld.MEthStorageTrie, multihash.KECCAK_256, storageNode.NodeValue)
		if err != nil {
			return err
		}
		storageModel := eth.StorageNodeModel{
			Path:       storageNode.Path,
			StorageKey: common.Bytes2Hex(storageNode.LeafKey),
			CID:        storageCIDStr,
			NodeType:   eth.ResolveFromNodeType(storageNode.NodeType),
		}
		if err := ssb.indexStorageTrieCID(tx, storageModel, stateID); err != nil {
			return err
		}
		ssb.txSize += 1
	}
	return nil
}

func (ssb *SnapshotBuilder) getHeader(height uint64) (eth.HeaderModel, error) {
	// if we already have a unequivocally valid header in our db at this height, use it
	header, err := ssb.retrieveHeader(height)
	if err == nil {
		logrus.Info("SnapshotBuilder using header found in local db")
		return header, nil
	}
	// otherwise fetch the header remotely, and publish and index it locally
	logrus.Info("SnapshotBuilder fetching remote header")
	return ssb.fetchHeader(height)
}

func (ssb *SnapshotBuilder) retrieveHeader(height uint64) (eth.HeaderModel, error) {
	headers := make([]eth.HeaderModel, 0)
	pgStr := `SELECT * FROM btc.header_cids
				WHERE block_number = $1
				AND times_validation > 0`
	if err := ssb.db.Select(&headers, pgStr, height); err != nil {
		return eth.HeaderModel{}, err
	}
	switch len(headers) {
	case 0:
		return eth.HeaderModel{}, fmt.Errorf("no valid header at height %d", height)
	case 1:
		return headers[0], nil
	default:
		return eth.HeaderModel{}, fmt.Errorf("more than one valid header at height %d", height)
	}
}

func (ssb *SnapshotBuilder) fetchHeader(height uint64) (eth.HeaderModel, error) {
	var head *types.Header
	err := ssb.http.CallContext(context.Background(), &head, "eth_getBlockByNumber", hexutil.EncodeUint64(height), false)
	if err == nil && head == nil {
		return eth.HeaderModel{}, ethereum.NotFound
	}
	headerNode, err := ipld.NewEthHeader(head)
	if err != nil {
		return eth.HeaderModel{}, err
	}

	tx, err := ssb.db.Beginx()
	if err != nil {
		return eth.HeaderModel{}, err
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
		return eth.HeaderModel{}, err
	}
	header := eth.HeaderModel{
		CID:             headerNode.Cid().String(),
		ParentHash:      head.ParentHash.String(),
		BlockNumber:     head.Number.String(),
		BlockHash:       head.Hash().String(),
		TotalDifficulty: "0",
		Reward:          "0",
		Bloom:           head.Bloom.Bytes(),
		StateRoot:       head.Root.String(),
		RctRoot:         head.ReceiptHash.String(),
		TxRoot:          head.TxHash.String(),
		UncleRoot:       head.UncleHash.String(),
		Timestamp:       head.Time,
	}
	headerID, err := ssb.indexHeaderCID(tx, header)
	if err != nil {
		return eth.HeaderModel{}, err
	}
	header.ID = headerID
	return header, err // return err explicitly so it is assigned to in the defer
}

// calls StreamTrie(ctx context.Context, hash common.Hash, done chan bool) (*rpc.Subscription, error)
func (ssb *SnapshotBuilder) stream(ctx context.Context, stateNodes chan []byte, hash common.Hash) (shared.ClientSubscription, error) {
	return ssb.ws.Subscribe(ctx, "statediff", stateNodes, "streamTrie", hash)
}

// header shares same table as regularly indexed headers
// but we are careful not to overwrite anything and leave the validation level at 0
// so that the regular processes will still collect diff data at that height
// and fill in things we miss by this process such as tx, uncles, receipts, and miner rewards
func (ssb *SnapshotBuilder) indexHeaderCID(tx *sqlx.Tx, header eth.HeaderModel) (int64, error) {
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
func (ssb *SnapshotBuilder) indexStateTrieCID(tx *sqlx.Tx, stateNode eth.StateNodeModel, headerID int64) (int64, error) {
	var stateID int64
	var stateKey string
	if stateNode.StateKey != nullHash.String() {
		stateKey = stateNode.StateKey
	}
	err := tx.QueryRowx(`INSERT INTO eth.state_cids (header_id, state_leaf_key, cid, state_path, node_type, eventual) VALUES ($1, $2, $3, $4, $5, $6)
									ON CONFLICT (header_id, state_path) DO UPDATE SET header_id = state_cids.header_id
									RETURNING id`,
		headerID, stateKey, stateNode.CID, stateNode.Path, stateNode.NodeType, true).Scan(&stateID)
	return stateID, err
}

// we write storage trie nodes collected in this fashion to a different table since they represent something different
func (ssb *SnapshotBuilder) indexStorageTrieCID(tx *sqlx.Tx, storageCID eth.StorageNodeModel, stateID int64) error {
	var storageKey string
	if storageCID.StorageKey != nullHash.String() {
		storageKey = storageCID.StorageKey
	}
	_, err := tx.Exec(`INSERT INTO eth.storage_trie_cids (state_id, storage_leaf_key, cid, storage_path, node_type, eventual) VALUES ($1, $2, $3, $4, $5, $6) 
							  ON CONFLICT (state_id, storage_path) DO NOTHING`,
		stateID, storageKey, storageCID.CID, storageCID.Path, storageCID.NodeType)
	return err
}
