package svc

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"log"
)

// AppState - used for
type AppState struct {
	Db *sql.DB
}

// DbOptions - used for
type DbOptions struct {
	DB       string `mapstructure:"db"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"hostname"`
	Port     string `mapstructure:"port"`
	Schema   string `mapstructure:"database"`
}

// dbInit - used for database initialization
func dbInit() (*AppState, error) {

	var dbObj DbOptions

	v := viper.New()
	v.AutomaticEnv()

	dbObj.DB = v.GetString("ETHBLOCKS_DB")
	dbObj.User = v.GetString("ETHBLOCKS_DBUSER")
	dbObj.Password = v.GetString("ETHBLOCKS_DBPASS")
	dbObj.Host = v.GetString("ETHBLOCKS_DBHOST")
	dbObj.Port = v.GetString("ETHBLOCKS_DBPORT")
	dbObj.Schema = v.GetString("ETHBLOCKS_DBNAME")

	db, err := sql.Open(dbObj.DB, fmt.Sprint(dbObj.User, ":", dbObj.Password, "@(", dbObj.Host,
		":", dbObj.Port, ")/", dbObj.Schema, "?charset=utf8mb4&parseTime=True"))
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
	appState.Db = db

	return appState, nil

}
