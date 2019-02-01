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
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/vulcanize/vulcanizedb/pkg/geth"
)

func GetEventSignature(solidityMethodSignature string) string {
	eventSignature := []byte(solidityMethodSignature)
	hash := crypto.Keccak256Hash(eventSignature)
	return hash.Hex()
}

func GetLogNoteSignature(solidityMethodSignature string) string {
	rawSignature := GetEventSignature(solidityMethodSignature)
	return rawSignature[:10] + "00000000000000000000000000000000000000000000000000000000"
}

func GetSolidityMethodSignature(abi, name string) string {
	parsedAbi, _ := geth.ParseAbi(abi)

	if method, ok := parsedAbi.Methods[name]; ok {
		return method.Sig()
	} else if event, ok := parsedAbi.Events[name]; ok {
		return getEventSignature(event)
	}
	panic("Error: could not get Solidity method signature for: " + name)
}

func getEventSignature(event abi.Event) string {
	types := make([]string, len(event.Inputs))
	for i, input := range event.Inputs {
		types[i] = input.Type.String()
		i++
	}

	return fmt.Sprintf("%v(%v)", event.Name, strings.Join(types, ","))
}
