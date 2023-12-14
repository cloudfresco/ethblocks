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
	AddTransactionLog(ctx context.Context, tx *sql.Tx, ethLog *types.Log, blockId uint, transactionId uint, transactionReceiptId uint) (*TransactionLog, error)
	GetTransactionLogs(ctx context.Context, transactionReceiptId uint) ([]*TransactionLog, error)
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
	Id                   uint
	Address              string
	LogData              []byte
	BlockNumber          uint64
	TxHash               string
	TxIndex              uint
	BlockHash            string
	LogIndex             uint
	Removed              bool
	BlockId              uint
	TransactionId        uint
	TransactionReceiptId uint
	Topics               []*TransactionLogTopic
}

// AddTransactionLog - add a transaction log to the db
func (tl *TransactionLogService) AddTransactionLog(ctx context.Context, tx *sql.Tx, ethLog *types.Log, blockId uint, transactionId uint, transactionReceiptId uint) (*TransactionLog, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		transLog := TransactionLog{}
		transLog.Address = ethLog.Address.Hex()
		transLog.LogData = ethLog.Data
		transLog.BlockNumber = ethLog.BlockNumber
		transLog.TxHash = ethLog.TxHash.Hex()
		transLog.TxIndex = ethLog.TxIndex
		transLog.BlockHash = ethLog.BlockHash.Hex()
		transLog.LogIndex = ethLog.Index
		transLog.Removed = ethLog.Removed
		transLog.BlockId = blockId
		transLog.TransactionId = transactionId
		transLog.TransactionReceiptId = transactionReceiptId
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
	  ( address,
			log_data,
			block_number,
      tx_hash,
			tx_index,
			block_hash,			
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
			transLog.Address,
			transLog.LogData,
			transLog.BlockNumber,
			transLog.TxHash,
			transLog.TxIndex,
			transLog.BlockHash,
			transLog.LogIndex,
			transLog.Removed,
			transLog.BlockId,
			transLog.TransactionId,
			transLog.TransactionReceiptId)
		if err != nil {
			log.Println(err)
			err = stmt.Close()
			return err
		}
		uId, err := res.LastInsertId()
		if err != nil {
			log.Println(err)
			err = stmt.Close()
			return err
		}
		transLog.Id = uint(uId)
		err = stmt.Close()
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	}
}

// GetTransactionLogs - used for getting logs by TransactionReceiptId
func (tl *TransactionLogService) GetTransactionLogs(ctx context.Context, transactionReceiptId uint) ([]*TransactionLog, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		transLogs := []*TransactionLog{}
		rows, err := tl.Db.QueryContext(ctx, `select
      id,
	    address,
			log_data,
			block_number,
      tx_hash,
			tx_index,
			block_hash,	
			log_index,
			removed,
			block_id,
			transaction_id,
			transaction_receipt_id from transaction_logs where transaction_receipt_id = ?`, transactionReceiptId)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		for rows.Next() {
			transLog := TransactionLog{}
			err = rows.Scan(
				&transLog.Id,
				&transLog.Address,
				&transLog.LogData,
				&transLog.BlockNumber,
				&transLog.TxHash,
				&transLog.TxIndex,
				&transLog.BlockHash,
				&transLog.LogIndex,
				&transLog.Removed,
				&transLog.BlockId,
				&transLog.TransactionId,
				&transLog.TransactionReceiptId)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			transactionLogTopicService := TransactionLogTopicService{Db: tl.Db}
			topics, err := transactionLogTopicService.GetTransactionLogTopics(ctx, transLog.Id)
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
