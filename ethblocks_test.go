package ethblocks

import (
	"log"
	"os"
	"testing"

	"gopkg.in/testfixtures.v2"
)

var (
	appState *AppState
	fixtures *testfixtures.Context
)

func TestMain(m *testing.M) {
	var err error

	appState, err = DbInitTest()
	if err != nil {
		log.Fatal(err)
	}

	fixtures, err = FixturesInit(appState.DbType, appState.Db)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}
