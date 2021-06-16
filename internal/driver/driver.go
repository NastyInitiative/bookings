package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// DB holds the database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifetime = 5 * time.Minute

// ConnectSQL - Creates database bool for Postgres;
//
// dsn - "database connection string"
func ConnectSQL(dsn string) (*DB, error) {
	newDb, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	newDb.SetMaxOpenConns(maxOpenDbConn)
	newDb.SetMaxIdleConns(maxIdleDbConn)
	newDb.SetConnMaxLifetime(maxDbLifetime)

	dbConn.SQL = newDb
	err = testDB(newDb)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

// NewDatabase - creates a new database for the application
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// testDb - tries to ping database
func testDB(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}

	return nil
}
