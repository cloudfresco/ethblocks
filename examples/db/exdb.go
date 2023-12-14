package main

import (
	"context"
	"errors"
	"log"
	"math/big"
	"reflect"

	etbcommon "github.com/cloudfresco/ethblocks/common"
	"github.com/cloudfresco/ethblocks/ethblocks"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	clientAddr := ethblocks.GetEthblocksClientAddr()
	client, err := ethblocks.GetClient(clientAddr)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	blockNumber := big.NewInt(7602500)
	block, err := ethblocks.GetBlockByNumber(ctx, client, blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	// create connection
	appState, err := etbcommon.DbInit()
	if err != nil {
		log.Fatal(err)
	}
	blockService := ethblocks.NewBlockService(appState.Db)
	blk1, err := blockService.AddBlock(ctx, client, block)
	if err != nil {
		log.Fatal(err)
	}
	blk2, err := blockService.GetBlock(ctx, blk1.Id)
	if err != nil {
		log.Fatal(err)
	}
	err = compareBlock(ctx, blk1, blk2)
	if err != nil {
		log.Fatal(err)
	}

	err = compareBlockUncles(ctx, appState, blk1)
	if err != nil {
		log.Fatal(err)
	}

	err = compareBlockTransactions(ctx, appState, blk1)
	if err != nil {
		log.Fatal(err)
	}

	err = compareReceiptsLogTopics(ctx, appState, blk1)
	if err != nil {
		log.Fatal(err)
	}
	err = etbcommon.DbClose(appState.Db)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Comparison of blocks done")
}

// compareBlock - Compare block
func compareBlock(ctx context.Context, blk1 *ethblocks.Block, blk2 *ethblocks.Block) error {
	if !reflect.DeepEqual(blk1, blk2) {
		return errors.New("Block Doesnt Match")
	}
	return nil
}

// compareBlockUncles - Compare Block Uncles
func compareBlockUncles(ctx context.Context, appState *etbcommon.AppState, blk1 *ethblocks.Block) error {
	blockUncleService := ethblocks.NewBlockUncleService(appState.Db)
	uncles, err := blockUncleService.GetBlockUncles(ctx, blk1.Id)
	if err != nil {
		log.Fatal(err)
	}
	if !reflect.DeepEqual(blk1.BlockUncles, uncles) {
		return errors.New("Block Uncles Doesnt Match")
	}
	return nil
}

// compareBlockTransactions - Compare Block Transactions
func compareBlockTransactions(ctx context.Context, appState *etbcommon.AppState, blk1 *ethblocks.Block) error {
	transactionService := ethblocks.NewTransactionService(appState.Db)
	transactions, err := transactionService.GetBlockTransactions(ctx, blk1.Id)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if !reflect.DeepEqual(blk1.Transactions, transactions) {
		return errors.New("Block Transactions Doesnt Match")
	}
	return nil
}

// compareReceiptsLogTopics - Compare Receipts Log Topics
func compareReceiptsLogTopics(ctx context.Context, appState *etbcommon.AppState, blk1 *ethblocks.Block) error {
	for _, trans := range blk1.Transactions {
		transactionReceiptService := ethblocks.NewTransactionReceiptService(appState.Db)
		receipts, err := transactionReceiptService.GetTransactionReceipts(ctx, trans.Id)
		if err != nil {
			log.Fatal(err)
			return err
		}
		if !reflect.DeepEqual(receipts, trans.TransactionReceipts) {
			return errors.New("Block Transaction Receipts Doesnt Match")
		}
		for _, receipt := range trans.TransactionReceipts {
			transactionLogService := ethblocks.NewTransactionLogService(appState.Db)
			logs, err := transactionLogService.GetTransactionLogs(ctx, receipt.Id)
			if err != nil {
				log.Fatal(err)
				return err
			}
			if !reflect.DeepEqual(logs, receipt.Logs) {
				return errors.New("Block Transaction Logs Doesnt Match")
			}

			for _, lg := range receipt.Logs {
				transactionLogTopicService := ethblocks.NewTransactionLogTopicService(appState.Db)
				topics, err := transactionLogTopicService.GetTransactionLogTopics(ctx, lg.Id)
				if err != nil {
					log.Fatal(err)
					return err
				}
				if !reflect.DeepEqual(topics, lg.Topics) {
					return errors.New("Block Transaction Topics Doesnt Match")
				}
			}
		}
	}
	return nil
}
