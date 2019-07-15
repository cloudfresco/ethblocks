package ethblocks

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/ethereum/go-ethereum/common"
)

// TransactionLogTopicIntf interface
type TransactionLogTopicIntf interface {
	AddTransactionLogTopic(ctx context.Context, tx *sql.Tx, s common.Hash, BlockID uint, TransactionID uint, TransactionReceiptID uint, TransactionLogID uint) (*TransactionLogTopic, error)
	InsertTransactionLogTopic(ctx context.Context, tx *sql.Tx, lt *TransactionLogTopic) error
	GetTransactionLogTopics(ctx context.Context, TransactionLogID uint) ([]*TransactionLogTopic, error)
}

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
func AddTransactionLogTopic(ctx context.Context, tx *sql.Tx, s common.Hash, BlockID uint, TransactionID uint, TransactionReceiptID uint, TransactionLogID uint) (*TransactionLogTopic, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		lt := TransactionLogTopic{}
		lt.Topic = s.Hex()
		lt.BlockID = BlockID
		lt.TransactionID = TransactionID
		lt.TransactionReceiptID = TransactionReceiptID
		lt.TransactionLogID = TransactionLogID
		err := InsertTransactionLogTopic(ctx, tx, &lt)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return &lt, nil
	}
}

// InsertTransactionLogTopic - insert transaction Topic details to db
func InsertTransactionLogTopic(ctx context.Context, tx *sql.Tx, lt *TransactionLogTopic) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return err
	default:
		stmt, err := tx.PrepareContext(ctx, `insert into transaction_log_topics
	  ( 
			topic,
			block_id,
			transaction_id,
			transaction_receipt_id,
      transaction_log_id)
  values (?,?,?,?,?);`)
		if err != nil {
			log.Println(err)
			return err
		}
		res, err := stmt.ExecContext(ctx,
			lt.Topic,
			lt.BlockID,
			lt.TransactionID,
			lt.TransactionReceiptID,
			lt.TransactionLogID)
		if err != nil {
			log.Println(err)
			err = stmt.Close()
			return err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.Println(err)
			err = stmt.Close()
			return err
		}
		lt.ID = uint(uID)
		err = stmt.Close()
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	}
}

// GetTransactionLogTopics - used for getting topics by TransactionLogID
func GetTransactionLogTopics(ctx context.Context, TransactionLogID uint) ([]*TransactionLogTopic, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		appState, err := dbInit()
		if err != nil {
			log.Println(err)
			return nil, err
		}
		db := appState.Db
		topics := []*TransactionLogTopic{}
		rows, err := db.QueryContext(ctx, `select
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
			if err != nil {
				log.Println(err)
				return nil, err
			}
			topics = append(topics, &lt)
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
		return topics, nil
	}
}
