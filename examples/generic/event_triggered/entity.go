// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package event_triggered

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type BurnEntity struct {
	TokenName    string
	TokenAddress common.Address
	Burner       common.Address
	Value        *big.Int
	Block        int64
	TxHash       string
}

type MintEntity struct {
	TokenName    string
	TokenAddress common.Address
	To           common.Address
	Amount       *big.Int
	Block        int64
	TxHash       string
}
