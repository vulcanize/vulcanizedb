package eth_block_header

import (
	"github.com/ethereum/go-ethereum/core/types"

	"gx/ipfs/QmWi2BYBL5gJ3CiAiQchg6rn1A8iBsrWy51EYxvHVjFvLb/go-ipld-format"
	"gx/ipfs/QmapdYm1b22Frv3k17fqrBYTFRxwiaVJkB299Mfn33edeB/go-cid"
)

type EthBlockHeaderNode struct {
	*types.Header

	cid     *cid.Cid
	rawdata []byte
}

func (ebh *EthBlockHeaderNode) RawData() []byte {
	return ebh.rawdata
}

func (ebh *EthBlockHeaderNode) Cid() *cid.Cid {
	return ebh.cid
}

func (EthBlockHeaderNode) String() string {
	return ""
}

func (EthBlockHeaderNode) Loggable() map[string]interface{} {
	panic("implement me")
}

func (EthBlockHeaderNode) Resolve(path []string) (interface{}, []string, error) {
	panic("implement me")
}

func (EthBlockHeaderNode) Tree(path string, depth int) []string {
	panic("implement me")
}

func (EthBlockHeaderNode) ResolveLink(path []string) (*format.Link, []string, error) {
	panic("implement me")
}

func (EthBlockHeaderNode) Copy() format.Node {
	panic("implement me")
}

func (EthBlockHeaderNode) Links() []*format.Link {
	panic("implement me")
}

func (EthBlockHeaderNode) Stat() (*format.NodeStat, error) {
	panic("implement me")
}

func (EthBlockHeaderNode) Size() (uint64, error) {
	panic("implement me")
}
