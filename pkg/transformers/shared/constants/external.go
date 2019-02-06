package constants

import (
	"fmt"

	"github.com/spf13/viper"
)

func getEnvironmentString(key string) string {
	value := viper.GetString(key)
	if value == "" {
		panic(fmt.Sprintf("No environment configuration variable set for key: \"%v\"", key))
	}
	return value
}

// Returns an int from the environment, defaulting to 0 if it does not exist
func getEnvironmentInt64(key string) int64 {
	value := viper.GetInt64(key)
	if value == -1 {
		return 0
	}
	return value
}

// Getters for contract addresses from environment files
func CatContractAddress() string     { return getEnvironmentString("contract.address.cat") }
func DripContractAddress() string    { return getEnvironmentString("contract.address.drip") }
func FlapperContractAddress() string { return getEnvironmentString("contract.address.mcd_flap") }
func FlipperContractAddress() string { return getEnvironmentString("contract.address.eth_flip") }
func FlopperContractAddress() string { return getEnvironmentString("contract.address.mcd_flop") }
func PepContractAddress() string     { return getEnvironmentString("contract.address.pep") }
func PipContractAddress() string     { return getEnvironmentString("contract.address.pip") }
func PitContractAddress() string     { return getEnvironmentString("contract.address.pit") }
func RepContractAddress() string     { return getEnvironmentString("contract.address.rep") }
func VatContractAddress() string     { return getEnvironmentString("contract.address.vat") }
func VowContractAddress() string     { return getEnvironmentString("contract.address.vow") }

func CatABI() string        { return getEnvironmentString("contract.abi.cat") }
func DripABI() string       { return getEnvironmentString("contract.abi.drip") }
func FlapperABI() string    { return getEnvironmentString("contract.abi.mcd_flap") }
func FlipperABI() string    { return getEnvironmentString("contract.abi.eth_flip") }
func FlopperABI() string    { return getEnvironmentString("contract.abi.mcd_flop") }
func MedianizerABI() string { return getEnvironmentString("contract.abi.medianizer") }
func PitABI() string        { return getEnvironmentString("contract.abi.pit") }
func VatABI() string        { return getEnvironmentString("contract.abi.vat") }
func VowABI() string        { return getEnvironmentString("contract.abi.vow") }

func CatDeploymentBlock() int64     { return getEnvironmentInt64("contract.deployment-block.cat") }
func DripDeploymentBlock() int64    { return getEnvironmentInt64("contract.deployment-block.drip") }
func FlapperDeploymentBlock() int64 { return getEnvironmentInt64("contract.deployment-block.mcd_flap") }
func FlipperDeploymentBlock() int64 { return getEnvironmentInt64("contract.deployment-block.eth_flip") }
func FlopperDeploymentBlock() int64 { return getEnvironmentInt64("contract.deployment-block.mcd_flop") }
func PepDeploymentBlock() int64     { return getEnvironmentInt64("contract.deployment-block.pep") }
func PipDeploymentBlock() int64     { return getEnvironmentInt64("contract.deployment-block.pip") }
func PitDeploymentBlock() int64     { return getEnvironmentInt64("contract.deployment-block.pit") }
func RepDeploymentBlock() int64     { return getEnvironmentInt64("contract.deployment-block.rep") }
func VatDeploymentBlock() int64     { return getEnvironmentInt64("contract.deployment-block.vat") }
func VowDeploymentBlock() int64     { return getEnvironmentInt64("contract.deployment-block.vow") }
func MedianizerDeploymentBlock() int64 {
	return getEnvironmentInt64("contract.deployment-block.medianizer")
}
