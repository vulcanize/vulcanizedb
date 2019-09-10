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

package mocks

import (
	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
)

type MockLogExtractor struct {
	AddedConfigs              []transformer.EventTransformerConfig
	AddTransformerConfigError error
	ExtractLogsCount          int
	ExtractLogsErrors         []error
	UncheckedHeadersExist     []bool
}

func (extractor *MockLogExtractor) AddTransformerConfig(config transformer.EventTransformerConfig) error {
	extractor.AddedConfigs = append(extractor.AddedConfigs, config)
	return extractor.AddTransformerConfigError
}

func (extractor *MockLogExtractor) ExtractLogs(recheckHeaders constants.TransformerExecution) (error, bool) {
	extractor.ExtractLogsCount++
	var errorThisRun error
	errorThisRun, extractor.ExtractLogsErrors = extractor.ExtractLogsErrors[0], extractor.ExtractLogsErrors[1:]
	if errorThisRun != nil {
		return errorThisRun, false
	}
	var missingHeadersExist bool
	missingHeadersExist, extractor.UncheckedHeadersExist = extractor.UncheckedHeadersExist[0], extractor.UncheckedHeadersExist[1:]
	return nil, missingHeadersExist
}
