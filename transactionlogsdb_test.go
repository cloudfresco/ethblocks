package ethblocks

import (
	"context"
	"database/sql"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
)

func TestTransactionLogService_AddTransactionLog(t *testing.T) {
	client, err := GetClient("https://mainnet.infura.io")
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()

	blockNumber := big.NewInt(7602500)
	block, err := GetBlockByNumber(ctx, client, blockNumber)
	if err != nil {
		t.Error(err)
	}
	transactions := GetTransactions(block)
	receipt, err := GetTransactionReceipt(ctx, client, transactions[1].Hash())
	if err != nil {
		t.Error(err)
	}
	tlogs := GetLogs(receipt)
	tlog := tlogs[0]

	err = LoadSQL(appState)
	if err != nil {
		t.Error(err)
		return
	}

	transactionLogService := NewTransactionLogService(appState.Db)

	tx, err := appState.Db.Begin()
	if err != nil {
		t.Error("err", err)
	}
	transLog := TransactionLog{}
	transLog.ID = uint(66)
	transLog.BlockNumber = uint64(7602500)
	transLog.BlockHash = "0xeecc6fa7bf2ae8d533854f44b22fd744a621bfb3520844240f0f67fa26c159c5"
	transLog.Address = "0x8E766F57F7d16Ca50B4A0b90b88f6468A09b0439"
	transLog.LogData = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 231, 169, 129, 250, 100, 76, 142, 128, 0}
	transLog.TxHash = "0x44d3a1bdb0e3de24eb5bc7f9cd601298503e22bb625641dd6f27dd13f64ef4f8"
	transLog.TxIndex = uint(1)
	transLog.LogIndex = uint(0)
	transLog.Removed = false
	transLog.BlockID = uint(1)
	transLog.TransactionID = uint(2)
	transLog.TransactionReceiptID = uint(2)

	type args struct {
		ctx                  context.Context
		tx                   *sql.Tx
		ethLog               *types.Log
		BlockID              uint
		TransactionID        uint
		TransactionReceiptID uint
	}
	tests := []struct {
		tl      *TransactionLogService
		args    args
		want    *TransactionLog
		wantErr bool
	}{
		{
			tl: transactionLogService,
			args: args{
				ctx:                  ctx,
				tx:                   tx,
				ethLog:               tlog,
				BlockID:              1,
				TransactionID:        2,
				TransactionReceiptID: 2,
			},
			want:    &transLog,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.tl.AddTransactionLog(tt.args.ctx, tt.args.tx, tt.args.ethLog, tt.args.BlockID, tt.args.TransactionID, tt.args.TransactionReceiptID)
		if (err != nil) != tt.wantErr {
			t.Errorf("TransactionLogService.AddTransactionLog() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("TransactionLogService.AddTransactionLog() = %v, want %v", got, tt.want)
		}
	}
	err = tx.Commit()
	if err != nil {
		t.Error(err)
	}
}

func TestTransactionLogService_GetTransactionLogs(t *testing.T) {
	ctx := context.Background()

	transactionLogService := NewTransactionLogService(appState.Db)

	transactionLogs := []*TransactionLog{}

	type args struct {
		ctx                  context.Context
		TransactionReceiptID uint
	}
	tests := []struct {
		tl      *TransactionLogService
		args    args
		want    []*TransactionLog
		wantErr bool
	}{
		{
			tl: transactionLogService,
			args: args{
				ctx:                  ctx,
				TransactionReceiptID: 2,
			},
			want:    transactionLogs,
			wantErr: false,
		},
	}
	for range tests {
		/*got, err := tt.tl.GetTransactionLogs(tt.args.ctx, tt.args.TransactionReceiptID)
		if (err != nil) != tt.wantErr {
			t.Errorf("TransactionLogService.GetTransactionLogs() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("TransactionLogService.GetTransactionLogs() = %v, want %v", got, tt.want)
		}*/
	}
}
