package service_test

import (
	"database/sql"
	"log"
	"reflect"
	"sort"
	"testing"
)

func TestCreateTables(t *testing.T) {

	dbase, err := sql.Open("mysql", "root:root@/FinanceTesting")
	if err != nil {
		log.Fatal("Failed to open connection to database")
	}

	var tables []string
	rows, err := dbase.Query("SHOW TABLES")
	if err != nil {
		log.Fatal("Failed to query tables from the database")
	}

	for rows.Next() {
		var tableName string
		rows.Scan(&tableName)
		tables = append(tables, tableName)
	}

	expectedTables := []string{"Goals", "MonthlyExpenses", "Transactions", "Variables"}

	sort.Strings(tables)
	sort.Strings(expectedTables)

	if reflect.DeepEqual(tables, expectedTables) {
		t.Fail()
	}

}
