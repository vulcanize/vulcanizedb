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
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ipfs/go-block-format"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/config"
)

// APIName is the namespace for the super node's eth api
const APIName = "eth"

// APIVersion is the version of the super node's eth api
const APIVersion = "0.0.1"

type PublicEthAPI struct {
	b *Backend
}

// NewPublicEthAPI creates a new PublicEthAPI with the provided underlying Backend
func NewPublicEthAPI(b *Backend) *PublicEthAPI {
	return &PublicEthAPI{
		b: b,
	}
}

// BlockNumber returns the block number of the chain head.
func (pea *PublicEthAPI) BlockNumber() hexutil.Uint64 {
	number, _ := pea.b.retriever.RetrieveLastBlockNumber()
	return hexutil.Uint64(number)
}

// GetLogs returns logs matching the given argument that are stored within the state.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getlogs
func (pea *PublicEthAPI) GetLogs(ctx context.Context, crit ethereum.FilterQuery) ([]*types.Log, error) {
	// Convert FilterQuery into ReceiptFilter
	addrStrs := make([]string, len(crit.Addresses))
	for i, addr := range crit.Addresses {
		addrStrs[i] = addr.String()
	}
	topicStrSets := make([][]string, 4)
	for i, topicSet := range crit.Topics {
		if i > 3 {
			break
		}
		for _, topic := range topicSet {
			topicStrSets[i] = append(topicStrSets[i], topic.String())
		}
	}
	filter := config.ReceiptFilter{
		Contracts: addrStrs,
		Topics:    topicStrSets,
	}
	tx, err := pea.b.db.Beginx()
	if err != nil {
		return nil, err
	}
	// If we have a blockhash to filter on, fire off single retrieval query
	if crit.BlockHash != nil {
		rctCIDs, err := pea.b.retriever.RetrieveRctCIDs(tx, filter, 0, crit.BlockHash, nil)
		if err != nil {
			return nil, err
		}
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		rctIPLDs, err := pea.b.fetcher.FetchRcts(rctCIDs)
		if err != nil {
			return nil, err
		}
		return extractLogsOfInterest(rctIPLDs, filter.Topics)
	}
	// Otherwise, create block range from criteria
	// nil values are filled in; to request a single block have both ToBlock and FromBlock equal that number
	startingBlock := crit.FromBlock
	endingBlock := crit.ToBlock
	if startingBlock == nil {
		startingBlockInt, err := pea.b.retriever.RetrieveFirstBlockNumber()
		if err != nil {
			return nil, err
		}
		startingBlock = big.NewInt(startingBlockInt)
	}
	if endingBlock == nil {
		endingBlockInt, err := pea.b.retriever.RetrieveLastBlockNumber()
		if err != nil {
			return nil, err
		}
		endingBlock = big.NewInt(endingBlockInt)
	}
	start := startingBlock.Int64()
	end := endingBlock.Int64()
	allRctCIDs := make([]ReceiptModel, 0)
	for i := start; i <= end; i++ {
		rctCIDs, err := pea.b.retriever.RetrieveRctCIDs(tx, filter, i, nil, nil)
		if err != nil {
			return nil, err
		}
		allRctCIDs = append(allRctCIDs, rctCIDs...)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	rctIPLDs, err := pea.b.fetcher.FetchRcts(allRctCIDs)
	if err != nil {
		return nil, err
	}
	return extractLogsOfInterest(rctIPLDs, filter.Topics)
}

func extractLogsOfInterest(rctIPLDs []blocks.Block, wantedTopics [][]string) ([]*types.Log, error) {
	var logs []*types.Log
	for _, rctIPLD := range rctIPLDs {
		rctRLP := rctIPLD.RawData()
		var rct types.ReceiptForStorage
		if err := rlp.DecodeBytes(rctRLP, &rct); err != nil {
			return nil, err
		}
		for _, log := range rct.Logs {
			if wanted := wantedLog(wantedTopics, log.Topics); wanted == true {
				logs = append(logs, log)
			}
		}
	}
	return logs, nil
}

// returns true if the log matches on the filter
func wantedLog(wantedTopics [][]string, actualTopics []common.Hash) bool {
	// actualTopics will always have length <= 4
	// wantedTopics will always have length == 4
	matches := 0
	for i, actualTopic := range actualTopics {
		// If we have topics in this filter slot, count as a match if the actualTopic matches one of the ones in this filter slot
		if len(wantedTopics[i]) > 0 {
			matches += sliceContainsHash(wantedTopics[i], actualTopic)
		} else {
			// Filter slot is empty, not matching any topics at this slot => counts as a match
			matches++
		}
	}
	if matches == len(actualTopics) {
		return true
	}
	return false
}

// returns 1 if the slice contains the hash, 0 if it does not
func sliceContainsHash(slice []string, hash common.Hash) int {
	for _, str := range slice {
		if str == hash.String() {
			return 1
		}
	}
	return 0
}

// GetHeaderByNumber returns the requested canonical block header.
// When blockNr is -1 the chain head is returned.
// We cannot support pending block calls since we do not have an active miner
func (pea *PublicEthAPI) GetHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (map[string]interface{}, error) {
	header, err := pea.b.HeaderByNumber(ctx, number)
	if header != nil && err == nil {
		return pea.rpcMarshalHeader(header)
	}
	return nil, err
}

// rpcMarshalHeader uses the generalized output filler, then adds the total difficulty field, which requires
// a `PublicEthAPI`.
func (pea *PublicEthAPI) rpcMarshalHeader(header *types.Header) (map[string]interface{}, error) {
	fields := RPCMarshalHeader(header)
	td, err := pea.b.GetTd(header.Hash())
	if err != nil {
		return nil, err
	}
	fields["totalDifficulty"] = (*hexutil.Big)(td)
	return fields, nil
}

// RPCMarshalHeader converts the given header to the RPC output .
func RPCMarshalHeader(head *types.Header) map[string]interface{} {
	return map[string]interface{}{
		"number":           (*hexutil.Big)(head.Number),
		"hash":             head.Hash(),
		"parentHash":       head.ParentHash,
		"nonce":            head.Nonce,
		"mixHash":          head.MixDigest,
		"sha3Uncles":       head.UncleHash,
		"logsBloom":        head.Bloom,
		"stateRoot":        head.Root,
		"miner":            head.Coinbase,
		"difficulty":       (*hexutil.Big)(head.Difficulty),
		"extraData":        hexutil.Bytes(head.Extra),
		"size":             hexutil.Uint64(head.Size()),
		"gasLimit":         hexutil.Uint64(head.GasLimit),
		"gasUsed":          hexutil.Uint64(head.GasUsed),
		"timestamp":        hexutil.Uint64(head.Time),
		"transactionsRoot": head.TxHash,
		"receiptsRoot":     head.ReceiptHash,
	}
}
