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

type TransferModel struct {
	TokenName    string `db:"token_name"`
	TokenAddress string `db:"token_address"`
	To           string `db:"to_address"`
	From         string `db:"from_address"`
	Tokens       string `db:"tokens"`
	Block        int64  `db:"block"`
	TxHash       string `db:"tx"`
}

type ApprovalModel struct {
	TokenName    string `db:"token_name"`
	TokenAddress string `db:"token_address"`
	Owner        string `db:"owner"`
	Spender      string `db:"spender"`
	Tokens       string `db:"tokens"`
	Block        int64  `db:"block"`
	TxHash       string `db:"tx"`
}
