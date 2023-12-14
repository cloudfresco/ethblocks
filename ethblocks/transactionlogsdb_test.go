package ethblocks

import (
	"context"
	"database/sql"
	"math/big"
	"reflect"
	"testing"

	"github.com/cloudfresco/ethblocks/test"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func TestTransactionLogService_AddTransactionLog(t *testing.T) {
	client, err := GetClient(clientAddr)
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

	err = test.LoadSQL(appState)
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
	transLog.Id = uint(66)
	transLog.Address = "0x8E766F57F7d16Ca50B4A0b90b88f6468A09b0439"
	transLog.LogData = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 231, 169, 129, 250, 100, 76, 142, 128, 0}
	transLog.BlockNumber = uint64(7602500)
	transLog.TxHash = "0x44d3a1bdb0e3de24eb5bc7f9cd601298503e22bb625641dd6f27dd13f64ef4f8"
	transLog.TxIndex = uint(1)
	transLog.BlockHash = "0xeecc6fa7bf2ae8d533854f44b22fd744a621bfb3520844240f0f67fa26c159c5"
	transLog.LogIndex = uint(0)
	transLog.Removed = false
	transLog.BlockId = uint(1)
	transLog.TransactionId = uint(2)
	transLog.TransactionReceiptId = uint(2)

	type args struct {
		ctx                  context.Context
		tx                   *sql.Tx
		ethLog               *types.Log
		BlockId              uint
		TransactionId        uint
		TransactionReceiptId uint
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
				BlockId:              1,
				TransactionId:        2,
				TransactionReceiptId: 2,
			},
			want:    &transLog,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		transLogResult, err := tt.tl.AddTransactionLog(tt.args.ctx, tt.args.tx, tt.args.ethLog, tt.args.BlockId, tt.args.TransactionId, tt.args.TransactionReceiptId)
		if (err != nil) != tt.wantErr {
			t.Errorf("TransactionLogService.AddTransactionLog() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(transLogResult, tt.want) {
			t.Errorf("TransactionLogService.AddTransactionLog() = %v, want %v", transLogResult, tt.want)
		}
		assert.NotNil(t, transLogResult)
		assert.Equal(t, transLogResult.Address, "0x8E766F57F7d16Ca50B4A0b90b88f6468A09b0439", "they should be equal")
		assert.Equal(t, transLogResult.LogData, transLog.LogData, "they should be equal")
		assert.Equal(t, transLogResult.BlockNumber, uint64(7602500), "they should be equal")
		assert.Equal(t, transLogResult.TxHash, "0x44d3a1bdb0e3de24eb5bc7f9cd601298503e22bb625641dd6f27dd13f64ef4f8", "they should be equal")
		assert.Equal(t, transLogResult.TxIndex, uint(1), "they should be equal")
		assert.Equal(t, transLogResult.BlockHash, "0xeecc6fa7bf2ae8d533854f44b22fd744a621bfb3520844240f0f67fa26c159c5", "they should be equal")
		assert.False(t, transLogResult.Removed, "Its false")

	}
	err = tx.Commit()
	if err != nil {
		t.Error(err)
	}
}
