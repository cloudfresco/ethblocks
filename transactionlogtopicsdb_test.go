package ethblocks

import (
	"context"
	"database/sql"
	"log"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestTransactionLogTopicService_AddTransactionLogTopic(t *testing.T) {
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
	receipt, err := GetTransactionReceipt(ctx, client, transactions[1].Hash())
	if err != nil {
		log.Fatal(err)
	}
	tlogs := GetLogs(receipt)
	tlog := tlogs[0]
	topics := GetTopics(tlog)
	topic := topics[0]
	err = LoadSQL(appState)
	if err != nil {
		log.Println(err)
		return
	}

	transactionLogTopicService := NewTransactionLogTopicService(appState.Db)

	tx, err := appState.Db.Begin()
	if err != nil {
		log.Println("err", err)
	}
	transLogTopic := TransactionLogTopic{}
	transLogTopic.ID = uint(175)
	transLogTopic.Topic = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	transLogTopic.BlockID = uint(1)
	transLogTopic.TransactionID = uint(2)
	transLogTopic.TransactionReceiptID = uint(2)
	transLogTopic.TransactionLogID = uint(1)

	type args struct {
		ctx                  context.Context
		tx                   *sql.Tx
		s                    common.Hash
		BlockID              uint
		TransactionID        uint
		TransactionReceiptID uint
		TransactionLogID     uint
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
				BlockID:              1,
				TransactionID:        2,
				TransactionReceiptID: 2,
				TransactionLogID:     1,
			},
			want:    &transLogTopic,
			wantErr: false,
		},
	}
	for _, tt := range tests {

		got, err := tt.t.AddTransactionLogTopic(tt.args.ctx, tt.args.tx, tt.args.s, tt.args.BlockID, tt.args.TransactionID, tt.args.TransactionReceiptID, tt.args.TransactionLogID)
		if (err != nil) != tt.wantErr {
			t.Errorf("TransactionLogTopicService.AddTransactionLogTopic() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("TransactionLogTopicService.AddTransactionLogTopic() = %v, want %v", got, tt.want)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Println(err)
	}
	err = DeleteSQL(appState)
	if err != nil {
		log.Println(err)
		return
	}
}

func TestTransactionLogTopicService_GetTransactionLogTopics(t *testing.T) {
	ctx := context.Background()

	transactionLogTopicService := NewTransactionLogTopicService(appState.Db)

	transactionLogTopics := []*TransactionLogTopic{}

	type args struct {
		ctx              context.Context
		TransactionLogID uint
	}
	tests := []struct {
		t       *TransactionLogTopicService
		args    args
		want    []*TransactionLogTopic
		wantErr bool
	}{
		{
			t: transactionLogTopicService,
			args: args{
				ctx:              ctx,
				TransactionLogID: 1,
			},
			want:    transactionLogTopics,
			wantErr: false,
		},
	}
	for range tests {

		/*got, err := tt.t.GetTransactionLogTopics(tt.args.ctx, tt.args.TransactionLogID)
		if (err != nil) != tt.wantErr {
			t.Errorf("TransactionLogTopicService.GetTransactionLogTopics() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("TransactionLogTopicService.GetTransactionLogTopics() = %v, want %v", got, tt.want)
		}*/

	}
}
