package cmd

import (
	"context"
	"fmt"

	"github.com/cloudfresco/ethblocks/ethblocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var TxHash string

// transactionbyhashCmd represents the transactionbyhash command
var transactionbyhashCmd = &cobra.Command{
	Use:   "transactionbyhash",
	Short: "Returns the transaction with the given hash",
	Long:  `Returns the transaction with the given hash`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := ethblocks.GetClient(ClientAddr)
		txHash := common.HexToHash(TxHash)
		ctx := context.Background()
		blocktransaction, _, err := ethblocks.GetTransactionByHash(ctx, client, txHash)
		if err != nil {
			fmt.Println(err)
			return
		}
		ethblocks.PrintTransaction(blocktransaction)
	},
}

func init() {
	rootCmd.AddCommand(transactionbyhashCmd)
	transactionbyhashCmd.Flags().StringVarP(&TxHash, "txHash", "t", "t", "Please Enter Hash")
	err := transactionbyhashCmd.MarkFlagRequired("txHash")
	if err != nil {
		fmt.Println(err)
		return
	}
}
