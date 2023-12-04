package ethblocks

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
)

// TransactionLogIntf - interface
type TransactionLogIntf interface {
	AddTransactionLog(ctx context.Context, tx *sql.Tx, ethLog *types.Log, BlockID uint, TransactionID uint, TransactionReceiptID uint) (*TransactionLog, error)
	GetTransactionLogs(ctx context.Context, TransactionReceiptID uint) ([]*TransactionLog, error)
}

// TransactionLogService - For accessing Transaction Log services
type TransactionLogService struct {
	Db *sql.DB
}

// NewTransactionLogService - Create Transaction Log service
func NewTransactionLogService(db *sql.DB) *TransactionLogService {
	return &TransactionLogService{
		Db: db,
	}
}

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
func (tl *TransactionLogService) AddTransactionLog(ctx context.Context, tx *sql.Tx, ethLog *types.Log, BlockID uint, TransactionID uint, TransactionReceiptID uint) (*TransactionLog, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		transLog := TransactionLog{}
		transLog.BlockNumber = ethLog.BlockNumber
		transLog.BlockHash = ethLog.BlockHash.Hex()
		transLog.Address = ethLog.Address.Hex()
		transLog.LogData = ethLog.Data
		transLog.TxHash = ethLog.TxHash.Hex()
		transLog.TxIndex = ethLog.TxIndex
		transLog.LogIndex = ethLog.Index
		transLog.Removed = ethLog.Removed
		transLog.BlockID = BlockID
		transLog.TransactionID = TransactionID
		transLog.TransactionReceiptID = TransactionReceiptID
		err := insertTransactionLog(ctx, tx, &transLog)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return &transLog, nil
	}
}

// insertTransactionLog - insert transaction Log details to db
func insertTransactionLog(ctx context.Context, tx *sql.Tx, transLog *TransactionLog) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return err
	default:
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
			return err
		}
		res, err := stmt.ExecContext(ctx,
			transLog.BlockNumber,
			transLog.BlockHash,
			transLog.Address,
			transLog.LogData,
			transLog.TxHash,
			transLog.TxIndex,
			transLog.LogIndex,
			transLog.Removed,
			transLog.BlockID,
			transLog.TransactionID,
			transLog.TransactionReceiptID)
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
		transLog.ID = uint(uID)
		err = stmt.Close()
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	}
}

// GetTransactionLogs - used for getting logs by TransactionReceiptID
func (tl *TransactionLogService) GetTransactionLogs(ctx context.Context, TransactionReceiptID uint) ([]*TransactionLog, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		transLogs := []*TransactionLog{}
		rows, err := tl.Db.QueryContext(ctx, `select
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
			transLog := TransactionLog{}
			err = rows.Scan(
				&transLog.ID,
				&transLog.BlockNumber,
				&transLog.BlockHash,
				&transLog.Address,
				&transLog.LogData,
				&transLog.TxHash,
				&transLog.TxIndex,
				&transLog.LogIndex,
				&transLog.Removed,
				&transLog.BlockID,
				&transLog.TransactionID,
				&transLog.TransactionReceiptID)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			transactionLogTopicService := TransactionLogTopicService{Db: tl.Db}
			topics, err := transactionLogTopicService.GetTransactionLogTopics(ctx, transLog.ID)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			transLog.Topics = topics
			transLogs = append(transLogs, &transLog)
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
		return transLogs, nil
	}
}
