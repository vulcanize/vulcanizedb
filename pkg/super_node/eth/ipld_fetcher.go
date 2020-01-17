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
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ipfs/go-block-format"
	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-cid"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
)

var (
	errUnexpectedNumberOfIPLDs = errors.New("ipfs batch fetch returned unexpected number of IPLDs")
)

// IPLDFetcher satisfies the IPLDFetcher interface for ethereum
type IPLDFetcher struct {
	BlockService blockservice.BlockService
}

// NewIPLDFetcher creates a pointer to a new IPLDFetcher
func NewIPLDFetcher(ipfsPath string) (*IPLDFetcher, error) {
	blockService, err := ipfs.InitIPFSBlockService(ipfsPath)
	if err != nil {
		return nil, err
	}
	return &IPLDFetcher{
		BlockService: blockService,
	}, nil
}

// Fetch is the exported method for fetching and returning all the IPLDS specified in the CIDWrapper
func (f *IPLDFetcher) Fetch(cids interface{}) (interface{}, error) {
	cidWrapper, ok := cids.(*CIDWrapper)
	if !ok {
		return nil, fmt.Errorf("eth fetcher: expected cids type %T got %T", &CIDWrapper{}, cids)
	}
	log.Debug("fetching iplds")
	iplds := new(IPLDWrapper)
	iplds.BlockNumber = cidWrapper.BlockNumber
	var err error
	iplds.Headers, err = f.FetchHeaders(cidWrapper.Headers)
	if err != nil {
		return nil, err
	}
	iplds.Uncles, err = f.FetchUncles(cidWrapper.Uncles)
	if err != nil {
		return nil, err
	}
	iplds.Transactions, err = f.FetchTrxs(cidWrapper.Transactions)
	if err != nil {
		return nil, err
	}
	iplds.Receipts, err = f.FetchRcts(cidWrapper.Receipts)
	if err != nil {
		return nil, err
	}
	iplds.StateNodes, err = f.FetchState(cidWrapper.StateNodes)
	if err != nil {
		return nil, err
	}
	iplds.StorageNodes, err = f.FetchStorage(cidWrapper.StorageNodes)
	if err != nil {
		return nil, err
	}
	return iplds, nil
}

// FetchHeaders fetches headers
// It uses the f.fetchBatch method
func (f *IPLDFetcher) FetchHeaders(cids []HeaderModel) ([]blocks.Block, error) {
	log.Debug("fetching header iplds")
	headerCids := make([]cid.Cid, 0, len(cids))
	for _, c := range cids {
		dc, err := cid.Decode(c.CID)
		if err != nil {
			return nil, err
		}
		headerCids = append(headerCids, dc)
	}
	headers := f.fetchBatch(headerCids)
	if len(headers) != len(headerCids) {
		log.Errorf("ipfs fetcher: number of header blocks returned (%d) does not match number expected (%d)", len(headers), len(headerCids))
		return headers, errUnexpectedNumberOfIPLDs
	}
	return headers, nil
}

// FetchUncles fetches uncles
// It uses the f.fetchBatch method
func (f *IPLDFetcher) FetchUncles(cids []HeaderModel) ([]blocks.Block, error) {
	log.Debug("fetching uncle iplds")
	uncleCids := make([]cid.Cid, 0, len(cids))
	for _, c := range cids {
		dc, err := cid.Decode(c.CID)
		if err != nil {
			return nil, err
		}
		uncleCids = append(uncleCids, dc)
	}
	uncles := f.fetchBatch(uncleCids)
	if len(uncles) != len(uncleCids) {
		log.Errorf("ipfs fetcher: number of uncle blocks returned (%d) does not match number expected (%d)", len(uncles), len(uncleCids))
		return uncles, errUnexpectedNumberOfIPLDs
	}
	return uncles, nil
}

