package cmd

import (
	"context"
	"fmt"
	"math/big"

	"github.com/cloudfresco/ethblocks/ethblocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

// codeatCmd represents the codeat command
var codeatCmd = &cobra.Command{
	Use:   "codeat",
	Short: "Returns the contract code of the given account",
	Long:  `Returns the contract code of the given account`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := ethblocks.GetClient(ClientAddr)
		address := common.HexToAddress(Address)
		blockNumber := big.NewInt(BlockNumber)
		ctx := context.Background()
		code, err := ethblocks.CodeAt(ctx, client, address, blockNumber)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("code", code)
	},
}

func init() {
	rootCmd.AddCommand(codeatCmd)
	codeatCmd.Flags().StringVarP(&Address, "address", "a", "0", "Please enter the Address")
	codeatCmd.Flags().Int64VarP(&BlockNumber, "blocknumber", "b", 1, "Please enter the Block Number")
	err := codeatCmd.MarkFlagRequired("address")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = codeatCmd.MarkFlagRequired("blocknumber")
	if err != nil {
		fmt.Println(err)
		return
	}
}
