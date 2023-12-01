package immudbGorm

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
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
		LastInsertIDReversed: true,
		CreateClauses:        []string{"INSERT", "VALUES"},
		UpdateClauses:        []string{"UPDATE", "SET", "WHERE", "ORDER BY", "LIMIT"},
		DeleteClauses:        []string{"DELETE", "FROM", "WHERE", "ORDER BY", "LIMIT"},
		QueryClauses:         []string{"SELECT", "FROM", "WHERE", "GROUP BY", "ORDER BY", "LIMIT"},
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
		dataType = "FLOAT"
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
	// Add string to UUID conversion, as a workaround for a bug in immudb.
	// ImmuDB expects raw UUIDs to be set as parameters, but does not have a
	// method to transfer raw UUID values. Only string encoded uuids can be
	// transferred, but these are not parsed to uuids currently.
	_, isUUID := v.(uuid.UUID)
	// If the type is not directly a uuid, check if it implements the valuer
	// interface and serialized itself to a uuid string.
	if !isUUID {
		valuer, ok := v.(driver.Valuer)
		if ok {
			value, err := valuer.Value()
			if err == nil {
				if str, ok := value.(string); ok {
					_, err = uuid.Parse(str)
					if err == nil {
						isUUID = true
					}
				}
			}
		}
	}
	// // If the value is not directly a UUID, check if it embeds a UUID value.
	// if !isUUID {
	// 	t := reflect.TypeOf(v)
	// 	for i := 0; i < t.NumField(); i++ {
	// 		if t.Field(i).Type.String() == "uuid.UUID" {
	// 			isUUID = true
	// 			break
	// 		}
	// 	}
	// }
	// Append an explicit cast to a uuid, if the value was determined to be a
	// a uuid value.
	if isUUID {
		writer.WriteString("::UUID")
	}
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
