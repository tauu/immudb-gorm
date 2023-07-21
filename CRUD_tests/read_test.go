package tests

import (
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {

	// Open connection
	db, err := OpenConnection(t)
	require.NoError(t, err, "There was an error opening connection")

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

	// Query the data previously inserted to the database
	var users []User

	err = db.Find(&users).Error
	if err != nil {
		log.Error().Err(err).Msg("No users found")
	}

	// Test cases
	assert.Equal(t, newUser.Name, users[0].Name, "The queried name does not match the defined name")
	assert.Equal(t, newUser.Age, users[0].Age, "The queried age does not match the defined age")
}
