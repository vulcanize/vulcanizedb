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

package btc

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

type CIDIndexer struct {
	db *postgres.DB
}

func NewCIDIndexer(db *postgres.DB) *CIDIndexer {
	return &CIDIndexer{
		db: db,
	}
}

func (in *CIDIndexer) Index(cids shared.CIDsForIndexing) error {
	cidWrapper, ok := cids.(*CIDPayload)
	if !ok {
		return fmt.Errorf("btc indexer expected cids type %T got %T", &CIDPayload{}, cids)
	}
	tx, err := in.db.Beginx()
	if err != nil {
		return err
	}
	headerID, err := in.indexHeaderCID(tx, cidWrapper.HeaderCID)
	if err != nil {
		logrus.Error("btc indexer error when indexing header")
		return err
	}
	if err := in.indexTransactionCIDs(tx, cidWrapper.TransactionCIDs, headerID); err != nil {
		logrus.Error("btc indexer error when indexing transactions")
		return err
	}
	return tx.Commit()
}

func (in *CIDIndexer) indexHeaderCID(tx *sqlx.Tx, header HeaderModel) (int64, error) {
	var headerID int64
	err := tx.QueryRowx(`INSERT INTO btc.header_cids (block_number, block_hash, parent_hash, cid, timestamp, bits, node_id, times_validated)
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
							ON CONFLICT (block_number, block_hash) DO UPDATE SET (parent_hash, cid, timestamp, bits, node_id, times_validated) = ($3, $4, $5, $6, $7, btc.header_cids.times_validated + 1)
							RETURNING id`,
		header.BlockNumber, header.BlockHash, header.ParentHash, header.CID, header.Timestamp, header.Bits, in.db.NodeID, 1).Scan(&headerID)
	return headerID, err
}

func (in *CIDIndexer) indexTransactionCIDs(tx *sqlx.Tx, transactions []TxModelWithInsAndOuts, headerID int64) error {
	for _, transaction := range transactions {
		txID, err := in.indexTransactionCID(tx, transaction, headerID)
		if err != nil {
			logrus.Error("btc indexer error when indexing header")
			return err
		}
		for _, input := range transaction.TxInputs {
			if err := in.indexTxInput(tx, input, txID); err != nil {
				logrus.Error("btc indexer error when indexing tx inputs")
				return err
			}
		}
		for _, output := range transaction.TxOutputs {
			if err := in.indexTxOutput(tx, output, txID); err != nil {
				logrus.Error("btc indexer error when indexing tx outputs")
				return err
			}
		}
	}
	return nil
}

func (in *CIDIndexer) indexTransactionCID(tx *sqlx.Tx, transaction TxModelWithInsAndOuts, headerID int64) (int64, error) {
	var txID int64
	err := tx.QueryRowx(`INSERT INTO btc.transaction_cids (header_id, tx_hash, index, cid, segwit, witness_hash)
							VALUES ($1, $2, $3, $4, $5, $6)
							ON CONFLICT (tx_hash) DO UPDATE SET (header_id, index, cid, segwit, witness_hash) = ($1, $3, $4, $5, $6)
							RETURNING id`,
		headerID, transaction.TxHash, transaction.Index, transaction.CID, transaction.SegWit, transaction.WitnessHash).Scan(&txID)
	return txID, err
}

func (in *CIDIndexer) indexTxInput(tx *sqlx.Tx, txInput TxInput, txID int64) error {
	_, err := tx.Exec(`INSERT INTO btc.tx_inputs (tx_id, index, witness, sig_script, outpoint_tx_hash, outpoint_index)
						VALUES ($1, $2, $3, $4, $5, $6)
						ON CONFLICT (tx_id, index) DO UPDATE SET (witness, sig_script, outpoint_tx_hash, outpoint_index) = ($3, $4, $5, $6)`,
		txID, txInput.Index, pq.Array(txInput.TxWitness), txInput.SignatureScript, txInput.PreviousOutPointHash, txInput.PreviousOutPointIndex)
	return err
}

func (in *CIDIndexer) indexTxOutput(tx *sqlx.Tx, txOuput TxOutput, txID int64) error {
	_, err := tx.Exec(`INSERT INTO btc.tx_outputs (tx_id, index, value, pk_script, script_class, addresses, required_sigs)
							VALUES ($1, $2, $3, $4, $5, $6, $7)
							ON CONFLICT (tx_id, index) DO UPDATE SET (value, pk_script, script_class, addresses, required_sigs) = ($3, $4, $5, $6, $7)`,
		txID, txOuput.Index, txOuput.Value, txOuput.PkScript, txOuput.ScriptClass, txOuput.Addresses, txOuput.RequiredSigs)
	return err
}
