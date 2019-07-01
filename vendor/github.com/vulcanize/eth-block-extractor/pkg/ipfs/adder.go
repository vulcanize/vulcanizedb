package ipfs

import ipld "github.com/ipfs/go-ipld-format"

type Adder interface {
	Add(node ipld.Node) error
}
