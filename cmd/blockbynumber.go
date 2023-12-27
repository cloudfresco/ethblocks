package cmd

import (
	"context"
	"fmt"
	"math/big"

	"github.com/cloudfresco/ethblocks/ethblocks"
	"github.com/spf13/cobra"
)

// blockbynumberCmd represents the blockbynumber command
var blockbynumberCmd = &cobra.Command{
	Use:   "blockbynumber",
	Short: "Returns the full block, given the number of the block",
	Long:  `Returns the full block, given the number of the block`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := ethblocks.GetClient(ClientAddr)
		blockNumber := big.NewInt(BlockNumber)
		ctx := context.Background()
		block, err := ethblocks.GetBlockByNumber(ctx, client, blockNumber)
		if err != nil {
			fmt.Println(err)
			return
		}
		ethblocks.PrintBlock(block)
	},
}

func init() {
	rootCmd.AddCommand(blockbynumberCmd)
	blockbynumberCmd.Flags().Int64VarP(&BlockNumber, "blocknumber", "b", 1, "Please enter the Block Number")
	err := blockbynumberCmd.MarkFlagRequired("blocknumber")
	if err != nil {
		fmt.Println(err)
		return
	}
}
