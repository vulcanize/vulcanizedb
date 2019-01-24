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

package fakes

import . "github.com/onsi/gomega"

type MockFsReader struct {
	readCalled      bool
	readPassedPath  string
	readReturnBytes []byte
	readReturnErr   error
}

func NewMockFsReader() *MockFsReader {
	return &MockFsReader{
		readCalled:      false,
		readPassedPath:  "",
		readReturnBytes: nil,
		readReturnErr:   nil,
	}
}

func (mfr *MockFsReader) SetReturnBytes(returnBytes []byte) {
	mfr.readReturnBytes = returnBytes
}

func (mfr *MockFsReader) SetReturnErr(err error) {
	mfr.readReturnErr = err
}

func (mfr *MockFsReader) Read(path string) ([]byte, error) {
	mfr.readCalled = true
	mfr.readPassedPath = path
	return mfr.readReturnBytes, mfr.readReturnErr
}

func (mfr *MockFsReader) AssertReadCalledWith(path string) {
	Expect(mfr.readCalled).To(BeTrue())
	Expect(mfr.readPassedPath).To(Equal(path))
}
