// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
