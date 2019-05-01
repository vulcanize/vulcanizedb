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

package ipfs

import (
	"context"

	"github.com/ipfs/go-block-format"
	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-cid"
)

// IPLDFetcher is the interface for fetching IPLD objects from IPFS
type IPLDFetcher interface {
	Fetch(cid cid.Cid) (blocks.Block, error)
	FetchBatch(cids []cid.Cid) []blocks.Block
}

// Fetcher is the underlying struct which supports the IPLDFetcher interface
type Fetcher struct {
	BlockService blockservice.BlockService
}

// NewIPLDFetcher creates a pointer to a new Fetcher which satisfies the IPLDFetcher interface
func NewIPLDFetcher(ipfsPath string) (*Fetcher, error) {
	blockService, err := InitIPFSBlockService(ipfsPath)
	if err != nil {
		return nil, err
	}
	return &Fetcher{
		BlockService: blockService,
	}, nil
}

// Fetch is used to fetch a batch of IPFS data blocks by cid
func (f *Fetcher) Fetch(cid cid.Cid) (blocks.Block, error) {
	return f.BlockService.GetBlock(context.Background(), cid)
}

// FetchBatch is used to fetch a batch of IPFS data blocks by cid
// There is no guarantee all are fetched, and no error in such a case, so
// downstream we will need to confirm which CIDs were fetched in the result set
func (f *Fetcher) FetchBatch(cids []cid.Cid) []blocks.Block {
	fetchedBlocks := make([]blocks.Block, 0, len(cids))
	blockChan := f.BlockService.GetBlocks(context.Background(), cids)
	for block := range blockChan {
		fetchedBlocks = append(fetchedBlocks, block)
	}
	return fetchedBlocks
}
