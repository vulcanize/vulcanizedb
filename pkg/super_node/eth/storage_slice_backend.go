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
	"bytes"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

// GetStorageSlice is the backend method for constructing and returning storage slice responses
func (b *Backend) GetStorageSlice(path string, depth int, root common.Hash) (*GetSliceResponse, error) {
	response := new(GetSliceResponse)
	response.init(path, depth, root)

	// "Trie Loading steps"
	trieLoadingStart := makeTimestamp()
	// Get the block height and state key for this storage root hash
	blockHeight, stateKey, err := b.getHeightAndStateKey(root)
	if err != nil {
		return nil, fmt.Errorf("GetStorageSlice blockheight and state key lookup error: %s", err.Error())
	}
	// Get all the paths
	headPath, stemPaths, slicePaths, err := getPaths(path, depth)
	if err != nil {
		return nil, fmt.Errorf("GetStorageSlice path generation error: %s", err.Error())
	}
	response.MetaData.TimeStats["00-trie-loading"] = strconv.Itoa(int(makeTimestamp() - trieLoadingStart))

	// Fetch stem nodes
	// some of the "stem" nodes can be leaf nodes (but not value nodes)
	stemNodes, stemLeafCIDs, _, timeSpent, err := b.getStorageNodes(blockHeight, stateKey, stemPaths)
	if err != nil {
		return nil, fmt.Errorf("GetStorageSlice stem node lookup error: %s", err.Error())
	}
	response.TrieNodes.Stem = stemNodes
	response.MetaData.TimeStats["01-fetch-stem-keys"] = timeSpent

	// Fetch slice nodes
	sliceNodes, sliceLeafCIDs, deepestPath, timeSpent, err := b.getStorageNodes(blockHeight, stateKey, slicePaths)
	if err != nil {
		return nil, fmt.Errorf("GetStorageSlice slice node lookup error: %s", err.Error())
	}
	response.TrieNodes.Slice = sliceNodes
	response.MetaData.TimeStats["02-fetch-slice-keys"] = timeSpent

	// Fetch head node
	headNode, headLeafCID, err := b.getStorageHeadNode(blockHeight, headPath)
	if err != nil {
		return nil, fmt.Errorf("GetStorageSlice head node lookup error: %s", err.Error())
	}
	response.TrieNodes.Head = headNode

	// Fill in metadata
	leafFetchStart := makeTimestamp()
	response.MetaData.NodeStats["03-leaves"] = strconv.Itoa(len(stemLeafCIDs) + len(sliceLeafCIDs) + len(headLeafCID))
	response.MetaData.TimeStats["03-fetch-leaves-info"] = strconv.Itoa(int(makeTimestamp() - leafFetchStart))
	response.MetaData.NodeStats["00-stem-and-head-nodes"] = strconv.Itoa(len(response.TrieNodes.Stem) + 1)
	maxDepth := deepestPath - len(headPath)
	if maxDepth < 0 {
		maxDepth = 0
	}
	response.MetaData.NodeStats["01-max-depth"] = strconv.Itoa(maxDepth)
	response.MetaData.NodeStats["02-total-trie-nodes"] = strconv.Itoa(len(response.TrieNodes.Stem) + len(response.TrieNodes.Slice) + 1)

	return response, nil
}

func (b *Backend) getHeightAndStateKey(root common.Hash) (int, string, error) {
	res := struct {
		Height       int    `db:"block_number"`
		StateLeafKey string `db:"state_leaf_key"`
	}{}
	pgStr := `SELECT block_number, state_leaf_key
			FROM eth.header_cids, eth.state_cids, eth.state_accounts
			WHERE state_cids.header_id = header_cids.id
			AND state_accounts.state_id = state_cids.id
			AND state_accounts.storage_root = $1
			ORDER BY block_number DESC
			LIMIT 1`
	return res.Height, res.StateLeafKey, b.DB.Get(&res, pgStr, root.String())
}

