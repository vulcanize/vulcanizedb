package eth_block_transactions

import (
	"github.com/ethereum/go-ethereum/core/types"
	"gx/ipfs/QmWi2BYBL5gJ3CiAiQchg6rn1A8iBsrWy51EYxvHVjFvLb/go-ipld-format"
	"gx/ipfs/QmapdYm1b22Frv3k17fqrBYTFRxwiaVJkB299Mfn33edeB/go-cid"
)

type EthTransactionNode struct {
	*types.Transaction

	cid     *cid.Cid
	rawdata []byte
}

func (etn *EthTransactionNode) RawData() []byte {
	return etn.rawdata
}

func (etn *EthTransactionNode) Cid() *cid.Cid {
	return etn.cid
}

func (EthTransactionNode) String() string {
	return ""
}

func (EthTransactionNode) Loggable() map[string]interface{} {
	panic("implement me")
}

func (EthTransactionNode) Resolve(path []string) (interface{}, []string, error) {
	panic("implement me")
}

func (EthTransactionNode) Tree(path string, depth int) []string {
	panic("implement me")
}

func (EthTransactionNode) ResolveLink(path []string) (*format.Link, []string, error) {
	panic("implement me")
}

func (EthTransactionNode) Copy() format.Node {
	panic("implement me")
}

func (EthTransactionNode) Links() []*format.Link {
	panic("implement me")
}

func (EthTransactionNode) Stat() (*format.NodeStat, error) {
	panic("implement me")
}

func (EthTransactionNode) Size() (uint64, error) {
	panic("implement me")
}
