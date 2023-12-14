package ethblocks

import (
	"context"
	"database/sql"
	"math/big"
	"reflect"
	"testing"

	"github.com/cloudfresco/ethblocks/test"
	"github.com/ethereum/go-ethereum/core/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
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
	err = test.LoadSQL(appState)
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
	bu.Id = uint(3)
	bu.ParentHash = "0xedc122f4cc5a34ef716487e81df0522d227e0425406357469e141cfbf772da31"
	bu.UncleHash = "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"
	bu.Coinbase = "0x52bc44d5378309EE2abF1539BF71dE1b7d7bE3b5"
	bu.BlockRoot = "0xee5d08e32fe8f2dda82accd4494a27bfd199da34d0c518aa2c2a595a4e423933"
	bu.TxHash = "0x817922d3ac2b30577a731ae6a0fa91496a87ce9e2bc64a6b0607c21788237b34"
	bu.ReceiptHash = "0xed61edfefe49ba3a263287b3ae6337257b094007f1d05b3f2af612e3080a37dc"
	bu.Bloom = []byte{16, 0, 1, 64, 0, 98, 34, 136, 0, 72, 140, 64, 16, 140, 32, 2, 0, 32, 16, 4, 0, 32, 0, 128, 5, 49, 17, 1, 4, 144, 33, 72, 26, 32, 1, 16, 16, 72, 90, 132, 197, 26, 160, 69, 0, 0, 1, 0, 0, 0, 2, 7, 128, 3, 193, 0, 35, 128, 13, 140, 0, 172, 0, 135, 4, 0, 8, 194, 33, 193, 224, 57, 32, 2, 131, 40, 34, 170, 0, 16, 72, 68, 64, 33, 0, 52, 16, 0, 7, 1, 104, 128, 16, 2, 0, 0, 13, 0, 0, 9, 6, 3, 65, 4, 8, 7, 17, 8, 0, 10, 8, 80, 82, 0, 131, 5, 8, 112, 0, 0, 8, 129, 8, 26, 65, 2, 0, 130, 0, 164, 16, 0, 0, 0, 48, 128, 0, 40, 6, 80, 145, 4, 4, 136, 0, 64, 65, 0, 0, 192, 0, 72, 160, 14, 144, 4, 6, 18, 1, 40, 130, 1, 10, 16, 0, 145, 1, 2, 150, 44, 148, 128, 136, 0, 0, 32, 146, 1, 4, 36, 130, 1, 148, 4, 80, 32, 0, 17, 12, 6, 169, 4, 64, 130, 32, 50, 16, 14, 144, 64, 36, 98, 128, 0, 20, 11, 50, 41, 128, 17, 128, 196, 240, 12, 129, 1, 33, 30, 0, 0, 16, 34, 112, 0, 8, 0, 66, 13, 166, 2, 8, 2, 0, 8, 0, 2, 5, 68, 10, 129, 36, 69, 128, 8, 40, 149, 130, 0, 10, 24, 97, 68, 16, 0, 8, 165}
	bu.Difficulty = uint64(1917036994703655)
	bu.BlockNumber = uint64(7602499)
	bu.GasLimit = uint64(8000029)
	bu.GasUsed = uint64(6557700)
	bu.BlockTime = uint64(1555735033)
	bu.Extra = []byte{80, 80, 89, 69, 32, 110, 97, 110, 111, 112, 111, 111, 108, 46, 111, 114, 103}
	bu.MixDigest = "0x4afc378aaceda9f84aaaca15050f9f7ea24eab2381dd190126ea87cd3e990909"
	bu.BlockNonce = uint64(10441032150657948177)
	bu.ParentBeaconRoot = ""
	bu.BlockId = uint(1)

	type args struct {
		ctx      context.Context
		tx       *sql.Tx
		blkuncle *types.Header
		BlockId  uint
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
				BlockId:  1,
			},
			want:    &bu,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		uncleResult, err := tt.bu.AddBlockUncle(tt.args.ctx, tt.args.tx, tt.args.blkuncle, tt.args.BlockId)
		if (err != nil) != tt.wantErr {
			t.Errorf("BlockUncleService.AddBlockUncle() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(uncleResult, tt.want) {
			t.Errorf("BlockUncleService.AddBlockUncle() = %v, want %v", uncleResult, tt.want)
		}
		assert.NotNil(t, uncleResult)
		assert.Equal(t, uncleResult.ParentHash, "0xedc122f4cc5a34ef716487e81df0522d227e0425406357469e141cfbf772da31", "they should be equal")
		assert.Equal(t, uncleResult.UncleHash, "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347", "they should be equal")
		assert.Equal(t, uncleResult.Coinbase, "0x52bc44d5378309EE2abF1539BF71dE1b7d7bE3b5", "they should be equal")
		assert.Equal(t, uncleResult.BlockRoot, "0xee5d08e32fe8f2dda82accd4494a27bfd199da34d0c518aa2c2a595a4e423933", "they should be equal")
		assert.Equal(t, uncleResult.TxHash, "0x817922d3ac2b30577a731ae6a0fa91496a87ce9e2bc64a6b0607c21788237b34", "they should be equal")
		assert.Equal(t, uncleResult.ReceiptHash, "0xed61edfefe49ba3a263287b3ae6337257b094007f1d05b3f2af612e3080a37dc", "they should be equal")
		assert.Equal(t, uncleResult.Bloom, bu.Bloom, "they should be equal")
		assert.Equal(t, uncleResult.Difficulty, uint64(1917036994703655), "they should be equal")
		assert.Equal(t, uncleResult.GasUsed, uint64(6557700), "they should be equal")
		assert.Equal(t, uncleResult.MixDigest, "0x4afc378aaceda9f84aaaca15050f9f7ea24eab2381dd190126ea87cd3e990909", "they should be equal")
		assert.Equal(t, uncleResult.BlockNonce, uint64(10441032150657948177), "they should be equal")
	}
	err = tx.Commit()
	if err != nil {
		t.Error(err)
	}
}

