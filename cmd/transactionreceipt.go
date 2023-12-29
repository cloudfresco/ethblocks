package cmd

import (
	"context"
	"fmt"

	"github.com/cloudfresco/ethblocks/ethblocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

// transactionreceiptCmd represents the transactionreceipt command
var transactionreceiptCmd = &cobra.Command{
	Use:   "transactionreceipt",
	Short: "Returns the receipt of a transaction by transaction hash",
	Long:  `Returns the receipt of a transaction by transaction hash`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := ethblocks.GetClient(ClientAddr)
		txHash := common.HexToHash(TxHash)
		ctx := context.Background()
		receipt, err := ethblocks.GetTransactionReceipt(ctx, client, txHash)
		if err != nil {
			fmt.Println(err)
			return
		}
		ethblocks.PrintReceipt(receipt)
	},
}

func init() {
	rootCmd.AddCommand(transactionreceiptCmd)
	transactionreceiptCmd.Flags().StringVarP(&TxHash, "txHash", "t", "t", "Please Enter Hash")
	err := transactionreceiptCmd.MarkFlagRequired("txHash")
	if err != nil {
		fmt.Println(err)
		return
	}
}
