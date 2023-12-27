package cmd

import (
	"context"
	"fmt"
	"math/big"

	"github.com/cloudfresco/ethblocks/ethblocks"
	"github.com/spf13/cobra"
)

// unclecountbyblocknumberCmd represents the unclecountbyblocknumber command
var unclecountbyblocknumberCmd = &cobra.Command{
	Use:   "unclecountbyblocknumber",
	Short: "Returns the number of uncles in a block from a block matching the given block number",
	Long:  `Returns the number of uncles in a block from a block matching the given block number`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := ethblocks.GetClient(ClientAddr)
		blockNumber := big.NewInt(BlockNumber)
		ctx := context.Background()
		count, err := ethblocks.GetUncleCountByBlockNumber(ctx, client, blockNumber)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("count", count)
	},
}

func init() {
	rootCmd.AddCommand(unclecountbyblocknumberCmd)
	unclecountbyblocknumberCmd.Flags().Int64VarP(&BlockNumber, "blocknumber", "b", 1, "Please enter the Block Number")
	err := unclecountbyblocknumberCmd.MarkFlagRequired("blocknumber")
	if err != nil {
		fmt.Println(err)
		return
	}
}