// FetchTrxs fetches transactions
// It uses the f.fetchBatch method
func (f *IPLDFetcher) FetchTrxs(cids []TxModel) ([]blocks.Block, error) {
	log.Debug("fetching transaction iplds")
	trxCids := make([]cid.Cid, 0, len(cids))
	for _, c := range cids {
		dc, err := cid.Decode(c.CID)
		if err != nil {
			return nil, err
		}
		trxCids = append(trxCids, dc)
	}
	trxs := f.fetchBatch(trxCids)
	if len(trxs) != len(trxCids) {
		log.Errorf("ipfs fetcher: number of transaction blocks returned (%d) does not match number expected (%d)", len(trxs), len(trxCids))
		return trxs, errUnexpectedNumberOfIPLDs
	}
	return trxs, nil
}

// FetchRcts fetches receipts
// It uses the f.fetchBatch method
func (f *IPLDFetcher) FetchRcts(cids []ReceiptModel) ([]blocks.Block, error) {
	log.Debug("fetching receipt iplds")
	rctCids := make([]cid.Cid, 0, len(cids))
	for _, c := range cids {
		dc, err := cid.Decode(c.CID)
		if err != nil {
			return nil, err
		}
		rctCids = append(rctCids, dc)
	}
	rcts := f.fetchBatch(rctCids)
	if len(rcts) != len(rctCids) {
		log.Errorf("ipfs fetcher: number of receipt blocks returned (%d) does not match number expected (%d)", len(rcts), len(rctCids))
		return rcts, errUnexpectedNumberOfIPLDs
	}
	return rcts, nil
}

// FetchState fetches state nodes
// It uses the single f.fetch method instead of the batch fetch, because it
// needs to maintain the data's relation to state keys
func (f *IPLDFetcher) FetchState(cids []StateNodeModel) (map[common.Hash]blocks.Block, error) {
	log.Debug("fetching state iplds")
	stateNodes := make(map[common.Hash]blocks.Block)
	for _, stateNode := range cids {
		if stateNode.CID == "" || stateNode.StateKey == "" {
			continue
		}
		dc, err := cid.Decode(stateNode.CID)
		if err != nil {
			return nil, err
		}
		state, err := f.fetch(dc)
		if err != nil {
			return nil, err
		}
		stateNodes[common.HexToHash(stateNode.StateKey)] = state
	}
	return stateNodes, nil
}

// FetchStorage fetches storage nodes
// It uses the single f.fetch method instead of the batch fetch, because it
// needs to maintain the data's relation to state and storage keys
func (f *IPLDFetcher) FetchStorage(cids []StorageNodeWithStateKeyModel) (map[common.Hash]map[common.Hash]blocks.Block, error) {
	log.Debug("fetching storage iplds")
	storageNodes := make(map[common.Hash]map[common.Hash]blocks.Block)
	for _, storageNode := range cids {
		if storageNode.CID == "" || storageNode.StorageKey == "" || storageNode.StateKey == "" {
			continue
		}
		dc, err := cid.Decode(storageNode.CID)
		if err != nil {
			return nil, err
		}
		storage, err := f.fetch(dc)
		if err != nil {
			return nil, err
		}
		if storageNodes[common.HexToHash(storageNode.StateKey)] == nil {
			storageNodes[common.HexToHash(storageNode.StateKey)] = make(map[common.Hash]blocks.Block)
		}
		storageNodes[common.HexToHash(storageNode.StateKey)][common.HexToHash(storageNode.StorageKey)] = storage
	}
	return storageNodes, nil
}

// fetch is used to fetch a single cid
func (f *IPLDFetcher) fetch(cid cid.Cid) (blocks.Block, error) {
	return f.BlockService.GetBlock(context.Background(), cid)
}

// fetchBatch is used to fetch a batch of IPFS data blocks by cid
// There is no guarantee all are fetched, and no error in such a case, so
// downstream we will need to confirm which CIDs were fetched in the result set
func (f *IPLDFetcher) fetchBatch(cids []cid.Cid) []blocks.Block {
	fetchedBlocks := make([]blocks.Block, 0, len(cids))
	blockChan := f.BlockService.GetBlocks(context.Background(), cids)
	for block := range blockChan {
		fetchedBlocks = append(fetchedBlocks, block)
	}
	return fetchedBlocks
}
