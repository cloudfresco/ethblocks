package ethblocks

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// BlockServiceIntf - BlockIntf interface for BlockService
type BlockServiceIntf interface {
	AddBlock(ctx context.Context, client *ethclient.Client, block *types.Block) (*Block, error)
	CreateBlockTransaction(ctx context.Context, client *ethclient.Client, tx *sql.Tx, blk *Block, block *types.Block) error
	GetBlock(ctx context.Context, id uint) (*Block, error)
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
	Id               uint
	ParentHash       string
	UncleHash        string
	Coinbase         string
	BlockRoot        string
	TxHash           string
	ReceiptHash      string
	Bloom            []byte
	Difficulty       uint64
	BlockNumber      uint64
	GasLimit         uint64
	GasUsed          uint64
	BlockTime        uint64
	Extra            []byte
	MixDigest        string
	BlockNonce       uint64
	BaseFee          uint64
	WithdrawalsHash  string
	BlobGasUsed      *uint64
	ExcessBlobGas    *uint64
	ParentBeaconRoot string
	BlockHash        string
	BlockSize        uint64
	ReceivedAt       time.Time
	BlockUncles      []*BlockUncle
	Transactions     []*Transaction
}

// AddBlock - add a block to the db
func (b *BlockService) AddBlock(ctx context.Context, client *ethclient.Client, block *types.Block) (*Block, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		blk := Block{}
		blk.ParentHash = block.ParentHash().Hex()
		blk.UncleHash = block.UncleHash().Hex()
		blk.Coinbase = block.Coinbase().Hex()
		blk.BlockRoot = block.Root().Hex()
		blk.TxHash = block.TxHash().Hex()
		blk.ReceiptHash = block.ReceiptHash().Hex()
		blk.Bloom = block.Bloom().Bytes()
		blk.Difficulty = block.Difficulty().Uint64()
		blk.BlockNumber = block.Number().Uint64()
		blk.GasLimit = block.GasLimit()
		blk.GasUsed = block.GasUsed()
		blk.BlockTime = block.Time()
		blk.Extra = block.Extra()
		blk.MixDigest = block.MixDigest().Hex()
		blk.BlockNonce = block.Nonce()
		if block.BaseFee() != nil {
			blk.BaseFee = block.BaseFee().Uint64()
		}
		blk.BlobGasUsed = block.BlobGasUsed()
		blk.ExcessBlobGas = block.ExcessBlobGas()
		if block.BeaconRoot() != nil {
			blk.ParentBeaconRoot = block.BeaconRoot().Hex()
		}
		blk.BlockHash = block.Hash().Hex()
		blk.BlockSize = block.Size()
		blk.ReceivedAt = block.ReceivedAt

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
			uncle, err := blockUncleService.AddBlockUncle(ctx, tx, blockuncle, blk.Id)
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
			transaction, err := transactionService.AddTransaction(ctx, tx, tns, blk.Id, blk.BlockNumber)
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
			treceipt, err := transactionReceiptService.AddTransactionReceipt(ctx, tx, receipt, blk.Id, blk.BlockNumber, block.Hash().Hex(), transaction.Id)
			if err != nil {
				log.Println(err)
				err = tx.Rollback()
				return err
			}
			tlogs := []*TransactionLog{}
			for _, lg := range GetLogs(receipt) {
				transactionLogService := TransactionLogService{Db: b.Db}
				tlg, err := transactionLogService.AddTransactionLog(ctx, tx, lg, blk.Id, transaction.Id, treceipt.Id)
				if err != nil {
					log.Println(err)
					err = tx.Rollback()
					return err
				}
				topics := []*TransactionLogTopic{}
				for _, tpc := range GetTopics(lg) {
					transactionLogTopicService := TransactionLogTopicService{Db: b.Db}
					topic, err := transactionLogTopicService.AddTransactionLogTopic(ctx, tx, tpc, blk.Id, transaction.Id, treceipt.Id, tlg.Id)
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
			parent_hash,
      uncle_hash,
      coinbase,
      block_root,
      tx_hash,
      receipt_hash,
      bloom,
      difficulty,
      block_number,
      gas_limit,
      gas_used,
      block_time,
      extra,
      mix_digest,
      block_nonce,
      base_fee,
      withdrawals_hash,
      blob_gas_used,
      excess_blob_gas,
      parent_beacon_root,
      block_hash,
      block_size,
      received_at)
  values (?,?,?,?,?,?,?,?,?,?,
          ?,?,?,?,?,?,?,?,?,?,
          ?,?,?);`)
		if err != nil {
			log.Println(err)
			return err
		}
		res, err := stmt.ExecContext(ctx,
			blk.ParentHash,
			blk.UncleHash,
			blk.Coinbase,
			blk.BlockRoot,
			blk.TxHash,
			blk.ReceiptHash,
			blk.Bloom,
			blk.Difficulty,
			blk.BlockNumber,
			blk.GasLimit,
			blk.GasUsed,
			blk.BlockTime,
			blk.Extra,
			blk.MixDigest,
			blk.BlockNonce,
			blk.BaseFee,
			blk.WithdrawalsHash,
			blk.BlobGasUsed,
			blk.ExcessBlobGas,
			blk.ParentBeaconRoot,
			blk.BlockHash,
			blk.BlockSize,
			blk.ReceivedAt)
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
		blk.Id = uint(uId)
		err = stmt.Close()
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	}
}

// GetBlock - used for
func (b *BlockService) GetBlock(ctx context.Context, id uint) (*Block, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		blk := Block{}
		row := b.Db.QueryRowContext(ctx, `select
      id,
			parent_hash,
      uncle_hash,
      coinbase,
      block_root,
      tx_hash,
      receipt_hash,
      bloom,
      difficulty,
      block_number,
      gas_limit,
      gas_used,
      block_time,
      extra,
      mix_digest,
      block_nonce,
      base_fee,
      withdrawals_hash,
      blob_gas_used,
      excess_blob_gas,
      parent_beacon_root,
      block_hash,
      block_size,
      received_at from blocks where id = ?;`, id)

		err := row.Scan(
			&blk.Id,
			&blk.ParentHash,
			&blk.UncleHash,
			&blk.Coinbase,
			&blk.BlockRoot,
			&blk.TxHash,
			&blk.ReceiptHash,
			&blk.Bloom,
			&blk.Difficulty,
			&blk.BlockNumber,
			&blk.GasLimit,
			&blk.GasUsed,
			&blk.BlockTime,
			&blk.Extra,
			&blk.MixDigest,
			&blk.BlockNonce,
			&blk.BaseFee,
			&blk.WithdrawalsHash,
			&blk.BlobGasUsed,
			&blk.ExcessBlobGas,
			&blk.ParentBeaconRoot,
			&blk.BlockHash,
			&blk.BlockSize,
			&blk.ReceivedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		blockUncleService := BlockUncleService{Db: b.Db}
		uncles, err := blockUncleService.GetBlockUncles(ctx, id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		transactionsService := TransactionService{Db: b.Db}
		transactions, err := transactionsService.GetBlockTransactions(ctx, id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		blk.BlockUncles = uncles
		blk.Transactions = transactions

		return &blk, nil
	}
}
