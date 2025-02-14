package pgdb

import (
	"database/sql"
	"log"

	"github.com/govargo/go-sqlcommenter/core"
	gosql "github.com/govargo/go-sqlcommenter/database/sql"
	_ "github.com/lib/pq"
)

func ConnectPG(connection string) *sql.DB {
	db, err := gosql.Open("postgres", connection, core.CommenterOptions{
		Config: core.CommenterConfig{EnableDBDriver: true, EnableRoute: true, EnableAction: true, EnableFramework: true, EnableTraceparent: true, EnableApplication: true},
	})
	if err != nil {
		log.Fatalf("Failed to connect to PG(%q), error: %v", connection, err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the database, error: %v", err)
	}

	return db
}
