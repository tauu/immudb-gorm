package immudbGorm

import (
	"database/sql"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
)

type Config struct {
	DriverName string
	DSN        string
}

type dialector struct {
	*Config
}

// Open creates a new dialector for connecting to the database at dsn.
func Open(dsn string) gorm.Dialector {
	return &dialector{Config: &Config{DSN: dsn, DriverName: "immudb"}}
}

// New creates a new dialector using the given configuration.
func New(config Config) gorm.Dialector {
	// Set the default driver name if the user has not set one.
	if config.DriverName == "" {
		config.DriverName = "immudb"
	}
	return &dialector{Config: &config}
}

// -- Dialector interface --

// Name is the name oif the sql dialect.
func (dialector dialector) Name() string {
	return "immudb"
}

// Initialize sets up the dialector for a database.
func (dialector dialector) Initialize(db *gorm.DB) (err error) {
	db.ConnPool, err = sql.Open(dialector.DriverName, dialector.DSN)
	if err != nil {
		return err
	}
	// Register default callbacks for insert and delete.
	// The default update callback is not useable,
	// as immudb uses the upsert clause instead of update.
	callbacks.RegisterDefaultCallbacks(db, &callbacks.Config{
		CreateClauses: []string{"INSERT", "VALUES"},
		UpdateClauses: []string{"UPDATE", "SET", "WHERE", "ORDER BY", "LIMIT"},
		DeleteClauses: []string{"DELETE", "FROM", "WHERE", "ORDER BY", "LIMIT"},
	})
	return nil
}

// Migrator creates a new migrator for the gorm database.
func (dialector dialector) Migrator(db *gorm.DB) gorm.Migrator {
	m := Migrator{
		migrator.Migrator{Config: migrator.Config{
			DB:                          db,
			Dialector:                   dialector,
			CreateIndexAfterCreateTable: true,
		}},
	}
	return m
}

// DataTypeOf creates a datatype definition matching the configuration of a field.
func (dialector dialector) DataTypeOf(field *schema.Field) string {
	dataType := ""
	switch field.DataType {
	case schema.Bool:
		dataType = "BOOLEAN"
	case schema.Int, schema.Uint:
		dataType = "INTEGER"
	case schema.Float:
		// Floats are not yet supported by immudb.
		// Therefore a blob is used to to store the binary representation.
		dataType = "BLOB[8]"
	case schema.String:
		dataType = "VARCHAR"
	case schema.Time:
		dataType = "TIMESTAMP"
	case schema.Bytes:
		dataType = "BLOB"
	default:
		dataType = string(field.DataType)
	}

	// Add a size constraint for the field if one is set.
	// Currentlty size constraints are only supported for BLOB and VARCHAR.
	if (dataType == "BLOB" || dataType == "VARCHAR") && field.Size > 0 {
		dataType = dataType + fmt.Sprintf("[%d]", field.Size)
	}

	// Set nullable constraint.
	if field.NotNull {
		dataType = dataType + " NOT NULL"
	}

	// Set auto increment for integer fields, if the field is also a primary key.
	if field.AutoIncrement && field.PrimaryKey && dataType == "INTEGER" {
		dataType = dataType + " AUTO_INCREMENT"
	}

	return dataType
}

// DefaultValueOf creates an sql expression to set a default value for a column.
// As immudb does not support default values at the moment,
// just an empty expression is returned for now.
func (dialector dialector) DefaultValueOf(*schema.Field) clause.Expression {
	return clause.Expr{}
}

// BindVarTo adds a placeholder for a variable in a SQL query.
func (dialector dialector) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v interface{}) {
	writer.WriteByte('?')
}

// QuoteTo quotes an identifier in a SQL query.
// As immudb does not support quoting identifiers at the moment,
// the dialector does not perform any quoting so far.
func (dialector dialector) QuoteTo(writer clause.Writer, str string) {
	writer.WriteString(str)
}

// Explain creates a string describing the SQL query.
func (dialector dialector) Explain(sql string, vars ...interface{}) string {
	return logger.ExplainSQL(sql, nil, `'`, vars...)
}
