package svc

import (
	"context"
	"database/sql"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
)

// TransactionReceipt - Used for
type TransactionReceipt struct {
	ID                uint
	BlockNumber       uint64
	BlockHash         string
	TxHash            string
	TxStatus          uint64
	CumulativeGasUsed uint64
	GasUsed           uint64
	ContractAddress   string
	PostState         []byte
	BlockID           uint
	TransactionID     uint
	Logs              []*TransactionLog
}

// AddTransactionReceipt - add a transaction to the db
func AddTransactionReceipt(ctx context.Context, tx *sql.Tx, receipt *types.Receipt, BlockID uint, BlockNumber uint64, BlockHash string, TransactionID uint) (*TransactionReceipt, error) {
	bl := TransactionReceipt{}
	bl.BlockNumber = BlockNumber
	bl.BlockHash = BlockHash
	bl.TxHash = receipt.TxHash.Hex()
	bl.TxStatus = receipt.Status
	bl.CumulativeGasUsed = receipt.CumulativeGasUsed
	bl.GasUsed = receipt.GasUsed
	bl.ContractAddress = receipt.ContractAddress.Hex()
	bl.PostState = receipt.PostState
	bl.BlockID = BlockID
	bl.TransactionID = TransactionID
	transactionReceipt, err := InsertTransactionReceipt(ctx, tx, bl)
	if err != nil {
		log.Println(err)
		err = tx.Rollback()
		return nil, err
	}
	return transactionReceipt, nil
}

// InsertTransactionReceipt - insert transaction receipt details to db
func InsertTransactionReceipt(ctx context.Context, tx *sql.Tx, receipt TransactionReceipt) (*TransactionReceipt, error) {
	stmt, err := tx.PrepareContext(ctx, `insert into transaction_receipts
	  ( 
			block_number,
			block_hash,
			tx_hash,
			tx_status,
			cumulative_gas_used,
			gas_used,
			contract_address,
			post_state,
			block_id,
			transaction_id)
  values (?,?,?,?,?,?,?,?,?,?);`)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	res, err := stmt.ExecContext(ctx,
		receipt.BlockNumber,
		receipt.BlockHash,
		receipt.TxHash,
		receipt.TxStatus,
		receipt.CumulativeGasUsed,
		receipt.GasUsed,
		receipt.ContractAddress,
		receipt.PostState,
		receipt.BlockID,
		receipt.TransactionID)
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
	receipt.ID = uint(uID)
	err = stmt.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &receipt, nil
}

// GetTransactionReceipts - used for getting receipts by TransactionID
func GetTransactionReceipts(ctx context.Context, TransactionID uint) ([]*TransactionReceipt, error) {
	appState, err := dbInit()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	db := appState.Db
	receipts := []*TransactionReceipt{}
	rows, err := db.QueryContext(ctx, `select
      id,
			block_number,
			block_hash,
			tx_hash,
			tx_status,
			cumulative_gas_used,
			gas_used,
			contract_address,
			post_state,
			block_id,
			transaction_id from transaction_receipts where transaction_id = ?`, TransactionID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		receipt := TransactionReceipt{}
		err = rows.Scan(
			&receipt.ID,
			&receipt.BlockNumber,
			&receipt.BlockHash,
			&receipt.TxHash,
			&receipt.TxStatus,
			&receipt.CumulativeGasUsed,
			&receipt.GasUsed,
			&receipt.ContractAddress,
			&receipt.PostState,
			&receipt.BlockID,
			&receipt.TransactionID)
		tlogs, err := GetTransactionLogs(ctx, receipt.ID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		receipt.Logs = tlogs
		receipts = append(receipts, &receipt)
	}
	err = rows.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = db.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return receipts, nil
}
