// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package geth

import (
	"errors"
	"strings"

	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/vulcanize/vulcanizedb/pkg/fs"
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
	reader := fs.FsReader{}
	filesBytes, err := reader.Read(abiFilePath)
	if err != nil {
		return "", ErrMissingAbiFile
	}
	return string(filesBytes), nil
}
