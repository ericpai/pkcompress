package main

import (
	"database/sql"
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testAnalyzer *Analyzer

var expectedStaffTable = &Table{
	dbName:       "pktest",
	name:         "staff",
	pkColumnName: "id",
	foreignKeys: map[string]*ForeignKey{
		"fk_staff_department_id": {
			tableName:     "staff",
			columnName:    "department_id",
			refTableName:  "department",
			refColumnName: "id",
			updateRule:    "RESTRICT",
			deleteRule:    "RESTRICT",
		},
	},
}

var expectedDepartmentTable = &Table{
	dbName:       "pktest",
	name:         "department",
	pkColumnName: "id",
	foreignKeys: map[string]*ForeignKey{
		"fk_department_region_id": {
			tableName:     "department",
			columnName:    "region_id",
			refTableName:  "region",
			refColumnName: "id",
			updateRule:    "RESTRICT",
			deleteRule:    "RESTRICT",
		},
	},
}

var expectedRegionTable = &Table{
	dbName:       "pktest",
	name:         "region",
	pkColumnName: "id",
	foreignKeys:  make(map[string]*ForeignKey),
}

func TestMain(m *testing.M) {

	connStr := fmt.Sprintf("pktest:dsA12(djks@tcp(127.0.0.1:3306)/pktest?interpolateParams=true&charset=utf8&timeout=1s")
	var conn *sql.DB
	var err error
	if conn, err = sql.Open("mysql", connStr); err != nil {
		panic(err)
	}

	testAnalyzer = newAnalyzer(conn, "pktest")
	m.Run()
}

func TestGetTables(t *testing.T) {
	tables, err := getTables(testAnalyzer.db)
	assert.Nil(t, err)
	sort.Strings(tables)
	expectedTables := []string{
		"staff",
		"department",
		"region",
	}
	sort.Strings(expectedTables)
	assert.EqualValues(t, expectedTables, tables)
}

func TestAnalyseColumnType(t *testing.T) {
	colType1, err := analyseColumnType(testAnalyzer.db, "pktest", "staff", "id")
	assert.Nil(t, err)
	assert.Equal(t, "int", colType1)
	colType2, err := analyseColumnType(testAnalyzer.db, "pktest", "staff", "name")
	assert.Nil(t, err)
	assert.Equal(t, "varchar", colType2)
	colType3, err := analyseColumnType(testAnalyzer.db, "pktest", "staff", "unknown_column")
	assert.NotNil(t, err)
	assert.IsType(t, UnknownColumnNameError{}, err)
	assert.Empty(t, colType3)
}

func TestGetForeignKeys(t *testing.T) {
	foreignKeys, err := getForeignKeys(testAnalyzer.db, "pktest", "staff")
	assert.Nil(t, err)
	assert.EqualValues(t, expectedStaffTable.foreignKeys, foreignKeys)
}
