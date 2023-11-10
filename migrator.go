package immudbGorm

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
)

type ErrMissingImmuDBsupport struct {
	Operation string
}

func (err *ErrMissingImmuDBsupport) Error() string {
	return "the migrator operation " + err.Operation + " is currently not implemented as immudb does not provide the required functionality for it."
}

// Use the default migrator of gorm as a base and define only new functions
// whenever the default ones use features not supported by immudb.
type Migrator struct {
	migrator.Migrator
}

// -- Migrator interface --

// AddColumn creates a column with the given name in the table referenced by value.
// This function is an almost identical copy of the default add column function.
// The only difference is the string "COLUMN" inserted after "ADD" in the sql query.
func (m Migrator) AddColumn(value interface{}, name string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		// avoid using the same name field
		f := stmt.Schema.LookUpField(name)
		if f == nil {
			return fmt.Errorf("failed to look up field with name: %s", name)
		}

		if !f.IgnoreMigration {
			return m.DB.Exec(
				"ALTER TABLE ? ADD COLUMN ? ?",
				m.CurrentTable(stmt), clause.Column{Name: f.DBName}, m.DB.Migrator().FullDataTypeOf(f),
			).Error
		}

		return nil
	})
}

// CreateConstraint creates a constraint on a table.
//
// Not implemented as immudb does not support constraints.
func (m Migrator) CreateConstraint(value interface{}, name string) error {
	return &ErrMissingImmuDBsupport{"CreateConstraint"}
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

// CreateView creates a view.
//
// Not implemented as immudb does not support views.
func (m Migrator) CreateView(name string, option gorm.ViewOption) error {
	return &ErrMissingImmuDBsupport{"CreateView"}
}

// CurrentDatabase returns the currently selected database.
func (m Migrator) CurrentDatabase() (name string) {
	// TODO add support in the driver to determine the current database.
	return ""
}

// DropColumn does not have a custom implementation as the default one is
// compatible with immudb.
//func (m Migrator) DropColumn(value interface{}, name string) error

// DropConstraint removes a constraint from a table.
//
// Not implemented as immudb does not support constraints.
func (m Migrator) DropConstraint(value interface{}, name string) error {
	return &ErrMissingImmuDBsupport{"DropConstraint"}
}

// DropIndex removes an index from a table.
//
// Not implemented as immudb does not support dropping indexes.
func (m Migrator) DropIndex(value interface{}, name string) error {
	return &ErrMissingImmuDBsupport{"DropIndex"}
}

// DropTable removes a table from a database.
//
// Not implemented as immudb does not support dropping tables.
func (m Migrator) DropTable(values ...interface{}) error {
	return &ErrMissingImmuDBsupport{"DropTable"}
}

// DropView removes a view from a database.
//
// Not implemented as immudb does not support views.
func (m Migrator) DropView(name string) error {
	return &ErrMissingImmuDBsupport{"DropView"}
}

// FullDataTypeOf returns field's db full data type
func (m Migrator) FullDataTypeOf(field *schema.Field) (expr clause.Expr) {
	expr.SQL = m.DataTypeOf(field)

	if field.NotNull {
		expr.SQL += " NOT NULL"
	}

	// In contrast to the default FullDataTypeOf implementation of the migrator,
	// this implementation ignores the unique and default setting, as both are
	// no supported by immudb.
	return
}

// GetIndexes returns all indexes for the table referenced by dst.
func (m Migrator) GetIndexes(dst interface{}) ([]gorm.Index, error) {
	var indexes []gorm.Index
	err := m.RunWithValue(dst, func(stmt *gorm.Statement) error {
		namer := m.DB.NamingStrategy
		// Retrieve the database connector.
		db, err := m.DB.DB()
		if err != nil {
			return err
		}
		// Query all indexes of the table.
		rows, err := db.Query("SELECT \"table\", name, \"unique\", \"primary\" FROM INDEXES( ? )", stmt.Table)
		if err != nil {
			return err
		}
		// Check if the result contains a row.
		for rows.Next() {
			// Get the name of the index.
			var table string
			var indexName string
			var unique bool
			var primary bool
			err = rows.Scan(&table, &indexName, &unique, &primary)
			if err != nil {
				return err
			}
			// Retrieve the columns string from the index name.
			column, err := extractColumnFromIndexName(indexName, table)
			if err != nil {
				return err
			}
			// Create a new index.
			indexes = append(indexes, ImmuDBindex{
				immudbName: indexName,
				gormName:   namer.IndexName(table, column),
				table:      table,
				column:     column,
				primary:    primary,
				unique:     unique,
			})
		}
		rows.Close()
		return nil
	})
	return indexes, err
}

// GetTables returns a list of all tables in the current database.
func (m Migrator) GetTables() (tableList []string, err error) {
	tables := []string{}
	// Retrieve the database connector.
	db, err := m.DB.DB()
	if err != nil {
		return tables, err
	}
	// Query all indexes of the table.
	rows, err := db.Query("SELECT name FROM TABLES()")
	if err != nil {
		return tables, err
	}
	// Check if the result contains a row.
	for rows.Next() {
		// Get the name of the index.
		var table string
		err = rows.Scan(&table)
		if err != nil {
			return tables, err
		}
		tables = append(tables, table)
	}
	rows.Close()
	return tables, nil
}

// HasColumn returns true if the table referenced by value contains a
// column with the given name.
func (m Migrator) HasColumn(value interface{}, field string) bool {
	columnExists := false
	err := m.RunWithValue(value, func(stmt *gorm.Statement) error {
		// Retrieve the database connector.
		db, err := m.DB.DB()
		if err != nil {
			return err
		}
		// Query all columns of the table.
		rows, err := db.Query("SELECT name FROM COLUMNS(?)", stmt.Table)
		if err != nil {
			return err
		}
		// Check if the result contains a row.
		for rows.Next() {
			// Get the name of the column.
			var name string
			err = rows.Scan(&name)
			if err != nil {
				return err
			}
			// Check if it matches the given column name.
			if field == name {
				columnExists = true
				return nil
			}
		}
		rows.Close()
		return nil
	})
	if err != nil {
		return false
	}
	return columnExists
}

// HasConstraints returns true if the table referenced by value contains a
// constraint with the given name.
//
// As immudb does not support constraints, this will always return false.
func (m Migrator) HasConstraint(value interface{}, name string) bool {
	return false
}

// HasTable determines if a specific tables exists in the current database.
func (m Migrator) HasTable(value interface{}) bool {
	tableExists := false
	err := m.RunWithValue(value, func(stmt *gorm.Statement) error {
		// Retrieve the database connector.
		db, err := m.DB.DB()
		if err != nil {
			return err
		}
		// Query for tables with this name.
		rows, err := db.Query("SELECT name FROM TABLES() WHERE name = ?", stmt.Table)
		if err != nil {
			return err
		}
		// Check if the result contains a row.
		tableExists = rows.Next()
		rows.Close()
		return nil
	})
	if err != nil {
		log.Printf("error checking if a table exists: %v", err)
		return false
	}
	return tableExists
}

// HasTable determines if an index with a specific name exists for a table.
func (m Migrator) HasIndex(value interface{}, name string) bool {
	indexExists := false
	err := m.RunWithValue(value, func(stmt *gorm.Statement) error {
		namer := m.DB.NamingStrategy
		// Retrieve the database connector.
		db, err := m.DB.DB()
		if err != nil {
			return err
		}
		// Query all indexes of the table.
		rows, err := db.Query("SELECT \"table\", name FROM INDEXES(?)", stmt.Table)
		if err != nil {
			return err
		}
		// Check if the result contains a row.
		for rows.Next() {
			// Get the name of the index.
			var table string
			var indexName string
			err = rows.Scan(&table, &indexName)
			if err != nil {
				return err
			}
			column, err := extractColumnFromIndexName(indexName, table)
			if err != nil {
				return err
			}
			// Check if either the name which immudb assigned to the index
			// matches the given name or if the immudb index name converted to
			// name which gorm would assign to this index matches the given name.
			if name == indexName || name == namer.IndexName(table, column) {
				indexExists = true
				return nil
			}
		}
		rows.Close()
		return nil
	})
	if err != nil {
		log.Printf("error checking if an index exists: %v", err)
		return false
	}
	return indexExists
}

// MigrateColumn does not have a custom implementation as the default one is
// compatible with immudb.
// func (m Migrator) MigrateColumn(value interface{}, field *schema.Field, columnType gorm.ColumnType) error

// RenameColumn does not have a custom implementation as the default one is
// compatible with immudb.
// func (m Migrator) RenameColumn(value interface{}, oldName, newName string) error

// RenameIndex alters the name of an index in the specified table.
//
// Not implemented as immudb does not support custom names for indexes.
func (m Migrator) RenameIndex(value interface{}, oldName, newName string) error {
	return &ErrMissingImmuDBsupport{"RenameIndex"}
}

// RenameTable alters the definition of a column in the specified table.
//
// Not implemented as immudb does not support changing the name of a table.
func (m Migrator) RenameTable(oldName, newName interface{}) error {
	return &ErrMissingImmuDBsupport{"RenameTable"}
}

type ImmuDBindex struct {
	immudbName string
	gormName   string
	table      string
	column     string
	primary    bool
	unique     bool
}

// Table returns table on which the index is defined.
func (ind ImmuDBindex) Table() string {
	return ind.table
}

// Table returns name gorm would assign to index, not the one immudb uses
// internally for the index.
func (ind ImmuDBindex) Name() string {
	return ind.gormName
}

// Columns returns the columns of the index.
func (ind ImmuDBindex) Columns() []string {
	// Split the column string into individual columns.
	return strings.Split(ind.column, ",")
}

// PrimaryKey return true is the index is for the primary key.
func (ind ImmuDBindex) PrimaryKey() (isPrimaryKey bool, ok bool) {
	return ind.primary, true
}

// Unique return true every entry for this key is unique.
func (ind ImmuDBindex) Unique() (unique bool, ok bool) {
	return ind.unique, true
}

// Option will always return an empty string, as immudb does not support options
// for indexes.
func (ind ImmuDBindex) Option() string {
	return ""
}

func extractColumnFromIndexName(name string, table string) (string, error) {
	// The format immudb uses for the index names is: table[columns] .
	// Remove the table name from the index name.
	columnsPart := strings.TrimPrefix(name, table)
	if columnsPart[0] != '[' || columnsPart[len(columnsPart)-1] != ']' {
		return "", errors.New("index name does not use the expected immudb format")
	}
	// Return the string between the brackets.
	return columnsPart[1 : len(columnsPart)-1], nil
}
