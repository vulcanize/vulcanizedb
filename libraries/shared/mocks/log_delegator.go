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
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
)

type MockLogDelegator struct {
	AddedTransformers []transformer.EventTransformer
	DelegateCallCount int
	DelegateErrors    []error
}

func (delegator *MockLogDelegator) AddTransformer(t transformer.EventTransformer) {
	delegator.AddedTransformers = append(delegator.AddedTransformers, t)
}

func (delegator *MockLogDelegator) DelegateLogs() error {
	delegator.DelegateCallCount++
	if len(delegator.DelegateErrors) > 1 {
		var delegateErrorThisRun error
		delegateErrorThisRun, delegator.DelegateErrors = delegator.DelegateErrors[0], delegator.DelegateErrors[1:]
		return delegateErrorThisRun
	} else if len(delegator.DelegateErrors) == 1 {
		thisErr := delegator.DelegateErrors[0]
		delegator.DelegateErrors = []error{}
		return thisErr
	}
	return nil
}
