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
	ExBlock()
	ExAccount()
	ExTransaction()
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
	blk1, err := svc.AddBlock(block)
	if err != nil {
		log.Fatal(err)
	}
	blk2, err := svc.GetBlock(blk1.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(reflect.DeepEqual(blk1, blk2))

	log.Println("GetBlockByHash")
	h := block.Hash()
	block, err = svc.GetBlockByHash(ctx, client, h)
	if err != nil {
		log.Fatal(err)
	}
	svc.PrintBlock(block)

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
	log.Println("GetUncleCountByBlockNumber:", count)

	count, err = svc.GetUncleCountByBlockHash(ctx, client, h)
	log.Println("GetUncleCountByBlockHash:", count)

	log.Println("GetBlocksByMiner:")
	blocks, err = svc.GetBlocksByMiner(ctx, client, "0x5a0b54d5dc17e0aadc383d2db43b0a0d3e029c4c", big.NewInt(7602500), big.NewInt(7602509))
	for _, b := range blocks {
		svc.PrintBlock(b)
	}
}

// ExAccount - Account Examples
func ExAccount() {
	clientaddr := "https://rinkeby.infura.io"
	client, err := svc.GetClient(clientaddr)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	account := "0x7AF3A1f8F9864F8E8B6A465F4eaeFa15321029f4"
	balance, err := svc.GetBalance(ctx, client, account)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("GetBalance :", balance)

	pendingbalance, err := svc.GetPendingBalanceAt(ctx, client, account)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("GetPendingBalanceAt:", pendingbalance)

}

// ExTransaction - Transaction Examples
func ExTransaction() {
	client, err := svc.GetClient("https://mainnet.infura.io")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	count, err := svc.GetBlockTransactionCountByNumber(ctx, client)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("GetBlockTransactionCountByNumber :", count)

	blockNumber := big.NewInt(7602500)
	block, err := svc.GetBlockByNumber(ctx, client, blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	count, err = svc.GetBlockTransactionCountByHash(ctx, client, block.Hash())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("GetBlockTransactionCountByHash:", count)

	txs := block.Transactions()
	blockhash := txs[1].Hash()
	tx, _, err := svc.GetTransactionByHash(ctx, client, blockhash)
	if err != nil {
		log.Fatal(err)
	}
	svc.PrintTransaction(tx)

	tx, err = svc.GetTransactionByBlockHashAndIndex(ctx, client, block.Hash(), uint(0))
	if err != nil {
		log.Fatal(err)
	}
	svc.PrintTransaction(tx)

	txs, err = svc.GetTransactionsByAddress(ctx, client, "0xEec606A66edB6f497662Ea31b5eb1610da87AB5f", big.NewInt(7602500), big.NewInt(7602509))
	if err != nil {
		log.Fatal(err)
	}
	for _, tx := range txs {
		svc.PrintTransaction(tx)
	}

	receipt, err := svc.GetTransactionReceipt(ctx, client, tx.Hash())
	svc.PrintReceipt(receipt)
}
