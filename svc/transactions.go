package svc

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// PrintTransaction - Print Transaction
func PrintTransaction(tx *types.Transaction) {
	log.Println("hash            : ", tx.Hash().Hex())
	log.Println("AccountNonce    : ", tx.Nonce())
	log.Println("Price           : ", tx.GasPrice().Uint64())
	log.Println("GasLimit        : ", tx.Gas())
	log.Println("TxAmount        : ", tx.Value().Uint64())
	log.Println("from            : ", getSender(tx))
	log.Println("to              : ", tx.To().Hex())
	log.Println("Size            : ", tx.Size())
}

// getSender - Get sender details
func getSender(tx *types.Transaction) string {
	msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()))
	if err != nil {
		return ""
	}

	return msg.From().Hex()
}

// GetTransactions - Get Transaction by block
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblockbyhash
func GetTransactions(block *types.Block) []*types.Transaction {
	transactions := block.Transactions()
	return transactions
}

// GetBlockTransactionCountByNumber - Get Block Transaction Count By Number
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblocktransactioncountbynumber
func GetBlockTransactionCountByNumber(ctx context.Context, client *ethclient.Client) (uint, error) {
	count, err := client.PendingTransactionCount(ctx)
	if err != nil {
		return uint(0), err
	}
	return count, nil
}

// GetBlockTransactionCountByHash - Get Block Transaction Count By Hash
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblocktransactioncountbyhash
func GetBlockTransactionCountByHash(ctx context.Context, client *ethclient.Client, blockHash common.Hash) (uint, error) {
	count, err := client.TransactionCount(ctx, blockHash)
	if err != nil {
		return uint(0), err
	}
	return count, nil
}

// GetTransactionByHash - Get Transaction By Hash
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactionbyhash
func GetTransactionByHash(ctx context.Context, client *ethclient.Client, hash common.Hash) (*types.Transaction, bool, error) {
	tx, isPending, err := client.TransactionByHash(ctx, hash)
	if err != nil {
		return nil, false, err
	}
	return tx, isPending, nil
}

// GetTransactionByBlockHashAndIndex - Get Transaction By BlockHash And Index
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactionbyblockhashandindex
func GetTransactionByBlockHashAndIndex(ctx context.Context, client *ethclient.Client, hash common.Hash, index uint) (*types.Transaction, error) {
	tx, err := client.TransactionInBlock(ctx, hash, index)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// GetTransactionByBlockNumberAndIndex - Get Transaction By BlockNumber And Index
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactionbyblocknumberandindex
func GetTransactionByBlockNumberAndIndex() {

}

// GetTransactionsByAddress - Get Transactions By Address
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
			frmaddress1 := common.HexToHash(getSender(tx))
			if frmaddress == frmaddress1 {
				txs = append(txs, tx)
			}
		}

	}
	return txs, nil
}

// CreateRawTransaction - Create Raw Transaction
func CreateRawTransaction(nonce uint64, toAddress common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *types.Transaction {
	tx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, data)
	return tx
}

// SendRawTransaction - Send Raw Transaction
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sendrawtransaction
func SendRawTransaction(ctx context.Context, client *ethclient.Client, tx *types.Transaction, prv *ecdsa.PrivateKey) error {
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), prv)
	if err != nil {
		return err
	}
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return err
	}
	return nil
}
