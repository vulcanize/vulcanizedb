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
