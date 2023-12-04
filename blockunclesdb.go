package ethblocks

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// BlockUncleIntf - interface
type BlockUncleIntf interface {
	AddBlockUncle(ctx context.Context, tx *sql.Tx, blkuncle *types.Header, BlockID uint) (*BlockUncle, error)
	GetBlockUncles(ctx context.Context, BlockID uint) ([]*BlockUncle, error)
}

// BlockUncleService - For accessing Block Uncle services
type BlockUncleService struct {
	Db *sql.DB
}

// NewBlockUncleService - Create Block Uncle service
func NewBlockUncleService(db *sql.DB) *BlockUncleService {
	return &BlockUncleService{
		Db: db,
	}
}

// BlockUncle - Used for
type BlockUncle struct {
	ID          uint
	BlockNumber uint64
	BlockTime   uint64
	ParentHash  string
	UncleHash   string
	BlockRoot   string
	TxHash      string
	ReceiptHash string
	MixDigest   string
	BlockNonce  uint64
	Coinbase    string
	GasLimit    uint64
	GasUsed     uint64
	Difficulty  uint64
	BlockSize   common.StorageSize
	BlockID     uint
}

// AddBlockUncle - add a block uncle to the db
func (bu *BlockUncleService) AddBlockUncle(ctx context.Context, tx *sql.Tx, blkuncle *types.Header, BlockID uint) (*BlockUncle, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		bl := BlockUncle{}
		bl.BlockNumber = blkuncle.Number.Uint64()
		bl.BlockTime = blkuncle.Time
		bl.ParentHash = blkuncle.ParentHash.Hex()
		bl.UncleHash = blkuncle.UncleHash.Hex()
		bl.BlockRoot = blkuncle.Root.Hex()
		bl.TxHash = blkuncle.TxHash.Hex()
		bl.ReceiptHash = blkuncle.ReceiptHash.Hex()
		bl.MixDigest = blkuncle.MixDigest.Hex()
		bl.BlockNonce = blkuncle.Nonce.Uint64()
		bl.Coinbase = blkuncle.Coinbase.Hex()
		bl.GasLimit = blkuncle.GasLimit
		bl.GasUsed = blkuncle.GasUsed
		bl.Difficulty = blkuncle.Difficulty.Uint64()
		bl.BlockSize = blkuncle.Size()
		bl.BlockID = BlockID
		err := insertBlockUncle(ctx, tx, &bl)
		if err != nil {
			log.Println(err)
			err = tx.Rollback()
			return nil, err
		}
		return &bl, nil
	}
}

// insertBlockUncle - insert block uncle details to db
func insertBlockUncle(ctx context.Context, tx *sql.Tx, blk *BlockUncle) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return err
	default:
		stmt, err := tx.PrepareContext(ctx, `insert into block_uncles
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
			block_size, 
      block_id)
  values (?,?,?,?,?,?,?,?,?,?,
          ?,?,?,?,?);`)
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
			blk.BlockSize,
			blk.BlockID)
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

// GetBlockUncles - used for
func (bu *BlockUncleService) GetBlockUncles(ctx context.Context, BlockID uint) ([]*BlockUncle, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		blockuncles := []*BlockUncle{}
		rows, err := bu.Db.QueryContext(ctx, `select 
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
		block_size, 
		block_id from block_uncles where block_id = ?`, BlockID)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		for rows.Next() {
			blk := BlockUncle{}
			err = rows.Scan(
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
				&blk.BlockSize,
				&blk.BlockID)
			if err != nil {
				log.Println(err)
				return nil, err
			}

			blockuncles = append(blockuncles, &blk)
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

		return blockuncles, nil
	}
}
