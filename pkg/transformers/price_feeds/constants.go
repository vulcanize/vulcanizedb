package price_feeds

import (
	"errors"
	"math/big"
)

var (
	ErrMultipleLogs  = errors.New("multiple matching logs found in block")
	ErrNoMatchingLog = errors.New("no matching log")
	Ether            = big.NewFloat(1e18)
	PeekMethodName   = "peek"
	PepAddress       = "0x99041F808D598B782D5a3e498681C2452A31da08"
	PepLogTopic0     = "0x296ba4ca62c6c21c95e828080cb8aec7481b71390585605300a8a76f9e95b527"
	PipAddress       = "0x729D19f657BD0614b4985Cf1D82531c67569197B"
	PipLogTopic0     = "0x1817835800000000000000000000000000000000000000000000000000000000"
	PipMedianizerABI = `[{"constant":false,"inputs":[{"name":"owner_","type":"address"}],"name":"setOwner","outputs":[],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"","type":"bytes32"}],"name":"poke","outputs":[],"payable":false,"type":"function"},{"constant":false,"inputs":[],"name":"poke","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"compute","outputs":[{"name":"","type":"bytes32"},{"name":"","type":"bool"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"wat","type":"address"}],"name":"set","outputs":[],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"wat","type":"address"}],"name":"unset","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"indexes","outputs":[{"name":"","type":"bytes12"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"next","outputs":[{"name":"","type":"bytes12"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"read","outputs":[{"name":"","type":"bytes32"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"peek","outputs":[{"name":"","type":"bytes32"},{"name":"","type":"bool"}],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"","type":"bytes12"}],"name":"values","outputs":[{"name":"","type":"address"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"min_","type":"uint96"}],"name":"setMin","outputs":[],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"authority_","type":"address"}],"name":"setAuthority","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"type":"function"},{"constant":false,"inputs":[],"name":"void","outputs":[],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"pos","type":"bytes12"},{"name":"wat","type":"address"}],"name":"set","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"authority","outputs":[{"name":"","type":"address"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"pos","type":"bytes12"}],"name":"unset","outputs":[],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"next_","type":"bytes12"}],"name":"setNext","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"min","outputs":[{"name":"","type":"uint96"}],"payable":false,"type":"function"},{"anonymous":true,"inputs":[{"indexed":true,"name":"sig","type":"bytes4"},{"indexed":true,"name":"guy","type":"address"},{"indexed":true,"name":"foo","type":"bytes32"},{"indexed":true,"name":"bar","type":"bytes32"},{"indexed":false,"name":"wad","type":"uint256"},{"indexed":false,"name":"fax","type":"bytes"}],"name":"LogNote","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"authority","type":"address"}],"name":"LogSetAuthority","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"}],"name":"LogSetOwner","type":"event"}]`
	Ray              = big.NewFloat(1e27)
	RepAddress       = "0xF5f94b7F9De14D43112e713835BCef2d55b76c1C"
	RepLogTopic0     = "0x296ba4ca62c6c21c95e828080cb8aec7481b71390585605300a8a76f9e95b527"
)
