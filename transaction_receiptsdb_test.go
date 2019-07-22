package ethblocks

import (
	"context"
	"database/sql"
	"log"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
)

func TestTransactionReceiptService_AddTransactionReceipt(t *testing.T) {
	client, err := GetClient("https://mainnet.infura.io")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	blockNumber := big.NewInt(7602500)
	block1, err := GetBlockByNumber(ctx, client, blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	transactions := GetTransactions(block1)
	receipt, err := GetTransactionReceipt(ctx, client, transactions[0].Hash())
	if err != nil {
		log.Println("err", err)
	}
	// load data into the test db
	err = fixtures.Load()
	if err != nil {
		log.Println("err", err)
	}
	transactionReceiptService := NewTransactionReceiptService(appState.Db)
	tx, err := appState.Db.Begin()
	if err != nil {
		log.Println("err", err)
	}
	transReceipt := TransactionReceipt{}
	transReceipt.ID = uint(103)
	transReceipt.BlockNumber = uint64(7602500)
	transReceipt.BlockHash = "0xeecc6fa7bf2ae8d533854f44b22fd744a621bfb3520844240f0f67fa26c159c5"
	transReceipt.TxHash = "0x02e8467c3c439e0e6f129be99c4006609c27bbc6eef8f881c211cc571d77ab27"
	transReceipt.TxStatus = uint64(1)
	transReceipt.CumulativeGasUsed = uint64(21000)
	transReceipt.GasUsed = uint64(21000)
	transReceipt.ContractAddress = "0x0000000000000000000000000000000000000000"
	transReceipt.PostState = receipt.PostState
	transReceipt.BlockID = uint(1)
	transReceipt.TransactionID = uint(1)

	type args struct {
		ctx           context.Context
		tx            *sql.Tx
		ethReceipt    *types.Receipt
		BlockID       uint
		BlockNumber   uint64
		BlockHash     string
		TransactionID uint
	}
	tests := []struct {
		t       *TransactionReceiptService
		args    args
		want    *TransactionReceipt
		wantErr bool
	}{
		{
			t: transactionReceiptService,
			args: args{
				ctx:           ctx,
				tx:            tx,
				ethReceipt:    receipt,
				BlockID:       1,
				BlockNumber:   block1.Number().Uint64(),
				BlockHash:     block1.Hash().Hex(),
				TransactionID: 1,
			},
			want:    &transReceipt,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.t.AddTransactionReceipt(tt.args.ctx, tt.args.tx, tt.args.ethReceipt, tt.args.BlockID, tt.args.BlockNumber, tt.args.BlockHash, tt.args.TransactionID)
		if (err != nil) != tt.wantErr {
			t.Errorf("TransactionReceiptService.AddTransactionReceipt() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("TransactionReceiptService.AddTransactionReceipt() = %v, want %v", got, tt.want)
		}
	}
}

func TestTransactionReceiptService_GetTransactionReceipts(t *testing.T) {
	ctx := context.Background()
	// load data into the test db
	err := fixtures.Load()
	if err != nil {
		log.Fatal(err)
	}

	transactionReceiptService := NewTransactionReceiptService(appState.Db)

	transactionReceipts := []*TransactionReceipt{}

	type args struct {
		ctx           context.Context
		TransactionID uint
	}
	tests := []struct {
		t       *TransactionReceiptService
		args    args
		want    []*TransactionReceipt
		wantErr bool
	}{
		{
			t: transactionReceiptService,
			args: args{
				ctx:           ctx,
				TransactionID: 1,
			},
			want:    transactionReceipts,
			wantErr: false,
		},
	}
	for range tests {
		/*got, err := tt.t.GetTransactionReceipts(tt.args.ctx, tt.args.TransactionID)
		if (err != nil) != tt.wantErr {
			t.Errorf("TransactionReceiptService.GetTransactionReceipts() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("TransactionReceiptService.GetTransactionReceipts() = %v, want %v", got, tt.want)
		}*/
	}
}
