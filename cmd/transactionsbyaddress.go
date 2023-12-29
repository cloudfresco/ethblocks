package cmd

import (
	"context"
	"fmt"
	"math/big"

	"github.com/cloudfresco/ethblocks/ethblocks"
	"github.com/spf13/cobra"
)

// transactionsbyaddressCmd represents the transactionsbyaddress command
var transactionsbyaddressCmd = &cobra.Command{
	Use:   "transactionsbyaddress",
	Short: "Returns the transactions in a range of blocks",
	Long:  `Returns the transactions in a range of blocks`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := ethblocks.GetClient(ClientAddr)
		startBlockNumber := big.NewInt(StartBlockNumber)
		endBlockNumber := big.NewInt(EndBlockNumber)
		ctx := context.Background()
		blocktransactions, err := ethblocks.GetTransactionsByAddress(ctx, client, Address, startBlockNumber, endBlockNumber)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, blocktransaction := range blocktransactions {
			ethblocks.PrintTransaction(blocktransaction)
		}
	},
}

func init() {
	rootCmd.AddCommand(transactionsbyaddressCmd)
	transactionsbyaddressCmd.Flags().Int64VarP(&StartBlockNumber, "startblocknumber", "s", 1, "Please enter the Start Block Number")
	transactionsbyaddressCmd.Flags().Int64VarP(&EndBlockNumber, "endblocknumber", "e", 1, "Please enter the End Block Number")
	transactionsbyaddressCmd.Flags().StringVarP(&Address, "address", "a", "a", "Please Enter Address")
	err := transactionsbyaddressCmd.MarkFlagRequired("startblocknumber")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = transactionsbyaddressCmd.MarkFlagRequired("endblocknumber")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = transactionsbyaddressCmd.MarkFlagRequired("address")
	if err != nil {
		fmt.Println(err)
		return
	}
}
