// VulcanizeDB
// Copyright © 2019 Vulcanize

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

package test_helpers

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/makerdao/vulcanizedb/pkg/config"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/constants"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/contract"
	"github.com/makerdao/vulcanizedb/pkg/contract_watcher/helpers/test_helpers/mocks"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/eth"
	"github.com/makerdao/vulcanizedb/pkg/eth/client"
	"github.com/makerdao/vulcanizedb/pkg/eth/converters"
	"github.com/makerdao/vulcanizedb/pkg/eth/node"
	"github.com/makerdao/vulcanizedb/test_config"
	. "github.com/onsi/gomega"
)

type TransferLog struct {
	ID       int64  `db:"id"`
	HeaderID int64  `db:"header_id"`
	LogIndex int64  `db:"log_idx"`
	TxIndex  int64  `db:"tx_idx"`
	From     string `db:"from_"`
	To       string `db:"to_"`
	Value    string `db:"value_"`
	RawLog   []byte `db:"raw_log"`
}

type NewOwnerLog struct {
	ID       int64  `db:"id"`
	HeaderID int64  `db:"header_id"`
	LogIndex int64  `db:"log_idx"`
	TxIndex  int64  `db:"tx_idx"`
	Node     string `db:"node_"`
	Label    string `db:"label_"`
	Owner    string `db:"owner_"`
	RawLog   []byte `db:"raw_log"`
}

func SetupDBandBC() (*postgres.DB, core.BlockChain) {
	con := test_config.TestClient
	testIPC := con.IPCPath
	rawRPCClient, err := rpc.Dial(testIPC)
	Expect(err).NotTo(HaveOccurred())
	rpcClient := client.NewRpcClient(rawRPCClient, testIPC)
	ethClient := ethclient.NewClient(rawRPCClient)
	blockChainClient := client.NewEthClient(ethClient)
	madeNode := node.MakeNode(rpcClient)
	transactionConverter := converters.NewTransactionConverter(ethClient)
	blockChain := eth.NewBlockChain(blockChainClient, rpcClient, madeNode, transactionConverter)

	db, err := postgres.NewDB(config.Database{
		Hostname: "localhost",
		Name:     "vulcanize_testing",
		Port:     5432,
	}, blockChain.Node())
	Expect(err).NotTo(HaveOccurred())

	return db, blockChain
}

func SetupTusdRepo(wantedEvents []string) (*postgres.DB, *contract.Contract) {
	db, err := postgres.NewDB(config.Database{
		Hostname: "localhost",
		Name:     "vulcanize_testing",
		Port:     5432,
	}, core.Node{})
	Expect(err).NotTo(HaveOccurred())

	info := SetupTusdContract(wantedEvents)

	return db, info
}

func SetupTusdContract(wantedEvents []string) *contract.Contract {
	p := mocks.NewParser(constants.TusdAbiString)
	err := p.Parse()
	Expect(err).ToNot(HaveOccurred())

	return contract.Contract{
		Address:       constants.TusdContractAddress,
		Abi:           p.Abi(),
		ParsedAbi:     p.ParsedAbi(),
		StartingBlock: 6194634,
		Events:        p.GetEvents(wantedEvents),
		FilterArgs:    map[string]bool{},
	}.Init()
}

func SetupMarketPlaceContract(wantedEvents []string) *contract.Contract {
	p := mocks.NewParser(constants.MarketPlaceAbiString)
	err := p.Parse()
	Expect(err).NotTo(HaveOccurred())

	return contract.Contract{
		Address:       constants.MarketPlaceContractAddress,
		StartingBlock: 6496012,
		Abi:           p.Abi(),
		ParsedAbi:     p.ParsedAbi(),
		Events:        p.GetEvents(wantedEvents),
		FilterArgs:    map[string]bool{},
	}.Init()
}

func SetupMolochContract(wantedEvents []string) *contract.Contract {
	p := mocks.NewParser(constants.MolochAbiString)
	err := p.Parse()
	Expect(err).NotTo(HaveOccurred())

	return contract.Contract{
		Address:       constants.MolochContractAddress,
		StartingBlock: 7218566,
		Abi:           p.Abi(),
		ParsedAbi:     p.ParsedAbi(),
		Events:        p.GetEvents(wantedEvents),
		FilterArgs:    map[string]bool{},
	}.Init()
}

func SetupOasisContract(wantedEvents []string) *contract.Contract {
	p := mocks.NewParser(constants.OasisAbiString)
	err := p.Parse()
	Expect(err).NotTo(HaveOccurred())

	return contract.Contract{
		Address:       constants.OasisContractAddress,
		StartingBlock: 7183773,
		Abi:           p.Abi(),
		ParsedAbi:     p.ParsedAbi(),
		Events:        p.GetEvents(wantedEvents),
		FilterArgs:    map[string]bool{},
	}.Init()
}

// TODO: tear down/setup DB from migrations so this doesn't alter the schema between tests
func TearDown(db *postgres.DB) {
	tx, err := db.Beginx()
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DELETE FROM public.addresses`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DELETE FROM public.headers`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec("DELETE FROM public.transactions")
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DELETE FROM public.receipts`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DROP TABLE public.checked_headers`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`CREATE TABLE public.checked_headers (
    	id SERIAL PRIMARY KEY,
    	header_id INTEGER UNIQUE NOT NULL REFERENCES headers (id) ON DELETE CASCADE);`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DROP SCHEMA IF EXISTS cw_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e CASCADE`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DROP SCHEMA IF EXISTS cw_0x314159265dd8dbb310642f98f50c066173c1259b CASCADE`)
	Expect(err).NotTo(HaveOccurred())

	err = tx.Commit()
	Expect(err).NotTo(HaveOccurred())

	_, err = db.Exec(`VACUUM public.checked_headers`)
	Expect(err).NotTo(HaveOccurred())
}
