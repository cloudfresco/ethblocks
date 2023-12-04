package ethblocks

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/spf13/viper"
)

const (
	ethblocksDbsqlMysqlTest     = "fixtures/ethblocks_mysql_test.sql"
	ethblocksDbsqlMysqlSchema   = "sql/mysql/ethblocks_mysql_schema.sql"
	ethblocksDbsqlMysqlTruncate = "fixtures/ethblocks_mysql_truncate.sql"
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
	dbOpt.MySQLTestFilePath = ethblocksDbsqlMysqlTest
	dbOpt.MySQLSchemaFilePath = ethblocksDbsqlMysqlSchema
	dbOpt.MySQLTruncateFilePath = ethblocksDbsqlMysqlTruncate
	dbOpt.PgSQLTestFilePath = ""
	dbOpt.PgSQLSchemaFilePath = ""
	dbOpt.PgSQLTruncateFilePath = ""
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
	appState.MySQLSchemaFilePath = dbOpt.MySQLSchemaFilePath
	appState.MySQLTruncateFilePath = dbOpt.MySQLTruncateFilePath
	appState.PgSQLTestFilePath = dbOpt.PgSQLTestFilePath
	appState.PgSQLSchemaFilePath = dbOpt.PgSQLSchemaFilePath
	appState.PgSQLTruncateFilePath = dbOpt.PgSQLTruncateFilePath
	return appState, nil
}

func execSQLFile(ctx context.Context, sqlFilePath string, db *sql.DB) error {
	content, err := ioutil.ReadFile(sqlFilePath)
	if err != nil {
		log.Println(err)
		return err
	}

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Fatal(err)
	}

	sqlLines := strings.Split(string(content), ";\n")

	for _, sqlLine := range sqlLines {
		if sqlLine != "" {
			_, err := tx.ExecContext(ctx, sqlLine)
			if err != nil {
				log.Println(err)
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					log.Printf("Load SQL failed: %v, unable to rollback: %v\n", err, rollbackErr)
					return err
				}
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// LoadSQL -- truncate, load data
func LoadSQL(appState *AppState) error {
	var err error
	ctx := context.Background()

	if appState.DbType == DbMysql {
		err = execSQLFile(ctx, appState.MySQLTruncateFilePath, appState.Db)
		if err != nil {
			log.Println(err)
			return err
		}

		err = execSQLFile(ctx, appState.MySQLTestFilePath, appState.Db)
		if err != nil {
			log.Println(err)
			return err
		}

	} else if appState.DbType == DbPgsql {
		err = execSQLFile(ctx, appState.PgSQLTruncateFilePath, appState.Db)
		if err != nil {
			log.Println(err)
			return err
		}

		err = execSQLFile(ctx, appState.PgSQLTestFilePath, appState.Db)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}
