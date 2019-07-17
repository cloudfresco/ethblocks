package ethblocks

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// BlockServiceIntf - BlockIntf interface for BlockService
type BlockServiceIntf interface {
	AddBlock(ctx context.Context, client *ethclient.Client, block *types.Block) (*Block, error)
	CreateBlockTransaction(ctx context.Context, client *ethclient.Client, tx *sql.Tx, blk *Block, block *types.Block) error
	GetBlock(ctx context.Context, ID uint) (*Block, error)
}

// BlockService - For accessing block services
type BlockService struct {
	Db *sql.DB
}

// NewBlockService - Create block service
func NewBlockService(db *sql.DB) *BlockService {
	return &BlockService{
		Db: db,
	}
}

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
func (b *BlockService) AddBlock(ctx context.Context, client *ethclient.Client, block *types.Block) (*Block, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		blk := Block{}
		blk.BlockNumber = block.Number().Uint64()
		blk.BlockTime = block.Time()
		blk.ParentHash = block.ParentHash().Hex()
		blk.UncleHash = block.UncleHash().Hex()
		blk.BlockRoot = block.Root().Hex()
		blk.TxHash = block.TxHash().Hex()
		blk.ReceiptHash = block.ReceiptHash().Hex()
		blk.MixDigest = block.MixDigest().Hex()
		blk.BlockNonce = block.Nonce()
		blk.Coinbase = block.Coinbase().Hex()
		blk.GasLimit = block.GasLimit()
		blk.GasUsed = block.GasUsed()
		blk.Difficulty = block.Difficulty().Uint64()
		blk.BlockSize = block.Size()
		tx, err := b.Db.Begin()
		if err != nil {
			log.Println(err)
			err = tx.Rollback()
			return nil, err
		}
		err = insertBlock(ctx, tx, &blk)
		if err != nil {
			log.Println(err)
			err = tx.Rollback()
			return nil, err
		}
		uncles := []*BlockUncle{}
		for _, blockuncle := range GetUncles(block) {
			blockUncleService := BlockUncleService{Db: b.Db}
			uncle, err := blockUncleService.AddBlockUncle(ctx, tx, blockuncle, blk.ID)
			if err != nil {
				log.Println(err)
				err = tx.Rollback()
				return nil, err
			}
			uncles = append(uncles, uncle)
		}
		blk.BlockUncles = uncles
		err = b.CreateBlockTransaction(ctx, client, tx, &blk, block)
		if err != nil {
			log.Println(err)
			err = tx.Rollback()
			return nil, err
		}
		err = tx.Commit()
		if err != nil {
			log.Println(err)
			err = tx.Rollback()
			return nil, err
		}
		return &blk, err
	}
}

// CreateBlockTransaction - add a block transaction to the db
func (b *BlockService) CreateBlockTransaction(ctx context.Context, client *ethclient.Client, tx *sql.Tx, blk *Block, block *types.Block) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return err
	default:
		transactions := []*Transaction{}
		for _, tns := range GetTransactions(block) {
			transactionService := TransactionService{Db: b.Db}
			transaction, err := transactionService.AddTransaction(ctx, tx, tns, blk.ID, blk.BlockNumber)
			if err != nil {
				log.Println(err)
				err = tx.Rollback()
				return err
			}

			receipt, err := GetTransactionReceipt(ctx, client, tns.Hash())
			if err != nil {
				log.Println(err)
				err = tx.Rollback()
				return err
			}
			receipts := []*TransactionReceipt{}
			transactionReceiptService := TransactionReceiptService{Db: b.Db}
			treceipt, err := transactionReceiptService.AddTransactionReceipt(ctx, tx, receipt, blk.ID, blk.BlockNumber, block.Hash().Hex(), transaction.ID)
			if err != nil {
				log.Println(err)
				err = tx.Rollback()
				return err
			}
			tlogs := []*TransactionLog{}
			for _, lg := range GetLogs(receipt) {
				transactionLogService := TransactionLogService{Db: b.Db}
				tlg, err := transactionLogService.AddTransactionLog(ctx, tx, lg, blk.ID, transaction.ID, treceipt.ID)
				if err != nil {
					log.Println(err)
					err = tx.Rollback()
					return err
				}
				topics := []*TransactionLogTopic{}
				for _, tpc := range GetTopics(lg) {
					transactionLogTopicService := TransactionLogTopicService{Db: b.Db}
					topic, err := transactionLogTopicService.AddTransactionLogTopic(ctx, tx, tpc, blk.ID, transaction.ID, treceipt.ID, tlg.ID)
					if err != nil {
						log.Println(err)
						err = tx.Rollback()
						return err
					}
					topics = append(topics, topic)
				}
				tlg.Topics = topics
				tlogs = append(tlogs, tlg)
			}
			treceipt.Logs = tlogs
			receipts = append(receipts, treceipt)
			transaction.TransactionReceipts = receipts
			transactions = append(transactions, transaction)
		}
		blk.Transactions = transactions
		return nil
	}
}

// insertBlock - insert block details to db
func insertBlock(ctx context.Context, tx *sql.Tx, blk *Block) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return err
	default:
		stmt, err := tx.PrepareContext(ctx, `insert into blocks
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
			return err
		}
		res, err := stmt.ExecContext(ctx,
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
			return err
		}
		uID, err := res.LastInsertId()
		if err != nil {
			log.Println(err)
			err = stmt.Close()
			return err
		}
		blk.ID = uint(uID)
		err = stmt.Close()
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	}
}

// GetBlock - used for
func (b *BlockService) GetBlock(ctx context.Context, ID uint) (*Block, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		blk := Block{}
		row := b.Db.QueryRowContext(ctx, `select
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

		err := row.Scan(
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
		blockUncleService := BlockUncleService{Db: b.Db}
		uncles, err := blockUncleService.GetBlockUncles(ctx, ID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		transactionsService := TransactionService{Db: b.Db}
		transactions, err := transactionsService.GetBlockTransactions(ctx, ID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		blk.BlockUncles = uncles
		blk.Transactions = transactions

		return &blk, nil
	}
}
