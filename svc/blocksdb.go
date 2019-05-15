package svc

import (
	"context"
	"database/sql"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Block - Used for
type Block struct {
	ID           uint
	BlockNumber  uint64
	BlockTime    uint64
	ParentHash   string
	UncleHash    string
	BlockRoot    string
	TxHash       string
	ReceiptHash  string
	MixDigest    string
	BlockNonce   uint64
	Coinbase     string
	GasLimit     uint64
	GasUsed      uint64
	Difficulty   uint64
	BlockSize    common.StorageSize
	BlockUncles  []*BlockUncle
	Transactions []*Transaction
}

// AddBlock - add a block to the db
func AddBlock(ctx context.Context, client *ethclient.Client, block *types.Block) (*Block, error) {

	appState, err := dbInit()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	bl := Block{}
	bl.BlockNumber = block.Number().Uint64()
	bl.BlockTime = block.Time()
	bl.ParentHash = block.ParentHash().Hex()
	bl.UncleHash = block.UncleHash().Hex()
	bl.BlockRoot = block.Root().Hex()
	bl.TxHash = block.TxHash().Hex()
	bl.ReceiptHash = block.ReceiptHash().Hex()
	bl.MixDigest = block.MixDigest().Hex()
	bl.BlockNonce = block.Nonce()
	bl.Coinbase = block.Coinbase().Hex()
	bl.GasLimit = block.GasLimit()
	bl.GasUsed = block.GasUsed()
	bl.Difficulty = block.Difficulty().Uint64()
	bl.BlockSize = block.Size()
	db := appState.Db
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		err = tx.Rollback()
		return nil, err
	}
	blk, err := InsertBlock(tx, bl)
	if err != nil {
		log.Println(err)
		err = tx.Rollback()
		return nil, err
	}
	uncles := []*BlockUncle{}
	for _, blockuncle := range GetUncles(block) {
		uncle, err := AddBlockUncle(tx, blockuncle, blk.ID)
		if err != nil {
			log.Println(err)
			err = tx.Rollback()
			return nil, err
		}
		uncles = append(uncles, uncle)
	}
	transactions := []*Transaction{}
	for _, tns := range GetTransactions(block) {
		transaction, err := AddTransaction(tx, tns, blk.ID, bl.BlockNumber)
		if err != nil {
			log.Println(err)
			err = tx.Rollback()
			return nil, err
		}

		receipt, err := GetTransactionReceipt(ctx, client, tns.Hash())
		if err != nil {
			log.Println(err)
			err = tx.Rollback()
			return nil, err
		}
		receipts := []*TransactionReceipt{}
		treceipt, err := AddTransactionReceipt(tx, receipt, blk.ID, blk.BlockNumber, block.Hash().Hex(), transaction.ID)
		if err != nil {
			log.Println(err)
			err = tx.Rollback()
			return nil, err
		}
		receipts = append(receipts, treceipt)
		transaction.TransactionReceipts = receipts
		transactions = append(transactions, transaction)
	}
	blk.BlockUncles = uncles
	blk.Transactions = transactions
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		err = tx.Rollback()
		return nil, err
	}
	return blk, nil
}

// InsertBlock - insert block details to db
func InsertBlock(tx *sql.Tx, blk Block) (*Block, error) {
	stmt, err := tx.Prepare(`insert into blocks
	  ( 
			block_number,
			block_time,
			parent_hash,
			uncle_hash,
			block_root,
			tx_hash,
			receipt_hash,
			mix_digest,
			block_nonce,
			coinbase,
			gas_limit,
			gas_used,
			difficulty,
			block_size)
  values (?,?,?,?,?,?,?,?,?,?,
          ?,?,?,?);`)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	res, err := stmt.Exec(
		blk.BlockNumber,
		blk.BlockTime,
		blk.ParentHash,
		blk.UncleHash,
		blk.BlockRoot,
		blk.TxHash,
		blk.ReceiptHash,
		blk.MixDigest,
		blk.BlockNonce,
		blk.Coinbase,
		blk.GasLimit,
		blk.GasUsed,
		blk.Difficulty,
		blk.BlockSize)
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
	blk.ID = uint(uID)
	err = stmt.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &blk, nil
}

// GetBlock - used for
func GetBlock(ID uint) (*Block, error) {
	appState, err := dbInit()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	db := appState.Db
	blk := Block{}
	row := db.QueryRow(`select
      id,
			block_number,
			block_time,
			parent_hash,
			uncle_hash,
			block_root,
			tx_hash,
			receipt_hash,
			mix_digest,
			block_nonce,
			coinbase,
			gas_limit,
			gas_used,
			difficulty,
			block_size from blocks where id = ?;`, ID)

	err = row.Scan(
		&blk.ID,
		&blk.BlockNumber,
		&blk.BlockTime,
		&blk.ParentHash,
		&blk.UncleHash,
		&blk.BlockRoot,
		&blk.TxHash,
		&blk.ReceiptHash,
		&blk.MixDigest,
		&blk.BlockNonce,
		&blk.Coinbase,
		&blk.GasLimit,
		&blk.GasUsed,
		&blk.Difficulty,
		&blk.BlockSize)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	uncles, err := GetBlockUncles(ID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	transactions, err := GetBlockTransactions(ID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	blk.BlockUncles = uncles
	blk.Transactions = transactions
	err = db.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &blk, nil
}
