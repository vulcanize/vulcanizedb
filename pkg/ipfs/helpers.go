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

	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
)

// InitIPFSBlockService is used to configure and return a BlockService using an ipfs repo path (e.g. ~/.ipfs)
func InitIPFSBlockService(ipfsPath string) (blockservice.BlockService, error) {
	r, err := fsrepo.Open(ipfsPath)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	cfg := &core.BuildCfg{
		Online: false,
		Repo:   r,
	}
	ipfsNode, err := core.NewNode(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return ipfsNode.Blocks, nil
}
