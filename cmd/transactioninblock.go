package cmd

import (
	"context"
	"fmt"

	"github.com/cloudfresco/ethblocks/ethblocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var (
	Index     uint
	BlockHash string
)

// transactioninblockCmd represents the transactioninblock command
var transactioninblockCmd = &cobra.Command{
	Use:   "transactioninblock",
	Short: "Returns a single transaction at index in the given block",
	Long:  `Returns a single transaction at index in the given block`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := ethblocks.GetClient(ClientAddr)
		blockHash := common.HexToHash(BlockHash)
		ctx := context.Background()
		blocktransaction, err := ethblocks.TransactionInBlock(ctx, client, blockHash, Index)
		if err != nil {
			fmt.Println(err)
			return
		}
		ethblocks.PrintTransaction(blocktransaction)
	},
}

func init() {
	rootCmd.AddCommand(transactioninblockCmd)
	transactioninblockCmd.Flags().StringVarP(&BlockHash, "blockHash", "k", "k", "Please Enter BlockHash")
	err := transactioninblockCmd.MarkFlagRequired("blockHash")
	if err != nil {
		fmt.Println(err)
		return
	}
	transactioninblockCmd.Flags().UintVarP(&Index, "index", "i", 1, "Please enter the Index")
	err = transactioninblockCmd.MarkFlagRequired("index")
	if err != nil {
		fmt.Println(err)
		return
	}
}
