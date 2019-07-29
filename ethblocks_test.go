package ethblocks

import (
	"log"
	"os"
	"testing"
)

var (
	appState *AppState
)

func TestMain(m *testing.M) {
	var err error
	appState, err = DbInitTest()
	if err != nil {
		log.Fatal(err)
	}
	exitVal := m.Run()
	os.Exit(exitVal)
}
