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

package flip_kick

import (
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/flip_kick"
)

type MockFlipKickRepository struct {
	CreateRecordError               error
	EndingBlockNumber               int64
	FlipKicksCreated                []flip_kick.FlipKickModel
	HeaderIds                       []int64
	HeadersToReturn                 []core.Header
	MissingHeadersError             error
	StartingBlockNumber             int64
	markHeaderCheckedErr            error
	markHeaderCheckedPassedHeaderId int64
}

func (mfkr *MockFlipKickRepository) Create(headerId int64, flipKick []flip_kick.FlipKickModel) error {
	mfkr.HeaderIds = append(mfkr.HeaderIds, headerId)
	mfkr.FlipKicksCreated = append(mfkr.FlipKicksCreated, flipKick...)

	return mfkr.CreateRecordError
}

func (mfkr *MockFlipKickRepository) MarkHeaderChecked(headerId int64) error {
	mfkr.markHeaderCheckedPassedHeaderId = headerId
	return mfkr.markHeaderCheckedErr
}

func (mfkr *MockFlipKickRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	mfkr.StartingBlockNumber = startingBlockNumber
	mfkr.EndingBlockNumber = endingBlockNumber

	return mfkr.HeadersToReturn, mfkr.MissingHeadersError
}

func (mfkr *MockFlipKickRepository) SetHeadersToReturn(headers []core.Header) {
	mfkr.HeadersToReturn = headers
}

func (mfkr *MockFlipKickRepository) SetCreateRecordError(err error) {
	mfkr.CreateRecordError = err
}

func (mfkr *MockFlipKickRepository) SetMarkHeaderCheckedErr(err error) {
	mfkr.markHeaderCheckedErr = err
}

func (mfkr *MockFlipKickRepository) SetMissingHeadersError(err error) {
	mfkr.MissingHeadersError = err
}

func (mfkr *MockFlipKickRepository) AssertMarkHeaderCheckedCalledWith(headerId int64) {
	Expect(mfkr.markHeaderCheckedPassedHeaderId).To(Equal(headerId))
}
