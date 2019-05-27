package svc

import (
	"context"
	"database/sql"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
)

// TransactionLog - Used for
type TransactionLog struct {
	ID                   uint
	BlockNumber          uint64
	BlockHash            string
	Address              string
	LogData              []byte
	TxHash               string
	TxIndex              uint
	LogIndex             uint
	Removed              bool
	BlockID              uint
	TransactionID        uint
	TransactionReceiptID uint
	Topics               []*TransactionLogTopic
}

// AddTransactionLog - add a transaction log to the db
func AddTransactionLog(ctx context.Context, tx *sql.Tx, lg *types.Log, BlockID uint, TransactionID uint, TransactionReceiptID uint) (*TransactionLog, error) {
	bl := TransactionLog{}
	bl.BlockNumber = lg.BlockNumber
	bl.BlockHash = lg.BlockHash.Hex()
	bl.Address = lg.Address.Hex()
	bl.LogData = lg.Data
	bl.TxHash = lg.TxHash.Hex()
	bl.TxIndex = lg.TxIndex
	bl.LogIndex = lg.Index
	bl.Removed = lg.Removed
	bl.BlockID = BlockID
	bl.TransactionID = TransactionID
	bl.TransactionReceiptID = TransactionReceiptID
	transactionLog, err := InsertTransactionLog(ctx, tx, bl)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return transactionLog, nil
}

// InsertTransactionLog - insert transaction Log details to db
func InsertTransactionLog(ctx context.Context, tx *sql.Tx, lg TransactionLog) (*TransactionLog, error) {
	stmt, err := tx.PrepareContext(ctx, `insert into transaction_logs
	  ( 
			block_number,
			block_hash,
			address,
			log_data,
			tx_hash,
			tx_index,
			log_index,
			removed,
			block_id,
			transaction_id,
			transaction_receipt_id)
  values (?,?,?,?,?,?,?,?,?,?,
          ?);`)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	res, err := stmt.ExecContext(ctx,
		lg.BlockNumber,
		lg.BlockHash,
		lg.Address,
		lg.LogData,
		lg.TxHash,
		lg.TxIndex,
		lg.LogIndex,
		lg.Removed,
		lg.BlockID,
		lg.TransactionID,
		lg.TransactionReceiptID)
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
	lg.ID = uint(uID)
	err = stmt.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &lg, nil
}

// GetTransactionLogs - used for getting logs by TransactionReceiptID
func GetTransactionLogs(ctx context.Context, TransactionReceiptID uint) ([]*TransactionLog, error) {
	appState, err := dbInit()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	db := appState.Db
	tlogs := []*TransactionLog{}
	rows, err := db.QueryContext(ctx, `select
      id,
			block_number,
			block_hash,
			address,
			log_data,
			tx_hash,
			tx_index,
			log_index,
			removed,
			block_id,
			transaction_id,
			transaction_receipt_id from transaction_logs where transaction_receipt_id = ?`, TransactionReceiptID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		lg := TransactionLog{}
		err = rows.Scan(
			&lg.ID,
			&lg.BlockNumber,
			&lg.BlockHash,
			&lg.Address,
			&lg.LogData,
			&lg.TxHash,
			&lg.TxIndex,
			&lg.LogIndex,
			&lg.Removed,
			&lg.BlockID,
			&lg.TransactionID,
			&lg.TransactionReceiptID)
		topics, err := GetTransactionLogTopics(ctx, lg.ID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		lg.Topics = topics
		tlogs = append(tlogs, &lg)
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
	return tlogs, nil
}
