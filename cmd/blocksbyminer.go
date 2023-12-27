package cmd

import (
	"context"
	"fmt"
	"math/big"

	"github.com/cloudfresco/ethblocks/ethblocks"
	"github.com/spf13/cobra"
)

var Miner string

// blocksbyminerCmd represents the blocksbyminer command
var blocksbyminerCmd = &cobra.Command{
	Use:   "blocksbyminer",
	Short: "Returns the blocks between StartBlockNumber and EndBlockNumber created by this miner",
	Long:  `Returns the blocks between StartBlockNumber and EndBlockNumber created by this miner.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := ethblocks.GetClient(ClientAddr)
		startBlockNumber := big.NewInt(StartBlockNumber)
		endBlockNumber := big.NewInt(EndBlockNumber)
		ctx := context.Background()
		blocks, err := ethblocks.GetBlocksByMiner(ctx, client, Miner, startBlockNumber, endBlockNumber)
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
	rootCmd.AddCommand(blocksbyminerCmd)
	blocksbyminerCmd.Flags().Int64VarP(&StartBlockNumber, "startblocknumber", "s", 1, "Please enter the Start Block Number")
	blocksbyminerCmd.Flags().Int64VarP(&EndBlockNumber, "endblocknumber", "e", 1, "Please enter the End Block Number")
	blocksbyminerCmd.Flags().StringVarP(&Miner, "miner", "m", "m", "Please Enter Miner")
	err := blocksbyminerCmd.MarkFlagRequired("startblocknumber")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = blocksbyminerCmd.MarkFlagRequired("endblocknumber")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = blocksbyminerCmd.MarkFlagRequired("miner")
	if err != nil {
		fmt.Println(err)
		return
	}
}
