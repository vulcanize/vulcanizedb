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
