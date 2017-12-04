package core

type Contract struct {
	Attributes ContractAttributes
	Hash       string
}

type ContractAttribute struct {
	Name string
	Type string
}

type ContractAttributes []ContractAttribute

func (attributes ContractAttributes) Len() int {
	return len(attributes)
}

func (attributes ContractAttributes) Swap(i, j int) {
	attributes[i], attributes[j] = attributes[j], attributes[i]
}

func (attributes ContractAttributes) Less(i, j int) bool {
	return attributes[i].Name < attributes[j].Name
}
