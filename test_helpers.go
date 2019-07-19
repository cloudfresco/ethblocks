package ethblocks

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/testfixtures.v2"
	"log"
)

// DbInitTest - used for database initialization
func DbInitTest() (*AppState, error) {

	var dbOpt DbOptions

	v := viper.New()
	v.AutomaticEnv()

	dbOpt.DB = v.GetString("ETHBLOCKS_DB")
	dbOpt.User = v.GetString("ETHBLOCKS_DBUSER_TEST")
	dbOpt.Password = v.GetString("ETHBLOCKS_DBPASS_TEST")
	dbOpt.Host = v.GetString("ETHBLOCKS_DBHOST")
	dbOpt.Port = v.GetString("ETHBLOCKS_DBPORT")
	dbOpt.Schema = v.GetString("ETHBLOCKS_DBNAME_TEST")

	db, err := sql.Open(dbOpt.DB, fmt.Sprint(dbOpt.User, ":", dbOpt.Password, "@(", dbOpt.Host,
		":", dbOpt.Port, ")/", dbOpt.Schema, "?charset=utf8mb4&parseTime=True"))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	appState := &AppState{}
	appState.DbType = dbOpt.DB
	appState.Db = db

	return appState, nil

}

// FixturesInit initialize fixtures
func FixturesInit(dbType string, db *sql.DB) (*testfixtures.Context, error) {
	var err error
	var fixtures *testfixtures.Context
	if dbType == "mysql" {
		fixtures, err = testfixtures.NewFolder(db, &testfixtures.MySQL{}, "fixtures")
		if err != nil {
			log.Println(err)
			return fixtures, err
		}
	} else if dbType == "pgsql" {
		fixtures, err = testfixtures.NewFolder(db, &testfixtures.PostgreSQL{}, "testdata/fixtures")
		if err != nil {
			log.Println(err)
			return fixtures, err
		}
	}
	return fixtures, nil
}
