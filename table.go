package main

import (
	"database/sql"
	"fmt"
	"strconv"
)

type Table struct {
	name         string
	dbName       string
	pkColumnName string
	foreignKeys  map[string]*ForeignKey
}

type ForeignKey struct {
	tableName     string
	columnName    string
	refTableName  string
	refColumnName string
	updateRule    string
	deleteRule    string
}

func (t *Table) cascadeForeginKeys(db *sql.DB) error {
	fmt.Printf("Start to cascade foregin keys of table %s\n", t.name)
	if len(t.foreignKeys) == 0 {
		fmt.Println("> [Skipped] No foreign keys found")
		return nil
	}

	for constraint_name, fk := range t.foreignKeys {
		if _, err := db.Exec(fmt.Sprintf(
			"ALTER TABLE `%s` DROP FOREIGN KEY `%s`", fk.tableName, constraint_name)); err != nil {
			return err
		}
		if _, err := db.Exec(fmt.Sprintf(
			"ALTER TABLE `%s` ADD CONSTRAINT `%s` FOREIGN KEY (`%s`) REFERENCES `%s`(`%s`) on UPDATE CASCADE on DELETE CASCADE",
			fk.tableName, constraint_name, fk.columnName, fk.refTableName, fk.refColumnName,
		)); err != nil {
			return err
		}
	}
	return nil
}

func (t *Table) resumeForeignKeys(db *sql.DB) error {
	fmt.Printf("Start to resume foregin keys of table %s\n", t.name)
	if len(t.foreignKeys) == 0 {
		fmt.Println("> [Skipped] No foreign keys found")
		return nil
	}

	for constraint_name, fk := range t.foreignKeys {
		if _, err := db.Exec(fmt.Sprintf(
			"ALTER TABLE `%s` DROP FOREIGN KEY `%s`", fk.tableName, constraint_name)); err != nil {
			return err
		}
		if _, err := db.Exec(fmt.Sprintf(
			"ALTER TABLE `%s` ADD CONSTRAINT `%s` FOREIGN KEY (`%s`) REFERENCES `%s`(`%s`) on UPDATE %s on DELETE %s",
			fk.tableName, constraint_name, fk.columnName, fk.refTableName, fk.refColumnName, fk.updateRule, fk.deleteRule,
		)); err != nil {
			return err
		}
	}
	return nil
}

func (t *Table) compressPrimaryKey(db *sql.DB) (int, error) {
	fmt.Printf("Start to compress primary key of table %s\n", t.name)
	pks, err := readDataSet(db,
		fmt.Sprintf("SELECT `%s` from `%s` order by `%s`", t.pkColumnName, t.name, t.pkColumnName))
	if err != nil {
		return 0, err
	}
	count := len(pks)
	fmt.Printf("> Rows count: %d. Compressing in a single transaction, please wait...\n", count)
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for i := 0; i < count; i++ {
		pkValue, err := strconv.Atoi(pks[i][t.pkColumnName])
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		if _, err := db.Exec(
			fmt.Sprintf("UPDATE `%s` SET `%s`=? WHERE `%s`=?", t.name, t.pkColumnName, t.pkColumnName),
			i+1, pkValue); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return 0, err
	}
	fmt.Println("> Compressing successfully")

	return count, nil
}

func (t *Table) updateAutoIncrement(db *sql.DB, nextIncrement int) error {
	fmt.Printf("Update auto increment of table %s to %d\n", t.name, nextIncrement)
	_, err := db.Exec(fmt.Sprintf("ALTER TABLE `%s` AUTO_INCREMENT=?", t.name), nextIncrement)
	return err
}

func (t *Table) getAutoIncrement(db *sql.DB) (int, error) {
	var autoIncrement int
	data, err := readDataSet(db,
		"SELECT auto_increment FROM information_schema.tables WHERE table_schema=? AND table_name=?",
		t.dbName, t.name,
	)
	if err != nil {
		return 0, err
	}
	if len(data) != 1 {
		return 0, UnknownTableError{t.dbName, t.name}
	}
	if autoIncrement, err = strconv.Atoi(data[0]["auto_increment"]); err != nil {
		return 0, err
	}
	return autoIncrement, nil

}
