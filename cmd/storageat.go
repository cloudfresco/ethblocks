package cmd

import (
	"context"
	"fmt"
	"math/big"

	"github.com/cloudfresco/ethblocks/ethblocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var TestSlot string

// storageatCmd represents the storageat command
var storageatCmd = &cobra.Command{
	Use:   "storageat",
	Short: "Returns the value of key in the contract storage of the given account",
	Long:  `Returns the value of key in the contract storage of the given account`,
	Run: func(cmd *cobra.Command, args []string) {
		client, _ := ethblocks.GetClient(ClientAddr)
		address := common.HexToAddress(Address)
		blockNumber := big.NewInt(BlockNumber)
		testSlot := common.HexToHash(TestSlot)
		ctx := context.Background()
		slotValue, err := ethblocks.StorageAt(ctx, client, address, testSlot, blockNumber)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("slotValue", slotValue)
	},
}

func init() {
	rootCmd.AddCommand(storageatCmd)
	storageatCmd.Flags().StringVarP(&Address, "address", "a", "0", "Please enter the Address")
	storageatCmd.Flags().StringVarP(&TestSlot, "testslot", "c", "c", "Please enter TestSlot")
	storageatCmd.Flags().Int64VarP(&BlockNumber, "blocknumber", "b", 1, "Please enter the Block Number")
	err := storageatCmd.MarkFlagRequired("address")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = storageatCmd.MarkFlagRequired("testslot")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = storageatCmd.MarkFlagRequired("blocknumber")
	if err != nil {
		fmt.Println(err)
		return
	}
}
