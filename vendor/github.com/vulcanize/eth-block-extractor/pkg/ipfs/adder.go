package ipfs

import ipld "gx/ipfs/QmWi2BYBL5gJ3CiAiQchg6rn1A8iBsrWy51EYxvHVjFvLb/go-ipld-format"

type Adder interface {
	Add(node ipld.Node) error
}
