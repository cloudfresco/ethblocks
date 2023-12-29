package ethblocks

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
)

// PrintTransaction - Print Transaction
func PrintTransaction(tx *types.Transaction) {
	txv, txr, txs := tx.RawSignatureValues()
	log.Println("From            : ", getSender(tx))
	log.Println("To              : ", tx.To().Hex())
	log.Println("AccountNonce    : ", tx.Nonce())
	log.Println("Hash            : ", tx.Hash().Hex())
	log.Println("Size            : ", tx.Size())
	log.Println("TxAmount        : ", tx.Value().Uint64())
	log.Println("TxType          : ", tx.Type())
	log.Println("ChainId         : ", tx.ChainId().Uint64())
	log.Println("GasLimit        : ", tx.Gas())
	log.Println("GasPrice        : ", tx.GasPrice().Uint64())
	log.Println("GasTipCap       : ", tx.GasTipCap().Uint64())
	log.Println("GasFeeCap       : ", tx.GasFeeCap().Uint64())
	log.Println("TxV             : ", txv.Uint64())
	log.Println("TxR             : ", txr.Uint64())
	log.Println("TxS             : ", txs.Uint64())
	log.Println("Data            : ", tx.Data())
}

// getSender - Get sender details
// https://github.com/ethereum/go-ethereum/issues/22918
func getSender(tx *types.Transaction) string {
	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		return ""
	}

	return from.Hex()
}

// GetTransactions - return the Transactions in this block
func GetTransactions(block *types.Block) []*types.Transaction {
	transactions := block.Transactions()
	return transactions
}

// TransactionInBlock - return a single transaction at index in the given block
func TransactionInBlock(ctx context.Context, client *ethclient.Client, blockHash common.Hash, index uint) (*types.Transaction, error) {
	tx, err := client.TransactionInBlock(ctx, blockHash, index)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// GetBlockTransactionCountByNumber - returns the total number of transactions in the pending state.
// https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_getblocktransactioncountbynumber
func GetBlockTransactionCountByNumber(ctx context.Context, client *ethclient.Client) (uint, error) {
	count, err := client.PendingTransactionCount(ctx)
	if err != nil {
		return uint(0), err
	}
	return count, nil
}

// GetBlockTransactionCountByHash - Returns the number of transactions in
// a block from a block matching the given block hash.
// https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_getblocktransactioncountbyhash
func GetBlockTransactionCountByHash(ctx context.Context, client *ethclient.Client, blockHash common.Hash) (uint, error) {
	count, err := client.TransactionCount(ctx, blockHash)
	if err != nil {
		return uint(0), err
	}
	return count, nil
}

// GetTransactionByHash - Returns the information about a transaction
// requested by transaction hash.
// https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_gettransactionbyhash
func GetTransactionByHash(ctx context.Context, client *ethclient.Client, hash common.Hash) (*types.Transaction, bool, error) {
	tx, isPending, err := client.TransactionByHash(ctx, hash)
	if err != nil {
		return nil, false, err
	}
	return tx, isPending, nil
}

// GetTransactionByBlockHashAndIndex - Returns information about a
// transaction by block hash and transaction index position.
// https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_gettransactionbyblockhashandindex
func GetTransactionByBlockHashAndIndex(ctx context.Context, client *ethclient.Client, hash common.Hash, index uint) (*types.Transaction, error) {
	tx, err := client.TransactionInBlock(ctx, hash, index)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// GetTransactionsByAddress - returns the transactions in a range of blocks
func GetTransactionsByAddress(ctx context.Context, client *ethclient.Client, frmaddr string, startBlockNumber *big.Int, endBlockNumber *big.Int) ([]*types.Transaction, error) {
	frmaddress := common.HexToHash(frmaddr)
	txs := []*types.Transaction{}
	one := big.NewInt(1)
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

// SendRawTransaction - Send a signed Transaction
// https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_sendrawtransaction
func SendRawTransaction(ctx context.Context, client *ethclient.Client, tx *types.Transaction, prv *ecdsa.PrivateKey) error {
	chainId, err := client.NetworkID(ctx)
	if err != nil {
		return err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), prv)
	if err != nil {
		return err
	}
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return err
	}
	return nil
}

// SendTransaction1 -- create and send a transaction
func SendTransaction1(client *ethclient.Client, fromAddr common.Address, toAddress *common.Address, privateKey *ecdsa.PrivateKey, value *big.Int, gas uint64, gasPrice *big.Int) error {
	nonce, err := client.PendingNonceAt(context.Background(), fromAddr)
	if err != nil {
		return err
	}

	return SendTransaction2(client, nonce, toAddress, privateKey, value, gas, gasPrice)
}

// SendTransaction2 -- create and send a transaction
func SendTransaction2(client *ethclient.Client, nonce uint64, toAddress *common.Address, privateKey *ecdsa.PrivateKey, value *big.Int, gas uint64, gasPrice *big.Int) error {
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return err
	}

	signer := types.LatestSignerForChainID(chainID)
	tx, err := types.SignNewTx(privateKey, signer, &types.LegacyTx{
		Nonce:    nonce,
		To:       toAddress,
		Value:    value,
		Gas:      gas,
		GasPrice: gasPrice,
	})
	if err != nil {
		return err
	}
	return client.SendTransaction(context.Background(), tx)
}

// SubscribePendingTransactions - Subscribe to Transactions
func SubscribePendingTransactions(ctx context.Context, client *ethclient.Client, gclient *gethclient.Client, nonce uint64, to *common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte, testKey *ecdsa.PrivateKey) error {
	// Subscribe to Transactions
	ch := make(chan common.Hash)
	_, err := gclient.SubscribePendingTransactions(ctx, ch)
	if err != nil {
		return err
	}
	// send transaction
	err = SendTransaction2(client, nonce, to, testKey, amount, gasLimit, gasPrice)
	if err != nil {
		return err
	}

	return nil
}

// SubscribeFullPendingTransactions - subscribe full pending transactions
func SubscribeFullPendingTransactions(ctx context.Context, client *ethclient.Client, gclient *gethclient.Client, nonce uint64, to *common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte, testKey *ecdsa.PrivateKey) error {
	// Subscribe to Transactions
	ch := make(chan *types.Transaction)
	_, err := gclient.SubscribeFullPendingTransactions(ctx, ch)
	if err != nil {
		return err
	}
	err = SendTransaction2(client, nonce, to, testKey, amount, gasLimit, gasPrice)
	if err != nil {
		return err
	}
	return nil
}
