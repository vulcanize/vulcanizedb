package eth_state_trie

import (
	"gx/ipfs/QmWi2BYBL5gJ3CiAiQchg6rn1A8iBsrWy51EYxvHVjFvLb/go-ipld-format"
	"gx/ipfs/QmapdYm1b22Frv3k17fqrBYTFRxwiaVJkB299Mfn33edeB/go-cid"
)

type EthStateTrieNode struct {
	cid     *cid.Cid
	rawdata []byte
}

func (estn EthStateTrieNode) RawData() []byte {
	return estn.rawdata
}

func (estn EthStateTrieNode) Cid() *cid.Cid {
	return estn.cid
}

func (EthStateTrieNode) String() string {
	panic("implement me")
}

func (EthStateTrieNode) Loggable() map[string]interface{} {
	panic("implement me")
}

func (EthStateTrieNode) Resolve(path []string) (interface{}, []string, error) {
	panic("implement me")
}

func (EthStateTrieNode) Tree(path string, depth int) []string {
	panic("implement me")
}

func (EthStateTrieNode) ResolveLink(path []string) (*format.Link, []string, error) {
	panic("implement me")
}

func (EthStateTrieNode) Copy() format.Node {
	panic("implement me")
}

func (EthStateTrieNode) Links() []*format.Link {
	panic("implement me")
}

func (EthStateTrieNode) Stat() (*format.NodeStat, error) {
	panic("implement me")
}

func (EthStateTrieNode) Size() (uint64, error) {
	panic("implement me")
}
