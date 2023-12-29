package cmd

import (
	"context"
	"fmt"

	"github.com/cloudfresco/ethblocks/ethblocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

// blocktransactioncountbyhashCmd represents the blocktransactioncountbyhash command
var blocktransactioncountbyhashCmd = &cobra.Command{
	Use:   "blocktransactioncountbyhash",
	Short: "Returns the number of transactions in a block from a block matching the given block hash",
	Long:  `Returns the number of transactions in a block from a block matching the given block hash`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := ethblocks.GetClient(ClientAddr)
		hash := common.HexToHash(Hash)
		ctx := context.Background()
		count, err := ethblocks.GetBlockTransactionCountByHash(ctx, client, hash)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("count", count)
	},
}

func init() {
	rootCmd.AddCommand(blocktransactioncountbyhashCmd)
	blocktransactioncountbyhashCmd.Flags().StringVarP(&Hash, "hash", "s", "s", "Please Enter Hash")
	err := blocktransactioncountbyhashCmd.MarkFlagRequired("hash")
	if err != nil {
		fmt.Println(err)
		return
	}
}
