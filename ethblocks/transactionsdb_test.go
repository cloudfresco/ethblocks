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

func TestTransactionService_AddTransaction(t *testing.T) {
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

	err = test.LoadSQL(appState)
	if err != nil {
		t.Error(err)
		return
	}

	transactionService := NewTransactionService(appState.Db)
	tx, err := appState.Db.Begin()
	if err != nil {
		t.Error(err)
	}
	trans := Transaction{}
	trans.Id = uint(103)
	trans.TxType = uint8(0)
	trans.ChainId = uint64(1)
	trans.TxData = []byte{}
	trans.Gas = uint64(200000)
	trans.GasPrice = uint64(30000000000)
	trans.GasTipCap = uint64(30000000000)
	trans.GasFeeCap = uint64(30000000000)
	trans.TxValue = uint64(2036190900000000000)
	trans.AccountNonce = uint64(75544)
	trans.TxTo = "0x594038b1359D33D750Bdcb448352Cc577a475B81"
	trans.TxV = uint64(37)
	trans.TxR = uint64(7696038145190415973)
	trans.TxS = uint64(11190351604279239082)
	trans.BlockNumber = uint64(7602500)
	trans.BlockHash = "0x02e8467c3c439e0e6f129be99c4006609c27bbc6eef8f881c211cc571d77ab27"
	trans.BlockId = uint(1)

	type args struct {
		ctx         context.Context
		tx          *sql.Tx
		ethTrans    *types.Transaction
		BlockId     uint
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
				BlockId:     1,
				BlockNumber: block.Number().Uint64(),
			},
			want:    &trans,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		transResult, err := tt.t.AddTransaction(tt.args.ctx, tt.args.tx, tt.args.ethTrans, tt.args.BlockId, tt.args.BlockNumber)
		if (err != nil) != tt.wantErr {
			t.Errorf("TransactionService.AddTransaction() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(transResult, tt.want) {
			t.Errorf("TransactionService.AddTransaction() = %v, want %v", transResult, tt.want)
		}
		assert.NotNil(t, transResult)
		assert.Equal(t, transResult.ChainId, uint64(1), "they should be equal")
		assert.Equal(t, transResult.Gas, uint64(200000), "they should be equal")
		assert.Equal(t, transResult.GasPrice, uint64(30000000000), "they should be equal")
		assert.Equal(t, transResult.GasTipCap, uint64(30000000000), "they should be equal")
		assert.Equal(t, transResult.GasFeeCap, uint64(30000000000), "they should be equal")
		assert.Equal(t, transResult.TxValue, uint64(2036190900000000000), "they should be equal")
		assert.Equal(t, transResult.TxTo, "0x594038b1359D33D750Bdcb448352Cc577a475B81", "they should be equal")
		assert.Equal(t, transResult.BlockNumber, uint64(7602500), "they should be equal")
		assert.Equal(t, transResult.BlockHash, "0x02e8467c3c439e0e6f129be99c4006609c27bbc6eef8f881c211cc571d77ab27", "they should be equal")
	}

	err = tx.Commit()
	if err != nil {
		t.Error(err)
	}
}
