package main

import (
	"context"
	"log"

	"github.com/cloudfresco/ethblocks"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	ExAccount()
}

// ExAccount - Account Examples
func ExAccount() {

  clientAddr := ethblocks.GetEthblocksClientAddr()
  client, err := ethblocks.GetClient(clientAddr)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	account := "0x7AF3A1f8F9864F8E8B6A465F4eaeFa15321029f4"
	balance, err := ethblocks.GetBalance(ctx, client, account)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("GetBalance :", balance)

	pendingbalance, err := ethblocks.GetPendingBalanceAt(ctx, client, account)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("GetPendingBalanceAt:", pendingbalance)

}
