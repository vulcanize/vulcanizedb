package geth

import (
	"errors"
	"io/ioutil"
	"strings"

	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	ErrInvalidAbiFile   = errors.New("invalid abi")
	ErrMissingAbiFile   = errors.New("missing abi")
	ErrApiRequestFailed = errors.New("etherscan api request failed")
)

type Response struct {
	Status  string
	Message string
	Result  string
}

type EtherScanAPI struct {
	client *http.Client
	url    string
}

func NewEtherScanClient(url string) *EtherScanAPI {
	return &EtherScanAPI{
		client: &http.Client{Timeout: 10 * time.Second},
		url:    url,
	}
}

func GenURL(network string) string {
	switch network {
	case "ropsten":
		return "https://ropsten.etherscan.io"
	case "kovan":
		return "https://kovan.etherscan.io"
	case "rinkeby":
		return "https://rinkeby.etherscan.io"
	default:
		return "https://api.etherscan.io"
	}
}

//https://api.etherscan.io/api?module=contract&action=getabi&address=%s
func (e *EtherScanAPI) GetAbi(contractHash string) (string, error) {
	target := new(Response)
	request := fmt.Sprintf("%s/api?module=contract&action=getabi&address=%s", e.url, contractHash)
	r, err := e.client.Get(request)
	if err != nil {
		return "", ErrApiRequestFailed
	}
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&target)
	return target.Result, err
}

func ParseAbiFile(abiFilePath string) (abi.ABI, error) {
	abiString, err := ReadAbiFile(abiFilePath)
	if err != nil {
		return abi.ABI{}, ErrMissingAbiFile
	}
	return ParseAbi(abiString)
}

func ParseAbi(abiString string) (abi.ABI, error) {
	parsedAbi, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return abi.ABI{}, ErrInvalidAbiFile
	}
	return parsedAbi, nil
}

func ReadAbiFile(abiFilePath string) (string, error) {
	filesBytes, err := ioutil.ReadFile(abiFilePath)
	if err != nil {
		return "", ErrMissingAbiFile
	}
	return string(filesBytes), nil
}
