package main

import (
	"database/sql"
	"fmt"
)

var intType = map[string]bool{
	"tinyint":   false,
	"smallint":  false,
	"mediumint": false,
	"int":       false,
	"bigint":    false,
}

// readDataSet executes the query string with placeholders replaced by args and returns the dataset.
func readDataSet(db *sql.DB, query string, args ...interface{}) ([]map[string]string, error) {
	var err error
	var result *sql.Rows
	var columnName []string
	if result, err = db.Query(query, args...); err != nil {
		return nil, err
	}
	defer result.Close()
	if columnName, err = result.Columns(); err != nil {
		return nil, err
	}
	columnCount := len(columnName)
	columnValue := make([]interface{}, columnCount)
	var dataset []map[string]string
	for result.Next() {
		for i := 0; i < columnCount; i++ {
			columnValue[i] = new([]byte)
		}
		if err = result.Scan(columnValue...); err != nil {
			return nil, err
		}
		row := make(map[string]string)
		for i := 0; i < columnCount; i++ {
			row[columnName[i]] = string(*columnValue[i].(*[]byte))
		}
		dataset = append(dataset, row)
	}
	return dataset, nil
}

func getTables(db *sql.DB) ([]string, error) {
	var tables []string
	tableData, err := readDataSet(db, "SHOW TABLES")
	if err != nil {
		return tables, err
	}
	for _, table := range tableData {
		for _, tableName := range table {
			tables = append(tables, tableName)
			fmt.Printf("Find table: %s\n", tableName)
		}
	}
	return tables, nil
}

func analyseTable(db *sql.DB, dbName, tableName string) (*Table, error) {
	fmt.Printf("Start to analyse table: %s.%s\n", dbName, tableName)

	result, err := readDataSet(db,
		`select kc.column_name as column_name from information_schema.KEY_COLUMN_USAGE
			as kc where kc.CONSTRAINT_NAME='PRIMARY' and table_schema=? and table_name=?`, dbName, tableName)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, NoSinglePrimaryKeyError{dbName}
	}
	pkName := result[0]["column_name"]
	fmt.Printf("> Find the single primary key: %s\n", pkName)

	pkType, err := analyseColumnType(db, dbName, tableName, pkName)

	if _, exist := intType[pkType]; !exist {
		fmt.Printf(">> The primary key is not integer type, but %s\n", pkType)
		pkName = ""
	}

	fmt.Println("> Start to find all foreign keys")
	foreginKeys, err := getForeignKeys(db, dbName, tableName)
	if err != nil {
		return nil, err
	}
	if len(foreginKeys) == 0 {
		fmt.Println(">> No foreigns are found")
	}

	table := &Table{
		name:         tableName,
		dbName:       dbName,
		pkColumnName: pkName,
		foreignKeys:  foreginKeys,
	}
	return table, err
}

func analyseColumnType(db *sql.DB, dbName, tableName, columnName string) (string, error) {
	fmt.Printf(">> Start to analyse column %s.%s", tableName, columnName)
	result, err := readDataSet(db,
		`select data_type from information_schema.columns where table_schema=? and table_name=? and column_name=?`,
		dbName, tableName, columnName,
	)
	if err != nil {
		return "", err
	}
	if len(result) != 1 {
		return "", UnknownColumnNameError{dbName, tableName, columnName}
	}
	dataType := result[0]["data_type"]
	fmt.Printf(" with data type: %s\n", dataType)
	return dataType, nil
}

func getForeignKeys(db *sql.DB, dbName, tableName string) (map[string]*ForeignKey, error) {
	result, err := readDataSet(db,
		`select
			t.CONSTRAINT_NAME as constraint_name,
			k.COLUMN_NAME as column_name,
			k.REFERENCED_TABLE_NAME as ref_table_name,
			k.REFERENCED_COLUMN_NAME as ref_column_name,
			r.UPDATE_RULE as update_rule,
			r.DELETE_RULE as delete_rule
		from
			information_schema.table_constraints as t,
			information_schema.referential_constraints as r,
			information_schema.KEY_COLUMN_USAGE as k
		where
			t.TABLE_SCHEMA=? and
			t.TABLE_NAME=? and
			t.constraint_type='FOREIGN KEY' and 
			t.CONSTRAINT_NAME = r.CONSTRAINT_NAME and
			k.CONSTRAINT_NAME = r.CONSTRAINT_NAME`,
		dbName, tableName,
	)
	foreignKeys := make(map[string]*ForeignKey)
	if err != nil {
		return foreignKeys, nil
	}
	for _, row := range result {
		fk := &ForeignKey{
			tableName:     tableName,
			columnName:    row["column_name"],
			refTableName:  row["ref_table_name"],
			refColumnName: row["ref_column_name"],
			updateRule:    row["update_rule"],
			deleteRule:    row["delete_rule"],
		}
		foreignKeys[row["constraint_name"]] = fk
		fmt.Printf(">> Constraint name %s with update rule: %s\n", row["constraint_name"], fk.updateRule)
	}
	return foreignKeys, nil
}
