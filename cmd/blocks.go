package cmd

import (
	"context"
	"fmt"
	"math/big"

	"github.com/cloudfresco/ethblocks/ethblocks"
	"github.com/spf13/cobra"
)

var (
	StartBlockNumber int64
	EndBlockNumber   int64
)

// blocksCmd represents the blocks command
var blocksCmd = &cobra.Command{
	Use:   "blocks",
	Short: "Returns the blocks between StartBlockNumber and EndBlockNumber",
	Long:  `Returns the blocks between StartBlockNumber and EndBlockNumber`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := ethblocks.GetClient(ClientAddr)
		startBlockNumber := big.NewInt(StartBlockNumber)
		endBlockNumber := big.NewInt(EndBlockNumber)
		ctx := context.Background()
		blocks, err := ethblocks.GetBlocks(ctx, client, startBlockNumber, endBlockNumber)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, block := range blocks {
			ethblocks.PrintBlock(block)
		}
	},
}

func init() {
	rootCmd.AddCommand(blocksCmd)
	blocksCmd.Flags().Int64VarP(&StartBlockNumber, "startblocknumber", "s", 1, "Please enter the Start Block Number")
	blocksCmd.Flags().Int64VarP(&EndBlockNumber, "endblocknumber", "e", 1, "Please enter the End Block Number")
	err := blocksCmd.MarkFlagRequired("startblocknumber")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = blocksCmd.MarkFlagRequired("endblocknumber")
	if err != nil {
		fmt.Println(err)
		return
	}
}
