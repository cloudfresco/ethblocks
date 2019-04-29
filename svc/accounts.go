package svc

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GetClient - Get client
func GetClient(addr string) (*ethclient.Client, error) {
	client, err := ethclient.Dial(addr)
	if err != nil {
		return nil, err
	}
	return client, err
}

// GetBalance - Get Balance
func GetBalance(ctx context.Context, client *ethclient.Client, addr string) (*big.Int, error) {
	account := common.HexToAddress(addr)
	balance, err := client.BalanceAt(ctx, account, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

// GetPendingBalanceAt - Get Pending Balance
func GetPendingBalanceAt(ctx context.Context, client *ethclient.Client, addr string) (*big.Int, error) {
	account := common.HexToAddress(addr)
	balance, err := client.PendingBalanceAt(ctx, account)
	if err != nil {
		return nil, err
	}
	return balance, nil
}
