package mysqldb

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/govargo/go-sqlcommenter/core"
	gosql "github.com/govargo/go-sqlcommenter/database/sql"
)

func ConnectMySQL(connection string) *sql.DB {
	db, err := gosql.Open("mysql", connection, core.CommenterOptions{
		Config: core.CommenterConfig{EnableDBDriver: true, EnableRoute: true, EnableAction: true, EnableFramework: true, EnableTraceparent: true, EnableApplication: true},
	})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL(%q), error: %v", connection, err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the database, error: %v", err)
	}

	return db
}
