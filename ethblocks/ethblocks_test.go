package ethblocks

import (
	"log"
	"os"
	"testing"

	etbcommon "github.com/cloudfresco/ethblocks/common"
	"github.com/cloudfresco/ethblocks/test"
)

var (
	appState   *etbcommon.AppState
	clientAddr string
)

func TestMain(m *testing.M) {
	var err error
	appState, err = test.DbInitTest()
	if err != nil {
		log.Fatal(err)
	}
	clientAddr = GetEthblocksClientAddr()
	exitVal := m.Run()
	os.Exit(exitVal)
}
