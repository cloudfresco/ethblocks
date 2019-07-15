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
	AddTransactionReceipt(ctx context.Context, tx *sql.Tx, ethReceipt *types.Receipt, BlockID uint, BlockNumber uint64, BlockHash string, TransactionID uint) (*TransactionReceipt, error)
	GetTransactionReceipts(ctx context.Context, TransactionID uint) ([]*TransactionReceipt, error)
}

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
func AddTransactionReceipt(ctx context.Context, tx *sql.Tx, ethReceipt *types.Receipt, BlockID uint, BlockNumber uint64, BlockHash string, TransactionID uint) (*TransactionReceipt, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		transReceipt := TransactionReceipt{}
		transReceipt.BlockNumber = BlockNumber
		transReceipt.BlockHash = BlockHash
		transReceipt.TxHash = ethReceipt.TxHash.Hex()
		transReceipt.TxStatus = ethReceipt.Status
		transReceipt.CumulativeGasUsed = ethReceipt.CumulativeGasUsed
		transReceipt.GasUsed = ethReceipt.GasUsed
		transReceipt.ContractAddress = ethReceipt.ContractAddress.Hex()
		transReceipt.PostState = ethReceipt.PostState
		transReceipt.BlockID = BlockID
		transReceipt.TransactionID = TransactionID
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
			return err
		}
		res, err := stmt.ExecContext(ctx,
			transReceipt.BlockNumber,
			transReceipt.BlockHash,
			transReceipt.TxHash,
			transReceipt.TxStatus,
			transReceipt.CumulativeGasUsed,
			transReceipt.GasUsed,
			transReceipt.ContractAddress,
			transReceipt.PostState,
			transReceipt.BlockID,
			transReceipt.TransactionID)
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
		transReceipt.ID = uint(uID)
		err = stmt.Close()
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	}
}

// GetTransactionReceipts - used for getting receipts by TransactionID
func GetTransactionReceipts(ctx context.Context, TransactionID uint) ([]*TransactionReceipt, error) {
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
		transreceipts := []*TransactionReceipt{}
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
			transReceipt := TransactionReceipt{}
			err = rows.Scan(
				&transReceipt.ID,
				&transReceipt.BlockNumber,
				&transReceipt.BlockHash,
				&transReceipt.TxHash,
				&transReceipt.TxStatus,
				&transReceipt.CumulativeGasUsed,
				&transReceipt.GasUsed,
				&transReceipt.ContractAddress,
				&transReceipt.PostState,
				&transReceipt.BlockID,
				&transReceipt.TransactionID)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			tlogs, err := GetTransactionLogs(ctx, transReceipt.ID)
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

		err = db.Close()
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return transreceipts, nil
	}
}
