package fakes

type BlockRetriever struct {
	FirstBlock    int64
	FirstBlockErr error
}

func (retriever *BlockRetriever) RetrieveFirstBlock() (int64, error) {
	return retriever.FirstBlock, retriever.FirstBlockErr
}

func (retriever *BlockRetriever) RetrieveMostRecentBlock() (int64, error) {
	return 0, nil
}
