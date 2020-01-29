package fakes

type MockFullSyncBlockRetriever struct {
	FirstBlock      int64
	FirstBlockErr   error
	MostRecentBlock int64
}

func (retriever *MockFullSyncBlockRetriever) RetrieveFirstBlock(contractAddr string) (int64, error) {
	return retriever.FirstBlock, retriever.FirstBlockErr
}

func (retriever *MockFullSyncBlockRetriever) RetrieveMostRecentBlock() (int64, error) {
	return retriever.MostRecentBlock, nil
}
