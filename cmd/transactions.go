package cmd

import (
	"context"
	"fmt"
	"math/big"

	"github.com/cloudfresco/ethblocks/ethblocks"

	"github.com/spf13/cobra"
)

// transactionsCmd represents the transactions command
var transactionsCmd = &cobra.Command{
	Use:   "transactions",
	Short: "Returns the transactions in this block",
	Long:  `Returns the transactions in this block`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := ethblocks.GetClient(ClientAddr)
		blockNumber := big.NewInt(BlockNumber)
		ctx := context.Background()
		block, err := ethblocks.GetBlockByNumber(ctx, client, blockNumber)
		if err != nil {
			fmt.Println(err)
			return
		}
		blocktransactions := ethblocks.GetTransactions(block)
		for _, blocktransaction := range blocktransactions {
			ethblocks.PrintTransaction(blocktransaction)
		}
	},
}

func init() {
	rootCmd.AddCommand(transactionsCmd)
	transactionsCmd.Flags().Int64VarP(&BlockNumber, "blocknumber", "b", 1, "Please enter the Block Number")
	err := transactionsCmd.MarkFlagRequired("blocknumber")
	if err != nil {
		fmt.Println(err)
		return
	}
}
