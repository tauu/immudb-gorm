package test_migrator

import (
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestAutoMigrate(t *testing.T) {

	// Open connection
	db, err := OpenConnection()
	if !assert.NoError(t, err, "An error ocurred while openning connection") {
		t.FailNow()
	}

	// Deletes the test directory
	defer DeleteTestDir()

	// Check if the table employees exist before creating it
	isTableCreatedBefore := db.Migrator().HasTable(&Employee{})
	assert.Equal(t, isTableCreatedBefore, false, "Table employees should not exist before creating it")

	// Create an employees table
	err = db.AutoMigrate(&Employee{})
	if err != nil {
		log.Error().Err(err).Msg("An error occurred while creating a new table")
	}

	// Check if the table employees exist after creating it
	isTableCreatedAfter := db.Migrator().HasTable(&Employee{})
	assert.Equal(t, isTableCreatedAfter, true, "Table employees should exist after creating it")

}
