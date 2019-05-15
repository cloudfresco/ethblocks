package main

import (
	"context"
	"log"
	"math/big"

	"github.com/cloudfresco/ethblocks/svc"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	ExTransaction()
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

	blocktransactions := svc.GetTransactions(block)
	for _, blocktransaction := range blocktransactions {
		svc.PrintTransaction(blocktransaction)
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

	for _, blocktransaction := range blocktransactions {
		receipt, err := svc.GetTransactionReceipt(ctx, client, blocktransaction.Hash())
		if err != nil {
			log.Fatal(err)
		}
		svc.PrintReceipt(receipt)
		for _, lg := range receipt.Logs {
			svc.PrintReceiptLog(lg)
			for _, topic := range lg.Topics {
				svc.PrintTopic(topic)
			}
		}
	}
}
