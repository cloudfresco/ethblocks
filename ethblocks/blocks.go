package ethblocks

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// PrintBlock - Print block
func PrintBlock(block *types.Block) {
	PrintBlockHeader(block.Header())
	PrintBlockUncles(block.Uncles())
	PrintBlockTransactions(block.Transactions())
}

// PrintBlockHeader - Print a single block Header
func PrintBlockHeader(header *types.Header) {
	log.Println(" Number          : ", header.Number.Uint64())
	log.Println(" ParentHash      : ", header.ParentHash.Hex())
	log.Println(" Sha3Uncles      : ", header.UncleHash.Hex())
	log.Println(" Hash            : ", header.Hash().Hex())
	log.Println(" Miner           : ", header.Coinbase.Hex())
	log.Println(" StateRoot       : ", header.Root.Hex())
	log.Println(" TransactionsRoot: ", header.TxHash.Hex())
	log.Println(" ReceiptsRoot    : ", header.ReceiptHash.Hex())
	log.Println(" LogsBloom       : ", header.Bloom)
	log.Println(" Difficulty      : ", header.Difficulty.Uint64())
	log.Println(" GasLimit        : ", header.GasLimit)
	log.Println(" GasUsed         : ", header.GasUsed)
	log.Println(" Timestamp       : ", header.Time)
	log.Println(" ExtraData       : ", header.Extra)
	log.Println(" MixHash         : ", header.MixDigest)
	log.Println(" Nonce           : ", header.Nonce)
	log.Println(" Size            : ", header.Size())
}

// PrintBlockUncles - Loop over the uncles and print each one
func PrintBlockUncles(uncles []*types.Header) {
	for _, uncle := range uncles {
		PrintBlockUncle(uncle)
	}
}

// PrintBlockUncle - Print a single block Uncle
func PrintBlockUncle(uncle *types.Header) {
	log.Println(" Number          : ", uncle.Number.Uint64())
	log.Println(" ParentHash      : ", uncle.ParentHash.Hex())
	log.Println(" Sha3Uncles      : ", uncle.UncleHash.Hex())
	log.Println(" Hash            : ", uncle.Hash().Hex())
	log.Println(" Miner           : ", uncle.Coinbase.Hex())
	log.Println(" StateRoot       : ", uncle.Root.Hex())
	log.Println(" TransactionsRoot: ", uncle.TxHash.Hex())
	log.Println(" ReceiptsRoot    : ", uncle.ReceiptHash.Hex())
	log.Println(" LogsBloom       : ", uncle.Bloom)
	log.Println(" Difficulty      : ", uncle.Difficulty.Uint64())
	log.Println(" GasLimit        : ", uncle.GasLimit)
	log.Println(" GasUsed         : ", uncle.GasUsed)
	log.Println(" Timestamp       : ", uncle.Time)
	log.Println(" ExtraData       : ", uncle.Extra)
	log.Println(" MixHash         : ", uncle.MixDigest)
	log.Println(" Nonce           : ", uncle.Nonce)
	log.Println(" Size            : ", uncle.Size())
}

// PrintBlockTransactions Loop over the transactions and print each one
func PrintBlockTransactions(blocktransactions []*types.Transaction) {
	for _, blocktransaction := range blocktransactions {
		PrintTransaction(blocktransaction)
	}
}

// GetBlockByNumber - Get block by block number
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblockbynumber
func GetBlockByNumber(ctx context.Context, client *ethclient.Client, blockNumber *big.Int) (*types.Block, error) {
	block, err := client.BlockByNumber(ctx, blockNumber)
	if err != nil {
		return nil, err
	}
	return block, nil
}

// BlockNumber - returns the number of most recent block
// https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_blocknumber
func BlockNumber(ctx context.Context, client *ethclient.Client) (string, error) {
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return "", err
	}
	return header.Number.String(), nil
}

// GetBlockByHash - Returns information about a block by hash
// https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_getblockbyhash
func GetBlockByHash(ctx context.Context, client *ethclient.Client, hash common.Hash) (*types.Block, error) {
	block, err := client.BlockByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return block, nil
}

// GetBlocks - return blocks between start and end
func GetBlocks(ctx context.Context, client *ethclient.Client, startBlockNumber *big.Int, endBlockNumber *big.Int) ([]*types.Block, error) {
	blocks := []*types.Block{}
	one := big.NewInt(1)
	for i := new(big.Int).Set(startBlockNumber); i.Cmp(endBlockNumber) <= 0; i.Add(i, one) {
		block, err := GetBlockByNumber(ctx, client, i)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}

// GetUncles - return the Uncles in the block
func GetUncles(block *types.Block) []*types.Header {
	blockuncles := block.Uncles()
	return blockuncles
}

// GetUncleCountByBlockNumber - Returns the number of uncles in a block
// from a block matching the given block number.
// https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_getunclecountbyblocknumber
func GetUncleCountByBlockNumber(ctx context.Context, client *ethclient.Client, blockNumber *big.Int) (int, error) {
	block, err := GetBlockByNumber(ctx, client, blockNumber)
	if err != nil {
		return 0, err
	}
	return len(block.Uncles()), err
}

// GetUncleCountByBlockHash - Returns the number of uncles in a block
// from a block matching the given block number.
// https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_getunclecountbyblocknumber
func GetUncleCountByBlockHash(ctx context.Context, client *ethclient.Client, hash common.Hash) (int, error) {
	block, err := GetBlockByHash(ctx, client, hash)
	if err != nil {
		return 0, err
	}
	return len(block.Uncles()), err
}

// GetBlocksByMiner - Get blocks by miner
func GetBlocksByMiner(ctx context.Context, client *ethclient.Client, miner string, startBlockNumber *big.Int, endBlockNumber *big.Int) ([]*types.Block, error) {
	mineraddr := common.HexToAddress(miner)
	blocks := []*types.Block{}
	one := big.NewInt(1)
	for i := new(big.Int).Set(startBlockNumber); i.Cmp(endBlockNumber) <= 0; i.Add(i, one) {
		block, err := GetBlockByNumber(ctx, client, i)
		if err != nil {
			return nil, err
		}
		if block.Coinbase() == mineraddr {
			blocks = append(blocks, block)
		}
	}
	return blocks, nil
}
