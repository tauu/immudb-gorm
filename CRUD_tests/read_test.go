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

func TestLimit(t *testing.T) {

	// Open connection
	db, err := OpenConnection(t)
	require.NoError(t, err, "There was an error opening connection")

	// Create a users table
	err = db.AutoMigrate(&User{})
	require.NoError(t, err, "There was an error creating users table")

	// Define a new record
	var newUser1 = User{Name: "Jose", Age: 33}
	var newUser2 = User{Name: "Dave", Age: 35}

	// Create a new user record
	err = db.Create(&newUser1).Error
	require.NoError(t, err, "An error occurred while creating the first user")
	err = db.Create(&newUser2).Error
	require.NoError(t, err, "An error occurred while creating the second user")

	// Query the data previously inserted to the database
	var users []User

	err = db.Debug().Limit(1).Find(&users).Error
	if err != nil {
		log.Error().Err(err).Msg("No users found")
	}

	// Test cases
	assert.Len(t, users, 1, "The query should only return one row if the limit is 1")
}
