package svc

import (
	"database/sql"
	"log"

	"github.com/ethereum/go-ethereum/common"
)

// TransactionLogTopic - Used for
type TransactionLogTopic struct {
	ID                   uint
	Topic                string
	BlockID              uint
	TransactionID        uint
	TransactionReceiptID uint
	TransactionLogID     uint
}

// AddTransactionLogTopic - add a transaction Topic to the db
func AddTransactionLogTopic(tx *sql.Tx, s common.Hash, BlockID uint, TransactionID uint, TransactionReceiptID uint, TransactionLogID uint) (*TransactionLogTopic, error) {
	bl := TransactionLogTopic{}
	bl.Topic = s.Hex()
	bl.BlockID = BlockID
	bl.TransactionID = TransactionID
	bl.TransactionReceiptID = TransactionReceiptID
	bl.TransactionLogID = TransactionLogID
	transactiontopic, err := InsertTransactionLogTopic(tx, bl)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return transactiontopic, nil
}

// InsertTransactionLogTopic - insert transaction Topic details to db
func InsertTransactionLogTopic(tx *sql.Tx, lt TransactionLogTopic) (*TransactionLogTopic, error) {
	stmt, err := tx.Prepare(`insert into transaction_log_topics
	  ( 
			topic,
			block_id,
			transaction_id,
			transaction_receipt_id,
      transaction_log_id)
  values (?,?,?,?,?);`)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	res, err := stmt.Exec(
		lt.Topic,
		lt.BlockID,
		lt.TransactionID,
		lt.TransactionReceiptID,
		lt.TransactionLogID)
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
	lt.ID = uint(uID)
	err = stmt.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &lt, nil
}

// GetTransactionLogTopics - used for getting topics by TransactionLogID
func GetTransactionLogTopics(TransactionLogID uint) ([]*TransactionLogTopic, error) {
	appState, err := dbInit()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	db := appState.Db
	topics := []*TransactionLogTopic{}
	rows, err := db.Query(`select
      id,
      topic,
			block_id,
			transaction_id,
			transaction_receipt_id,
      transaction_log_id from transaction_log_topics where transaction_log_id = ?`, TransactionLogID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		lt := TransactionLogTopic{}
		err = rows.Scan(
			&lt.ID,
			&lt.Topic,
			&lt.BlockID,
			&lt.TransactionID,
			&lt.TransactionReceiptID,
			&lt.TransactionLogID)
		topics = append(topics, &lt)
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
	return topics, nil
}
