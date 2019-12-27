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

	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/sirupsen/logrus"
)

func PopulateMissingHeaders(blockChain core.BlockChain, headerRepository datastore.HeaderRepository, startingBlockNumber int64) (int, error) {
	lastBlock, err := blockChain.LastBlock()
	if err != nil {
		logrus.Error("PopulateMissingHeaders: Error getting last block: ", err)
		return 0, err
	}

	blockNumbers, err := headerRepository.MissingBlockNumbers(startingBlockNumber, lastBlock.Int64())
	if err != nil {
		logrus.Error("PopulateMissingHeaders: Error getting missing block numbers: ", err)
		return 0, err
	} else if len(blockNumbers) == 0 {
		return 0, nil
	}

	logrus.Debug(getBlockRangeString(blockNumbers))
	_, err = RetrieveAndUpdateHeaders(blockChain, headerRepository, blockNumbers)
	if err != nil {
		logrus.Error("PopulateMissingHeaders: Error getting/updating headers: ", err)
		return 0, err
	}
	return len(blockNumbers), nil
}

func RetrieveAndUpdateHeaders(blockChain core.BlockChain, headerRepository datastore.HeaderRepository, blockNumbers []int64) (int, error) {
	headers, err := blockChain.GetHeadersByNumbers(blockNumbers)
	for _, header := range headers {
		_, err = headerRepository.CreateOrUpdateHeader(header)
		if err != nil {
			if err == repositories.ErrValidHeaderExists {
				continue
			}
			return 0, err
		}
	}
	return len(blockNumbers), nil
}

func getBlockRangeString(blockRange []int64) string {
	return fmt.Sprintf("Backfilling |%v| blocks", len(blockRange))
}
