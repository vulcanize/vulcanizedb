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

package erc20_test_helpers

type TokenSupplyDBRow struct {
	ID           int64
	Supply       int64
	BlockID      int64  `db:"block_id"`
	TokenAddress string `db:"token_address"`
}

type TransferDBRow struct {
	ID             int64 `db:"id"`
	VulcanizeLogID int64 `db:"vulcanize_log_id"`
}
