// Copyright © 2019 Vulcanize, Inc
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipfs/go-ipfs-ds-help"
	"github.com/multiformats/go-multihash"
	"gitlab.com/microo8/plgo"
)

func MultihashKeyFromCIDString(c string) string {
	logger := plgo.NewErrorLogger("", log.Ltime|log.Lshortfile)
	dc, err := cid.Decode(c)
	if err != nil {
		logger.Fatalf("cid decode error: %s", err.Error())
	}
	dbKey := dshelp.CidToDsKey(dc)
	return blockstore.BlockPrefix.String() + dbKey.String()
}

func CIDStringFromKeccak256(codec uint64, hash []byte) string {
	logger := plgo.NewErrorLogger("", log.Ltime|log.Lshortfile)
	buf, err := multihash.Encode(hash, multihash.KECCAK_256)
	if err != nil {
		logger.Fatalf("multihash encode error: %s", err.Error())
	}
	return cid.NewCidV1(codec, multihash.Multihash(buf)).String()
}

func MultihashKeyFromKeccak256(codec uint64, hash []byte) string {
	logger := plgo.NewErrorLogger("", log.Ltime|log.Lshortfile)
	buf, err := multihash.Encode(hash, multihash.KECCAK_256)
	if err != nil {
		logger.Fatalf("multihash encode error: %s", err.Error())
	}
	dc := cid.NewCidV1(codec, multihash.Multihash(buf))
	dbKey := dshelp.CidToDsKey(dc)
	return blockstore.BlockPrefix.String() + dbKey.String()
}

func BytesToKeccak256Bytes(raw []byte) []byte {
	return crypto.Keccak256(raw)
}

func HexToKeccak256Bytes(hex string) []byte {
	return crypto.Keccak256(common.Hex2Bytes(hex))
}

func BytesToKeccak256Hex(raw []byte) string {
	return common.Bytes2Hex(crypto.Keccak256(raw))
}

func HexToKeccak256Hex(hex string) string {
	return common.Bytes2Hex(crypto.Keccak256(common.Hex2Bytes(hex)))
}