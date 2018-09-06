package shared

import "github.com/ethereum/go-ethereum/crypto"

func GetEventSignature(solidityMethodSignature string) string {
	eventSignature := []byte(solidityMethodSignature)
	hash := crypto.Keccak256Hash(eventSignature)
	return hash.Hex()
}

func GetLogNoteSignature(solidityMethodSignature string) string {
	rawSignature := GetEventSignature(solidityMethodSignature)
	return rawSignature[:10] + "00000000000000000000000000000000000000000000000000000000"
}
