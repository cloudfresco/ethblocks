package cmd

import (
	"context"
	"fmt"

	"github.com/cloudfresco/ethblocks/ethblocks"
	"github.com/spf13/cobra"
)

// blocktransactioncountbynumberCmd represents the blocktransactioncountbynumber command
var blocktransactioncountbynumberCmd = &cobra.Command{
	Use:   "blocktransactioncountbynumber",
	Short: "Returns the total number of transactions in the pending state",
	Long:  `Returns the total number of transactions in the pending state`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := ethblocks.GetClient(ClientAddr)
		ctx := context.Background()
		count, err := ethblocks.GetBlockTransactionCountByNumber(ctx, client)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("count", count)
	},
}

func init() {
	rootCmd.AddCommand(blocktransactioncountbynumberCmd)
}
