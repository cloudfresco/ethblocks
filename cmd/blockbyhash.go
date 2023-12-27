package cmd

import (
	"context"
	"fmt"

	"github.com/cloudfresco/ethblocks/ethblocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var Hash string

// blockbyhashCmd represents the blockbyhash command
var blockbyhashCmd = &cobra.Command{
	Use:   "blockbyhash",
	Short: "Returns the full block, given the hash of the block",
	Long:  `Returns the full block, given the hash of the block`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := ethblocks.GetClient(ClientAddr)
		hash := common.HexToHash(Hash)
		ctx := context.Background()
		block, err := ethblocks.GetBlockByHash(ctx, client, hash)
		if err != nil {
			fmt.Println(err)
			return
		}
		ethblocks.PrintBlock(block)
	},
}

func init() {
	rootCmd.AddCommand(blockbyhashCmd)
	blockbyhashCmd.Flags().StringVarP(&Hash, "hash", "s", "s", "Please Enter Hash")
	err := blockbyhashCmd.MarkFlagRequired("hash")
	if err != nil {
		fmt.Println(err)
		return
	}
}
