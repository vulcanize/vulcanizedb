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

// GetStateSlice is the backend method for constructing and returning state slice responses
func (b *Backend) GetStateSlice(path string, depth int, root common.Hash) (*GetSliceResponse, error) {
	response := new(GetSliceResponse)
	response.init(path, depth, root)

	// "Trie Loading steps"
	trieLoadingStart := makeTimestamp()
	blockHeight, err := b.getHeight(root)
	if err != nil {
		return nil, fmt.Errorf("GetStateSlice blockheight lookup error: %s", err.Error())
	}
	headPath, stemPaths, slicePaths, err := getPaths(path, depth)
	if err != nil {
		return nil, fmt.Errorf("GetStateSlice path generation error: %s", err.Error())
	}
	response.MetaData.TimeStats["00-trie-loading"] = strconv.Itoa(int(makeTimestamp() - trieLoadingStart))

	// Fetch stem nodes
	// some of the "stem" nodes can be leaf nodes (but not value nodes)
	stemNodes, stemLeafCIDs, _, timeSpent, err := b.getStateNodes(blockHeight, stemPaths)
	if err != nil {
		return nil, fmt.Errorf("GetStateSlice stem node lookup error: %s", err.Error())
	}
	response.TrieNodes.Stem = stemNodes
	response.MetaData.TimeStats["01-fetch-stem-keys"] = timeSpent

	// Fetch slice nodes
	sliceNodes, sliceLeafCIDs, deepestPath, timeSpent, err := b.getStateNodes(blockHeight, slicePaths)
	if err != nil {
		return nil, fmt.Errorf("GetStateSlice slice node lookup error: %s", err.Error())
	}
	response.TrieNodes.Slice = sliceNodes
	response.MetaData.TimeStats["02-fetch-slice-keys"] = timeSpent

	// Fetch head node
	headNode, headLeafCID, err := b.getStateHeadNode(blockHeight, headPath)
	if err != nil {
		return nil, fmt.Errorf("GetStateSlice head node lookup error: %s", err.Error())
	}
	response.TrieNodes.Head = headNode

	// Fetch leaf contract data and fill in remaining metadata
	leafFetchStart := makeTimestamp()
	leafNodes := make([]cid.Cid, 0, len(stemLeafCIDs)+len(sliceLeafCIDs)+len(headLeafCID))
	leafNodes = append(leafNodes, stemLeafCIDs...)
	leafNodes = append(leafNodes, sliceLeafCIDs...)
	leafNodes = append(leafNodes, headLeafCID...)
	// TODO: fill in contract data `response.Leaves`

	response.MetaData.NodeStats["03-leaves"] = strconv.Itoa(len(leafNodes))
	response.MetaData.NodeStats["04-smart-contracts"] = "" // TODO: count # of contracts
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

func (b *Backend) getHeight(root common.Hash) (int, error) {
	var blockHeight int
	pgStr := `SELECT block_number
			FROM eth.header_cids
			WHERE state_root = $1`
	return blockHeight, b.DB.Get(&blockHeight, pgStr, root.String())
}

func (b Backend) getStateNodes(blockHeight int, paths [][]byte) (map[string]string, []cid.Cid, int, string, error) {
	nodes := make(map[string]string)
	fetchStart := makeTimestamp()

	// Get CIDs for all nodes at the provided paths
	leafCIDs, intermediateCIDs, deepestPath, err := b.getStateNodeCIDs(blockHeight, paths)
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

func (b *Backend) getStateNodeCIDs(blockHeight int, paths [][]byte) ([]cid.Cid, []cid.Cid, int, error) {
	tx, err := b.DB.Beginx()
	if err != nil {
		return nil, nil, 0, err
	}
	deepestPath := 0
	intermediateNodes := make([]cid.Cid, 0)
	leafNodes := make([]cid.Cid, 0)
	pgStr := `SELECT state_cids.cid, state_cids.node_type
			FROM eth.state_cids, eth.header_cids
			WHERE state_cids.header_id = header_cids.id
			AND header_cids.block_number <= $1
			AND state_cids.state_path = $2
			ORDER BY block_number DESC LIMIT 1`
	for _, path := range paths {
		var node nodeDBResponse
		if err := tx.Get(&node, pgStr, blockHeight, path); err != nil {
			if err == sql.ErrNoRows { // we will not find a node for each path
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

func (b *Backend) getStateHeadNode(blockHeight int, headPath []byte) (map[string]string, []cid.Cid, error) {
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
