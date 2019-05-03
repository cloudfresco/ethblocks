package main

import (
	"context"
	"log"
	"math/big"
	"reflect"

	"github.com/cloudfresco/ethblocks/svc"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	ExDb()
}

// ExDb - Block DB Examples
func ExDb() {

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
	blk1, err := svc.AddBlock(ctx, client, block)
	if err != nil {
		log.Fatal(err)
	}
	blk2, err := svc.GetBlock(blk1.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(reflect.DeepEqual(blk1, blk2))

	uncles, err := svc.GetBlockUncles(blk1.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(reflect.DeepEqual(blk1.BlockUncles, uncles))

	transactions, err := svc.GetBlockTransactions(blk1.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(reflect.DeepEqual(blk1.Transactions, transactions))
}
