package cmd

import (
	"context"
	"fmt"
	"math/big"

	"github.com/cloudfresco/ethblocks/ethblocks"
	"github.com/spf13/cobra"
)

// unclesCmd represents the uncles command
var unclesCmd = &cobra.Command{
	Use:   "uncles",
	Short: "Returns the uncle blocks from the given block number",
	Long:  `Returns the uncle blocks from the given block number`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := ethblocks.GetClient(ClientAddr)
		blockNumber := big.NewInt(BlockNumber)
		ctx := context.Background()
		block, err := ethblocks.GetBlockByNumber(ctx, client, blockNumber)
		if err != nil {
			fmt.Println(err)
			return
		}
		blockUncles := ethblocks.GetUncles(block)
		for _, blockuncle := range blockUncles {
			ethblocks.PrintBlockUncle(blockuncle)
		}
	},
}

func init() {
	rootCmd.AddCommand(unclesCmd)
	unclesCmd.Flags().Int64VarP(&BlockNumber, "blocknumber", "b", 1, "Please enter the Block Number")
	err := unclesCmd.MarkFlagRequired("blocknumber")
	if err != nil {
		fmt.Println(err)
		return
	}
}
