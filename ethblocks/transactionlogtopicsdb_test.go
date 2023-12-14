package ethblocks

import (
	"context"
	"database/sql"
	"math/big"
	"reflect"
	"testing"

	"github.com/cloudfresco/ethblocks/test"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestTransactionLogTopicService_AddTransactionLogTopic(t *testing.T) {
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
	topics := GetTopics(tlog)
	topic := topics[0]
	err = test.LoadSQL(appState)
	if err != nil {
		t.Error(err)
		return
	}

	transactionLogTopicService := NewTransactionLogTopicService(appState.Db)

	tx, err := appState.Db.Begin()
	if err != nil {
		t.Error("err", err)
	}
	transLogTopic := TransactionLogTopic{}
	transLogTopic.Id = uint(175)
	transLogTopic.Topic = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	transLogTopic.BlockId = uint(1)
	transLogTopic.TransactionId = uint(2)
	transLogTopic.TransactionReceiptId = uint(2)
	transLogTopic.TransactionLogId = uint(1)

	type args struct {
		ctx                  context.Context
		tx                   *sql.Tx
		s                    common.Hash
		BlockId              uint
		TransactionId        uint
		TransactionReceiptId uint
		TransactionLogId     uint
	}
	tests := []struct {
		t       *TransactionLogTopicService
		args    args
		want    *TransactionLogTopic
		wantErr bool
	}{
		{
			t: transactionLogTopicService,
			args: args{
				ctx:                  ctx,
				tx:                   tx,
				s:                    topic,
				BlockId:              1,
				TransactionId:        2,
				TransactionReceiptId: 2,
				TransactionLogId:     1,
			},
			want:    &transLogTopic,
			wantErr: false,
		},
	}
	for _, tt := range tests {

		transLogTopicResult, err := tt.t.AddTransactionLogTopic(tt.args.ctx, tt.args.tx, tt.args.s, tt.args.BlockId, tt.args.TransactionId, tt.args.TransactionReceiptId, tt.args.TransactionLogId)
		if (err != nil) != tt.wantErr {
			t.Errorf("TransactionLogTopicService.AddTransactionLogTopic() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(transLogTopicResult, tt.want) {
			t.Errorf("TransactionLogTopicService.AddTransactionLogTopic() = %v, want %v", transLogTopicResult, tt.want)
		}
		assert.NotNil(t, transLogTopicResult)
		assert.Equal(t, transLogTopicResult.Topic, "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", "they should be equal")
	}
	err = tx.Commit()
	if err != nil {
		t.Error(err)
	}
}
