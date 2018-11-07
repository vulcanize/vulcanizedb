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

package history

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
)

type HeaderValidator struct {
	blockChain       core.BlockChain
	headerRepository datastore.HeaderRepository
	windowSize       int
}

func NewHeaderValidator(blockChain core.BlockChain, repository datastore.HeaderRepository, windowSize int) HeaderValidator {
	return HeaderValidator{
		blockChain:       blockChain,
		headerRepository: repository,
		windowSize:       windowSize,
	}
}

func (validator HeaderValidator) ValidateHeaders() ValidationWindow {
	window := MakeValidationWindow(validator.blockChain, validator.windowSize)
	blockNumbers := MakeRange(window.LowerBound, window.UpperBound)
	RetrieveAndUpdateHeaders(validator.blockChain, validator.headerRepository, blockNumbers)
	return window
}
