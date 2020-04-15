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
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-cid"
	api "github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/interface-go-ipfs-core/options"
	ma "github.com/multiformats/go-multiaddr"
)

type Adder interface {
	Add(node ipld.Node) error
}

// InitIPFSPlugins is used to initialized IPFS plugins before creating a new IPFS node
// This should only be called once
func InitIPFSPlugins() error {
	l, err := loader.NewPluginLoader("")
	if err != nil {
		return err
	}
	err = l.Initialize()
	if err != nil {
		return err
	}
	return l.Inject()
}

// InitIPFSBlockService is used to configure and return a BlockService using an ipfs repo path (e.g. ~/.ipfs)
func InitIPFSBlockService(ipfsPath string) (blockservice.BlockService, error) {
	r, openErr := fsrepo.Open(ipfsPath)
	if openErr != nil {
		return nil, openErr
	}
	ctx := context.Background()
	cfg := &core.BuildCfg{
		Online: false,
		Repo:   r,
	}
	ipfsNode, newNodeErr := core.NewNode(ctx, cfg)
	if newNodeErr != nil {
		return nil, newNodeErr
	}
	return ipfsNode.Blocks, nil
}

type IPFS struct {
	n   *core.IpfsNode
	ctx context.Context
}

func (ipfs *IPFS) Add(node ipld.Node) error {
	return ipfs.n.DAG.Add(ipfs.n.Context(), node)
}

func InitIPFSNode(repoPath string) (*IPFS, error) {
	r, err := fsrepo.Open(repoPath)
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
	return &IPFS{n: ipfsNode, ctx: ctx}, nil
}

type IPFSClient struct {
	c   *api.HttpApi
	ctx context.Context
}

func (ipfs *IPFSClient) Add(node ipld.Node) error {
	c := node.Cid()
	prefix := c.Prefix()
	format := cid.CodecToStr[prefix.Codec]
	if prefix.Version == 0 {
		format = "v0"
	}
	res, err := ipfs.c.Block().Put(ipfs.ctx, bytes.NewReader(node.RawData()),
		options.Block.Hash(prefix.MhType, prefix.MhLength),
		options.Block.Format(format),
		options.Block.Pin(true))
	if err != nil {
		return err
	}
	if !res.Path().Cid().Equals(c) {
		return fmt.Errorf("cids didn't match - local %s, remote %s", c.String(), res.Path().Cid().String())
	}
	return nil
}

func InitIPFSClient(multiAddr ma.Multiaddr) (*IPFSClient, error) {
	c, err := api.NewApiWithClient(multiAddr, &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyFromEnvironment,
			DisableKeepAlives: true,
		},
	})
	if err != nil {
		return nil, err
	}
	return &IPFSClient{
		c:   c,
		ctx: context.Background(),
	}, nil
}
