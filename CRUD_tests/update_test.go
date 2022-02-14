package tests

import (
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	// Open connection
	db, err := OpenConnection()
	if !assert.NoError(t, err, "There was an error openning connection") {
		t.FailNow()
	}

	// Delete the test directory
	defer DeleteTestDir()

	// Checking if users table exists after creating it
	usersBeforeExists := TableChecker(db, "users")
	assert.Equal(t, false, usersBeforeExists, "Table users should not exists before creating it")

	// Create a users table
	err = db.AutoMigrate(&User{})
	assert.NoError(t, err, "There was an error creating users table")

	// Checking if users table exists after creating it
	usersAfterExists := TableChecker(db, "users")
	assert.Equal(t, true, usersAfterExists, "Table users should exists after creating it")

	// Define a new record
	var newUser = User{Name: "Jose", Age: 33}

	// Create a new user record
	res := db.Create(&newUser)
	assert.NoError(t, res.Error, "An error occurred while creating a new record")

	// Update a particular record
	result := db.Model(&User{}).Where("id = ?", newUser.ID).Updates(User{Name: "Joel", Age: 100})
	assert.NoError(t, result.Error, "An error occurred while updating")

	// Query the database to check that the data has been updated correctly
	var user User

	err = db.First(&user, newUser.ID).Error
	if err != nil {
		log.Error().Err(err).Msg("No user found was with that ID")
	}

	// Test cases
	assert.Equal(t, "Joel", user.Name, "The update failed")
	assert.Equal(t, 100, user.Age, "The update failed")

}
