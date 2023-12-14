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

func TestTransactionReceiptService_AddTransactionReceipt(t *testing.T) {
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
	receipt, err := GetTransactionReceipt(ctx, client, transactions[0].Hash())
	if err != nil {
		t.Error("err", err)
	}
	// load data into the test db
	err = test.LoadSQL(appState)
	if err != nil {
		t.Error(err)
		return
	}

	transactionReceiptService := NewTransactionReceiptService(appState.Db)
	tx, err := appState.Db.Begin()
	if err != nil {
		t.Error("err", err)
	}
	transReceipt := TransactionReceipt{}
	transReceipt.Id = uint(103)
	transReceipt.ReceiptType = uint8(0)
	transReceipt.PostState = receipt.PostState
	transReceipt.TxStatus = uint64(1)
	transReceipt.CumulativeGasUsed = uint64(21000)
	transReceipt.Bloom = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	transReceipt.TxHash = "0x02e8467c3c439e0e6f129be99c4006609c27bbc6eef8f881c211cc571d77ab27"
	transReceipt.ContractAddress = "0x0000000000000000000000000000000000000000"
	transReceipt.GasUsed = uint64(21000)
	transReceipt.EffectiveGasPrice = uint64(30000000000)
	transReceipt.BlobGasUsed = uint64(0)
	transReceipt.BlobGasPrice = uint64(0)
	transReceipt.BlockHash = "0xeecc6fa7bf2ae8d533854f44b22fd744a621bfb3520844240f0f67fa26c159c5"
	transReceipt.BlockNumber = uint64(7602500)
	transReceipt.TransactionIndex = uint(0)
	transReceipt.BlockId = uint(1)
	transReceipt.TransactionId = uint(1)

	type args struct {
		ctx           context.Context
		tx            *sql.Tx
		ethReceipt    *types.Receipt
		BlockId       uint
		BlockNumber   uint64
		BlockHash     string
		TransactionId uint
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
				BlockId:       1,
				BlockNumber:   block.Number().Uint64(),
				BlockHash:     block.Hash().Hex(),
				TransactionId: 1,
			},
			want:    &transReceipt,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		transactionReceiptResult, err := tt.t.AddTransactionReceipt(tt.args.ctx, tt.args.tx, tt.args.ethReceipt, tt.args.BlockId, tt.args.BlockNumber, tt.args.BlockHash, tt.args.TransactionId)
		if (err != nil) != tt.wantErr {
			t.Errorf("TransactionReceiptService.AddTransactionReceipt() error = %v, wantErr %v", err, tt.wantErr)
			return
		}

		if !reflect.DeepEqual(transactionReceiptResult, tt.want) {
			t.Errorf("TransactionReceiptService.AddTransactionReceipt() = %v, want %v", transactionReceiptResult, tt.want)
		}
		assert.NotNil(t, transactionReceiptResult)
		assert.Equal(t, transactionReceiptResult.TxStatus, uint64(1), "they should be equal")
		assert.Equal(t, transactionReceiptResult.CumulativeGasUsed, uint64(21000), "they should be equal")
		assert.Equal(t, transactionReceiptResult.Bloom, transReceipt.Bloom, "they should be equal")
		assert.Equal(t, transactionReceiptResult.TxHash, "0x02e8467c3c439e0e6f129be99c4006609c27bbc6eef8f881c211cc571d77ab27", "they should be equal")
		assert.Equal(t, transactionReceiptResult.ContractAddress, "0x0000000000000000000000000000000000000000", "they should be equal")
		assert.Equal(t, transactionReceiptResult.GasUsed, uint64(21000), "they should be equal")
		assert.Equal(t, transactionReceiptResult.EffectiveGasPrice, uint64(30000000000), "they should be equal")
		assert.Equal(t, transactionReceiptResult.BlockNumber, uint64(7602500), "they should be equal")

	}
	err = tx.Commit()
	if err != nil {
		t.Error(err)
	}
}
