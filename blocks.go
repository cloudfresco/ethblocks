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
	log.Println("Number:", block.Number().Uint64())
	log.Println(" Hash            : ", block.Hash().Hex())
	log.Println(" ParentHash      : ", block.ParentHash().Hex())
	log.Println(" Nonce           : ", block.Nonce())
	log.Println(" Sha3Uncles      : ", block.UncleHash().Hex())
	log.Println(" LogsBloom       : ", block.Bloom())
	log.Println(" StateRoot       : ", block.Root().Hex())
	log.Println(" TransactionsRoot: ", block.TxHash().Hex())
	log.Println(" Miner           : ", block.Coinbase().Hex())
	log.Println(" Difficulty      : ", block.Difficulty().Uint64())
	log.Println(" ExtraData       : ", block.Extra())
	log.Println(" Size           : ", block.Size())
	log.Println(" GasLimit        : ", block.GasLimit())
	log.Println(" GasUsed         : ", block.GasUsed())
	log.Println(" Timestamp       : ", block.Time())
	log.Println(" Transactions    : ", block.Transactions())
	log.Println(" Length of transactions    : ", len(block.Transactions()))
	log.Println(" Uncles          : ", block.Uncles())
}

// PrintBlockUncle - Print block Uncle
func PrintBlockUncle(uncle *types.Header) {
	log.Println("Number:", uncle.Number.Uint64())
	log.Println(" ParentHash      : ", uncle.ParentHash.Hex())
	log.Println(" Nonce           : ", uncle.Nonce.Uint64())
	log.Println(" Sha3Uncles      : ", uncle.UncleHash.Hex())
	log.Println(" LogsBloom       : ", uncle.Bloom)
	log.Println(" StateRoot       : ", uncle.Root.Hex())
	log.Println(" TransactionsRoot: ", uncle.TxHash.Hex())
	log.Println(" Miner           : ", uncle.Coinbase.Hex())
	log.Println(" Difficulty      : ", uncle.Difficulty.Uint64())
	log.Println(" Size           : ", uncle.Size())
	log.Println(" GasLimit        : ", uncle.GasLimit)
	log.Println(" GasUsed         : ", uncle.GasUsed)
	log.Println(" Timestamp       : ", uncle.Time)
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

// GetBlockByHash - Get block by block hash
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblockbyhash
func GetBlockByHash(ctx context.Context, client *ethclient.Client, hash common.Hash) (*types.Block, error) {
	block, err := client.BlockByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return block, nil
}

// BlockNumber - Get the latest block number
func BlockNumber(ctx context.Context, client *ethclient.Client) (string, error) {
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return "", err
	}
	return header.Number.String(), nil
}

// GetBlocks - Get blocks between start and end
func GetBlocks(ctx context.Context, client *ethclient.Client, startBlockNumber *big.Int, endBlockNumber *big.Int) ([]*types.Block, error) {
	blocks := []*types.Block{}
	var one = big.NewInt(1)
	for i := new(big.Int).Set(startBlockNumber); i.Cmp(endBlockNumber) <= 0; i.Add(i, one) {
		block, err := GetBlockByNumber(ctx, client, i)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}

// GetUncles - Get Uncles by block
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblockbyhash
func GetUncles(block *types.Block) []*types.Header {
	blockuncles := block.Uncles()
	return blockuncles
}

// GetUncleCountByBlockNumber - Get Uncle Count By Block Number
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getunclecountbyblocknumber
func GetUncleCountByBlockNumber(ctx context.Context, client *ethclient.Client, blockNumber *big.Int) (int, error) {
	block, err := GetBlockByNumber(ctx, client, blockNumber)
	if err != nil {
		return 0, err
	}
	return len(block.Uncles()), err
}

// GetUncleCountByBlockHash - Get Uncle Count By Block Hash
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getunclecountbyblockhash
func GetUncleCountByBlockHash(ctx context.Context, client *ethclient.Client, hash common.Hash) (int, error) {
	block, err := GetBlockByHash(ctx, client, hash)
	if err != nil {
		return 0, err
	}
	return len(block.Uncles()), err
}

// GetUncleByBlockHashAndIndex - Get Uncle By Block Hash and Index
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getunclebyblockhashandindex
func GetUncleByBlockHashAndIndex() {

}

// GetUncleByBlockNumberAndIndex - Get Uncle  By Block Number and Index
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getunclebyblocknumberandindex
func GetUncleByBlockNumberAndIndex() {

}

// GetBlocksByMiner - Get blocks by miner
func GetBlocksByMiner(ctx context.Context, client *ethclient.Client, miner string, startBlockNumber *big.Int, endBlockNumber *big.Int) ([]*types.Block, error) {
	mineraddr := common.HexToAddress(miner)
	blocks := []*types.Block{}
	var one = big.NewInt(1)
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
