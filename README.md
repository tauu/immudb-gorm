# gorm driver for immudb

This gorm driver uses immusql driver for database/sql to connect to a immudb database instance.

## Usage

The following example connects to an immudb instance running on localhost with port 3322 and selects the database test.

```golang
package main

import (
    "log"

	immudbGorm "github.com/tauu/immudb-gorm"
	"gorm.io/gorm"
)

func main() {
    dsn := "immudb://immudb:immudb@localhost:3322/test"
	db, err := gorm.Open(immudbGorm.Open(dsn), &gorm.Config{})
	if err != nil {
        log.Fatal("opening database connection failed")
	}
}
```

## Features

### dialector interface
- [x] Name
- [x] Initialize
- [x] Migrator
- [x] DataTypeOf
      schema.Float is only supported by a workaround using a BLOB[8] column type, due to a lack of float support in immudb.
- [ ] DefaultValueOf
      Default values for columns are currently not supported by immudb.
- [x] BindVarTo
- [x] QuoteTo
- [x] Explain

### migrator interface
At the moment the migrator is able to create tables for a database schema, but it cannot modify existing tables. Functions marked with *, cannot be supported due to limitations of immudb.

- [-] AutoMigrate
  Partial support is implemented, creating table works but altering does not.
- [ ] CurrentDatabase
- [ ] FullDataTypeOf
- [x] CreateTable
- [ ] DropTable*
- [x] HasTable
- [ ] RenameTable
- [ ] AddColumn
- [ ] DropColumn*
- [ ] AlterColumn*
- [ ] HasColumn
- [ ] RenameColumn*
- [ ] MigrateColumn
- [ ] ColumnTypes
- [ ] CreateConstraint*
- [ ] DropConstraint*
- [ ] HasConstraint*
- [x] CreateIndex
- [ ] DropIndex*
- [ ] HasIndex
- [ ] RenameIndex*
