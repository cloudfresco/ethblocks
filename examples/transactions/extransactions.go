package main

import (
	"context"
	"log"
	"math/big"

	"github.com/cloudfresco/ethblocks"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	ExTransaction()
}

// ExTransaction - Transaction Examples
func ExTransaction() {
	clientAddr := ethblocks.GetEthblocksClientAddr()
	client, err := ethblocks.GetClient(clientAddr)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	count, err := ethblocks.GetBlockTransactionCountByNumber(ctx, client)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("GetBlockTransactionCountByNumber :", count)

	blockNumber := big.NewInt(7602500)
	block, err := ethblocks.GetBlockByNumber(ctx, client, blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	blocktransactions := ethblocks.GetTransactions(block)
	for _, blocktransaction := range blocktransactions {
		ethblocks.PrintTransaction(blocktransaction)
	}

	count, err = ethblocks.GetBlockTransactionCountByHash(ctx, client, block.Hash())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("GetBlockTransactionCountByHash:", count)

	txs := block.Transactions()
	blockhash := txs[1].Hash()
	tx, _, err := ethblocks.GetTransactionByHash(ctx, client, blockhash)
	if err != nil {
		log.Fatal(err)
	}
	ethblocks.PrintTransaction(tx)

	tx, err = ethblocks.GetTransactionByBlockHashAndIndex(ctx, client, block.Hash(), uint(0))
	if err != nil {
		log.Fatal(err)
	}
	ethblocks.PrintTransaction(tx)

	txs, err = ethblocks.GetTransactionsByAddress(ctx, client, "0xEec606A66edB6f497662Ea31b5eb1610da87AB5f", big.NewInt(7602500), big.NewInt(7602509))
	if err != nil {
		log.Fatal(err)
	}
	for _, tx := range txs {
		ethblocks.PrintTransaction(tx)
	}

	for _, blocktransaction := range blocktransactions {
		receipt, err := ethblocks.GetTransactionReceipt(ctx, client, blocktransaction.Hash())
		if err != nil {
			log.Fatal(err)
		}
		ethblocks.PrintReceipt(receipt)
		for _, lg := range receipt.Logs {
			ethblocks.PrintReceiptLog(lg)
			for _, topic := range lg.Topics {
				ethblocks.PrintTopic(topic)
			}
		}
	}
}
