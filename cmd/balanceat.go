package cmd

import (
	"context"
	"fmt"
	"math/big"

	"github.com/cloudfresco/ethblocks/ethblocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var (
	Address     string
	BlockNumber int64
)

// balanceatCmd represents the balanceat command
var balanceatCmd = &cobra.Command{
	Use:   "balanceat",
	Short: "Returns the wei balance of the given account",
	Long:  `Returns the wei balance of the given account`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := ethblocks.GetClient(ClientAddr)
		address := common.HexToAddress(Address)
		blockNumber := big.NewInt(BlockNumber)
		ctx := context.Background()
		balance, err := ethblocks.BalanceAt(ctx, client, address, blockNumber)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("balance", balance)
	},
}

func init() {
	rootCmd.AddCommand(balanceatCmd)
	balanceatCmd.Flags().StringVarP(&Address, "address", "a", "0", "Please enter the Address")
	balanceatCmd.Flags().Int64VarP(&BlockNumber, "blocknumber", "b", 1, "Please enter the Block Number")
	err := balanceatCmd.MarkFlagRequired("address")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = balanceatCmd.MarkFlagRequired("blocknumber")
	if err != nil {
		fmt.Println(err)
		return
	}
}
