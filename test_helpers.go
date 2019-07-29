package ethblocks

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"io/ioutil"
	"strings"
)

// DbInitTest - used for database initialization
func DbInitTest() (*AppState, error) {
	var dbOpt DbOptions
	var db *sql.DB
	var err error

	v := viper.New()
	v.AutomaticEnv()

	dbOpt.DB = v.GetString("ETHBLOCKS_DB")
	dbOpt.User = v.GetString("ETHBLOCKS_DBUSER_TEST")
	dbOpt.Password = v.GetString("ETHBLOCKS_DBPASS_TEST")
	dbOpt.Host = v.GetString("ETHBLOCKS_DBHOST")
	dbOpt.Port = v.GetString("ETHBLOCKS_DBPORT")
	dbOpt.Schema = v.GetString("ETHBLOCKS_DBNAME_TEST")
	dbOpt.MySQLTestFilePath = v.GetString("ETHBLOCKS_DBSQL_MYSQL_TEST")
	dbOpt.PgSQLTestFilePath = v.GetString("ETHBLOCKS_DBSQL_PGSQL_TEST")

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
	appState.Schema = dbOpt.Schema
	appState.MySQLTestFilePath = dbOpt.MySQLTestFilePath
	appState.PgSQLTestFilePath = dbOpt.PgSQLTestFilePath
	return appState, nil

}

// LoadSQL -- load data into all tables
func LoadSQL(appState *AppState) error {
	ctx := context.Background()
	content, err := ioutil.ReadFile(appState.MySQLTestFilePath)

	if err != nil {
		log.Println(err)
		return err
	}

	sqlLines := strings.Split(string(content), ";\n")

	for _, sqlLine := range sqlLines {

		if sqlLine != "" {
			_, err := appState.Db.ExecContext(ctx, sqlLine)
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}
	return nil
}

// DeleteSQL -- delete data from all tables
func DeleteSQL(appState *AppState) error {
	ctx := context.Background()
	tables := []string{}
	tableSchema := appState.Schema
	sql := "select table_name from information_schema.tables where table_schema = " + " '" + tableSchema + "' " + ";"
	rows, err := appState.Db.QueryContext(ctx, sql)
	if err != nil {
		log.Println(err)
		return err
	}
	var tableName string
	for rows.Next() {
		err = rows.Scan(&tableName)
		if err != nil {
			log.Println(err)
			err = rows.Close()
			if err != nil {
				log.Println(err)
				return err
			}
			return err
		}
		tables = append(tables, tableName)
	}
	err = rows.Close()
	if err != nil {
		log.Println(err)
		return err
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
		return err
	}

	for _, tableName := range tables {
		sql = "truncate " + tableName
		_, err := appState.Db.ExecContext(ctx, sql)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}
