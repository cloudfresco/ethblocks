package ethblocks

import (
	"log"
	"os"
	"strings"
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

// skipIfPruned skips the test when the RPC node no longer serves the requested
// historical block (public nodes prune old state), and fails on any other error.
func skipIfPruned(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		return
	}
	if strings.Contains(err.Error(), "pruned history unavailable") {
		t.Skip(err)
	}
	t.Fatal(err)
}
