package ethblocks

import (
	"log"
	"os"
	"testing"
)

var (
	appState *AppState
  clientAddr string
)

func TestMain(m *testing.M) {
	var err error
	appState, err = DbInitTest()
	if err != nil {
		log.Fatal(err)
	}
  clientAddr = GetEthblocksClientAddr()
	exitVal := m.Run()
	os.Exit(exitVal)
}
