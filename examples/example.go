package main

import (
	"context"
	"log"
	"math/big"
	//"crypto/ecdsa"

	"github.com/cloudfresco/ethblocks/svc"
	//"github.com/ethereum/go-ethereum/crypto"
	//"github.com/ethereum/go-ethereum/common"
)

func main() {
	log.Println("Account Examples")
	ExAccount()
	log.Println("Block Examples")
	ExBlock()
	log.Println("Transaction Examples")
	CreateTx()
}

// ExAccount - used to test account examples
func ExAccount() {

	clientaddr := "https://rinkeby.infura.io"
	client, err := svc.GetClient(clientaddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(client)
	ctx := context.Background()
	account := "0x7AF3A1f8F9864F8E8B6A465F4eaeFa15321029f4"
	balance, err := svc.GetBalance(ctx, client, account)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("GetBalance :", balance)

	balance2, err := svc.GetBalance2(ctx, clientaddr, account)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("GetBalance2 :", balance2)

	balance3, err := svc.GetPendingBalanceAt(ctx, client, account)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("GetPendingBalanceAt:", balance3)

	balance4, err := svc.GetPendingBalanceAt2(ctx, clientaddr, account)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("GetPendingBalanceAt2:", balance4)

}

// ExBlock - used to test block examples
func ExBlock() {
	client, err := svc.GetClient("https://mainnet.infura.io")
	//client, err := ethclient.Dial("https://rinkeby.infura.io")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	blockNumber := big.NewInt(7602500)
	log.Println("GetBlockByNumber")
	block, err := svc.GetBlockByNumber(ctx, client, blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("GetBlockByHash")
	h := block.Hash()
	_, err = svc.GetBlockByHash(ctx, client, h)
	if err != nil {
		log.Fatal(err)
	}

	blocknumber, err := svc.BlockNumber(ctx, client)
	log.Println("Blocknumber:", blocknumber)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("GetBlocks")
	blocks, err := svc.GetBlocks(ctx, client, big.NewInt(7602500), big.NewInt(7602509))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(blocks)

	count, err := svc.GetUncleCountByBlockNumber(ctx, client, blockNumber)
	log.Println("GetUncleCountByBlockNumber:", count)

	count, err = svc.GetUncleCountByBlockHash(ctx, client, h)
	log.Println("GetUncleCountByBlockHash:", count)

	log.Println("GetBlocksByMiner:")
	blocks, err = svc.GetBlocksByMiner(ctx, client, "0x5a0b54d5dc17e0aadc383d2db43b0a0d3e029c4c", big.NewInt(7602500), big.NewInt(7602509))
}

// CreateTx - used to test Trasaction examples
func CreateTx() {
	client, err := svc.GetClient("https://mainnet.infura.io")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(client)
	ctx := context.Background()

	blockNumber := big.NewInt(7602500)
	log.Println("GetBlockByNumber")
	block, err := svc.GetBlockByNumber(ctx, client, blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	txs := block.Transactions()
	blockhash := txs[0].Hash()
	tx, _, err := svc.GetTransactionByHash(ctx, client, blockhash)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("tx:", tx)

	txs1, err := svc.GetTransactionsByAddress(ctx, client, "0xEec606A66edB6f497662Ea31b5eb1610da87AB5f", big.NewInt(7602500), big.NewInt(7602509))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("txs:", txs1)
	/*client, err := svc.GetClient("https://rinkeby.infura.io")
	  if err != nil {
	      log.Fatal(err)
	  }

		//log.Println("client:", client)
	  //transaction
		privateKeyStr1 :=   "8cbf2c86bf186e170a67aa6bbdb55a7b36c0080f698eab810d9f45bdf4c9e4dd"

	  //address2 := 0x54320380fD5b5e0389A1Ea23d25689e354dD0Dd9
		privateKeyStr2 := "0d7ed2be1a81713eecfbdabb660a835893fbf1e9550e12c3f902c2562046fcb6"



	  privateKey1, err := crypto.HexToECDSA(privateKeyStr1)
	  if err != nil {
	      log.Fatal(err)
	  }

	  privateKey2, err := crypto.HexToECDSA(privateKeyStr2)
	  if err != nil {
	      log.Fatal(err)
	  }


	  publicKey1 := privateKey1.Public()
	  publicKeyECDSA1, ok := publicKey1.(*ecdsa.PublicKey)
	  if !ok {
	      log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	  }

	  fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA1)
		//log.Println("fromAddress:", fromAddress)

	  publicKey2 := privateKey2.Public()
	  publicKeyECDSA2, ok := publicKey2.(*ecdsa.PublicKey)
	  if !ok {
	      log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	  }

	  toAddress := crypto.PubkeyToAddress(*publicKeyECDSA2)
		//log.Println("toAddress:", toAddress)


	  nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
		//log.Println("nonce:", nonce)
	  if err != nil {
	      log.Fatal(err)
	  }

	  value := big.NewInt(1000000000000000000) // in wei (1 eth)
	  gasLimit := uint64(21000)                // in units
	  gasPrice, err := client.SuggestGasPrice(context.Background())
		//log.Println("gasPrice:", gasPrice)
	  if err != nil {
	      log.Fatal(err)
	  }

	  var data []byte
	  tx := svc.CreateRawTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

		log.Println("tx:", tx)*/

}
