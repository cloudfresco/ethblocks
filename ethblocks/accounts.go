package ethblocks

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
)

// GetClient - connects a client to the given URL
func GetClient(addr string) (*ethclient.Client, error) {
	client, err := ethclient.Dial(addr)
	if err != nil {
		return nil, err
	}
	return client, err
}

// GetBalance - returns the wei balance of the given account.
func GetBalance(ctx context.Context, client *ethclient.Client, addr string) (*big.Int, error) {
	account := common.HexToAddress(addr)
	balance, err := client.BalanceAt(ctx, account, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

// GetGethClient - Get gethclient
func GetGethClient(addr string) (*gethclient.Client, error) {
	client, err := ethclient.Dial(addr)
	if err != nil {
		return nil, err
	}
	gclient := gethclient.New(client.Client())
	return gclient, err
}

// BalanceAt - returns the wei balance of the given account.
func BalanceAt(ctx context.Context, client *ethclient.Client, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	balance, err := client.BalanceAt(ctx, account, blockNumber)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

// GetPendingBalanceAt - returns the wei balance of the given account in the pending state
func GetPendingBalanceAt(ctx context.Context, client *ethclient.Client, addr string) (*big.Int, error) {
	account := common.HexToAddress(addr)
	balance, err := client.PendingBalanceAt(ctx, account)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

// StorageAt - returns the value of key in the contract storage of the given account.
func StorageAt(ctx context.Context, client *ethclient.Client, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	slotValue, err := client.StorageAt(ctx, account, key, blockNumber)
	if err != nil {
		return nil, err
	}
	return slotValue, nil
}

// CodeAt - returns the contract code of the given account.
func CodeAt(ctx context.Context, client *ethclient.Client, account common.Address, blockNumber *big.Int) ([]byte, error) {
	code, err := client.CodeAt(ctx, account, blockNumber)
	if err != nil {
		return nil, err
	}
	return code, nil
}

// GetProof - returns the account and storage values of the specified account including the Merkle-proof
func GetProof(ctx context.Context, gclient *gethclient.Client, account common.Address, keys []string, blockNumber *big.Int) (*gethclient.AccountResult, error) {
	result, err := gclient.GetProof(ctx, account, keys, blockNumber)
	if err != nil {
		return nil, err
	}
	if result.Address != account {
		err = fmt.Errorf("unexpected address, have: %v want: %v", result.Address, account)
		return nil, err
	}
	return result, nil
}

// NonceAt - returns the account nonce of the given account
func NonceAt(ctx context.Context, client *ethclient.Client, account common.Address, blockNumber *big.Int) (uint64, error) {
	nonce, err := client.NonceAt(ctx, account, blockNumber)
	if err != nil {
		return 0, err
	}
	return nonce, nil
}
