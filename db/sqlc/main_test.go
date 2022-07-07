package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQuery *Queries
var testDB *sql.DB

//we create a func call testMain, this is the main entry point
// for all test funcs

const (
	dbDriver = "postgres"
	dbSource = "postgresql://dbusername:dbpassword@localhost:5432/simpleBank?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal(err)
	}
	testQuery = New(testDB)
	os.Exit(m.Run())
}
