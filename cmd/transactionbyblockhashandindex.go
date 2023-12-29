package cmd

import (
	"context"
	"fmt"

	"github.com/cloudfresco/ethblocks/ethblocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

// transactionbyblockhashandindexCmd represents the transactionbyblockhashandindex command
var transactionbyblockhashandindexCmd = &cobra.Command{
	Use:   "transactionbyblockhashandindex",
	Short: "Returns information about a transaction by block hash and transaction index position",
	Long:  `Returns information about a transaction by block hash and transaction index position`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := ethblocks.GetClient(ClientAddr)
		blockHash := common.HexToHash(BlockHash)
		ctx := context.Background()
		blocktransaction, err := ethblocks.GetTransactionByBlockHashAndIndex(ctx, client, blockHash, Index)
		if err != nil {
			fmt.Println(err)
			return
		}
		ethblocks.PrintTransaction(blocktransaction)
	},
}

func init() {
	rootCmd.AddCommand(transactionbyblockhashandindexCmd)
	transactionbyblockhashandindexCmd.Flags().StringVarP(&BlockHash, "blockHash", "k", "k", "Please Enter BlockHash")
	err := transactionbyblockhashandindexCmd.MarkFlagRequired("blockHash")
	if err != nil {
		fmt.Println(err)
		return
	}
	transactionbyblockhashandindexCmd.Flags().UintVarP(&Index, "index", "i", 1, "Please enter the Index")
	err = transactionbyblockhashandindexCmd.MarkFlagRequired("index")
	if err != nil {
		fmt.Println(err)
		return
	}
}
