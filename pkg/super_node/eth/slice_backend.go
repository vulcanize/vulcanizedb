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

/*
import (
	"bytes"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
)

var pathSteps = []byte{'\x00', '\x01', '\x02', '\x03', '\x04', '\x05', '\x06', '\x07', '\x08', '\x09', '\x0a', '\x0b', '\x0c', '\x0d', '\x0e', '\x0f'}

// GetStateSlice is the backend method for constructing and returning state slice responses
func (b *Backend) GetStateSlice(path string, depth int, root common.Hash) (*GetSliceResponse, error) {
	response := new(GetSliceResponse)
	response.SliceID = fmt.Sprintf("%s-%d-%s", path, depth, root.String())

	// Get all the paths
	trieLoadingStart := makeTimestamp() // this doesn't really represent trie loading time for our process
	// Get the block height for this state root hash
	var blockHeight int
	pgStr := `SELECT block_number
			FROM eth.header_cids
			WHERE state_root = $1`
	if err := b.DB.Get(&blockHeight, pgStr, root.String()); err != nil {
		return nil, err
	}
	// Get head path and stem and slice paths
	headPath, stemPaths, slicePaths, err := getPaths(path, depth)
	if err != nil {
		return nil, err
	}
	response.MetaData.TimeStats["00-trie-loading"] =  strconv.Itoa(int(makeTimestamp() - trieLoadingStart))

	// Get the stem nodes
	stemFetchStart := makeTimestamp()
	stemCIDs, _, err := b.getStateNodeCIDs(blockHeight, stemPaths)
	if err != nil {
		return nil, err
	}
	stemIPLDs := b.Fetcher.fetchBatch(stemCIDs)
	if len(stemIPLDs) != len(stemCIDs) {
		return nil, fmt.Errorf("expected %d stem IPLDs, got %d", len(stemCIDs), len(stemIPLDs))
	}
	for _, stemIPLD := range stemIPLDs {
		decodedMh, err := multihash.Decode(stemIPLD.Cid().Hash())
		if err != nil {
			return nil, err
		}
		hash := crypto.Keccak256Hash(stemIPLD.RawData())
		if !bytes.Equal(hash.Bytes(), decodedMh.Digest) {
			panic("multihash digest should equal keccak of raw data")
		}
		response.TrieNodes.Stem[common.Bytes2Hex(decodedMh.Digest)] = common.Bytes2Hex(stemIPLD.RawData())
	}
	response.MetaData.TimeStats["01-fetch-stem-keys"] = strconv.Itoa(int(makeTimestamp() - stemFetchStart))

	// Get the slice nodes
	sliceFetchStart := makeTimestamp()
	sliceCIDs, deepestPath, err := b.getStateNodeCIDs(blockHeight, slicePaths)
	if err != nil {
		return nil, err
	}
	sliceIPLDs := b.Fetcher.fetchBatch(sliceCIDs)
	if len(sliceIPLDs) != len(sliceCIDs) {
		return nil, fmt.Errorf("expected %d slice IPLDs, got %d", len(sliceCIDs), len(sliceIPLDs))
	}
	for _, sliceIPLD := range sliceIPLDs {
		decodedMh, err := multihash.Decode(sliceIPLD.Cid().Hash())
		if err != nil {
			return nil, err
		}
		hash := crypto.Keccak256Hash(sliceIPLD.RawData())
		if !bytes.Equal(hash.Bytes(), decodedMh.Digest) {
			panic("multihash digest should equal keccak of raw data")
		}
		response.TrieNodes.SliceNodes[common.Bytes2Hex(decodedMh.Digest)] = common.Bytes2Hex(sliceIPLD.RawData())
	}
	response.MetaData.TimeStats["02-fetch-slice-keys"] = strconv.Itoa(int(makeTimestamp() - sliceFetchStart))

	// Get head node
	pgStr = `SELECT cid
			FROM eth.state_cids, eth.header_cids
			WHERE state_cids.header_id = header_cids.id
			AND header_cids.block_number <= $1
			AND state_cids.state_path = $2
			ORDER BY block_number DESC LIMIT 1`
	var cidStr string
	if err := b.DB.Get(&cidStr, pgStr, blockHeight, headPath); err != nil {
		return nil, err
	}
	headCID, err := cid.Decode(cidStr)
	if err != nil {
		return nil, err
	}
	headIPLD, err := b.Fetcher.fetch(headCID)
	if err != nil {
		return nil, err
	}
	decodedMh, err := multihash.Decode(headIPLD.RawData())
	if err != nil {
		return nil, err
	}
	hash := crypto.Keccak256Hash(headIPLD.RawData())
	if !bytes.Equal(hash.Bytes(), decodedMh.Digest) {
		panic("multihash digest should equal keccak of raw data")
	}
	response.TrieNodes.Head[common.Bytes2Hex(decodedMh.Digest)] = common.Bytes2Hex(headIPLD.RawData())

	// Get leafs
	leafFetchStart := makeTimestamp() // do we really need to split the leaf fetch into a seperate step just so we can time it...?
	allCids := append(append(stemCIDs, sliceCIDs...), headCID)

	// Fill in metadata
	response.MetaData.NodeStats["00-stem-and-head-nodes"] = strconv.Itoa(len(response.TrieNodes.Stem) + 1)
	response.MetaData.NodeStats["01-max-depth"] = strconv.Itoa(deepestPath - len(headPath))
	response.MetaData.NodeStats["02-total-trie-nodes"] = strconv.Itoa(len(response.TrieNodes.Stem) + len(response.TrieNodes.SliceNodes) + 1)
	response.MetaData.NodeStats["03-leaves"] = ""
	response.MetaData.NodeStats["04-smart-contracs"] = ""

	return response, nil
}

func (b *Backend) getStateNodeLeafs(cids []cid.Cid) ([]cid.Cid, error) {
	tx, err := b.DB.Beginx()
	if err != nil {
		return nil, err
	}
	leafCids := make([]cid.Cid, 0)
	pgStr := `SELECT cid
			FROM eth.state_cids
			WHERE state_cids.cid = $1
			AND state_cids.node_type = 2`
	for _, c := range cids {
		var cidStr string
		if err := tx.Get(&cidStr, pgStr, c.String()); err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			return nil, err
		}
		dc, err := cid.Decode(cidStr)
		if err != nil {
			return nil, err
		}
		leafCids = append(cids, dc)
	}
	return leafCids, nil
}

func (b *Backend) getStateNodeCIDs(blockHeight int, paths [][]byte) ([]cid.Cid, int, error) {
	tx, err := b.DB.Beginx()
	if err != nil {
		return nil, 0, err
	}
	deepestPath := 0
	cids := make([]cid.Cid, 0, len(paths))
	pgStr := `SELECT state_cids.cid
			FROM eth.state_cids, eth.header_cids
			WHERE state_cids.header_id = header_cids.id
			AND header_cids.block_number <= $1
			AND state_cids.state_path = $2
			ORDER BY block_number DESC LIMIT 1`
	for _, path := range paths {
		var cidStr string
		if err := tx.Get(&cidStr, pgStr, blockHeight, path); err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			return nil, 0, err
		}
		pathLen := len(path)
		if pathLen > deepestPath {
			deepestPath = pathLen
		}
		dc, err := cid.Decode(cidStr)
		if err != nil {
			return nil, 0, err
		}
		cids = append(cids, dc)
	}
	return cids, deepestPath, tx.Commit()
}

// GetStorageSlice is the backend method for constructing and returning storage slice responses
func (b *Backend) GetStorageSlice(path string, depth int, root common.Hash) (*GetSliceResponse, error) {
	response := new(GetSliceResponse)
	response.SliceID = fmt.Sprintf("%s-%d-%s", path, depth, root.String())

	// Get all the paths
	trieLoadingStart := makeTimestamp()
	res := struct{
		Height int `db:"block_number"`
		StateLeafKey string `db:"state_leaf_key"`
	}{}
	pgStr := `SELECT block_number, state_leaf_key
			FROM eth.header_cids, eth.state_cids, eth.state_accounts
			WHERE state_cids.header_id = header_cids.id
			AND state_accounts.state_id = state_cids.id
			AND state_accounts.storage_root = $1
			ORDER BY block_number DESC
			LIMIT 1`
	if err := b.DB.Get(&res, pgStr, root.String()); err != nil {
		return nil, err
	}
	// Get head path and stem and slice paths
	headPath, stemPaths, slicePaths, err := getPaths(path, depth)
	if err != nil {
		return nil, err
	}
	response.MetaData.TimeStats["00-trie-loading"] =  strconv.Itoa(int(makeTimestamp() - trieLoadingStart))
	stemFetchStart := makeTimestamp()
	stemCIDs, _, err := b.getStorageNodeCIDs(res.Height, res.StateLeafKey, stemPaths)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *Backend) getStorageNodeCIDs(blockHeight int, stateKey string, paths [][]byte) ([]cid.Cid, int, error) {
	tx, err := b.DB.Beginx()
	if err != nil {
		return nil, 0, err
	}
	deepestPath := 0
	cids := make([]cid.Cid, 0, len(paths))
	pgStr := `SELECT storage_cids.cid
			FROM eth.storage_cids, eth.state_cids, eth.header_cids
			WHERE storage_cids.state_id = state_cids.id
			AND state_cids.header_id = header_cids.id
			AND state_cids.state_leaf_key = $1
			AND header_cids.block_number <= $2
			AND storage_cids.storage_path = $3
			ORDER BY block_number DESC LIMIT 1`
	for _, path := range paths {
		var cidStr string
		if err := tx.Get(&cidStr, pgStr, stateKey, blockHeight, path); err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			return nil, 0, err
		}
		pathLen := len(path)
		if pathLen > deepestPath {
			deepestPath = pathLen
		}
		dc, err := cid.Decode(cidStr)
		if err != nil {
			return nil, 0, err
		}
		cids = append(cids, dc)
	}
	return cids, deepestPath, tx.Commit()
}

func getPaths(path string, depth int) ([]byte, [][]byte, [][]byte, error) {
	// Convert the head hex path to a decoded byte path
	headPath := keyStringToHexBytes(path)
	// Collect all of the stem paths
	pathLen := len(headPath)
	if pathLen > 64 { // max path len is 64
		return nil, nil, nil, fmt.Errorf("path length cannot exceed 64; got %d", pathLen)
	}
	maxDepth := 64 - pathLen
	if depth > maxDepth {
		return nil, nil, nil, fmt.Errorf("max depth for path %s is %d; got %d", path, maxDepth, depth)
	}
	stemPaths := make([][]byte, 0, pathLen)
	for i := 0; i < pathLen; i++ {
		stemPaths = append(stemPaths, headPath[:i])
	}
	// Generate all of the slice paths
	slicePaths := make([][]byte, 0, iPow(16, depth))
	makeSlicePaths(headPath, depth, &slicePaths)
	return headPath, stemPaths, slicePaths, nil
}

// iterative function to generate the set of slice paths
func makeSlicePaths(path []byte, depth int, slicePaths *[][]byte) {
	depth-- // decrement the depth
	nextPaths := make([][]byte, 16) // slice to hold the next 16 paths
	for i, step := range pathSteps {  // iterate through steps
		nextPath := append(path, step) // create next paths by adding steps to current path
		nextPaths[i] = nextPath
		newSlicePaths := append(*slicePaths, nextPath) // add next paths to the collection of all slice paths
		slicePaths = &newSlicePaths
	}
	if depth == 0 { // if depth has reach 0, return
		return
	}
	for _, nextPath := range nextPaths {  // if not, then we iterate over the next paths
		makeSlicePaths(nextPath, depth, slicePaths) // and repeat the process for each one
	}
}

// converts a hex string path to a decoded hex byte path
func keyStringToHexBytes(str string) []byte {
	path := common.Hex2Bytes(str)
	l := len(path)*2
	var nibbles = make([]byte, l)
	for i, b := range path {
		nibbles[i*2] = b / 16
		nibbles[i*2+1] = b % 16
	}
	return nibbles
}

// you have to write your own integer exponentiation function in Go...
func iPow(a, b int) int {
	result := 1
	for 0 != b {
		if 0 != (b & 1) {
			result *= a

		}
		b >>= 1
		a *= a
	}
	return result
}

// use to return timestamp in milliseconds
func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

*/
