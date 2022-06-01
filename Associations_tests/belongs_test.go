package test_associations

import (
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestBelongs(t *testing.T) {

	// Open connection
	db, err := OpenConnection()
	if !assert.NoError(t, err, "An error occurred opening connection") {
		t.FailNow()
	}

	// Delete test directory
	defer DeleteTestDir()

	// Check if tables employees and companies exists before creating them
	employeesExists := TableChecker("employees", db)
	companiesExists := TableChecker("companies", db)

	// Test cases
	assert.Equal(t, false, employeesExists, "Table employees already exists before creating it")
	assert.Equal(t, false, companiesExists, "Table companies already exists before creating it")

	// Create employees table
	err = db.AutoMigrate(&Employee{})
	assert.NoError(t, err, "An error occurred while creating tables")

	// Check if tables employees and companies were created
	employeesExists = TableChecker("employees", db)
	companiesExists = TableChecker("companies", db)

	// Test cases
	assert.Equal(t, true, employeesExists, "There was a problem creating the table employees")
	assert.Equal(t, true, companiesExists, "There was a problem creating the table companies")

	// Define a new record
	newEmployee := Employee{
		Name: "Joel",
		Company: Company{
			Name: "Net`Q GmbH",
		},
	}

	// Insert a new record to the data base
	result := db.Create(&newEmployee)
	assert.NoError(t, result.Error, "There was an error inserting a new record")

	// Make sure that the data previously inserted is there
	var employee Employee

	err = db.Preload("Company").First(&employee).Error
	if err != nil {
		log.Error().Err(err).Msg("There was an error quering the first employee")
	}

	// Test cases
	assert.Equal(t, 1, int(employee.ID), "There was an error quering employee's ID")
	assert.Equal(t, "Joel", employee.Name, "There was an error quering employee's name")
	assert.Equal(t, 1, int(employee.CompanyID), "There was an error quering employee's companyID")
	assert.Equal(t, "Net`Q GmbH", employee.Company.Name, "There was an error quering the company's name")

}
