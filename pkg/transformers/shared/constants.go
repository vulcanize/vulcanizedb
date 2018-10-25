// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shared

import (
	"github.com/spf13/viper"
)

func getContractValue(key string, fallback string) string {
	value := viper.GetString(key)
	if value == "" {
		return fallback
	}
	return value
}

var (
	DataItemLength = 32

	CatABI        = `[{"constant":true,"inputs":[],"name":"vat","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function","signature":"0x36569e77"},{"constant":true,"inputs":[],"name":"vow","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function","signature":"0x626cb3c5"},{"constant":true,"inputs":[{"name":"","type":"uint256"}],"name":"flips","outputs":[{"name":"ilk","type":"bytes32"},{"name":"urn","type":"bytes32"},{"name":"ink","type":"uint256"},{"name":"tab","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function","signature":"0x70d9235a"},{"constant":true,"inputs":[],"name":"nflip","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function","signature":"0x76181a51"},{"constant":true,"inputs":[],"name":"live","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function","signature":"0x957aa58c"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"wards","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function","signature":"0xbf353dbb"},{"constant":true,"inputs":[{"name":"","type":"bytes32"}],"name":"ilks","outputs":[{"name":"flip","type":"address"},{"name":"chop","type":"uint256"},{"name":"lump","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function","signature":"0xd9638d36"},{"constant":true,"inputs":[],"name":"pit","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function","signature":"0xf03c7c6e"},{"inputs":[{"name":"vat_","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor","signature":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"name":"ilk","type":"bytes32"},{"indexed":true,"name":"urn","type":"bytes32"},{"indexed":false,"name":"ink","type":"uint256"},{"indexed":false,"name":"art","type":"uint256"},{"indexed":false,"name":"tab","type":"uint256"},{"indexed":false,"name":"flip","type":"uint256"},{"indexed":false,"name":"iArt","type":"uint256"}],"name":"Bite","type":"event","signature":"0x99b5620489b6ef926d4518936cfec15d305452712b88bd59da2d9c10fb0953e8"},{"anonymous":true,"inputs":[{"indexed":true,"name":"sig","type":"bytes4"},{"indexed":true,"name":"guy","type":"address"},{"indexed":true,"name":"foo","type":"bytes32"},{"indexed":true,"name":"bar","type":"bytes32"},{"indexed":false,"name":"wad","type":"uint256"},{"indexed":false,"name":"fax","type":"bytes"}],"name":"LogNote","type":"event","signature":"0x644843f351d3fba4abcd60109eaff9f54bac8fb8ccf0bab941009c21df21cf31"},{"constant":false,"inputs":[{"name":"guy","type":"address"}],"name":"rely","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x65fae35e"},{"constant":false,"inputs":[{"name":"guy","type":"address"}],"name":"deny","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x9c52a7f1"},{"constant":false,"inputs":[{"name":"ilk","type":"bytes32"},{"name":"what","type":"bytes32"},{"name":"data","type":"uint256"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x1a0b287e"},{"constant":false,"inputs":[{"name":"what","type":"bytes32"},{"name":"data","type":"address"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0xd4e8be83"},{"constant":false,"inputs":[{"name":"ilk","type":"bytes32"},{"name":"what","type":"bytes32"},{"name":"flip","type":"address"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0xebecb39d"},{"constant":false,"inputs":[{"name":"ilk","type":"bytes32"},{"name":"urn","type":"bytes32"}],"name":"bite","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x72f7b593"},{"constant":false,"inputs":[{"name":"n","type":"uint256"},{"name":"wad","type":"uint256"}],"name":"flip","outputs":[{"name":"id","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0xe6f95917"}]`
	DripABI       = `[{"constant":true,"inputs":[],"name":"vat","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function","signature":"0x36569e77"},{"constant":true,"inputs":[],"name":"repo","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function","signature":"0x56ff3122"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"wards","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function","signature":"0xbf353dbb"},{"constant":true,"inputs":[{"name":"","type":"bytes32"}],"name":"ilks","outputs":[{"name":"vow","type":"bytes32"},{"name":"tax","type":"uint256"},{"name":"rho","type":"uint48"}],"payable":false,"stateMutability":"view","type":"function","signature":"0xd9638d36"},{"inputs":[{"name":"vat_","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor","signature":"constructor"},{"anonymous":true,"inputs":[{"indexed":true,"name":"sig","type":"bytes4"},{"indexed":true,"name":"guy","type":"address"},{"indexed":true,"name":"foo","type":"bytes32"},{"indexed":true,"name":"bar","type":"bytes32"},{"indexed":false,"name":"wad","type":"uint256"},{"indexed":false,"name":"fax","type":"bytes"}],"name":"LogNote","type":"event","signature":"0x644843f351d3fba4abcd60109eaff9f54bac8fb8ccf0bab941009c21df21cf31"},{"constant":false,"inputs":[{"name":"guy","type":"address"}],"name":"rely","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x65fae35e"},{"constant":false,"inputs":[{"name":"guy","type":"address"}],"name":"deny","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x9c52a7f1"},{"constant":true,"inputs":[],"name":"era","outputs":[{"name":"","type":"uint48"}],"payable":false,"stateMutability":"view","type":"function","signature":"0x143e55e0"},{"constant":false,"inputs":[{"name":"ilk","type":"bytes32"},{"name":"vow","type":"bytes32"},{"name":"tax","type":"uint256"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x1a0b287e"},{"constant":false,"inputs":[{"name":"what","type":"bytes32"},{"name":"data","type":"uint256"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x29ae8114"},{"constant":false,"inputs":[{"name":"ilk","type":"bytes32"}],"name":"drip","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x44e2a5a8"}]`
	FlipperABI    = `[{"constant":true,"inputs":[{"name":"","type":"uint256"}],"name":"bids","outputs":[{"name":"bid","type":"uint256"},{"name":"lot","type":"uint256"},{"name":"guy","type":"address"},{"name":"tic","type":"uint48"},{"name":"end","type":"uint48"},{"name":"urn","type":"bytes32"},{"name":"gal","type":"address"},{"name":"tab","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function","signature":"0x4423c5f1"},{"constant":true,"inputs":[],"name":"ttl","outputs":[{"name":"","type":"uint48"}],"payable":false,"stateMutability":"view","type":"function","signature":"0x4e8b1dd5"},{"constant":true,"inputs":[],"name":"gem","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function","signature":"0x7bd2bea7"},{"constant":true,"inputs":[],"name":"beg","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function","signature":"0x7d780d82"},{"constant":true,"inputs":[],"name":"tau","outputs":[{"name":"","type":"uint48"}],"payable":false,"stateMutability":"view","type":"function","signature":"0xcfc4af55"},{"constant":true,"inputs":[],"name":"kicks","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function","signature":"0xcfdd3302"},{"constant":true,"inputs":[],"name":"dai","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function","signature":"0xf4b9fa75"},{"inputs":[{"name":"dai_","type":"address"},{"name":"gem_","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor","signature":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"name":"id","type":"uint256"},{"indexed":false,"name":"lot","type":"uint256"},{"indexed":false,"name":"bid","type":"uint256"},{"indexed":false,"name":"gal","type":"address"},{"indexed":false,"name":"end","type":"uint48"},{"indexed":true,"name":"urn","type":"bytes32"},{"indexed":false,"name":"tab","type":"uint256"}],"name":"Kick","type":"event","signature":"0xbac86238bdba81d21995024470425ecb370078fa62b7271b90cf28cbd1e3e87e"},{"anonymous":true,"inputs":[{"indexed":true,"name":"sig","type":"bytes4"},{"indexed":true,"name":"guy","type":"address"},{"indexed":true,"name":"foo","type":"bytes32"},{"indexed":true,"name":"bar","type":"bytes32"},{"indexed":false,"name":"wad","type":"uint256"},{"indexed":false,"name":"fax","type":"bytes"}],"name":"LogNote","type":"event","signature":"0x644843f351d3fba4abcd60109eaff9f54bac8fb8ccf0bab941009c21df21cf31"},{"constant":true,"inputs":[],"name":"era","outputs":[{"name":"","type":"uint48"}],"payable":false,"stateMutability":"view","type":"function","signature":"0x143e55e0"},{"constant":false,"inputs":[{"name":"urn","type":"bytes32"},{"name":"gal","type":"address"},{"name":"tab","type":"uint256"},{"name":"lot","type":"uint256"},{"name":"bid","type":"uint256"}],"name":"kick","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0xeae19d9e"},{"constant":false,"inputs":[{"name":"id","type":"uint256"}],"name":"tick","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0xfc7b6aee"},{"constant":false,"inputs":[{"name":"id","type":"uint256"},{"name":"lot","type":"uint256"},{"name":"bid","type":"uint256"}],"name":"tend","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x4b43ed12"},{"constant":false,"inputs":[{"name":"id","type":"uint256"},{"name":"lot","type":"uint256"},{"name":"bid","type":"uint256"}],"name":"dent","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x5ff3a382"},{"constant":false,"inputs":[{"name":"id","type":"uint256"}],"name":"deal","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0xc959c42b"}]`
	FlopperABI    = `[{"constant":true,"inputs":[],"name":"era","outputs":[{"name":"","type":"uint48"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"uint256"}],"name":"bids","outputs":[{"name":"bid","type":"uint256"},{"name":"lot","type":"uint256"},{"name":"guy","type":"address"},{"name":"tic","type":"uint48"},{"name":"end","type":"uint48"},{"name":"vow","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"ttl","outputs":[{"name":"","type":"uint48"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"id","type":"uint256"},{"name":"lot","type":"uint256"},{"name":"bid","type":"uint256"}],"name":"dent","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"guy","type":"address"}],"name":"rely","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"gem","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"beg","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"guy","type":"address"}],"name":"deny","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"gal","type":"address"},{"name":"lot","type":"uint256"},{"name":"bid","type":"uint256"}],"name":"kick","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"wards","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"id","type":"uint256"}],"name":"deal","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"tau","outputs":[{"name":"","type":"uint48"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"kicks","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"dai","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"dai_","type":"address"},{"name":"gem_","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"name":"id","type":"uint256"},{"indexed":false,"name":"lot","type":"uint256"},{"indexed":false,"name":"bid","type":"uint256"},{"indexed":false,"name":"gal","type":"address"},{"indexed":false,"name":"end","type":"uint48"}],"name":"Kick","type":"event"},{"anonymous":true,"inputs":[{"indexed":true,"name":"sig","type":"bytes4"},{"indexed":true,"name":"guy","type":"address"},{"indexed":true,"name":"foo","type":"bytes32"},{"indexed":true,"name":"bar","type":"bytes32"},{"indexed":false,"name":"wad","type":"uint256"},{"indexed":false,"name":"fax","type":"bytes"}],"name":"LogNote","type":"event"}]`
	MedianizerABI = `[{"constant":false,"inputs":[{"name":"owner_","type":"address"}],"name":"setOwner","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"","type":"bytes32"}],"name":"poke","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"poke","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"compute","outputs":[{"name":"","type":"bytes32"},{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"wat","type":"address"}],"name":"set","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"wat","type":"address"}],"name":"unset","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"indexes","outputs":[{"name":"","type":"bytes12"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"next","outputs":[{"name":"","type":"bytes12"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"read","outputs":[{"name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"peek","outputs":[{"name":"","type":"bytes32"},{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"bytes12"}],"name":"values","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"min_","type":"uint96"}],"name":"setMin","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"authority_","type":"address"}],"name":"setAuthority","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"void","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"pos","type":"bytes12"},{"name":"wat","type":"address"}],"name":"set","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"authority","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"pos","type":"bytes12"}],"name":"unset","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"next_","type":"bytes12"}],"name":"setNext","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"min","outputs":[{"name":"","type":"uint96"}],"payable":false,"stateMutability":"view","type":"function"},{"anonymous":false,"inputs":[{"indexed":false,"name":"val","type":"bytes32"}],"name":"LogValue","type":"event"},{"anonymous":true,"inputs":[{"indexed":true,"name":"sig","type":"bytes4"},{"indexed":true,"name":"guy","type":"address"},{"indexed":true,"name":"foo","type":"bytes32"},{"indexed":true,"name":"bar","type":"bytes32"},{"indexed":false,"name":"wad","type":"uint256"},{"indexed":false,"name":"fax","type":"bytes"}],"name":"LogNote","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"authority","type":"address"}],"name":"LogSetAuthority","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"}],"name":"LogSetOwner","type":"event"}]]`
	PitABI        = `[{"constant":true,"inputs":[],"name":"vat","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function","signature":"0x36569e77"},{"constant":true,"inputs":[],"name":"live","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function","signature":"0x957aa58c"},{"constant":true,"inputs":[],"name":"drip","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function","signature":"0x9f678cca"},{"constant":true,"inputs":[],"name":"Line","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function","signature":"0xbabe8a3f"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"wards","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function","signature":"0xbf353dbb"},{"constant":true,"inputs":[{"name":"","type":"bytes32"}],"name":"ilks","outputs":[{"name":"spot","type":"uint256"},{"name":"line","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function","signature":"0xd9638d36"},{"inputs":[{"name":"vat_","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor","signature":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"name":"ilk","type":"bytes32"},{"indexed":true,"name":"urn","type":"bytes32"},{"indexed":false,"name":"ink","type":"uint256"},{"indexed":false,"name":"art","type":"uint256"},{"indexed":false,"name":"dink","type":"int256"},{"indexed":false,"name":"dart","type":"int256"},{"indexed":false,"name":"iArt","type":"uint256"}],"name":"Frob","type":"event","signature":"0xb2afa28318bcc689926b52835d844de174ef8de97e982a85c0199d584920791b"},{"anonymous":true,"inputs":[{"indexed":true,"name":"sig","type":"bytes4"},{"indexed":true,"name":"guy","type":"address"},{"indexed":true,"name":"foo","type":"bytes32"},{"indexed":true,"name":"bar","type":"bytes32"},{"indexed":false,"name":"wad","type":"uint256"},{"indexed":false,"name":"fax","type":"bytes"}],"name":"LogNote","type":"event","signature":"0x644843f351d3fba4abcd60109eaff9f54bac8fb8ccf0bab941009c21df21cf31"},{"constant":false,"inputs":[{"name":"guy","type":"address"}],"name":"rely","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x65fae35e"},{"constant":false,"inputs":[{"name":"guy","type":"address"}],"name":"deny","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x9c52a7f1"},{"constant":false,"inputs":[{"name":"ilk","type":"bytes32"},{"name":"what","type":"bytes32"},{"name":"data","type":"uint256"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x1a0b287e"},{"constant":false,"inputs":[{"name":"what","type":"bytes32"},{"name":"data","type":"uint256"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x29ae8114"},{"constant":false,"inputs":[{"name":"what","type":"bytes32"},{"name":"data","type":"address"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0xd4e8be83"},{"constant":false,"inputs":[{"name":"ilk","type":"bytes32"},{"name":"dink","type":"int256"},{"name":"dart","type":"int256"}],"name":"frob","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x5a984ded"}]`
	VatABI        = `[{"constant":true,"inputs":[],"name":"debt","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function","signature":"0x0dca59c1"},{"constant":true,"inputs":[{"name":"","type":"bytes32"},{"name":"","type":"bytes32"}],"name":"urns","outputs":[{"name":"ink","type":"uint256"},{"name":"art","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function","signature":"0x26e27482"},{"constant":true,"inputs":[],"name":"vice","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function","signature":"0x2d61a355"},{"constant":true,"inputs":[{"name":"","type":"bytes32"}],"name":"sin","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function","signature":"0xa60f1d3e"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"wards","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function","signature":"0xbf353dbb"},{"constant":true,"inputs":[{"name":"","type":"bytes32"},{"name":"","type":"bytes32"}],"name":"gem","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function","signature":"0xc0912683"},{"constant":true,"inputs":[{"name":"","type":"bytes32"}],"name":"ilks","outputs":[{"name":"take","type":"uint256"},{"name":"rate","type":"uint256"},{"name":"Ink","type":"uint256"},{"name":"Art","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function","signature":"0xd9638d36"},{"constant":true,"inputs":[{"name":"","type":"bytes32"}],"name":"dai","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function","signature":"0xf53e4e69"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor","signature":"constructor"},{"anonymous":true,"inputs":[{"indexed":true,"name":"sig","type":"bytes4"},{"indexed":true,"name":"foo","type":"bytes32"},{"indexed":true,"name":"bar","type":"bytes32"},{"indexed":true,"name":"too","type":"bytes32"},{"indexed":false,"name":"fax","type":"bytes"}],"name":"Note","type":"event","signature":"0x8c2dbbc2b33ffaa77c104b777e574a8a4ff79829dfee8b66f4dc63e3f8067152"},{"constant":false,"inputs":[{"name":"guy","type":"address"}],"name":"rely","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x65fae35e"},{"constant":false,"inputs":[{"name":"guy","type":"address"}],"name":"deny","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x9c52a7f1"},{"constant":false,"inputs":[{"name":"ilk","type":"bytes32"}],"name":"init","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x3b663195"},{"constant":false,"inputs":[{"name":"ilk","type":"bytes32"},{"name":"guy","type":"bytes32"},{"name":"rad","type":"int256"}],"name":"slip","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x42066cbb"},{"constant":false,"inputs":[{"name":"ilk","type":"bytes32"},{"name":"src","type":"bytes32"},{"name":"dst","type":"bytes32"},{"name":"rad","type":"int256"}],"name":"flux","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0xa6e41821"},{"constant":false,"inputs":[{"name":"src","type":"bytes32"},{"name":"dst","type":"bytes32"},{"name":"rad","type":"int256"}],"name":"move","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x78f19470"},{"constant":false,"inputs":[{"name":"i","type":"bytes32"},{"name":"u","type":"bytes32"},{"name":"v","type":"bytes32"},{"name":"w","type":"bytes32"},{"name":"dink","type":"int256"},{"name":"dart","type":"int256"}],"name":"tune","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x5dd6471a"},{"constant":false,"inputs":[{"name":"i","type":"bytes32"},{"name":"u","type":"bytes32"},{"name":"v","type":"bytes32"},{"name":"w","type":"bytes32"},{"name":"dink","type":"int256"},{"name":"dart","type":"int256"}],"name":"grab","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x3690ae4c"},{"constant":false,"inputs":[{"name":"u","type":"bytes32"},{"name":"v","type":"bytes32"},{"name":"rad","type":"int256"}],"name":"heal","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x990a5f63"},{"constant":false,"inputs":[{"name":"i","type":"bytes32"},{"name":"u","type":"bytes32"},{"name":"rate","type":"int256"}],"name":"fold","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0xe6a6a64d"},{"constant":false,"inputs":[{"name":"i","type":"bytes32"},{"name":"u","type":"bytes32"},{"name":"take","type":"int256"}],"name":"toll","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x09b7a0b5"}]`
	VowABI        = `[{"constant":true,"inputs":[],"name":"Awe","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"Joy","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"flap","outputs":[{"name":"id","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"hump","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"wad","type":"uint256"}],"name":"kiss","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"what","type":"bytes32"},{"name":"data","type":"uint256"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"Ash","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"era","type":"uint48"}],"name":"flog","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"vat","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"Woe","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"wait","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"guy","type":"address"}],"name":"rely","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"bump","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"tab","type":"uint256"}],"name":"fess","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"row","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"uint48"}],"name":"sin","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"guy","type":"address"}],"name":"deny","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"flop","outputs":[{"name":"id","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"wards","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"sump","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"Sin","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"what","type":"bytes32"},{"name":"addr","type":"address"}],"name":"file","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"cow","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"wad","type":"uint256"}],"name":"heal","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":true,"inputs":[{"indexed":true,"name":"sig","type":"bytes4"},{"indexed":true,"name":"guy","type":"address"},{"indexed":true,"name":"foo","type":"bytes32"},{"indexed":true,"name":"bar","type":"bytes32"},{"indexed":false,"name":"wad","type":"uint256"},{"indexed":false,"name":"fax","type":"bytes"}],"name":"LogNote","type":"event"}]`

	// temporary addresses from Kovan deployment
	CatContractAddress     = getContractValue("contract.cat", "0x2f34f22a00ee4b7a8f8bbc4eaee1658774c624e0")
	DripContractAddress    = getContractValue("contract.drip", "0x891c04639a5edcae088e546fa125b5d7fb6a2b9d")
	FlipperContractAddress = getContractValue("contract.eth_flip", "0x32D496Ad866D110060866B7125981C73642cc509") // ETH FLIP Contract
	FlopperContractAddress = getContractValue("contract.mcd_flop", "0x6191C9b0086c2eBF92300cC507009b53996FbFFa") // MCD FLOP Contract
	PepContractAddress     = getContractValue("contract.pep", "0xB1997239Cfc3d15578A3a09730f7f84A90BB4975")
	PipContractAddress     = getContractValue("contract.pip", "0x9FfFE440258B79c5d6604001674A4722FfC0f7Bc")
	PitContractAddress     = getContractValue("contract.pit", "0xe7cf3198787c9a4daac73371a38f29aaeeced87e")
	RepContractAddress     = getContractValue("contract.rep", "0xf88bBDc1E2718F8857F30A180076ec38d53cf296")
	VatContractAddress     = getContractValue("contract.vat", "0xcd726790550afcd77e9a7a47e86a3f9010af126b")
	VowContractAddress     = getContractValue("contract.vow", "0x3728e9777B2a0a611ee0F89e00E01044ce4736d1")

	//TODO: get pit and drip file method signatures directly from the ABI
	biteMethod                = GetSolidityMethodSignature(CatABI, "Bite")
	catFileChopLumpMethod     = "file(bytes32,bytes32,uint256)"
	catFileFlipMethod         = GetSolidityMethodSignature(CatABI, "file")
	catFilePitVowMethod       = "file(bytes32,address)"
	dealMethod                = GetSolidityMethodSignature(FlipperABI, "deal")
	dentMethod                = GetSolidityMethodSignature(FlipperABI, "dent")
	dripDripMethod            = GetSolidityMethodSignature(DripABI, "drip")
	dripFileIlkMethod         = "file(bytes32,bytes32,uint256)"
	dripFileRepoMethod        = GetSolidityMethodSignature(DripABI, "file")
	dripFileVowMethod         = "file(bytes32,bytes32)"
	flipKickMethod            = GetSolidityMethodSignature(FlipperABI, "Kick")
	flogMethod                = GetSolidityMethodSignature(VowABI, "flog")
	flopKickMethod            = GetSolidityMethodSignature(FlopperABI, "Kick")
	frobMethod                = GetSolidityMethodSignature(PitABI, "Frob")
	logValueMethod            = GetSolidityMethodSignature(MedianizerABI, "LogValue")
	pitFileDebtCeilingMethod  = "file(bytes32,uint256)"
	pitFileIlkMethod          = "file(bytes32,bytes32,uint256)"
	pitFileStabilityFeeMethod = GetSolidityMethodSignature(PitABI, "file")
	tendMethod                = GetSolidityMethodSignature(FlipperABI, "tend")
	vatHealMethod             = GetSolidityMethodSignature(VatABI, "heal")
	vatGrabMethod             = GetSolidityMethodSignature(VatABI, "grab")
	vatInitMethod             = GetSolidityMethodSignature(VatABI, "init")
	vatMoveMethod             = GetSolidityMethodSignature(VatABI, "move")
	vatFoldMethod             = GetSolidityMethodSignature(VatABI, "fold")
	vatSlipMethod             = GetSolidityMethodSignature(VatABI, "slip")
	vatTollMethod             = GetSolidityMethodSignature(VatABI, "toll")
	vatTuneMethod             = GetSolidityMethodSignature(VatABI, "tune")
	vatFluxMethod             = GetSolidityMethodSignature(VatABI, "flux")

	BiteSignature                = GetEventSignature(biteMethod)
	DealSignature                = GetLogNoteSignature(dealMethod)
	CatFileChopLumpSignature     = GetLogNoteSignature(catFileChopLumpMethod)
	CatFileFlipSignature         = GetLogNoteSignature(catFileFlipMethod)
	CatFilePitVowSignature       = GetLogNoteSignature(catFilePitVowMethod)
	DentFunctionSignature        = GetLogNoteSignature(dentMethod)
	DripDripSignature            = GetLogNoteSignature(dripDripMethod)
	DripFileIlkSignature         = GetLogNoteSignature(dripFileIlkMethod)
	DripFileRepoSignature        = GetLogNoteSignature(dripFileRepoMethod)
	DripFileVowSignature         = GetLogNoteSignature(dripFileVowMethod)
	FlipKickSignature            = GetEventSignature(flipKickMethod)
	FlogSignature                = GetLogNoteSignature(flogMethod)
	FlopKickSignature            = GetEventSignature(flopKickMethod)
	FrobSignature                = GetEventSignature(frobMethod)
	LogValueSignature            = GetEventSignature(logValueMethod)
	PitFileDebtCeilingSignature  = GetLogNoteSignature(pitFileDebtCeilingMethod)
	PitFileIlkSignature          = GetLogNoteSignature(pitFileIlkMethod)
	PitFileStabilityFeeSignature = GetLogNoteSignature(pitFileStabilityFeeMethod)
	TendFunctionSignature        = GetLogNoteSignature(tendMethod)
	VatHealSignature             = GetLogNoteSignature(vatHealMethod)
	VatGrabSignature             = GetLogNoteSignature(vatGrabMethod)
	VatInitSignature             = GetLogNoteSignature(vatInitMethod)
	VatMoveSignature             = GetLogNoteSignature(vatMoveMethod)
	VatFoldSignature             = GetLogNoteSignature(vatFoldMethod)
	VatSlipSignature             = GetLogNoteSignature(vatSlipMethod)
	VatTollSignature             = GetLogNoteSignature(vatTollMethod)
	VatTuneSignature             = GetLogNoteSignature(vatTuneMethod)
	VatFluxSignature             = GetLogNoteSignature(vatFluxMethod)

	BiteLabel                = "bite"
	DealLabel                = "deal"
	CatFileChopLumpLabel     = "catFileChopLump"
	CatFileFlipLabel         = "catFileFlip"
	CatFilePitVowLabel       = "catFilePitVow"
	DentLabel                = "dent"
	DripDripLabel            = "dripDrip"
	DripFileIlkLabel         = "dripFileIlk"
	DripFileRepoLabel        = "dripFileRepo"
	DripFileVowLabel         = "dripFileVow"
	FlipKickLabel            = "flipKick"
	FlogLabel                = "flog"
	FlopKickLabel            = "flopKick"
	FrobLabel                = "frob"
	PitFileDebtCeilingLabel  = "pitFileDebtCeiling"
	PitFileIlkLabel          = "pitFileIlk"
	PitFileStabilityFeeLabel = "pitFileStabilityFee"
	PriceFeedLabel           = "priceFeed"
	TendLabel                = "tend"
	VatHealLabel             = "vatHeal"
	VatGrabLabel             = "vatGrab"
	VatInitLabel             = "vatInit"
	VatMoveLabel             = "vatMove"
	VatFoldLabel             = "vatFold"
	VatSlipLabel             = "vatSlip"
	VatTollLabel             = "vatToll"
	VatTuneLabel             = "vatTune"
	VatFluxLabel             = "vatFlux"
)
