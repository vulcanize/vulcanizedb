package eth_storage_trie

import (
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/util"
)

const (
	EthStorageTrieNodeCode = 0x98
)

type StorageTrieDagPutter struct {
	adder ipfs.Adder
}

func NewStorageTrieDagPutter(adder ipfs.Adder) *StorageTrieDagPutter {
	return &StorageTrieDagPutter{adder: adder}
}

func (stdp StorageTrieDagPutter) DagPut(raw interface{}) ([]string, error) {
	input := raw.([]byte)
	cid, err := util.RawToCid(EthStorageTrieNodeCode, input)
	if err != nil {
		return nil, err
	}
	node := &EthStorageTrieNode{
		cid:     cid,
		rawdata: input,
	}
	err = stdp.adder.Add(node)
	if err != nil {
		return nil, err
	}
	return []string{node.Cid().String()}, nil
}
