package ethblocks

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"log"
)

// DbMysql for DbType is mysql
const DbMysql string = "mysql"

// DbPgsql for DbType is pgsql
const DbPgsql string = "pgsql"

// AppState - used for
type AppState struct {
	DbType                string
	Db                    *sql.DB
	Schema                string
	MySQLTestFilePath     string
	MySQLSchemaFilePath   string
	MySQLTruncateFilePath string
	PgSQLTestFilePath     string
	PgSQLSchemaFilePath   string
	PgSQLTruncateFilePath string
}

// DbOptions - used for
type DbOptions struct {
	DB                    string `mapstructure:"db"`
	User                  string `mapstructure:"user"`
	Password              string `mapstructure:"password"`
	Host                  string `mapstructure:"hostname"`
	Port                  string `mapstructure:"port"`
	Schema                string `mapstructure:"database"`
	MySQLTestFilePath     string `mapstructure:"mysql_test_file_path"`
	MySQLSchemaFilePath   string `mapstructure:"mysql_schema_file_path"`
	MySQLTruncateFilePath string `mapstructure:"mysql_truncate_file_path"`
	PgSQLTestFilePath     string `mapstructure:"pgsql_test_file_path"`
	PgSQLSchemaFilePath   string `mapstructure:"pgsql_schema_file_path"`
	PgSQLTruncateFilePath string `mapstructure:"pgsql_truncate_file_path"`
}

// DbInit - used for database initialization
func DbInit() (*AppState, error) {

	var dbOpt DbOptions
	var db *sql.DB
	var err error

	v := viper.New()
	v.AutomaticEnv()

	dbOpt.DB = v.GetString("ETHBLOCKS_DB")
	dbOpt.User = v.GetString("ETHBLOCKS_DBUSER")
	dbOpt.Password = v.GetString("ETHBLOCKS_DBPASS")
	dbOpt.Host = v.GetString("ETHBLOCKS_DBHOST")
	dbOpt.Port = v.GetString("ETHBLOCKS_DBPORT")
	dbOpt.Schema = v.GetString("ETHBLOCKS_DBNAME")

	if dbOpt.DB == DbMysql {
		db, err = sql.Open(dbOpt.DB, fmt.Sprint(dbOpt.User, ":", dbOpt.Password, "@(", dbOpt.Host,
			":", dbOpt.Port, ")/", dbOpt.Schema, "?charset=utf8mb4&parseTime=True"))
		if err != nil {
			log.Println(err)
			return nil, err
		}
	} else if dbOpt.DB == DbPgsql {

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

// DbClose - used for closing database
func DbClose(db *sql.DB) error {
	err := db.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
