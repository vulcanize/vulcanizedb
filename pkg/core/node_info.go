package core

import (
	"fmt"
)

type NodeType int

const (
	GETH NodeType = iota
	PARITY
	INFURA
)

type Node struct {
	GenesisBlock string
	NetworkID    float64
	ID           string
	ClientName   string
}

type ParityNodeInfo struct {
	Track         string
	ParityVersion `json:"version"`
	Hash          string
}

func (pn ParityNodeInfo) String() string {
	return fmt.Sprintf("Parity/v%d.%d.%d/", pn.Major, pn.Minor, pn.Patch)
}

type ParityVersion struct {
	Major int
	Minor int
	Patch int
}
