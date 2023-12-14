package ethblocks

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
)

// TransactionReceiptIntf interface
type TransactionReceiptIntf interface {
	AddTransactionReceipt(ctx context.Context, tx *sql.Tx, ethReceipt *types.Receipt, BlockId uint, BlockNumber uint64, BlockHash string, TransactionId uint) (*TransactionReceipt, error)
	GetTransactionReceipts(ctx context.Context, transactionId uint) ([]*TransactionReceipt, error)
}

// TransactionReceiptService - For accessing Transaction Receipt services
type TransactionReceiptService struct {
	Db *sql.DB
}

// NewTransactionReceiptService - Create Transaction Receipt service
func NewTransactionReceiptService(db *sql.DB) *TransactionReceiptService {
	return &TransactionReceiptService{
		Db: db,
	}
}

// TransactionReceipt - Used for
type TransactionReceipt struct {
	Id                uint
	ReceiptType       uint8
	PostState         []byte
	TxStatus          uint64
	CumulativeGasUsed uint64
	Bloom             []byte
	TxHash            string
	ContractAddress   string
	GasUsed           uint64
	EffectiveGasPrice uint64
	BlobGasUsed       uint64
	BlobGasPrice      uint64
	BlockHash         string
	BlockNumber       uint64
	TransactionIndex  uint
	BlockId           uint
	TransactionId     uint
	Logs              []*TransactionLog
}

// AddTransactionReceipt - add a transaction to the db
func (t *TransactionReceiptService) AddTransactionReceipt(ctx context.Context, tx *sql.Tx, ethReceipt *types.Receipt, blockId uint, blockNumber uint64, blockHash string, transactionId uint) (*TransactionReceipt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		transReceipt := TransactionReceipt{}
		transReceipt.ReceiptType = ethReceipt.Type
		transReceipt.PostState = ethReceipt.PostState
		transReceipt.TxStatus = ethReceipt.Status
		transReceipt.CumulativeGasUsed = ethReceipt.CumulativeGasUsed
		transReceipt.Bloom = ethReceipt.Bloom.Bytes()
		transReceipt.TxHash = ethReceipt.TxHash.Hex()
		transReceipt.ContractAddress = ethReceipt.ContractAddress.Hex()
		transReceipt.GasUsed = ethReceipt.GasUsed
		transReceipt.EffectiveGasPrice = ethReceipt.EffectiveGasPrice.Uint64()
		transReceipt.BlobGasUsed = ethReceipt.BlobGasUsed
		if ethReceipt.BlobGasPrice != nil {
			transReceipt.BlobGasPrice = ethReceipt.BlobGasPrice.Uint64()
		}
		transReceipt.BlockHash = blockHash
		transReceipt.BlockNumber = blockNumber
		transReceipt.TransactionIndex = ethReceipt.TransactionIndex
		transReceipt.BlockId = blockId
		transReceipt.TransactionId = transactionId
		err := insertTransactionReceipt(ctx, tx, &transReceipt)
		if err != nil {
			log.Println(err)
			err = tx.Rollback()
			return nil, err
		}
		return &transReceipt, nil
	}
}

// insertTransactionReceipt - insert transaction Receipt details to db
func insertTransactionReceipt(ctx context.Context, tx *sql.Tx, transReceipt *TransactionReceipt) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return err
	default:
		stmt, err := tx.PrepareContext(ctx, `insert into transaction_receipts
	  ( 
			receipt_type,
      post_state,
      tx_status,
      cumulative_gas_used,
      bloom,
      tx_hash,
      contract_address,
      gas_used,
      effective_gas_price,
      blob_gas_used,
      blob_gas_price,
      block_hash,
      block_number,
      transaction_index,
      block_id,
      transaction_id)
  values (?,?,?,?,?,?,?,?,?,?,
          ?,?,?,?,?,?);`)
		if err != nil {
			log.Println(err)
			return err
		}
		res, err := stmt.ExecContext(ctx,
			transReceipt.ReceiptType,
			transReceipt.PostState,
			transReceipt.TxStatus,
			transReceipt.CumulativeGasUsed,
			transReceipt.Bloom,
			transReceipt.TxHash,
			transReceipt.ContractAddress,
			transReceipt.GasUsed,
			transReceipt.EffectiveGasPrice,
			transReceipt.BlobGasUsed,
			transReceipt.BlobGasPrice,
			transReceipt.BlockHash,
			transReceipt.BlockNumber,
			transReceipt.TransactionIndex,
			transReceipt.BlockId,
			transReceipt.TransactionId)
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
		transReceipt.Id = uint(uId)
		err = stmt.Close()
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	}
}

// GetTransactionReceipts - used for getting receipts by TransactionId
func (t *TransactionReceiptService) GetTransactionReceipts(ctx context.Context, transactionId uint) ([]*TransactionReceipt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		transreceipts := []*TransactionReceipt{}
		rows, err := t.Db.QueryContext(ctx, `select
      id,
			receipt_type,
      post_state,
      tx_status,
      cumulative_gas_used,
      bloom,
      tx_hash,
      contract_address,
      gas_used,
      effective_gas_price,
      blob_gas_used,
      blob_gas_price,
      block_hash,
      block_number,
      transaction_index,
      block_id,
      transaction_id from transaction_receipts where transaction_id = ?`, transactionId)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		for rows.Next() {
			transReceipt := TransactionReceipt{}
			err = rows.Scan(
				&transReceipt.Id,
				&transReceipt.ReceiptType,
				&transReceipt.PostState,
				&transReceipt.TxStatus,
				&transReceipt.CumulativeGasUsed,
				&transReceipt.Bloom,
				&transReceipt.TxHash,
				&transReceipt.ContractAddress,
				&transReceipt.GasUsed,
				&transReceipt.EffectiveGasPrice,
				&transReceipt.BlobGasUsed,
				&transReceipt.BlobGasPrice,
				&transReceipt.BlockHash,
				&transReceipt.BlockNumber,
				&transReceipt.TransactionIndex,
				&transReceipt.BlockId,
				&transReceipt.TransactionId)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			transactionLogService := TransactionLogService{Db: t.Db}
			tlogs, err := transactionLogService.GetTransactionLogs(ctx, transReceipt.Id)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			transReceipt.Logs = tlogs
			transreceipts = append(transreceipts, &transReceipt)
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

		return transreceipts, nil
	}
}
