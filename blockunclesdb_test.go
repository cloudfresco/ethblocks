package ethblocks

import (
	"context"
	"database/sql"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	_ "github.com/go-sql-driver/mysql"
)

func TestBlockUncleService_AddBlockUncle(t *testing.T) {

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
	uncles := GetUncles(block)
	uncle := uncles[0]
	err = LoadSQL(appState)
	if err != nil {
		t.Error(err)
		return
	}

	blockUncleService := NewBlockUncleService(appState.Db)
	tx, err := appState.Db.Begin()
	if err != nil {
		t.Error(err)
	}
	bu := BlockUncle{}
	bu.ID = uint(3)
	bu.BlockNumber = uint64(7602499)
	bu.BlockTime = uint64(1555735033)
	bu.ParentHash = "0xedc122f4cc5a34ef716487e81df0522d227e0425406357469e141cfbf772da31"
	bu.UncleHash = "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"
	bu.BlockRoot = "0xee5d08e32fe8f2dda82accd4494a27bfd199da34d0c518aa2c2a595a4e423933"
	bu.TxHash = "0x817922d3ac2b30577a731ae6a0fa91496a87ce9e2bc64a6b0607c21788237b34"
	bu.ReceiptHash = "0xed61edfefe49ba3a263287b3ae6337257b094007f1d05b3f2af612e3080a37dc"
	bu.MixDigest = "0x4afc378aaceda9f84aaaca15050f9f7ea24eab2381dd190126ea87cd3e990909"
	bu.BlockNonce = uint64(10441032150657948177)
	bu.Coinbase = "0x52bc44d5378309EE2abF1539BF71dE1b7d7bE3b5"
	bu.GasLimit = uint64(8000029)
	bu.GasUsed = uint64(6557700)
	bu.Difficulty = uint64(1917036994703655)
	bu.BlockSize = uncle.Size()
	bu.BlockID = uint(1)

	type args struct {
		ctx      context.Context
		tx       *sql.Tx
		blkuncle *types.Header
		BlockID  uint
	}
	tests := []struct {
		bu      *BlockUncleService
		args    args
		want    *BlockUncle
		wantErr bool
	}{
		{
			bu: blockUncleService,
			args: args{
				ctx:      ctx,
				tx:       tx,
				blkuncle: uncle,
				BlockID:  1,
			},
			want:    &bu,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.bu.AddBlockUncle(tt.args.ctx, tt.args.tx, tt.args.blkuncle, tt.args.BlockID)
		if (err != nil) != tt.wantErr {
			t.Errorf("BlockUncleService.AddBlockUncle() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("BlockUncleService.AddBlockUncle() = %v, want %v", got, tt.want)
		}
	}
	err = tx.Commit()
	if err != nil {
		t.Error(err)
	}
}

func TestBlockUncleService_GetBlockUncles(t *testing.T) {
	ctx := context.Background()
	err := LoadSQL(appState)
	if err != nil {
		t.Error(err)
		return
	}

	blockUncleService := NewBlockUncleService(appState.Db)

	uncles := []*BlockUncle{}
	bu1 := BlockUncle{}
	bu1.ID = uint(1)
	bu1.BlockNumber = uint64(7602499)
	bu1.BlockTime = uint64(1555735033)
	bu1.ParentHash = "0xedc122f4cc5a34ef716487e81df0522d227e0425406357469e141cfbf772da31"
	bu1.UncleHash = "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"
	bu1.BlockRoot = "0xee5d08e32fe8f2dda82accd4494a27bfd199da34d0c518aa2c2a595a4e423933"
	bu1.TxHash = "0x817922d3ac2b30577a731ae6a0fa91496a87ce9e2bc64a6b0607c21788237b34"
	bu1.ReceiptHash = "0xed61edfefe49ba3a263287b3ae6337257b094007f1d05b3f2af612e3080a37dc"
	bu1.MixDigest = "0x4afc378aaceda9f84aaaca15050f9f7ea24eab2381dd190126ea87cd3e990909"
	bu1.BlockNonce = uint64(10441032150657948177)
	bu1.Coinbase = "0x52bc44d5378309EE2abF1539BF71dE1b7d7bE3b5"
	bu1.GasLimit = uint64(8000029)
	bu1.GasUsed = uint64(6557700)
	bu1.Difficulty = uint64(1917036994703655)
	bu1.BlockSize = common.StorageSize(570)
	bu1.BlockID = uint(1)

	bu2 := BlockUncle{}
	bu2.ID = uint(2)
	bu2.BlockNumber = uint64(7602499)
	bu2.BlockTime = uint64(1555735033)
	bu2.ParentHash = "0xedc122f4cc5a34ef716487e81df0522d227e0425406357469e141cfbf772da31"
	bu2.UncleHash = "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"
	bu2.BlockRoot = "0x021ae28620d5faa1b08c28af9531d03e73b7a712d5bd9861fce10c3333ab051d"
	bu2.TxHash = "0x7017dc17389685301bbb5c2b4420c2a0393f19e02f9944e676de7c95b17a99e8"
	bu2.ReceiptHash = "0x8ba50f482753e692429572561ed56de18c9337667e5d5666a6f2992c1b7685e2"
	bu2.MixDigest = "0x2063fae923c6e5b474e366426ff3b4d735ab67edef767914db411a3207a7efb7"
	bu2.BlockNonce = uint64(4565172828713654969)
	bu2.Coinbase = "0xb2930B35844a230f00E51431aCAe96Fe543a0347"
	bu2.GasLimit = uint64(8000029)
	bu2.GasUsed = uint64(6578700)
	bu2.Difficulty = uint64(1917036994703655)
	bu2.BlockSize = common.StorageSize(558)
	bu2.BlockID = uint(1)

	uncles = append(uncles, &bu1)
	uncles = append(uncles, &bu2)

	type args struct {
		ctx     context.Context
		BlockID uint
	}
	tests := []struct {
		bu      *BlockUncleService
		args    args
		want    []*BlockUncle
		wantErr bool
	}{
		{
			bu: blockUncleService,
			args: args{
				ctx:     ctx,
				BlockID: 1,
			},
			want:    uncles,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := tt.bu.GetBlockUncles(tt.args.ctx, tt.args.BlockID)
		if (err != nil) != tt.wantErr {
			t.Errorf("BlockUncleService.GetBlockUncles() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("BlockUncleService.GetBlockUncles() = %v, want %v", got, tt.want)
		}
	}

}
