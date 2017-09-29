package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectedStaffData = []map[string]string{
	{
		"id":            "1",
		"name":          "Mr. Little",
		"department_id": "1",
	},
	{
		"id":            "2",
		"name":          "Mrs. Medium",
		"department_id": "1",
	},
	{
		"id":            "3",
		"name":          "Mr. Big",
		"department_id": "2",
	},
}

var expectedDepartmentData = []map[string]string{
	{
		"id":        "1",
		"name":      "Sales",
		"region_id": "1",
	},
	{
		"id":        "2",
		"name":      "Tech",
		"region_id": "2",
	},
}

var expectedRegionData = []map[string]string{
	{
		"id":   "1",
		"name": "East",
	},
	{
		"id":   "2",
		"name": "West",
	},
}

func TestAnalyseAndCompress(t *testing.T) {
	assert.Nil(t, testAnalyzer.analyseAndCompress())
	actualStaffData, err := readDataSet(testAnalyzer.db, "SELECT * FROM `staff`")
	assert.Nil(t, err)
	assert.EqualValues(t, expectedStaffData, actualStaffData)
	actualDepartmentData, err := readDataSet(testAnalyzer.db, "SELECT * FROM `department`")
	assert.Nil(t, err)
	assert.EqualValues(t, expectedDepartmentData, actualDepartmentData)
	actualRegionData, err := readDataSet(testAnalyzer.db, "SELECT * FROM `region`")
	assert.Nil(t, err)
	assert.EqualValues(t, expectedRegionData, actualRegionData)

	actualStaffAutoIncrement, err := expectedStaffTable.getAutoIncrement(testAnalyzer.db)
	assert.Nil(t, err)
	assert.Equal(t, 4, actualStaffAutoIncrement)

	actualDepartmentAutoIncrement, err := expectedDepartmentTable.getAutoIncrement(testAnalyzer.db)
	assert.Nil(t, err)
	assert.Equal(t, 3, actualDepartmentAutoIncrement)

	actualRegionAutoDepartment, err := expectedRegionTable.getAutoIncrement(testAnalyzer.db)
	assert.Nil(t, err)
	assert.Equal(t, 3, actualRegionAutoDepartment)
}
