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
	AddedConfigs        []transformer.EventTransformerConfig
	ExtractLogsCount    int
	ExtractLogsErrors   []error
	MissingHeadersExist []bool
}

func (extractor *MockLogExtractor) AddTransformerConfig(config transformer.EventTransformerConfig) {
	extractor.AddedConfigs = append(extractor.AddedConfigs, config)
}

func (extractor *MockLogExtractor) ExtractLogs(recheckHeaders constants.TransformerExecution, errs chan error, missingHeadersFound chan bool) {
	extractor.ExtractLogsCount++
	var errorThisRun error
	errorThisRun, extractor.ExtractLogsErrors = extractor.ExtractLogsErrors[0], extractor.ExtractLogsErrors[1:]
	if errorThisRun != nil {
		errs <- errorThisRun
	}
	var missingHeadersExist bool
	missingHeadersExist, extractor.MissingHeadersExist = extractor.MissingHeadersExist[0], extractor.MissingHeadersExist[1:]
	missingHeadersFound <- missingHeadersExist
}
