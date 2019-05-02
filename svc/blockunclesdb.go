package svc

import (
	"database/sql"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

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
func AddBlockUncle(tx *sql.Tx, blkuncle *types.Header, BlockID uint) (*BlockUncle, error) {
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
	blockuncle, err := InsertBlockUncle(tx, bl)
	if err != nil {
		log.Println(err)
		err = tx.Rollback()
		return nil, err
	}
	return blockuncle, nil
}

// InsertBlockUncle - insert block uncle details to db
func InsertBlockUncle(tx *sql.Tx, blk BlockUncle) (*BlockUncle, error) {
	stmt, err := tx.Prepare(`insert into block_uncles
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
		blk.BlockSize,
		blk.BlockID)
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

// GetBlockUncles - used for
func GetBlockUncles(BlockID uint) (*[]BlockUncle, error) {
	appState, err := dbInit()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	db := appState.Db
	blk := BlockUncle{}
	blockuncles := []BlockUncle{}
	rows, err := db.Query(`select 
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

		blockuncles = append(blockuncles, blk)
	}
	err = rows.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &blockuncles, nil
}