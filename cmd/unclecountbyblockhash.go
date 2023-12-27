package cmd

import (
	"context"
	"fmt"

	"github.com/cloudfresco/ethblocks/ethblocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

// unclecountbyblockhashCmd represents the unclecountbyblockhash command
var unclecountbyblockhashCmd = &cobra.Command{
	Use:   "unclecountbyblockhash",
	Short: "Returns the number of uncles in a block from a block matching the given block hash",
	Long:  `Returns the number of uncles in a block from a block matching the given block hash`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := ethblocks.GetClient(ClientAddr)
		hash := common.HexToHash(Hash)
		ctx := context.Background()
		count, err := ethblocks.GetUncleCountByBlockHash(ctx, client, hash)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("count", count)
	},
}

func init() {
	rootCmd.AddCommand(unclecountbyblockhashCmd)
	unclecountbyblockhashCmd.Flags().StringVarP(&Hash, "hash", "s", "s", "Please Enter Hash")
	err := unclecountbyblockhashCmd.MarkFlagRequired("hash")
	if err != nil {
		fmt.Println(err)
		return
	}
}
