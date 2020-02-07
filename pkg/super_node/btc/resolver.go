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
	"fmt"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"

	"github.com/ipfs/go-block-format"
)

// IPLDResolver satisfies the IPLDResolver interface for bitcoin
type IPLDResolver struct{}

// NewIPLDResolver returns a pointer to an IPLDResolver which satisfies the IPLDResolver interface
func NewIPLDResolver() *IPLDResolver {
	return &IPLDResolver{}
}

// Resolve is the exported method for resolving all of the BTC IPLDs packaged in an IpfsBlockWrapper
func (eir *IPLDResolver) Resolve(iplds shared.FetchedIPLDs) (shared.ServerResponse, error) {
	ipfsBlocks, ok := iplds.(*IPLDWrapper)
	if !ok {
		return StreamResponse{}, fmt.Errorf("eth resolver expected iplds type %T got %T", &IPLDWrapper{}, iplds)
	}
	return StreamResponse{
		BlockNumber:       ipfsBlocks.BlockNumber,
		SerializedHeaders: eir.resolve(ipfsBlocks.Headers),
		SerializedTxs:     eir.resolve(ipfsBlocks.Transactions),
	}, nil
}

func (eir *IPLDResolver) resolve(iplds []blocks.Block) [][]byte {
	rlps := make([][]byte, 0, len(iplds))
	for _, ipld := range iplds {
		rlps = append(rlps, ipld.RawData())
	}
	return rlps
}
