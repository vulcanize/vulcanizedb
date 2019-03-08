package fakes

type MockLightBlockRetriever struct {
	FirstBlock    int64
	FirstBlockErr error
}

func (retriever *MockLightBlockRetriever) RetrieveFirstBlock() (int64, error) {
	return retriever.FirstBlock, retriever.FirstBlockErr
}

func (retriever *MockLightBlockRetriever) RetrieveMostRecentBlock() (int64, error) {
	return 0, nil
}
