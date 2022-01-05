package immudbGorm

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/tauu/immusql"
	"gorm.io/gorm"
	"gorm.io/gorm/migrator"
)

// Use the default migrator of gorm as a base and define only new functions
// whenever the default ones use features not supported by immudb.
type Migrator struct {
	migrator.Migrator
}

// -- Migrator interface --

// HasTable determines if a specific tables exists in the current database.
func (m Migrator) HasTable(value interface{}) bool {
	tableExists := false
	err := m.RunWithValue(value, func(stmt *gorm.Statement) error {
		// Retrieve the database connector.
		db, err := m.DB.DB()
		if err != nil {
			return err
		}
		// Get an actual connection to the database.
		conn, err := db.Conn(context.Background())
		err = conn.Raw(func(driverConn interface{}) error {
			// Check if the connection uses the immusql driver.
			if v, ok := driverConn.(immusql.ImmuDBconn); ok {
				// Use the ExistTable function of the driver
				// to determine if the table exists.
				tableExists, err = v.ExistTable(stmt.Table)
				return err
			}
			return errors.New("connection is not a ImmuDBConn")
		})
		return err
	})
	if err != nil {
		log.Printf("error checking if a table exists: %v", err)
		return false
	}
	return tableExists
}

// CreateIndex creates an index on a table. This code has been copied from the default CreateIndex
// function of the gorm migrator and adjusted to not define a name for an index. At the moment immudb
// does not support named indexes.
func (m Migrator) CreateIndex(value interface{}, name string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if idx := stmt.Schema.LookIndex(name); idx != nil {
			opts := m.DB.Migrator().(migrator.BuildIndexOptionsInterface).BuildIndexOptions(idx.Fields, stmt)
			values := []interface{}{m.CurrentTable(stmt), opts}

			createIndexSQL := "CREATE "
			// Classes are currently not suppored.
			//
			// if idx.Class != "" {
			// 	createIndexSQL += idx.Class + " "
			// }
			createIndexSQL += "INDEX ON ??"

			// Types are currently not suppored.
			//
			// if idx.Type != "" {
			// 	createIndexSQL += " USING " + idx.Type
			// }

			// Comments are not supported
			//
			// if idx.Comment != "" {
			// 	createIndexSQL += fmt.Sprintf(" COMMENT '%s'", idx.Comment)
			// }

			if idx.Option != "" {
				createIndexSQL += " " + idx.Option
			}

			return m.DB.Exec(createIndexSQL, values...).Error
		}

		return fmt.Errorf("failed to create index with name %s", name)
	})
}
