package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/hyson007/simpleBank/util"
	_ "github.com/lib/pq"
)

var testQuery *Queries
var testDB *sql.DB

//we create a func call testMain, this is the main entry point
// for all test funcs

func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("unable to load config:", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}
	testQuery = New(testDB)
	os.Exit(m.Run())
}
