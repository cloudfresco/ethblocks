package ethblocks

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// PrintReceipt - Print Receipt
func PrintReceipt(receipt *types.Receipt) {
	log.Println("TxStatus: ", receipt.Status)
	log.Println("CumulativeGasUsed: ", receipt.CumulativeGasUsed)
	log.Println("Bloom: ", receipt.Bloom)
	log.Println("Logs: ", receipt.Logs)
	log.Println("TxHash: ", receipt.TxHash.Hex())
	log.Println("ContractAddress :", receipt.ContractAddress.Hex())
	log.Println("GasUsed: ", receipt.GasUsed)
}

// GetTransactionReceipt - Returns the receipt of a transaction by
// transaction hash
// https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_gettransactionreceiptt
func GetTransactionReceipt(ctx context.Context, client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	receipt, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}
