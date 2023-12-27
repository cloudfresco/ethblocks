package cmd

import (
	"context"
	"fmt"

	"github.com/cloudfresco/ethblocks/ethblocks"
	"github.com/spf13/cobra"
)

// blocknumberCmd represents the blocknumber command
var blocknumberCmd = &cobra.Command{
	Use:   "blocknumber",
	Short: "Returns the most recent block number",
	Long:  `Returns the most recent block number`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := ethblocks.GetClient(ClientAddr)
		ctx := context.Background()
		blockNumber, err := ethblocks.BlockNumber(ctx, client)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("block Number", blockNumber)
	},
}

func init() {
	rootCmd.AddCommand(blocknumberCmd)
}