func TestBlockUncleService_GetBlockUncles(t *testing.T) {
	ctx := context.Background()
	err := test.LoadSQL(appState)
	if err != nil {
		t.Error(err)
		return
	}

	blockUncleService := NewBlockUncleService(appState.Db)

	uncles := []*BlockUncle{}
	bu1 := BlockUncle{}
	bu1.Id = uint(1)
	bu1.ParentHash = "0xedc122f4cc5a34ef716487e81df0522d227e0425406357469e141cfbf772da31"
	bu1.UncleHash = "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"
	bu1.Coinbase = "0x52bc44d5378309EE2abF1539BF71dE1b7d7bE3b5"
	bu1.BlockRoot = "0xee5d08e32fe8f2dda82accd4494a27bfd199da34d0c518aa2c2a595a4e423933"
	bu1.TxHash = "0x817922d3ac2b30577a731ae6a0fa91496a87ce9e2bc64a6b0607c21788237b34"
	bu1.ReceiptHash = "0xed61edfefe49ba3a263287b3ae6337257b094007f1d05b3f2af612e3080a37dc"
	bu1.Bloom = []byte{16, 0, 1, 64, 0, 98, 34, 136, 0, 72, 140, 64, 16, 140, 32, 2, 0, 32, 16, 4, 0, 32, 0, 128, 5, 49, 17, 1, 4, 144, 33, 72, 26, 32, 1, 16, 16, 72, 90, 132, 197, 26, 160, 69, 0, 0, 1, 0, 0, 0, 2, 7, 128, 3, 193, 0, 35, 128, 13, 140, 0, 172, 0, 135, 4, 0, 8, 194, 33, 193, 224, 57, 32, 2, 131, 40, 34, 170, 0, 16, 72, 68, 64, 33, 0, 52, 16, 0, 7, 1, 104, 128, 16, 2, 0, 0, 13, 0, 0, 9, 6, 3, 65, 4, 8, 7, 17, 8, 0, 10, 8, 80, 82, 0, 131, 5, 8, 112, 0, 0, 8, 129, 8, 26, 65, 2, 0, 130, 0, 164, 16, 0, 0, 0, 48, 128, 0, 40, 6, 80, 145, 4, 4, 136, 0, 64, 65, 0, 0, 192, 0, 72, 160, 14, 144, 4, 6, 18, 1, 40, 130, 1, 10, 16, 0, 145, 1, 2, 150, 44, 148, 128, 136, 0, 0, 32, 146, 1, 4, 36, 130, 1, 148, 4, 80, 32, 0, 17, 12, 6, 169, 4, 64, 130, 32, 50, 16, 14, 144, 64, 36, 98, 128, 0, 20, 11, 50, 41, 128, 17, 128, 196, 240, 12, 129, 1, 33, 30, 0, 0, 16, 34, 112, 0, 8, 0, 66, 13, 166, 2, 8, 2, 0, 8, 0, 2, 5, 68, 10, 129, 36, 69, 128, 8, 40, 149, 130, 0, 10, 24, 97, 68, 16, 0, 8, 165}
	bu1.Difficulty = uint64(1917036994703655)
	bu1.BlockNumber = uint64(7602499)
	bu1.GasLimit = uint64(8000029)
	bu1.GasUsed = uint64(6557700)
	bu1.BlockTime = uint64(1555735033)
	bu1.Extra = []byte{80, 80, 89, 69, 32, 110, 97, 110, 111, 112, 111, 111, 108, 46, 111, 114, 103}
	bu1.MixDigest = "0x4afc378aaceda9f84aaaca15050f9f7ea24eab2381dd190126ea87cd3e990909"
	bu1.BlockNonce = uint64(10441032150657948177)
	bu1.ParentBeaconRoot = ""
	bu1.BlockId = uint(1)

	bu2 := BlockUncle{}
	bu2.Id = uint(2)
	bu2.ParentHash = "0xedc122f4cc5a34ef716487e81df0522d227e0425406357469e141cfbf772da31"
	bu2.UncleHash = "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"
	bu2.Coinbase = "0xb2930B35844a230f00E51431aCAe96Fe543a0347"
	bu2.BlockRoot = "0x021ae28620d5faa1b08c28af9531d03e73b7a712d5bd9861fce10c3333ab051d"
	bu2.TxHash = "0x7017dc17389685301bbb5c2b4420c2a0393f19e02f9944e676de7c95b17a99e8"
	bu2.ReceiptHash = "0x8ba50f482753e692429572561ed56de18c9337667e5d5666a6f2992c1b7685e2"
	bu2.Bloom = []byte{16, 0, 1, 64, 0, 98, 34, 136, 0, 72, 140, 64, 16, 140, 32, 2, 0, 32, 16, 4, 0, 32, 0, 128, 5, 49, 17, 1, 4, 144, 33, 72, 26, 32, 1, 16, 16, 72, 90, 132, 197, 26, 160, 69, 0, 0, 1, 0, 0, 0, 2, 7, 128, 3, 193, 0, 35, 128, 13, 140, 0, 172, 0, 135, 4, 0, 8, 194, 33, 193, 224, 57, 32, 2, 131, 40, 34, 170, 0, 16, 72, 68, 64, 33, 0, 52, 16, 0, 7, 1, 104, 128, 16, 2, 0, 0, 13, 0, 0, 9, 6, 3, 65, 4, 8, 7, 17, 8, 0, 10, 8, 80, 82, 0, 131, 5, 8, 112, 0, 0, 8, 129, 8, 26, 65, 2, 0, 130, 0, 164, 16, 0, 0, 0, 48, 128, 0, 40, 6, 80, 145, 4, 4, 136, 0, 64, 65, 0, 0, 192, 0, 72, 160, 14, 144, 4, 6, 18, 1, 40, 130, 1, 10, 16, 0, 145, 1, 2, 150, 44, 148, 128, 136, 0, 0, 32, 146, 1, 4, 36, 130, 1, 148, 4, 80, 32, 0, 17, 12, 6, 169, 4, 64, 130, 32, 50, 16, 14, 144, 64, 36, 98, 128, 0, 20, 11, 50, 41, 128, 17, 128, 196, 240, 12, 129, 1, 33, 30, 0, 0, 16, 34, 112, 0, 8, 0, 66, 13, 166, 2, 8, 2, 0, 8, 0, 2, 5, 68, 10, 129, 36, 69, 128, 8, 40, 149, 130, 0, 10, 24, 97, 68, 16, 0, 8, 165}
	bu2.Difficulty = uint64(1917036994703655)
	bu2.BlockNumber = uint64(7602499)
	bu2.GasLimit = uint64(8000029)
	bu2.GasUsed = uint64(6578700)
	bu2.BlockTime = uint64(1555735033)
	bu2.Extra = []byte{115, 105, 110, 103, 50}
	bu2.MixDigest = "0x2063fae923c6e5b474e366426ff3b4d735ab67edef767914db411a3207a7efb7"
	bu2.BlockNonce = uint64(4565172828713654969)
	bu2.ParentBeaconRoot = ""
	bu2.BlockId = uint(1)

	uncles = append(uncles, &bu1)
	uncles = append(uncles, &bu2)

	type args struct {
		ctx     context.Context
		BlockId uint
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
				BlockId: 1,
			},
			want:    uncles,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		unclesResp, err := tt.bu.GetBlockUncles(tt.args.ctx, tt.args.BlockId)
		if (err != nil) != tt.wantErr {
			t.Errorf("BlockUncleService.GetBlockUncles() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		assert.NotNil(t, unclesResp)
		assert.Equal(t, unclesResp, tt.want)
		assert.NotNil(t, unclesResp)
		uncleResult1 := unclesResp[0]
		uncleResult2 := unclesResp[1]
		assert.Equal(t, uncleResult1.ParentHash, "0xedc122f4cc5a34ef716487e81df0522d227e0425406357469e141cfbf772da31", "they should be equal")
		assert.Equal(t, uncleResult1.UncleHash, "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347", "they should be equal")
		assert.Equal(t, uncleResult1.Coinbase, "0x52bc44d5378309EE2abF1539BF71dE1b7d7bE3b5", "they should be equal")
		assert.Equal(t, uncleResult1.BlockRoot, "0xee5d08e32fe8f2dda82accd4494a27bfd199da34d0c518aa2c2a595a4e423933", "they should be equal")
		assert.Equal(t, uncleResult1.TxHash, "0x817922d3ac2b30577a731ae6a0fa91496a87ce9e2bc64a6b0607c21788237b34", "they should be equal")
		assert.Equal(t, uncleResult1.ReceiptHash, "0xed61edfefe49ba3a263287b3ae6337257b094007f1d05b3f2af612e3080a37dc", "they should be equal")
		assert.Equal(t, uncleResult1.Bloom, bu1.Bloom, "they should be equal")
		assert.Equal(t, uncleResult1.Difficulty, uint64(1917036994703655), "they should be equal")
		assert.Equal(t, uncleResult1.GasUsed, uint64(6557700), "they should be equal")
		assert.Equal(t, uncleResult1.MixDigest, "0x4afc378aaceda9f84aaaca15050f9f7ea24eab2381dd190126ea87cd3e990909", "they should be equal")
		assert.Equal(t, uncleResult1.BlockNonce, uint64(10441032150657948177), "they should be equal")
		assert.Equal(t, uncleResult2.ParentHash, "0xedc122f4cc5a34ef716487e81df0522d227e0425406357469e141cfbf772da31", "they should be equal")
		assert.Equal(t, uncleResult2.UncleHash, "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347", "they should be equal")
		assert.Equal(t, uncleResult2.Coinbase, "0xb2930B35844a230f00E51431aCAe96Fe543a0347", "they should be equal")
		assert.Equal(t, uncleResult2.BlockRoot, "0x021ae28620d5faa1b08c28af9531d03e73b7a712d5bd9861fce10c3333ab051d", "they should be equal")
		assert.Equal(t, uncleResult2.TxHash, "0x7017dc17389685301bbb5c2b4420c2a0393f19e02f9944e676de7c95b17a99e8", "they should be equal")
		assert.Equal(t, uncleResult2.ReceiptHash, "0x8ba50f482753e692429572561ed56de18c9337667e5d5666a6f2992c1b7685e2", "they should be equal")
		assert.Equal(t, uncleResult2.Bloom, bu2.Bloom, "they should be equal")
		assert.Equal(t, uncleResult2.Difficulty, uint64(1917036994703655), "they should be equal")
		assert.Equal(t, uncleResult2.GasUsed, uint64(6578700), "they should be equal")
		assert.Equal(t, uncleResult2.MixDigest, "0x2063fae923c6e5b474e366426ff3b4d735ab67edef767914db411a3207a7efb7", "they should be equal")
		assert.Equal(t, uncleResult2.BlockNonce, uint64(4565172828713654969), "they should be equal")
		assert.Equal(t, uncleResult2.BlockNumber, uint64(7602499), "they should be equal")

	}
}
