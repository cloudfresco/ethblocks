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
	AddTransaction(ctx context.Context, tx *sql.Tx, ethTrans *types.Transaction, BlockID uint, BlockNumber uint64) (*Transaction, error)
	GetBlockTransactions(ctx context.Context, BlockID uint) ([]*Transaction, error)
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
	ID                  uint
	BlockNumber         uint64
	BlockHash           string
	AccountNonce        uint64
	Price               uint64
	GasLimit            uint64
	TxAmount            uint64
	Payload             []byte
	TxV                 uint64
	TxR                 uint64
	TxS                 uint64
	BlockID             uint
	TransactionReceipts []*TransactionReceipt
}

// AddTransaction - add a transaction to the db
func (t *TransactionService) AddTransaction(ctx context.Context, tx *sql.Tx, ethTrans *types.Transaction, BlockID uint, BlockNumber uint64) (*Transaction, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		txv, txr, txs := ethTrans.RawSignatureValues()
		trans := Transaction{}
		trans.BlockNumber = BlockNumber
		trans.BlockHash = ethTrans.Hash().Hex()
		trans.AccountNonce = ethTrans.Nonce()
		trans.Price = ethTrans.GasPrice().Uint64()
		trans.GasLimit = ethTrans.Gas()
		trans.TxAmount = ethTrans.Value().Uint64()
		trans.Payload = ethTrans.Data()
		trans.TxV = txv.Uint64()
		trans.TxR = txr.Uint64()
		trans.TxS = txs.Uint64()
		trans.BlockID = BlockID
		err := insertTransaction(ctx, tx, &trans)
		if err != nil {
			log.Println(err)
			err = tx.Rollback()
			return nil, err
		}
		return &trans, nil
	}
}

// insertTransaction - insert transaction details to db
func insertTransaction(ctx context.Context, tx *sql.Tx, trans *Transaction) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return err
	default:
		stmt, err := tx.PrepareContext(ctx, `insert into transactions
	  ( 
			block_number,
			block_hash,
			account_nonce,
			price,
			gas_limit,
			tx_amount,
			payload,
			tx_v,
			tx_r,
			tx_s,
			block_id)
  values (?,?,?,?,?,?,?,?,?,?,
          ?);`)
		if err != nil {
			log.Println(err)
			return err
		}
		res, err := stmt.ExecContext(ctx,
			trans.BlockNumber,
			trans.BlockHash,
			trans.AccountNonce,
			trans.Price,
			trans.GasLimit,
			trans.TxAmount,
			trans.Payload,
			trans.TxV,
			trans.TxR,
			trans.TxS,
			trans.BlockID)
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
		trans.ID = uint(uID)
		err = stmt.Close()
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	}
}

// GetBlockTransactions - used for
func (t *TransactionService) GetBlockTransactions(ctx context.Context, BlockID uint) ([]*Transaction, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		transactions := []*Transaction{}
		rows, err := t.Db.QueryContext(ctx, `select
      id,
			block_number,
			block_hash,
			account_nonce,
			price,
			gas_limit,
			tx_amount,
			payload,
			tx_v,
			tx_r,
			tx_s,
			block_id from transactions where block_id = ?`, BlockID)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		for rows.Next() {
			trans := Transaction{}
			err = rows.Scan(
				&trans.ID,
				&trans.BlockNumber,
				&trans.BlockHash,
				&trans.AccountNonce,
				&trans.Price,
				&trans.GasLimit,
				&trans.TxAmount,
				&trans.Payload,
				&trans.TxV,
				&trans.TxR,
				&trans.TxS,
				&trans.BlockID)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			transactionReceiptService := TransactionReceiptService{Db: t.Db}
			receipts, err := transactionReceiptService.GetTransactionReceipts(ctx, trans.ID)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			trans.TransactionReceipts = receipts

			transactions = append(transactions, &trans)

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
