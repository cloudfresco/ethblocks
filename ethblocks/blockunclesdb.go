package ethblocks

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
)

// BlockUncleIntf - interface
type BlockUncleIntf interface {
	AddBlockUncle(ctx context.Context, tx *sql.Tx, blkUncle *types.Header, blockId uint) (*BlockUncle, error)
	GetBlockUncles(ctx context.Context, blockId uint) ([]*BlockUncle, error)
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
	BlockId          uint
}

// AddBlockUncle - add a block uncle to the db
func (bu *BlockUncleService) AddBlockUncle(ctx context.Context, tx *sql.Tx, blkUncle *types.Header, blockId uint) (*BlockUncle, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		blockUncle := BlockUncle{}
		blockUncle.ParentHash = blkUncle.ParentHash.Hex()
		blockUncle.UncleHash = blkUncle.UncleHash.Hex()
		blockUncle.Coinbase = blkUncle.Coinbase.Hex()
		blockUncle.BlockRoot = blkUncle.Root.Hex()
		blockUncle.TxHash = blkUncle.TxHash.Hex()
		blockUncle.ReceiptHash = blkUncle.ReceiptHash.Hex()
		blockUncle.Bloom = blkUncle.Bloom.Bytes()
		blockUncle.Difficulty = blkUncle.Difficulty.Uint64()
		blockUncle.BlockNumber = blkUncle.Number.Uint64()
		blockUncle.GasLimit = blkUncle.GasLimit
		blockUncle.GasUsed = blkUncle.GasUsed
		blockUncle.BlockTime = blkUncle.Time
		blockUncle.Extra = blkUncle.Extra
		blockUncle.MixDigest = blkUncle.MixDigest.Hex()
		blockUncle.BlockNonce = blkUncle.Nonce.Uint64()
		if blkUncle.BaseFee != nil {
			blockUncle.BaseFee = blkUncle.BaseFee.Uint64()
		}
		if blkUncle.BlobGasUsed != nil {
			blockUncle.BlobGasUsed = blkUncle.BlobGasUsed
		}
		if blkUncle.ExcessBlobGas != nil {
			blockUncle.ExcessBlobGas = blkUncle.ExcessBlobGas
		}
		if blkUncle.ParentBeaconRoot != nil {
			blockUncle.ParentBeaconRoot = blkUncle.ParentBeaconRoot.Hex()
		}
		blockUncle.BlockId = blockId
		err := insertBlockUncle(ctx, tx, &blockUncle)
		if err != nil {
			log.Println(err)
			err = tx.Rollback()
			return nil, err
		}
		return &blockUncle, nil
	}
}

// insertBlockUncle - insert block uncle details to db
func insertBlockUncle(ctx context.Context, tx *sql.Tx, blockUncle *BlockUncle) error {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return err
	default:
		stmt, err := tx.PrepareContext(ctx, `insert into block_uncles
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
      block_id)
  values (?,?,?,?,?,?,?,?,?,?,
          ?,?,?,?,?,?,?,?,?,?,
          ?);`)
		if err != nil {
			log.Println(err)
			return err
		}
		res, err := stmt.ExecContext(ctx,
			blockUncle.ParentHash,
			blockUncle.UncleHash,
			blockUncle.Coinbase,
			blockUncle.BlockRoot,
			blockUncle.TxHash,
			blockUncle.ReceiptHash,
			blockUncle.Bloom,
			blockUncle.Difficulty,
			blockUncle.BlockNumber,
			blockUncle.GasLimit,
			blockUncle.GasUsed,
			blockUncle.BlockTime,
			blockUncle.Extra,
			blockUncle.MixDigest,
			blockUncle.BlockNonce,
			blockUncle.BaseFee,
			blockUncle.WithdrawalsHash,
			blockUncle.BlobGasUsed,
			blockUncle.ExcessBlobGas,
			blockUncle.ParentBeaconRoot,
			blockUncle.BlockId)
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
		blockUncle.Id = uint(uId)
		err = stmt.Close()
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	}
}

// GetBlockUncles - used for
func (bu *BlockUncleService) GetBlockUncles(ctx context.Context, blockId uint) ([]*BlockUncle, error) {
	select {
	case <-ctx.Done():
		err := errors.New("Client closed connection")
		return nil, err
	default:
		blockuncles := []*BlockUncle{}
		rows, err := bu.Db.QueryContext(ctx, `select 
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
		block_id from block_uncles where block_id = ?`, blockId)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		for rows.Next() {
			blockUncle := BlockUncle{}
			err = rows.Scan(
				&blockUncle.Id,
				&blockUncle.ParentHash,
				&blockUncle.UncleHash,
				&blockUncle.Coinbase,
				&blockUncle.BlockRoot,
				&blockUncle.TxHash,
				&blockUncle.ReceiptHash,
				&blockUncle.Bloom,
				&blockUncle.Difficulty,
				&blockUncle.BlockNumber,
				&blockUncle.GasLimit,
				&blockUncle.GasUsed,
				&blockUncle.BlockTime,
				&blockUncle.Extra,
				&blockUncle.MixDigest,
				&blockUncle.BlockNonce,
				&blockUncle.BaseFee,
				&blockUncle.WithdrawalsHash,
				&blockUncle.BlobGasUsed,
				&blockUncle.ExcessBlobGas,
				&blockUncle.ParentBeaconRoot,
				&blockUncle.BlockId)
			if err != nil {
				log.Println(err)
				return nil, err
			}

			blockuncles = append(blockuncles, &blockUncle)
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
