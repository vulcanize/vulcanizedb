package eth_state_trie

import (
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipld-format"
)

type EthStateTrieNode struct {
	cid     cid.Cid
	rawdata []byte
}

func (estn EthStateTrieNode) RawData() []byte {
	return estn.rawdata
}

func (estn EthStateTrieNode) Cid() cid.Cid {
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
