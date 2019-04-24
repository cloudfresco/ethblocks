package svc

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// PrintTransaction - used to print Transaction
func PrintTransaction(tx *types.Transaction) {
	log.Println("tx:", tx)
	log.Println("hash          : ", tx.Hash())
	log.Println("nonce           : ", tx.Nonce())
	log.Println("Size     : ", tx.Size())
	log.Println("value           : ", tx.Value())
	log.Println("gasPrice        : ", tx.GasPrice())
	log.Println("gas             : ", tx.Gas())
	log.Println("from            : ", GetSender(tx))
	log.Println("to              : ", tx.To())
}

// GetSender - used to get sender details
func GetSender(tx *types.Transaction) string {
	msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()))
	if err != nil {
		log.Fatal(err)
	}

	return msg.From().Hex()
}

// GetBlockTransactionCountByNumber - used to Get Block Transaction Count By Number
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblocktransactioncountbynumber
func GetBlockTransactionCountByNumber(ctx context.Context, client *ethclient.Client) (uint, error) {
	count, err := client.PendingTransactionCount(ctx)
	if err != nil {
		return uint(0), err
	}
	return count, nil
}

// GetBlockTransactionCountByHash - used to Get Block Transaction Count By Hash
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblocktransactioncountbyhash
func GetBlockTransactionCountByHash(ctx context.Context, client *ethclient.Client, blockHash common.Hash) (uint, error) {
	count, err := client.TransactionCount(ctx, blockHash)
	if err != nil {
		return uint(0), err
	}
	return count, nil
}

// GetTransactionByHash - used to Get Transaction By Hash
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactionbyhash
func GetTransactionByHash(ctx context.Context, client *ethclient.Client, hash common.Hash) (*types.Transaction, bool, error) {
	tx, isPending, err := client.TransactionByHash(ctx, hash)
	if err != nil {
		return nil, false, err
	}
	return tx, isPending, nil
}

// GetTransactionByBlockHashAndIndex - used to Get Transaction By BlockHash And Index
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactionbyblockhashandindex
func GetTransactionByBlockHashAndIndex(ctx context.Context, client *ethclient.Client, hash common.Hash, index uint) (*types.Transaction, error) {
	tx, err := client.TransactionInBlock(ctx, hash, index)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// GetTransactionByBlockNumberAndIndex - used to Get Transaction By BlockNumber And Index
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactionbyblocknumberandindex
func GetTransactionByBlockNumberAndIndex() {

}

// GetTransactionReceipt - used to Get Transaction Receipt
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactionreceipt
func GetTransactionReceipt(ctx context.Context, client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	receipt, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

// CreateRawTransaction - used to Create Raw Transaction
func CreateRawTransaction(nonce uint64, toAddress common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *types.Transaction {
	tx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, data)
	return tx
}

// SendRawTransaction - used to Send Raw Transaction
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sendrawtransaction
func SendRawTransaction(ctx context.Context, client *ethclient.Client, tx *types.Transaction) error {
	err := client.SendTransaction(ctx, tx)
	if err != nil {
		return err
	}
	return nil
}

// GetTransactionsByAddress - used to Get Transactions By Address
func GetTransactionsByAddress(ctx context.Context, client *ethclient.Client, frmaddr string, startBlockNumber *big.Int, endBlockNumber *big.Int) ([]*types.Transaction, error) {
	frmaddress := common.HexToHash(frmaddr)
	txs := []*types.Transaction{}
	var one = big.NewInt(1)
	for i := new(big.Int).Set(startBlockNumber); i.Cmp(endBlockNumber) <= 0; i.Add(i, one) {
		block, err := GetBlockByNumber(ctx, client, i)
		if err != nil {
			return nil, err
		}
		for _, tx := range block.Transactions() {
			frmaddress1 := common.HexToHash(GetSender(tx))
			if frmaddress == frmaddress1 {
				txs = append(txs, tx)
			}
		}

	}
	return txs, nil
}
