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
	"context"
	"errors"
	"fmt"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"

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
func (f *IPLDFetcher) Fetch(cids shared.CIDsForFetching) (shared.IPLDs, error) {
	cidWrapper, ok := cids.(*CIDWrapper)
	if !ok {
		return nil, fmt.Errorf("btc fetcher: expected cids type %T got %T", &CIDWrapper{}, cids)
	}
	log.Debug("fetching iplds")
	iplds := IPLDs{}
	iplds.BlockNumber = cidWrapper.BlockNumber
	var err error
	iplds.Headers, err = f.FetchHeaders(cidWrapper.Headers)
	if err != nil {
		return nil, err
	}
	iplds.Transactions, err = f.FetchTrxs(cidWrapper.Transactions)
	if err != nil {
		return nil, err
	}
	return iplds, nil
}

// FetchHeaders fetches headers
// It uses the f.fetchBatch method
func (f *IPLDFetcher) FetchHeaders(cids []HeaderModel) ([][]byte, error) {
	log.Debug("fetching header iplds")
	headerCids := make([]cid.Cid, len(cids))
	for i, c := range cids {
		dc, err := cid.Decode(c.CID)
		if err != nil {
			return nil, err
		}
		headerCids[i] = dc
	}
	headers := f.fetchBatch(headerCids)
	headersBytes := make([][]byte, len(headers))
	for i, header := range headers {
		headersBytes[i] = header.RawData()
	}
	if len(headersBytes) != len(headerCids) {
		log.Errorf("ipfs fetcher: number of header blocks returned (%d) does not match number expected (%d)", len(headers), len(headerCids))
		return headersBytes, errUnexpectedNumberOfIPLDs
	}
	return headersBytes, nil
}

// FetchTrxs fetches transactions
// It uses the f.fetchBatch method
func (f *IPLDFetcher) FetchTrxs(cids []TxModel) ([][]byte, error) {
	log.Debug("fetching transaction iplds")
	trxCids := make([]cid.Cid, len(cids))
	for i, c := range cids {
		dc, err := cid.Decode(c.CID)
		if err != nil {
			return nil, err
		}
		trxCids[i] = dc
	}
	trxs := f.fetchBatch(trxCids)
	trxBytes := make([][]byte, len(trxs))
	for i, trx := range trxs {
		trxBytes[i] = trx.RawData()
	}
	if len(trxBytes) != len(trxCids) {
		log.Errorf("ipfs fetcher: number of transaction blocks returned (%d) does not match number expected (%d)", len(trxs), len(trxCids))
		return trxBytes, errUnexpectedNumberOfIPLDs
	}
	return trxBytes, nil
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
