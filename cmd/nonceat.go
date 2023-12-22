package cmd

import (
	"context"
	"fmt"
	"math/big"

	"github.com/cloudfresco/ethblocks/ethblocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

// nonceatCmd represents the nonceat command
var nonceatCmd = &cobra.Command{
	Use:   "nonceat",
	Short: "Returns the account nonce of the given account",
	Long:  `Returns the account nonce of the given account`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := ethblocks.GetClient(ClientAddr)
		address := common.HexToAddress(Address)
		blockNumber := big.NewInt(BlockNumber)
		ctx := context.Background()
		nonce, err := ethblocks.NonceAt(ctx, client, address, blockNumber)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("nonce", nonce)
	},
}

func init() {
	rootCmd.AddCommand(nonceatCmd)
	nonceatCmd.Flags().StringVarP(&Address, "address", "a", "0", "Please enter the Address")
	nonceatCmd.Flags().Int64VarP(&BlockNumber, "blocknumber", "b", 1, "Please enter the Block Number")
	err := nonceatCmd.MarkFlagRequired("address")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = nonceatCmd.MarkFlagRequired("blocknumber")
	if err != nil {
		fmt.Println(err)
		return
	}
}
