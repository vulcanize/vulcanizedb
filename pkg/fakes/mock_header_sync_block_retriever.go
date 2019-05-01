package fakes

type MockHeaderSyncBlockRetriever struct {
	FirstBlock    int64
	FirstBlockErr error
}

func (retriever *MockHeaderSyncBlockRetriever) RetrieveFirstBlock() (int64, error) {
	return retriever.FirstBlock, retriever.FirstBlockErr
}

func (retriever *MockHeaderSyncBlockRetriever) RetrieveMostRecentBlock() (int64, error) {
	return 0, nil
}
