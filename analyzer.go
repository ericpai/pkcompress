package main

import (
	"database/sql"
	"fmt"
)

type Analyzer struct {
	db     *sql.DB
	dbName string
}

func newAnalyzer(db *sql.DB, dbName string) *Analyzer {
	return &Analyzer{
		db:     db,
		dbName: dbName,
	}
}

func (a *Analyzer) analyseAndCompress() error {
	tables, err := getTables(a.db)
	if err != nil {
		return err
	}
	analysedTables := make([]*Table, 0, len(tables))
	for _, tableName := range tables {
		if table, err := analyseTable(a.db, a.dbName, tableName); err != nil {
			fmt.Printf("! Analyse table %s error: %s\n", tableName, err.Error())
			return err
		} else {
			analysedTables = append(analysedTables, table)
		}
	}
	for _, table := range analysedTables {
		if err := table.cascadeForeginKeys(a.db); err != nil {
			fmt.Printf("! Cascade table foreign keys error: %s\n", err.Error())
			return err
		}
	}
	for _, table := range analysedTables {
		count, err := table.compressPrimaryKey(a.db)
		if err != nil {
			fmt.Printf("! Compress table primary key error: %s\n", err.Error())
			return err
		}
		if err := table.updateAutoIncrement(a.db, count+1); err != nil {
			fmt.Printf("! Update auto_increment error: %s\n", err.Error())
			return err
		}
	}
	for _, table := range analysedTables {
		if err := table.resumeForeignKeys(a.db); err != nil {
			fmt.Printf("! Resume table foreign keys error: %s\n", err.Error())
			return err
		}
	}
	fmt.Println("### WORK DONE! ###")
	return nil
}
