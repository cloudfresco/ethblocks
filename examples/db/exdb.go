package main

import (
	"context"
	"errors"
	"log"
	"math/big"
	"reflect"

	"github.com/cloudfresco/ethblocks/svc"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
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
	blk2, err := svc.GetBlock(ctx, blk1.ID)
	if err != nil {
		log.Fatal(err)
	}
	err = compareBlock(ctx, blk1, blk2)
	if err != nil {
		log.Fatal(err)
	}

	err = compareBlockUncles(ctx, blk1)
	if err != nil {
		log.Fatal(err)
	}

	err = compareBlockTransactions(ctx, blk1)
	if err != nil {
		log.Fatal(err)
	}

	err = compareReceiptsLogTopics(ctx, blk1)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Comparison of blocks done")
}

// compareBlock - Compare block
func compareBlock(ctx context.Context, blk1 *svc.Block, blk2 *svc.Block) error {
	if reflect.DeepEqual(blk1, blk2) == false {
		return errors.New("Block Doesnt Match")
	}
	return nil
}

// compareBlockUncles - Compare Block Uncles
func compareBlockUncles(ctx context.Context, blk1 *svc.Block) error {
	uncles, err := svc.GetBlockUncles(ctx, blk1.ID)
	if err != nil {
		log.Fatal(err)
	}
	if reflect.DeepEqual(blk1.BlockUncles, uncles) == false {
		return errors.New("Block Uncles Doesnt Match")
	}
	return nil
}

// compareBlockTransactions - Compare Block Transactions
func compareBlockTransactions(ctx context.Context, blk1 *svc.Block) error {
	transactions, err := svc.GetBlockTransactions(ctx, blk1.ID)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if reflect.DeepEqual(blk1.Transactions, transactions) == false {
		return errors.New("Block Transactions Doesnt Match")
	}
	return nil
}

// compareReceiptsLogTopics - Compare Receipts Log Topics
func compareReceiptsLogTopics(ctx context.Context, blk1 *svc.Block) error {
	for _, trans := range blk1.Transactions {
		receipts, err := svc.GetTransactionReceipts(ctx, trans.ID)
		if err != nil {
			log.Fatal(err)
			return err
		}
		if reflect.DeepEqual(receipts, trans.TransactionReceipts) == false {
			return errors.New("Block Transaction Receipts Doesnt Match")
		}
		for _, receipt := range trans.TransactionReceipts {
			logs, err := svc.GetTransactionLogs(ctx, receipt.ID)
			if err != nil {
				log.Fatal(err)
				return err
			}
			if reflect.DeepEqual(logs, receipt.Logs) == false {
				return errors.New("Block Transaction Logs Doesnt Match")
			}

			for _, lg := range receipt.Logs {
				topics, err := svc.GetTransactionLogTopics(ctx, lg.ID)
				if err != nil {
					log.Fatal(err)
					return err
				}
				if reflect.DeepEqual(topics, lg.Topics) == false {
					return errors.New("Block Transaction Topics Doesnt Match")
				}
			}
		}
	}
	return nil
}