func (b Backend) getStorageNodes(blockHeight int, stateKey string, slicePaths [][]byte) (map[string]string, []cid.Cid, int, string, error) {
	nodes := make(map[string]string)
	fetchStart := makeTimestamp()

	// Get CIDs for all nodes at the provided paths
	leafCIDs, intermediateCIDs, deepestPath, err := b.getStorageNodeCIDs(blockHeight, stateKey, slicePaths)
	if err != nil {
		return nil, nil, 0, "", err
	}

	// Fetch IPLDs for all CIDs
	nodeCIDs := append(leafCIDs, intermediateCIDs...)
	nodeIPLDs := b.Fetcher.fetchBatch(nodeCIDs)
	if len(nodeIPLDs) != len(nodeCIDs) {
		return nil, nil, 0, "", fmt.Errorf("expected %d IPLDs, got %d", len(nodeCIDs), len(nodeIPLDs))
	}

	// Pack info into response map
	for _, nodeIPLD := range nodeIPLDs {
		decodedMh, err := multihash.Decode(nodeIPLD.Cid().Hash())
		if err != nil {
			return nil, nil, 0, "", err
		}
		hash := crypto.Keccak256Hash(nodeIPLD.RawData())
		if !bytes.Equal(hash.Bytes(), decodedMh.Digest) {
			panic("multihash digest should equal keccak of raw data")
		}
		nodes[common.Bytes2Hex(decodedMh.Digest)] = common.Bytes2Hex(nodeIPLD.RawData())
	}

	return nodes, leafCIDs, deepestPath, strconv.Itoa(int(makeTimestamp() - fetchStart)), nil
}

func (b *Backend) getStorageNodeCIDs(blockHeight int, stateKey string, paths [][]byte) ([]cid.Cid, []cid.Cid, int, error) {
	tx, err := b.DB.Beginx()
	if err != nil {
		return nil, nil, 0, err
	}
	deepestPath := 0
	intermediateNodes := make([]cid.Cid, 0)
	leafNodes := make([]cid.Cid, 0)
	pgStr := `SELECT storage_cids.cid, storage_nodes.node_type
			FROM eth.storage_cids, eth.state_cids, eth.header_cids
			WHERE storage_cids.state_id = state_cids.id
			AND state_cids.header_id = header_cids.id
			AND state_cids.state_leaf_key = $1
			AND header_cids.block_number <= $2
			AND storage_cids.storage_path = $3
			ORDER BY block_number DESC LIMIT 1`
	for _, path := range paths {
		var node nodeDBResponse
		if err := tx.Get(&node, pgStr, stateKey, blockHeight, path); err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			shared.Rollback(tx)
			return nil, nil, 0, err
		}
		pathLen := len(path)
		if pathLen > deepestPath {
			deepestPath = pathLen
		}
		dc, err := cid.Decode(node.CID)
		if err != nil {
			shared.Rollback(tx)
			return nil, nil, 0, err
		}
		if node.NodeType == 2 {
			leafNodes = append(leafNodes, dc)
		} else {
			intermediateNodes = append(intermediateNodes, dc)
		}
	}
	return leafNodes, intermediateNodes, deepestPath, tx.Commit()
}

func (b *Backend) getStorageHeadNode(blockHeight int, headPath []byte) (map[string]string, []cid.Cid, error) {
	headNode := make(map[string]string)
	pgStr := `SELECT state_cids.cid, state_cids.node_type
			FROM eth.state_cids, eth.header_cids
			WHERE state_cids.header_id = header_cids.id
			AND header_cids.block_number <= $1
			AND state_cids.state_path = $2
			ORDER BY block_number DESC LIMIT 1`
	var node nodeDBResponse
	if err := b.DB.Get(&node, pgStr, blockHeight, headPath); err != nil {
		return nil, nil, err
	}
	headCID, err := cid.Decode(node.CID)
	if err != nil {
		return nil, nil, err
	}
	headIPLD, err := b.Fetcher.fetch(headCID)
	if err != nil {
		return nil, nil, err
	}
	decodedMh, err := multihash.Decode(headIPLD.RawData())
	if err != nil {
		return nil, nil, err
	}
	hash := crypto.Keccak256Hash(headIPLD.RawData())
	if !bytes.Equal(hash.Bytes(), decodedMh.Digest) {
		panic("multihash digest should equal keccak of raw data")
	}
	headNode[common.Bytes2Hex(decodedMh.Digest)] = common.Bytes2Hex(headIPLD.RawData())
	if node.NodeType == 2 {
		return headNode, []cid.Cid{headCID}, nil
	}
	return headNode, nil, nil
}
