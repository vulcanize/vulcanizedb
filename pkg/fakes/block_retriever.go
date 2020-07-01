package fakes

type MockBlockRetriever struct {
	FirstBlock    int64
	FirstBlockErr error
}

func (retriever *MockBlockRetriever) RetrieveFirstBlock() (int64, error) {
	return retriever.FirstBlock, retriever.FirstBlockErr
}

func (retriever *MockBlockRetriever) RetrieveMostRecentBlock() (int64, error) {
	return 0, nil
}
