package tests

import (
	"log"
	"net/url"
	"os"

	immudbGorm "github.com/tauu/immudb-gorm"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string
	Age  int
}

func OpenConnection() (*gorm.DB, error) {

	// Creates a database object
	// http://foo/asdfadf/test
	// "file:///test.html"
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// URI to storage location for the database.
	// Example format: immudbe:///folderA/folderB/databaseName
	url := url.URL{
		Scheme: "immudbe",
		Path:   path + "/test/testdb",
	}

	// Open a connection to immudb
	//dsn := "immudb://immudb:immudb@localhost:3322/test"
	db, err := gorm.Open(immudbGorm.Open(url.String()), &gorm.Config{})
	if err != nil {
		log.Fatal("opening database connection failed")
	}

	return db, nil
}

func DeleteTestDir() {
	// Delete the test directory
	err := os.RemoveAll("test")
	if err != nil {
		log.Fatal("deleting test directory failed")
	}
}

func TableChecker(db *gorm.DB, name string) bool {

	// Check if the table exists
	result := db.Migrator().HasTable(name)

	return result
}
