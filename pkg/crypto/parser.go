package crypto

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discv5"
)

type PublicKeyParser interface {
	ParsePublicKey(privateKey string) (string, error)
}

type EthPublicKeyParser struct{}

func (EthPublicKeyParser) ParsePublicKey(privateKey string) (string, error) {
	np, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", err
	}
	pubKey := discv5.PubkeyID(&np.PublicKey)
	return pubKey.String(), nil
}
