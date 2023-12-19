package main

import (
	"context"
	"log"
	"math/big"

	"github.com/cloudfresco/ethblocks/ethblocks"
	"github.com/ethereum/go-ethereum/common"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	ExAccount()
}

// ExAccount - Account Examples
func ExAccount() {
	clientAddr := ethblocks.GetEthblocksClientAddr()
	address := common.HexToAddress(clientAddr)
	client, err := ethblocks.GetClient(clientAddr)
	if err != nil {
		log.Fatal(err)
	}

	gclient, err := ethblocks.GetGethClient(clientAddr)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	account := "0x7AF3A1f8F9864F8E8B6A465F4eaeFa15321029f4"
	blockNumber := big.NewInt(7602500)
	// testKey, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testSlot := common.HexToHash("0xdeadbeef")

	balance, err := ethblocks.GetBalance(ctx, client, account)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("GetBalance :", balance)

	balanceAt, err := ethblocks.BalanceAt(ctx, client, address, blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("BalanceAt :", balanceAt)

	pendingbalance, err := ethblocks.GetPendingBalanceAt(ctx, client, account)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("GetPendingBalanceAt :", pendingbalance)

	slotValue, err := ethblocks.StorageAt(ctx, client, address, testSlot, blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("StorageAt :", slotValue)

	code, err := ethblocks.CodeAt(ctx, client, address, blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("CodeAt :", code)

	accountResult, err := ethblocks.GetProof(ctx, gclient, address, []string{testSlot.String()}, blockNumber)
	if err != nil {
		log.Println(err)
	}
	log.Println("GetProof :", accountResult)

	nonce, err := ethblocks.NonceAt(ctx, client, address, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("NonceAt", nonce)
}
