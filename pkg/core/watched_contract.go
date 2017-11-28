package core

type WatchedContract struct {
	Hash         string
	Transactions []Transaction
}

type ContractAttribute struct {
	Name string
	Type string
}

type ContractAttributes []ContractAttribute

func (s ContractAttributes) Len() int {
	return len(s)
}
func (s ContractAttributes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ContractAttributes) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}