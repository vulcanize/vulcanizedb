package eth_block_header

import (
	"github.com/ethereum/go-ethereum/core/types"
	ipld "gx/ipfs/QmWi2BYBL5gJ3CiAiQchg6rn1A8iBsrWy51EYxvHVjFvLb/go-ipld-format"

	"github.com/vulcanize/eth-block-extractor/pkg/ipfs"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/util"
	"github.com/vulcanize/eth-block-extractor/pkg/wrappers/rlp"
)

const (
	EthBlockHeaderCode = 0x90
)

type BlockHeaderDagPutter struct {
	adder   ipfs.Adder
	decoder rlp.Decoder
}

func NewBlockHeaderDagPutter(adder ipfs.Adder, decoder rlp.Decoder) *BlockHeaderDagPutter {
	return &BlockHeaderDagPutter{adder: adder, decoder: decoder}
}

func (bhdp *BlockHeaderDagPutter) DagPut(raw interface{}) ([]string, error) {
	input := raw.([]byte)
	nd, err := bhdp.getNodeForBlockHeader(input)
	if err != nil {
		return nil, err
	}
	err = bhdp.adder.Add(nd)
	if err != nil {
		return nil, err
	}
	return []string{nd.Cid().String()}, nil
}

func (bhdp *BlockHeaderDagPutter) getNodeForBlockHeader(raw []byte) (ipld.Node, error) {
	var blockHeader types.Header
	err := bhdp.decoder.Decode(raw, &blockHeader)
	if err != nil {
		return nil, err
	}
	blockHeaderCid, err := util.RawToCid(EthBlockHeaderCode, raw)
	if err != nil {
		return nil, err
	}
	return &EthBlockHeaderNode{
		Header:  &blockHeader,
		cid:     blockHeaderCid,
		rawdata: raw,
	}, nil
}
