package main

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/cloudfresco/ethblocks/ethblocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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

	gclient, err := ethblocks.GetGethClient(clientAddr)
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
	testKey, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	err = ethblocks.SubscribePendingTransactions(ctx, client, gclient, 16, &common.Address{1}, big.NewInt(1), 22000, big.NewInt(1), nil, testKey)
	if err != nil {
		log.Println(err)
	}

	err = ethblocks.SubscribeFullPendingTransactions(ctx, client, gclient, 1, &common.Address{1}, big.NewInt(1), 22000, big.NewInt(1), nil, testKey)
	if err != nil {
		log.Println(err)
	}

	clientAddr2, pKey, toAddr := ethblocks.GetEthblocksClient2Details()
	client2, err := ethblocks.GetClient(clientAddr2)
	if err != nil {
		log.Fatal(err)
	}
	privateKey, err := crypto.HexToECDSA(pKey)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	if err != nil {
		log.Println(err)
	}

	toAddress := common.HexToAddress(toAddr)

	err = ethblocks.SendTransaction1(client2, fromAddress, &toAddress, privateKey, big.NewInt(100000000000000000), uint64(21000), big.NewInt(2000000000))
	if err != nil {
		log.Fatal("SendTransaction1", err)
	}
}
