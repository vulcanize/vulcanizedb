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
	headerID, err := in.indexHeaderCID(tx, cidWrapper.HeaderCID, in.db.NodeID)
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

func (in *CIDIndexer) indexHeaderCID(tx *sqlx.Tx, header HeaderModel, nodeID int64) (int64, error) {
	var headerID int64
	err := tx.QueryRowx(`INSERT INTO btc.header_cids (block_number, block_hash, parent_hash, cid, timestamp, bits, node_id)
							VALUES ($1, $2, $3, $4, $5, $6, $7) 
							ON CONFLICT (block_number, block_hash) DO UPDATE SET (parent_hash, cid, timestamp, bits, node_id) = ($3, $4, $5, $6, $7)
							RETURNING id`,
		header.BlockNumber, header.BlockHash, header.ParentHash, header.CID, header.Timestamp, header.Bits, nodeID).Scan(&headerID)
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
	// resolve zero-value hash to null value (coinbase tx input with no referenced outputs)
	if txInput.PreviousOutPointHash == "0000000000000000000000000000000000000000000000000000000000000000" {
		_, err := tx.Exec(`INSERT INTO btc.tx_inputs (tx_id, index, witness, sig_script, outpoint_id)
						VALUES ($1, $2, $3, $4, $5)
						ON CONFLICT (tx_id, index) DO UPDATE SET (witness, sig_script, outpoint_id) = ($3, $4, $5)`,
			txID, txInput.Index, pq.Array(txInput.TxWitness), txInput.SignatureScript, nil)
		return err
	}
	var referencedOutPutID int64
	if err := tx.Get(&referencedOutPutID, `SELECT tx_outputs.id FROM btc.transaction_cids, btc.tx_outputs
						WHERE tx_outputs.tx_id = transaction_cids.id
						AND tx_hash = $1
						AND tx_outputs.index = $2`, txInput.PreviousOutPointHash, txInput.PreviousOutPointIndex); err != nil {
		logrus.Errorf("btc indexer could not find the tx output (tx hash %s, output index %d) referenced in tx input %d of tx id %d", txInput.PreviousOutPointHash, txInput.PreviousOutPointIndex, txInput.Index, txID)
		return err
	}
	if referencedOutPutID == 0 {
		return fmt.Errorf("btc indexer could not find the tx output (tx hash %s, output index %d) referenced in tx input %d of tx id %d", txInput.PreviousOutPointHash, txInput.PreviousOutPointIndex, txInput.Index, txID)
	}
	_, err := tx.Exec(`INSERT INTO btc.tx_inputs (tx_id, index, witness, sig_script, outpoint_id)
						VALUES ($1, $2, $3, $4, $5)
						ON CONFLICT (tx_id, index) DO UPDATE SET (witness, sig_script, outpoint_id) = ($3, $4, $5)`,
		txID, txInput.Index, pq.Array(txInput.TxWitness), txInput.SignatureScript, referencedOutPutID)
	return err
}

func (in *CIDIndexer) indexTxOutput(tx *sqlx.Tx, txOuput TxOutput, txID int64) error {
	_, err := tx.Exec(`INSERT INTO btc.tx_outputs (tx_id, index, value, pk_script, script_class, addresses, required_sigs)
							VALUES ($1, $2, $3, $4, $5, $6, $7)
							ON CONFLICT (tx_id, index) DO UPDATE SET (value, pk_script, script_class, addresses, required_sigs) = ($3, $4, $5, $6, $7)`,
		txID, txOuput.Index, txOuput.Value, txOuput.PkScript, txOuput.ScriptClass, txOuput.Addresses, txOuput.RequiredSigs)
	return err
}
