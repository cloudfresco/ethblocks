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

	client, err := svc.GetClient("https://mainnet.infura.io")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	blockNumber := big.NewInt(7602500)
	log.Println("GetBlockByNumber")
	block, err := svc.GetBlockByNumber(ctx, client, blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	svc.PrintBlock(block)

	log.Println("GetBlockByHash")
	h := block.Hash()
	block, err = svc.GetBlockByHash(ctx, client, h)
	if err != nil {
		log.Fatal(err)
	}
	svc.PrintBlock(block)

	blockuncles := svc.GetUncles(block)
	for _, blockuncle := range blockuncles {
		svc.PrintBlockUncle(blockuncle)
	}

	log.Println("GetBlockNumber")
	blocknumber, err := svc.BlockNumber(ctx, client)
	log.Println("Blocknumber:", blocknumber)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("GetBlocks")
	blocks, err := svc.GetBlocks(ctx, client, big.NewInt(7602500), big.NewInt(7602509))
	if err != nil {
		log.Fatal(err)
	}
	for _, b := range blocks {
		svc.PrintBlock(b)
	}

	count, err := svc.GetUncleCountByBlockNumber(ctx, client, blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("GetUncleCountByBlockNumber:", count)

	count, err = svc.GetUncleCountByBlockHash(ctx, client, h)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("GetUncleCountByBlockHash:", count)

	log.Println("GetBlocksByMiner:")
	blocks, err = svc.GetBlocksByMiner(ctx, client, "0x5a0b54d5dc17e0aadc383d2db43b0a0d3e029c4c", big.NewInt(7602500), big.NewInt(7602509))
	if err != nil {
		log.Fatal(err)
	}
	for _, b := range blocks {
		svc.PrintBlock(b)
	}
}
