package tests

import (
	"log"
	"net/url"
	"testing"

	immudbGorm "github.com/tauu/immudb-gorm"
	_ "github.com/tauu/immusql"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string
	Age  int
}

func OpenConnection(t *testing.T) (*gorm.DB, error) {

	// URI to storage location for the database.
	// Example format: immudbe:///folderA/folderB/databaseName
	url := url.URL{
		Scheme: "immudbe",
		Path:   t.TempDir(),
	}

	// Open a connection to immudb
	//dsn := "immudb://immudb:immudb@localhost:3322/test"
	db, err := gorm.Open(immudbGorm.Open(url.String()), &gorm.Config{})
	if err != nil {
		log.Fatal("opening database connection failed")
	}

	return db, nil
}

func TableChecker(db *gorm.DB, name string) bool {

	// Check if the table exists
	result := db.Migrator().HasTable(name)

	return result
}
