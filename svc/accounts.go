package svc

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GetClient - Used to get client
func GetClient(addr string) (*ethclient.Client, error) {
	client, err := ethclient.Dial(addr)
	if err != nil {
		return nil, err
	}
	return client, err
}

// GetBalance - Used to get Balance
func GetBalance(ctx context.Context, client *ethclient.Client, addr string) (*big.Int, error) {
	account := common.HexToAddress(addr)
	balance, err := client.BalanceAt(ctx, account, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

// GetBalance2 - Used to get Balance
func GetBalance2(ctx context.Context, clientAddr string, addr string) (*big.Int, error) {
	client, err := GetClient(clientAddr)
	account := common.HexToAddress(addr)
	balance, err := client.BalanceAt(ctx, account, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

// GetPendingBalanceAt - Used to get Pending Balance
func GetPendingBalanceAt(ctx context.Context, client *ethclient.Client, addr string) (*big.Int, error) {
	account := common.HexToAddress(addr)
	balance, err := client.PendingBalanceAt(ctx, account)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

// GetPendingBalanceAt2 - Used to get Pending Balance
func GetPendingBalanceAt2(ctx context.Context, clientAddr string, addr string) (*big.Int, error) {
	client, err := GetClient(clientAddr)
	account := common.HexToAddress(addr)
	balance, err := client.PendingBalanceAt(ctx, account)
	if err != nil {
		return nil, err
	}
	return balance, nil
}
