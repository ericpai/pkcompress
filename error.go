package main

import "fmt"

type NoSinglePrimaryKeyError struct {
	tableName string
}

func (e NoSinglePrimaryKeyError) Error() string {
	return fmt.Sprintf("table %s doesn't have a single primary key", e.tableName)
}

type NonIntegerColumnTypeError struct {
	columnName string
	typeName   string
}

func (e NonIntegerColumnTypeError) Error() string {
	return fmt.Sprintf("column %s is not an integer type, but %s", e.columnName, e.typeName)
}

type UnknownColumnNameError struct {
	dbName     string
	tableName  string
	columnName string
}

func (e UnknownColumnNameError) Error() string {
	return fmt.Sprintf("unknown column %s.%s of database %s", e.tableName, e.columnName, e.dbName)
}

type UnknownTableError struct {
	dbName    string
	tableName string
}

func (e UnknownTableError) Error() string {
	return fmt.Sprintf("unknown table %s of database %s", e.tableName, e.dbName)
}
