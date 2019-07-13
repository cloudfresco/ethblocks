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
	clientaddr := "https://rinkeby.infura.io"
	client, err := svc.GetClient(clientaddr)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	account := "0x7AF3A1f8F9864F8E8B6A465F4eaeFa15321029f4"
	balance, err := svc.GetBalance(ctx, client, account)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("GetBalance :", balance)

	pendingbalance, err := svc.GetPendingBalanceAt(ctx, client, account)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("GetPendingBalanceAt:", pendingbalance)

}
