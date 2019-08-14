package ethblocks

import (
	"context"
	"database/sql"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
)

func TestTransactionService_AddTransaction(t *testing.T) {
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

	err = LoadSQL(appState)
	if err != nil {
		t.Error(err)
		return
	}

	transactionService := NewTransactionService(appState.Db)
	tx, err := appState.Db.Begin()
	if err != nil {
		t.Error(err)
	}
	t1 := Transaction{}
	t1.ID = uint(103)
	t1.BlockNumber = uint64(7602500)
	t1.BlockHash = "0x02e8467c3c439e0e6f129be99c4006609c27bbc6eef8f881c211cc571d77ab27"
	t1.AccountNonce = uint64(75544)
	t1.Price = uint64(30000000000)
	t1.GasLimit = uint64(200000)
	t1.TxAmount = uint64(2036190900000000000)
	t1.Payload = []byte{}
	t1.TxV = uint64(37)
	t1.TxR = uint64(7696038145190415973)
	t1.TxS = uint64(11190351604279239082)
	t1.BlockID = uint(1)

	type args struct {
		ctx         context.Context
		tx          *sql.Tx
		ethTrans    *types.Transaction
		BlockID     uint
		BlockNumber uint64
	}
	tests := []struct {
		t       *TransactionService
		args    args
		want    *Transaction
		wantErr bool
	}{
		{
			t: transactionService,
			args: args{
				ctx:         ctx,
				tx:          tx,
				ethTrans:    transactions[0],
				BlockID:     1,
				BlockNumber: block.Number().Uint64(),
			},
			want:    &t1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.t.AddTransaction(tt.args.ctx, tt.args.tx, tt.args.ethTrans, tt.args.BlockID, tt.args.BlockNumber)
		if (err != nil) != tt.wantErr {
			t.Errorf("TransactionService.AddTransaction() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("TransactionService.AddTransaction() = %v, want %v", got, tt.want)
		}
	}
	err = tx.Commit()
	if err != nil {
		t.Error(err)
	}
}

func TestTransactionService_GetBlockTransactions(t *testing.T) {
	ctx := context.Background()

	transactionService := NewTransactionService(appState.Db)

	transactions := []*Transaction{}

	type args struct {
		ctx     context.Context
		BlockID uint
	}
	tests := []struct {
		t       *TransactionService
		args    args
		want    []*Transaction
		wantErr bool
	}{
		{
			t: transactionService,
			args: args{
				ctx:     ctx,
				BlockID: 1,
			},
			want:    transactions,
			wantErr: false,
		},
	}
	for range tests {

		/*got, err := tt.t.GetBlockTransactions(tt.args.ctx, tt.args.BlockID)
		if (err != nil) != tt.wantErr {
			t.Errorf("TransactionService.GetBlockTransactions() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			//t.Errorf("TransactionService.GetBlockTransactions() = %v, want %v", got, tt.want)
		}*/
	}
}
