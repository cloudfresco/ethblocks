package main

import (
	"context"
	"log"
	"math/big"

	"github.com/cloudfresco/ethblocks"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	ExBlock()
}

// ExBlock - Block Examples
func ExBlock() {
	clientAddr := ethblocks.GetEthblocksClientAddr()
	client, err := ethblocks.GetClient(clientAddr)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	blockNumber := big.NewInt(7602500)
	log.Println("GetBlockByNumber")
	block, err := ethblocks.GetBlockByNumber(ctx, client, blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	ethblocks.PrintBlock(block)

	log.Println("GetBlockByHash")
	h := block.Hash()
	block, err = ethblocks.GetBlockByHash(ctx, client, h)
	if err != nil {
		log.Fatal(err)
	}
	ethblocks.PrintBlock(block)

	blockuncles := ethblocks.GetUncles(block)
	for _, blockuncle := range blockuncles {
		ethblocks.PrintBlockUncle(blockuncle)
	}

	log.Println("GetBlockNumber")
	blocknumber, err := ethblocks.BlockNumber(ctx, client)
	log.Println("Blocknumber:", blocknumber)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("GetBlocks")
	blocks, err := ethblocks.GetBlocks(ctx, client, big.NewInt(7602500), big.NewInt(7602509))
	if err != nil {
		log.Fatal(err)
	}
	for _, b := range blocks {
		ethblocks.PrintBlock(b)
	}

	count, err := ethblocks.GetUncleCountByBlockNumber(ctx, client, blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("GetUncleCountByBlockNumber:", count)

	count, err = ethblocks.GetUncleCountByBlockHash(ctx, client, h)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("GetUncleCountByBlockHash:", count)

	log.Println("GetBlocksByMiner:")
	blocks, err = ethblocks.GetBlocksByMiner(ctx, client, "0x5a0b54d5dc17e0aadc383d2db43b0a0d3e029c4c", big.NewInt(7602500), big.NewInt(7602509))
	if err != nil {
		log.Fatal(err)
	}
	for _, b := range blocks {
		ethblocks.PrintBlock(b)
	}
}
