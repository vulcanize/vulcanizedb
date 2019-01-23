// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package constants

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Basic abi needed to check which interfaces are adhered to
var SupportsInterfaceABI = `[{"constant":true,"inputs":[{"name":"interfaceID","type":"bytes4"}],"name":"supportsInterface","outputs":[{"name":"","type":"bool"}],"payable":false,"type":"function"}]`

// Individual event interfaces for constructing ABI from
var SupportsInterace = `{"constant":true,"inputs":[{"name":"interfaceID","type":"bytes4"}],"name":"supportsInterface","outputs":[{"name":"","type":"bool"}],"payable":false,"type":"function"}`
var AddrChangeInterface = `{"anonymous":false,"inputs":[{"indexed":true,"name":"node","type":"bytes32"},{"indexed":false,"name":"a","type":"address"}],"name":"AddrChanged","type":"event"}`
var ContentChangeInterface = `{"anonymous":false,"inputs":[{"indexed":true,"name":"node","type":"bytes32"},{"indexed":false,"name":"hash","type":"bytes32"}],"name":"ContentChanged","type":"event"}`
var NameChangeInterface = `{"anonymous":false,"inputs":[{"indexed":true,"name":"node","type":"bytes32"},{"indexed":false,"name":"name","type":"string"}],"name":"NameChanged","type":"event"}`
var AbiChangeInterface = `{"anonymous":false,"inputs":[{"indexed":true,"name":"node","type":"bytes32"},{"indexed":true,"name":"contentType","type":"uint256"}],"name":"ABIChanged","type":"event"}`
var PubkeyChangeInterface = `{"anonymous":false,"inputs":[{"indexed":true,"name":"node","type":"bytes32"},{"indexed":false,"name":"x","type":"bytes32"},{"indexed":false,"name":"y","type":"bytes32"}],"name":"PubkeyChanged","type":"event"}`
var TextChangeInterface = `{"anonymous":false,"inputs":[{"indexed":true,"name":"node","type":"bytes32"},{"indexed":false,"name":"indexedKey","type":"string"},{"indexed":false,"name":"key","type":"string"}],"name":"TextChanged","type":"event"}`
var MultihashChangeInterface = `{"anonymous":false,"inputs":[{"indexed":true,"name":"node","type":"bytes32"},{"indexed":false,"name":"hash","type":"bytes"}],"name":"MultihashChanged","type":"event"}`
var ContenthashChangeInterface = `{"anonymous":false,"inputs":[{"indexed":true,"name":"node","type":"bytes32"},{"indexed":false,"name":"hash","type":"bytes"}],"name":"ContenthashChanged","type":"event"}`

var StartingBlock = int64(3648359)

// Resolver interface signatures
type Interface int

const (
	MetaSig Interface = iota
	AddrChangeSig
	ContentChangeSig
	NameChangeSig
	AbiChangeSig
	PubkeyChangeSig
	TextChangeSig
	MultihashChangeSig
	ContentHashChangeSig
)

func (e Interface) Hex() string {
	strings := [...]string{
		"0x01ffc9a7",
		"0x3b3b57de",
		"0xd8389dc5",
		"0x691f3431",
		"0x2203ab56",
		"0xc8690233",
		"0x59d1d43c",
		"0xe89401a1",
		"0xbc1c58d1",
	}

	if e < MetaSig || e > ContentHashChangeSig {
		return "Unknown"
	}

	return strings[e]
}

func (e Interface) Bytes() [4]uint8 {
	if e < MetaSig || e > ContentHashChangeSig {
		return [4]byte{}
	}

	str := e.Hex()
	by, _ := hexutil.Decode(str)
	var byArray [4]uint8
	for i := 0; i < 4; i++ {
		byArray[i] = by[i]
	}

	return byArray
}

func (e Interface) EventSig() string {
	strings := [...]string{
		"",
		"AddrChanged(bytes32,address)",
		"ContentChanged(bytes32,bytes32)",
		"NameChanged(bytes32,string)",
		"ABIChanged(bytes32,uint256)",
		"PubkeyChanged(bytes32,bytes32,bytes32)",
		"TextChanged(bytes32,string,string)",
		"MultihashChanged(bytes32,bytes)",
		"ContenthashChanged(bytes32,bytes)",
	}

	if e < MetaSig || e > ContentHashChangeSig {
		return "Unknown"
	}

	return strings[e]
}

func (e Interface) MethodSig() string {
	strings := [...]string{
		"supportsInterface(bytes4)",
		"addr(bytes32)",
		"content(bytes32)",
		"name(bytes32)",
		"ABI(bytes32,uint256)",
		"pubkey(bytes32)",
		"text(bytes32,string)",
		"multihash(bytes32)",
		"setContenthash(bytes32,bytes)",
	}

	if e < MetaSig || e > ContentHashChangeSig {
		return "Unknown"
	}

	return strings[e]
}
