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

package every_block

// Struct to hold token supply data
type TokenSupply struct {
	Value        string
	TokenAddress string
	BlockNumber  int64
}

// Struct to hold token holder address balance data
type TokenBalance struct {
	Value              string
	TokenAddress       string
	BlockNumber        int64
	TokenHolderAddress string
}

// Struct to hold token allowance data
type TokenAllowance struct {
	Value               string
	TokenAddress        string
	BlockNumber         int64
	TokenHolderAddress  string
	TokenSpenderAddress string
}
