package fakes

type MockFullBlockRetriever struct {
	FirstBlock      int64
	FirstBlockErr   error
	MostRecentBlock int64
}

func (retriever *MockFullBlockRetriever) RetrieveFirstBlock(contractAddr string) (int64, error) {
	return retriever.FirstBlock, retriever.FirstBlockErr
}

func (retriever *MockFullBlockRetriever) RetrieveMostRecentBlock() (int64, error) {
	return retriever.MostRecentBlock, nil
}
