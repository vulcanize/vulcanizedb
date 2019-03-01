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

package mocks

import (
	"github.com/ethereum/go-ethereum/core/types"
)

type MockLogNoteConverter struct {
	err                   error
	returnModels          []interface{}
	PassedLogs            []types.Log
	ToModelsCalledCounter int
}

func (converter *MockLogNoteConverter) ToModels(ethLogs []types.Log) ([]interface{}, error) {
	converter.PassedLogs = ethLogs
	converter.ToModelsCalledCounter++
	return converter.returnModels, converter.err
}

func (converter *MockLogNoteConverter) SetConverterError(e error) {
	converter.err = e
}

func (converter *MockLogNoteConverter) SetReturnModels(models []interface{}) {
	converter.returnModels = models
}
