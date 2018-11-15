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

type MockCryptoParser struct {
	parsePublicKeyCalled           bool
	parsePublicKeyPassedPrivateKey string
	parsePublicKeyReturnString     string
	parsePublicKeyReturnErr        error
}

func NewMockCryptoParser() *MockCryptoParser {
	return &MockCryptoParser{
		parsePublicKeyCalled:           false,
		parsePublicKeyPassedPrivateKey: "",
		parsePublicKeyReturnString:     "",
		parsePublicKeyReturnErr:        nil,
	}
}

func (mcp *MockCryptoParser) SetReturnVal(pubKey string) {
	mcp.parsePublicKeyReturnString = pubKey
}

func (mcp *MockCryptoParser) SetReturnErr(err error) {
	mcp.parsePublicKeyReturnErr = err
}

func (mcp *MockCryptoParser) ParsePublicKey(privateKey string) (string, error) {
	mcp.parsePublicKeyCalled = true
	mcp.parsePublicKeyPassedPrivateKey = privateKey
	return mcp.parsePublicKeyReturnString, mcp.parsePublicKeyReturnErr
}

func (mcp *MockCryptoParser) AssertParsePublicKeyCalledWith(privateKey string) {
	Expect(mcp.parsePublicKeyCalled).To(BeTrue())
	Expect(mcp.parsePublicKeyPassedPrivateKey).To(Equal(privateKey))
}
