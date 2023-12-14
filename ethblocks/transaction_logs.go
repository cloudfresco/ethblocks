package ethblocks

import (
	"log"

	"github.com/ethereum/go-ethereum/core/types"
)

// PrintReceiptLog - Print Receipt Log
func PrintReceiptLog(lg *types.Log) {
	log.Println("BlockNumber  : ", lg.BlockNumber)
	log.Println("BlockHash    : ", lg.BlockHash.Hex())
	log.Println("Address  : ", lg.Address.Hex())
	log.Println("TxHash    : ", lg.TxHash.Hex())
	log.Println("TxIndex  : ", lg.TxIndex)
	log.Println("LogIndex    : ", lg.Index)
	log.Println("Removed         : ", lg.Removed)
}

// GetLogs - Get Logs by receipt
func GetLogs(receipt *types.Receipt) []*types.Log {
	logs := receipt.Logs
	return logs
}
