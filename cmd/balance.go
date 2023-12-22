package cmd

import (
	"context"
	"fmt"

	"github.com/cloudfresco/ethblocks/ethblocks"
	"github.com/spf13/cobra"
)

// balanceCmd represents the balance command
var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Returns the wei balance of the given account",
	Long:  `Returns the wei balance of the given account`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := ethblocks.GetClient(ClientAddr)
		ctx := context.Background()
		balance, err := ethblocks.GetBalance(ctx, client, Address)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("balance", balance)
	},
}

func init() {
	rootCmd.AddCommand(balanceCmd)
	balanceCmd.Flags().StringVarP(&Address, "address", "a", "0", "Please enter the Address")
	err := balanceCmd.MarkFlagRequired("address")
	if err != nil {
		fmt.Println(err)
		return
	}
}
