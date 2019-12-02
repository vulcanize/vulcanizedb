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

package history

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type ValidationWindow struct {
	LowerBound int64
	UpperBound int64
}

func (window ValidationWindow) Size() int {
	return int(window.UpperBound - window.LowerBound)
}

func MakeValidationWindow(blockchain core.BlockChain, windowSize int) (ValidationWindow, error) {
	upperBound, err := blockchain.LastBlock()
	if err != nil {
		log.Error("MakeValidationWindow: error getting LastBlock: ", err)
		return ValidationWindow{}, err
	}
	lowerBound := upperBound.Int64() - int64(windowSize)
	return ValidationWindow{lowerBound, upperBound.Int64()}, nil
}

func MakeRange(min, max int64) []int64 {
	a := make([]int64, max-min+1)
	for i := range a {
		a[i] = min + int64(i)
	}
	return a
}

func (window ValidationWindow) GetString() string {
	return fmt.Sprintf("Validating Blocks |%v|-- Validation Window --|%v|",
		window.LowerBound, window.UpperBound)
}
