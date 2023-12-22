package cmd

import (
	"context"
	"fmt"

	"github.com/cloudfresco/ethblocks/ethblocks"
	"github.com/spf13/cobra"
)

// pendingbalanceatCmd represents the pendingbalanceat command
var pendingbalanceatCmd = &cobra.Command{
	Use:   "pendingbalanceat",
	Short: "Returns the wei balance of the given account in the pending state",
	Long:  `Returns the wei balance of the given account in the pending state`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := ethblocks.GetClient(ClientAddr)
		ctx := context.Background()
		balance, err := ethblocks.GetPendingBalanceAt(ctx, client, Address)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("balance", balance)
	},
}

func init() {
	rootCmd.AddCommand(pendingbalanceatCmd)
	pendingbalanceatCmd.Flags().StringVarP(&Address, "address", "a", "0", "Please enter the Address")
	err := pendingbalanceatCmd.MarkFlagRequired("address")
	if err != nil {
		fmt.Println(err)
		return
	}
}
