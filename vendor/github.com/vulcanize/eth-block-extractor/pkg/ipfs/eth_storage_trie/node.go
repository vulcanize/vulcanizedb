package eth_storage_trie

import (
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipld-format"
)

type EthStorageTrieNode struct {
	cid     cid.Cid
	rawdata []byte
}

func (estn *EthStorageTrieNode) RawData() []byte {
	return estn.rawdata
}

func (estn *EthStorageTrieNode) Cid() cid.Cid {
	return estn.cid
}

func (*EthStorageTrieNode) String() string {
	panic("implement me")
}

func (*EthStorageTrieNode) Loggable() map[string]interface{} {
	panic("implement me")
}

func (*EthStorageTrieNode) Resolve(path []string) (interface{}, []string, error) {
	panic("implement me")
}

func (*EthStorageTrieNode) Tree(path string, depth int) []string {
	panic("implement me")
}

func (*EthStorageTrieNode) ResolveLink(path []string) (*format.Link, []string, error) {
	panic("implement me")
}

func (*EthStorageTrieNode) Copy() format.Node {
	panic("implement me")
}

func (*EthStorageTrieNode) Links() []*format.Link {
	panic("implement me")
}

func (*EthStorageTrieNode) Stat() (*format.NodeStat, error) {
	panic("implement me")
}

func (*EthStorageTrieNode) Size() (uint64, error) {
	panic("implement me")
}
