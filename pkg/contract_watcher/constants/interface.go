// VulcanizeDB
// Copyright © 2019 Vulcanize

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

// Resolver interface signatures
type Interface int

// Interface enums
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

// Hex returns the hex signature for an interface
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

// Bytes returns the bytes signature for an interface
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

// EventSig returns the event signature for an interface
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
