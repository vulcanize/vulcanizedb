package crypto

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discover"
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
	pubKey := discover.PubkeyID(&np.PublicKey)
	return pubKey.String(), nil
}
