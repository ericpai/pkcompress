package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCascadeAndResumeForeignKeys(t *testing.T) {

	assert.Nil(t, expectedStaffTable.cascadeForeginKeys(testAnalyzer.db))
	foreignKeys, err := getForeignKeys(testAnalyzer.db, expectedStaffTable.dbName, expectedStaffTable.name)
	assert.Nil(t, err)
	for constraint_name, expectedFk := range expectedStaffTable.foreignKeys {
		actualFk := foreignKeys[constraint_name]
		assert.NotNil(t, actualFk)
		assert.Equal(t, expectedFk.tableName, actualFk.tableName)
		assert.Equal(t, expectedFk.columnName, actualFk.columnName)
		assert.Equal(t, expectedFk.refTableName, actualFk.refTableName)
		assert.Equal(t, expectedFk.refColumnName, actualFk.refColumnName)
		assert.Equal(t, "CASCADE", actualFk.deleteRule)
		assert.Equal(t, "CASCADE", actualFk.updateRule)
	}

	assert.Nil(t, expectedStaffTable.resumeForeignKeys(testAnalyzer.db))
	foreignKeys, err = getForeignKeys(testAnalyzer.db, expectedStaffTable.dbName, expectedStaffTable.name)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedStaffTable.foreignKeys, foreignKeys)
}

func TestUpdateAndGetAutoIncrement(t *testing.T) {

	expectedAutoIncrement := 100
	assert.Nil(t, expectedStaffTable.updateAutoIncrement(testAnalyzer.db, expectedAutoIncrement))
	actualAutoIncrement, err := expectedStaffTable.getAutoIncrement(testAnalyzer.db)
	assert.Nil(t, err)
	assert.Equal(t, expectedAutoIncrement, actualAutoIncrement)
}

func TestCompressPrimaryKey(t *testing.T) {

	expectedStaffTable.cascadeForeginKeys(testAnalyzer.db)
	count, err := expectedStaffTable.compressPrimaryKey(testAnalyzer.db)
	assert.Nil(t, err)
	assert.Equal(t, 3, count)
	actualData, err := readDataSet(testAnalyzer.db, "SELECT * FROM staff")
	assert.Nil(t, err)
	assert.Equal(t, 3, len(actualData))

	for i := 0; i < 3; i++ {
		assert.Equal(t, expectedStaffData[i]["id"], actualData[i]["id"])
		assert.Equal(t, expectedStaffData[i]["name"], actualData[i]["name"])
	}
	expectedStaffTable.resumeForeignKeys(testAnalyzer.db)
}
