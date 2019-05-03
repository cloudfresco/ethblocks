package svc

import (
	"database/sql"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
)

// Transaction - Used for
type Transaction struct {
	ID                  uint
	BlockNumber         uint64
	BlockHash           string
	AccountNonce        uint64
	Price               uint64
	GasLimit            uint64
	TxAmount            uint64
	Payload             []byte
	TxV                 uint64
	TxR                 uint64
	TxS                 uint64
	BlockID             uint
	TransactionReceipts []*TransactionReceipt
}

// AddTransaction - add a transaction to the db
func AddTransaction(tx *sql.Tx, tns *types.Transaction, BlockID uint, BlockNumber uint64) (*Transaction, error) {
	txv, txr, txs := tns.RawSignatureValues()
	bl := Transaction{}
	bl.BlockNumber = BlockNumber
	bl.BlockHash = tns.Hash().Hex()
	bl.AccountNonce = tns.Nonce()
	bl.Price = tns.GasPrice().Uint64()
	bl.GasLimit = tns.Gas()
	bl.TxAmount = tns.Value().Uint64()
	bl.Payload = tns.Data()
	bl.TxV = txv.Uint64()
	bl.TxR = txr.Uint64()
	bl.TxS = txs.Uint64()
	bl.BlockID = BlockID
	transaction, err := InsertTransaction(tx, bl)
	if err != nil {
		log.Println(err)
		err = tx.Rollback()
		return nil, err
	}
	return transaction, nil
}

// InsertTransaction - insert transaction details to db
func InsertTransaction(tx *sql.Tx, trans Transaction) (*Transaction, error) {
	stmt, err := tx.Prepare(`insert into transactions
	  ( 
			block_number,
			block_hash,
			account_nonce,
			price,
			gas_limit,
			tx_amount,
			payload,
			tx_v,
			tx_r,
			tx_s,
			block_id)
  values (?,?,?,?,?,?,?,?,?,?,
          ?);`)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	res, err := stmt.Exec(
		trans.BlockNumber,
		trans.BlockHash,
		trans.AccountNonce,
		trans.Price,
		trans.GasLimit,
		trans.TxAmount,
		trans.Payload,
		trans.TxV,
		trans.TxR,
		trans.TxS,
		trans.BlockID)
	if err != nil {
		log.Println(err)
		err = stmt.Close()
		return nil, err
	}
	uID, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		err = stmt.Close()
		return nil, err
	}
	trans.ID = uint(uID)
	err = stmt.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &trans, nil
}

// GetBlockTransactions - used for
func GetBlockTransactions(BlockID uint) ([]*Transaction, error) {
	appState, err := dbInit()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	db := appState.Db
	transactions := []*Transaction{}
	rows, err := db.Query(`select
      id,
			block_number,
			block_hash,
			account_nonce,
			price,
			gas_limit,
			tx_amount,
			payload,
			tx_v,
			tx_r,
			tx_s,
			block_id from transactions where block_id = ?`, BlockID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		trans := Transaction{}
		err = rows.Scan(
			&trans.ID,
			&trans.BlockNumber,
			&trans.BlockHash,
			&trans.AccountNonce,
			&trans.Price,
			&trans.GasLimit,
			&trans.TxAmount,
			&trans.Payload,
			&trans.TxV,
			&trans.TxR,
			&trans.TxS,
			&trans.BlockID)

		receipts, err := GetTransactionReceipts(trans.ID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		trans.TransactionReceipts = receipts

		transactions = append(transactions, &trans)

	}
	err = rows.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = db.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return transactions, nil
}
