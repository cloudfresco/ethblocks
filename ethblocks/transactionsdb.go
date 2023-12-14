package ethblocks

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
)

// TransactionIntf - interface
type TransactionIntf interface {
	AddTransaction(ctx context.Context, tx *sql.Tx, ethTrans *types.Transaction, blockId uint, blockNumber uint64) (*Transaction, error)
	GetBlockTransactions(ctx context.Context, blockId uint) ([]*Transaction, error)
}

// TransactionService - For accessing Transaction services
type TransactionService struct {
	Db *sql.DB
}

// NewTransactionService - Create Transaction service
func NewTransactionService(db *sql.DB) *TransactionService {
	return &TransactionService{
		Db: db,
	}
}

// Transaction - Used for
type Transaction struct {
	Id                  uint
	TxType              uint8
	ChainId             uint64
	TxData              []byte
	Gas                 uint64
	GasPrice            uint64
	GasTipCap           uint64
	GasFeeCap           uint64
	TxValue             uint64
	AccountNonce        uint64
	TxTo                string
	TxV                 uint64
	TxR                 uint64
	TxS                 uint64
	BlockNumber         uint64
	BlockHash           string
	BlockId             uint
	TransactionReceipts []*TransactionReceipt
}

// AddTransaction - add a transaction to the db
func (t *TransactionService) AddTransaction(ctx context.Context, tx *sql.Tx, ethTransaction *types.Transaction, blockId uint, blockNumber uint64) (*Transaction, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		txv, txr, txs := ethTransaction.RawSignatureValues()
		transaction := Transaction{}
		transaction.TxType = ethTransaction.Type()
		transaction.ChainId = ethTransaction.ChainId().Uint64()
		transaction.TxData = ethTransaction.Data()
		transaction.Gas = ethTransaction.Gas()
		transaction.GasPrice = ethTransaction.GasPrice().Uint64()
		transaction.GasTipCap = ethTransaction.GasTipCap().Uint64()
		transaction.GasFeeCap = ethTransaction.GasFeeCap().Uint64()
		transaction.TxValue = ethTransaction.Value().Uint64()
		transaction.AccountNonce = ethTransaction.Nonce()
		transaction.TxTo = ethTransaction.To().Hex()
		transaction.TxV = txv.Uint64()
		transaction.TxR = txr.Uint64()
		transaction.TxS = txs.Uint64()
		transaction.BlockNumber = blockNumber
		transaction.BlockHash = ethTransaction.Hash().Hex()
		transaction.BlockId = blockId
		err := insertTransaction(ctx, tx, &transaction)
		if err != nil {
			log.Println(err)
			err = tx.Rollback()
			return nil, err
		}
		return &transaction, nil
	}
}

// insertTransaction - insert transaction details to db
func insertTransaction(ctx context.Context, tx *sql.Tx, transaction *Transaction) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return err
	default:
		stmt, err := tx.PrepareContext(ctx, `insert into transactions
	  ( 
			tx_type,
      chain_id,
      tx_data,
      gas,
      gas_price,
      gas_tip_cap,
      gas_fee_cap,
      tx_value,
      account_nonce,
      tx_to,
      tx_v,
      tx_r,
      tx_s,
      block_number,
      block_hash,
			block_id)
  values (?,?,?,?,?,?,?,?,?,?,
          ?,?,?,?,?,?);`)
		if err != nil {
			log.Println(err)
			return err
		}
		res, err := stmt.ExecContext(ctx,
			transaction.TxType,
			transaction.ChainId,
			transaction.TxData,
			transaction.Gas,
			transaction.GasPrice,
			transaction.GasTipCap,
			transaction.GasFeeCap,
			transaction.TxValue,
			transaction.AccountNonce,
			transaction.TxTo,
			transaction.TxV,
			transaction.TxR,
			transaction.TxS,
			transaction.BlockNumber,
			transaction.BlockHash,
			transaction.BlockId)
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
		transaction.Id = uint(uId)
		err = stmt.Close()
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	}
}

// GetBlockTransactions - used for
func (t *TransactionService) GetBlockTransactions(ctx context.Context, blockId uint) ([]*Transaction, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		transactions := []*Transaction{}
		rows, err := t.Db.QueryContext(ctx, `select
      id,
			tx_type,
      chain_id,
      tx_data,
      gas,
      gas_price,
      gas_tip_cap,
      gas_fee_cap,
      tx_value,
      account_nonce,
      tx_to,
      tx_v,
      tx_r,
      tx_s,
      block_number,
      block_hash,
			block_id from transactions where block_id = ?`, blockId)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		for rows.Next() {
			transaction := Transaction{}
			err = rows.Scan(
				&transaction.Id,
				&transaction.TxType,
				&transaction.ChainId,
				&transaction.TxData,
				&transaction.Gas,
				&transaction.GasPrice,
				&transaction.GasTipCap,
				&transaction.GasFeeCap,
				&transaction.TxValue,
				&transaction.AccountNonce,
				&transaction.TxTo,
				&transaction.TxV,
				&transaction.TxR,
				&transaction.TxS,
				&transaction.BlockNumber,
				&transaction.BlockHash,
				&transaction.BlockId)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			transactionReceiptService := TransactionReceiptService{Db: t.Db}
			receipts, err := transactionReceiptService.GetTransactionReceipts(ctx, transaction.Id)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			transaction.TransactionReceipts = receipts

			transactions = append(transactions, &transaction)

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

		return transactions, nil
	}
}
