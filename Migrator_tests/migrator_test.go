package test_migrator

import (
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestHasTable(t *testing.T) {

	// Open connection
	db, err := OpenConnection(t)
	require.NoError(t, err, "An error ocurred while opening connection")

	// Check if the table employees exist before creating it
	isTableCreatedBefore := db.Migrator().HasTable(&Employee{})
	assert.Equal(t, isTableCreatedBefore, false, "Table employees should not exist before creating it")

	// Create an employees table
	err = db.Migrator().CreateTable(&Employee{})
	if err != nil {
		log.Error().Err(err).Msg("An error occurred while creating a new table")
	}

	// Check if the table employees exist after creating it
	isTableCreatedAfter := db.Migrator().HasTable(&Employee{})
	assert.Equal(t, isTableCreatedAfter, true, "Table employees should exist after creating it")

}

func TestRenameTable(t *testing.T) {

	// Open connection
	db, err := OpenConnection(t)
	require.NoError(t, err, "An error ocurred while opening connection")

	// Check if the table employees exist before creating it
	isTableCreatedBefore := db.Migrator().HasTable(&Employee{})
	assert.Equal(t, isTableCreatedBefore, false, "Table employees should not exist before creating it")

	// Create an employees table
	err = db.Migrator().CreateTable(&Employee{})
	if err != nil {
		log.Error().Err(err).Msg("An error occurred while creating a new table")
	}

	// Check if the table employees exist after creating it
	isTableCreatedAfter := db.Migrator().HasTable(&Employee{})
	assert.Equal(t, isTableCreatedAfter, true, "Table employees should exist after creating it")

	type EmployeeOfTheMonth struct {
		gorm.Model
		Promote bool
	}

	// Rename the table employees.
	err = db.Migrator().RenameTable(&Employee{}, &EmployeeOfTheMonth{})
	assert.NoError(t, err, "renaming a table should not cause an error")

	employeeTableExists := db.Migrator().HasTable(&Employee{})
	employeeOfTheMonthTableExists := db.Migrator().HasTable(&EmployeeOfTheMonth{})
	assert.False(t, employeeTableExists, "Table employees table should no longer exists after it was renamed")
	assert.True(t, employeeOfTheMonthTableExists, "Table employeeOfTheMonths table should exists after the employees table was renamed")
}

func TestHasIndex(t *testing.T) {
	// Open connection
	db, err := OpenConnection(t)
	require.NoError(t, err, "An error ocurred while opening connection")

	// Create an employees table
	err = db.Migrator().CreateTable(&Employee{})
	if err != nil {
		log.Error().Err(err).Msg("An error occurred while creating a new table")
	}

	// Check if the table employees has the default index for the deleted_at column.
	hasPrimaryIndex := db.Migrator().HasIndex(&Employee{}, "idx_employees_deleted_at")
	assert.Equal(t, true, hasPrimaryIndex, "Table employees should have an index for the deleted_at column.")
}

func TestGetIndexes(t *testing.T) {
	// Open connection
	db, err := OpenConnection(t)
	require.NoError(t, err, "An error ocurred while opening connection")

	// Create an employees table
	err = db.Migrator().CreateTable(&Employee{})
	require.NoError(t, err, "creating a table in an empty database should not cause an error")

	// Check if the table employees has the default index for the deleted_at column.
	indexes, err := db.Migrator().GetIndexes(&Employee{})
	assert.NoError(t, err, "getting the indexes of an existing table should not cause an error")
	require.Equal(t, 2, len(indexes), "the test table should have two indexes, one for the primary id and one for the deleted_at column")
	// The first index should be the primary key.
	primary, ok := indexes[0].PrimaryKey()
	assert.True(t, ok, "checking if an index is the primary key should never fail")
	assert.True(t, primary, "the first index should also be the primary one")
	unique, ok := indexes[0].Unique()
	assert.True(t, ok, "checking if an index is unique should never fail")
	assert.True(t, unique, "the primary index should be unique")
	assert.Equal(t, []string{"id"}, indexes[0].Columns(), "the primary index should only contain the id column")
	// The second index should be for the deleted_at column.
	primary, ok = indexes[1].PrimaryKey()
	assert.True(t, ok, "checking if an index is the primary key should never fail")
	assert.False(t, primary, "the second index should not be the primary key")
	unique, ok = indexes[1].Unique()
	assert.True(t, ok, "checking if an index is unique should never fail")
	assert.False(t, unique, "the second index should not be unique")
	assert.Equal(t, []string{"deleted_at"}, indexes[1].Columns(), "the second index should only contain the deleted_at column")
}

func TestGetTables(t *testing.T) {
	// Open connection
	db, err := OpenConnection(t)
	require.NoError(t, err, "An error ocurred while opening connection")

	// Create an employees table
	err = db.Migrator().CreateTable(&Employee{})
	require.NoError(t, err, "creating a table in an empty database should not cause an error")

	// Get a list of tables and verify that the just created table is part of it.
	tables, err := db.Migrator().GetTables()
	assert.NoError(t, err, "retrieving tables from a database should never fail")
	assert.Equal(t, tables, []string{"employees"}, "the database should contain only the beforehand created table")
}

func TestHasColumn(t *testing.T) {
	// Open connection
	db, err := OpenConnection(t)
	require.NoError(t, err, "An error ocurred while opening connection")

	// Create an employees table
	err = db.Migrator().CreateTable(&Employee{})
	require.NoError(t, err, "creating a table in an empty database should not cause an error")

	// Get a list of tables and verify that the just created table is part of it.
	result := db.Migrator().HasColumn(&Employee{}, "salary")
	assert.True(t, result, "the created table should have a column named salary")
}

func TestRenameColumn(t *testing.T) {
	// Open connection
	db, err := OpenConnection(t)
	require.NoError(t, err, "An error ocurred while opening connection")

	// Create an employees table
	err = db.Migrator().CreateTable(&Employee{})
	require.NoError(t, err, "creating a table in an empty database should not cause an error")

	// Rename the salary column.
	err = db.Migrator().RenameColumn(&Employee{}, "salary", "new_salary")
	assert.NoError(t, err, "renaming a column should not cause an error")

	hasNewColumn := db.Migrator().HasColumn(&Employee{}, "new_salary")
	assert.True(t, hasNewColumn, "the table should have the renamed column")
}

func TestAddColumn(t *testing.T) {
	// Open connection
	db, err := OpenConnection(t)
	require.NoError(t, err, "An error ocurred while opening connection")

	// Create an employees table
	err = db.Migrator().CreateTable(&Employee{})
	require.NoError(t, err, "creating a table in an empty database should not cause an error")

	type Employee struct {
		gorm.Model
		Promote bool
	}

	// Add a promote column.
	err = db.Migrator().AddColumn(&Employee{}, "promote")
	assert.NoError(t, err, "adding a column should not cause an error")

	hasNewColumn := db.Migrator().HasColumn(&Employee{}, "promote")
	assert.True(t, hasNewColumn, "the table should have the added column")
}

func TestDropColumn(t *testing.T) {
	// Open connection
	db, err := OpenConnection(t)
	require.NoError(t, err, "An error ocurred while opening connection")

	// Create an employees table
	err = db.Migrator().CreateTable(&Employee{})
	require.NoError(t, err, "creating a table in an empty database should not cause an error")

	// Remove salary column.
	err = db.Migrator().DropColumn(&Employee{}, "salary")
	assert.NoError(t, err, "dropping a column should not cause an error")

	hasNewColumn := db.Migrator().HasColumn(&Employee{}, "salary")
	assert.False(t, hasNewColumn, "the table should not have the dropped column")
}
