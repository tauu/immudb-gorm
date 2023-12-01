package tests

import (
	"database/sql/driver"
	"fmt"
	"log"
	"net/url"
	"testing"

	"github.com/google/uuid"
	immudbGorm "github.com/tauu/immudb-gorm"
	_ "github.com/tauu/immusql"

	"gorm.io/gorm"
)

// Define a custom uuid type.
type myUUID struct {
	uuid.UUID
}

func (u myUUID) Value() (driver.Value, error) {
	return u.UUID.Value()
}

// Define a custom nullable uuid type.
type myNullUUID struct {
	uuid.UUID
	valid bool
}

func (u myNullUUID) Value() (driver.Value, error) {
	if !u.valid {
		return nil, nil
	}
	return u.UUID.Value()
}

func (u *myNullUUID) Scan(src interface{}) error {
	switch src := src.(type) {
	// nil indicates that the column is NULL in the database.
	case nil:
		u.valid = false
		return nil

	case string:
		id, err := uuid.Parse(src)
		if err != nil {
			return err
		}
		u.valid = true
		u.UUID = id
		return nil
	default:
		return fmt.Errorf("scanning type %T into myNullUUID is not supported, must be string or nil", src)
	}
}

type User struct {
	gorm.Model
	Name       string
	Age        int
	Height     float64
	CompanyID  uuid.UUID  `gorm:"type:UUID"`
	GroupID    myUUID     `gorm:"type:UUID"`
	ContractID myNullUUID `gorm:"type:UUID;nullable:true"`
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
