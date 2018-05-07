package cold_import

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/crypto"
	"github.com/vulcanize/vulcanizedb/pkg/fs"
	"strings"
)

const (
	ColdImportClientName         = "LevelDbColdImport"
	ColdImportNetworkId  float64 = 1
)

var (
	NoChainDataErr = errors.New("Level DB path does not include chaindata extension.")
	NoGethRootErr  = errors.New("Level DB path does not include root path to geth.")
)

type ColdImportNodeBuilder struct {
	reader fs.Reader
	parser crypto.PublicKeyParser
}

func NewColdImportNodeBuilder(reader fs.Reader, parser crypto.PublicKeyParser) ColdImportNodeBuilder {
	return ColdImportNodeBuilder{reader: reader, parser: parser}
}

func (cinb ColdImportNodeBuilder) GetNode(genesisBlock []byte, levelPath string) (core.Node, error) {
	var coldNode core.Node
	nodeKeyPath, err := getNodeKeyPath(levelPath)
	if err != nil {
		return coldNode, err
	}
	nodeKey, err := cinb.reader.Read(nodeKeyPath)
	if err != nil {
		return coldNode, err
	}
	nodeId, err := cinb.parser.ParsePublicKey(string(nodeKey))
	if err != nil {
		return coldNode, err
	}
	genesisBlockHash := common.BytesToHash(genesisBlock).String()
	coldNode = core.Node{
		GenesisBlock: genesisBlockHash,
		NetworkID:    ColdImportNetworkId,
		ID:           nodeId,
		ClientName:   ColdImportClientName,
	}
	return coldNode, nil
}

func getNodeKeyPath(levelPath string) (string, error) {
	chaindataExtension := "chaindata"
	if !strings.Contains(levelPath, chaindataExtension) {
		return "", NoChainDataErr
	}
	chaindataExtensionLength := len(chaindataExtension)
	gethRootPathLength := len(levelPath) - chaindataExtensionLength
	if gethRootPathLength <= chaindataExtensionLength {
		return "", NoGethRootErr
	}
	gethRootPath := levelPath[:gethRootPathLength]
	nodeKeyPath := gethRootPath + "nodekey"
	return nodeKeyPath, nil
}
