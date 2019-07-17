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
	GetTransactionLogTopics(ctx context.Context, TransactionLogID uint) ([]*TransactionLogTopic, error)
}

// TransactionLogTopicService - For accessing Transaction Log Topic services
type TransactionLogTopicService struct {
	Db *sql.DB
}

// NewTransactionLogTopicService - Create Transaction Log Topic service
func NewTransactionLogTopicService(db *sql.DB) *TransactionLogTopicService {
	return &TransactionLogTopicService{
		Db: db,
	}
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
func (t *TransactionLogTopicService) AddTransactionLogTopic(ctx context.Context, tx *sql.Tx, s common.Hash, BlockID uint, TransactionID uint, TransactionReceiptID uint, TransactionLogID uint) (*TransactionLogTopic, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		transLogTopic := TransactionLogTopic{}
		transLogTopic.Topic = s.Hex()
		transLogTopic.BlockID = BlockID
		transLogTopic.TransactionID = TransactionID
		transLogTopic.TransactionReceiptID = TransactionReceiptID
		transLogTopic.TransactionLogID = TransactionLogID
		err := insertTransactionLogTopic(ctx, tx, &transLogTopic)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return &transLogTopic, nil
	}
}

// insertTransactionLogTopic - insert transaction Topic details to db
func insertTransactionLogTopic(ctx context.Context, tx *sql.Tx, transLogTopic *TransactionLogTopic) error {
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
			transLogTopic.Topic,
			transLogTopic.BlockID,
			transLogTopic.TransactionID,
			transLogTopic.TransactionReceiptID,
			transLogTopic.TransactionLogID)
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
		transLogTopic.ID = uint(uID)
		err = stmt.Close()
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	}
}

// GetTransactionLogTopics - used for getting topics by TransactionLogID
func (t *TransactionLogTopicService) GetTransactionLogTopics(ctx context.Context, TransactionLogID uint) ([]*TransactionLogTopic, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		topics := []*TransactionLogTopic{}
		rows, err := t.Db.QueryContext(ctx, `select
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
			transLogTopic := TransactionLogTopic{}
			err = rows.Scan(
				&transLogTopic.ID,
				&transLogTopic.Topic,
				&transLogTopic.BlockID,
				&transLogTopic.TransactionID,
				&transLogTopic.TransactionReceiptID,
				&transLogTopic.TransactionLogID)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			topics = append(topics, &transLogTopic)
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

		return topics, nil
	}
}
