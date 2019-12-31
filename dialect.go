package schema

import (
	"strings"
)

// query defines dialect query types.
type query int

// query type enum.
const (
	columnTypes query = iota // Index of query to get column type info.
	tableNames               // Index of query to get table names.
	viewNames                // Index of query to get view names.
)

// dialect describes how each database 'flavour' provides its metadata.
type dialect struct {
	// queries for fetching metadata: tableNames, viewNames, columnTypes.
	queries [3]string
}

// driverDialect is a registry, mapping database/sql driver names to database dialects.
// This is somewhat fragile.
var driverDialect map[string]*dialect = map[string]*dialect{
	"*sqlite3.SQLiteDriver":       &sqlite,   // github.com/mattn/go-sqlite3
	"*sqlite.impl":                &sqlite,   // github.com/gwenn/gosqlite
	"sqlite3.Driver":              &sqlite,   // github.com/mxk/go-sqlite - TODO(js) No datatypes.
	"*pq.Driver":                  &postgres, // github.com/lib/pq
	"*stdlib.Driver":              &postgres, // github.com/jackc/pgx
	"*pgsqldriver.postgresDriver": &postgres, // github.com/jbarham/gopgsqldriver - TODO(js) No datatypes.
	"*mysql.MySQLDriver":          &mysql,    // github.com/go-sql-driver/mysql
	"*godrv.Driver":               &mysql,    // github.com/ziutek/mymysql - TODO(js) No datatypes.
	"*mssql.MssqlDriver":          &mssql,    // github.com/denisenkom/go-mssqldb
	"*freetds.MssqlDriver":        &mssql,    // github.com/minus5/gofreetds - TODO(js) No datatypes. Error on create view.
	"*goracle.drv":                &oracle,   // gopkg.in/goracle.v2
	"*ora.Drv":                    &oracle,   // gopkg.in/rana/ora.v4 - TODO(js) Mismatched datatypes.
	"*oci8.OCI8Driver":            &oracle,   // github.com/mattn/go-oci8 - TODO(js) Mismatched datatypes.
}

// TODO Should we expose a method of registering a driver string/dialect in our registry?
// -- It would allow folk to work around the fragility. e.g.
//
// func Register(driver sql.Driver, d *Dialect) {}
//

// pack a string, normalising its whitespace.
func pack(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// Dialect interface contains behaviors that differ across SQL database
type Dialect interface {
	// GetName get dialect's name
	GetName() string

	// SetDB set db for dialect
	//SetDB(db SQLCommon)

	// BindVar return the placeholder for actual values in SQL statements, in many dbs it is "?", Postgres using $1
	BindVar(i int) string
	// Quote quotes field name to avoid SQL parsing exceptions by using a reserved word as a field name
	Quote(key string) string
	// DataTypeOf return data's sql type
	//DataTypeOf(field *StructField) string

	// HasIndex check has index or not
	HasIndex(tableName string, indexName string) bool
	// HasForeignKey check has foreign key or not
	HasForeignKey(tableName string, foreignKeyName string) bool
	// RemoveIndex remove index
	RemoveIndex(tableName string, indexName string) error
	// HasTable check has table or not
	HasTable(tableName string) bool
	// HasColumn check has column or not
	HasColumn(tableName string, columnName string) bool
	// ModifyColumn modify column's type
	ModifyColumn(tableName string, columnName string, typ string) error

	// LimitAndOffsetSQL return generated SQL with Limit and Offset, as mssql has special case
	LimitAndOffsetSQL(limit, offset interface{}) string
	// SelectFromDummyTable return select values, for most dbs, `SELECT values` just works, mysql needs `SELECT value FROM DUAL`
	SelectFromDummyTable() string
	// LastInsertIdReturningSuffix most dbs support LastInsertId, but postgres needs to use `RETURNING`
	LastInsertIDReturningSuffix(tableName, columnName string) string
	// DefaultValueStr
	DefaultValueStr() string

	// BuildKeyName returns a valid key name (foreign key, index key) for the given table, field and reference
	BuildKeyName(kind, tableName string, fields ...string) string

	// NormalizeIndexAndColumn returns valid index name and column name depending on each dialect
	NormalizeIndexAndColumn(indexName, columnName string) (string, string)

	// CurrentDatabase return current database name
	CurrentDatabase() string
}

// ParseFieldStructForDialect get field's sql data type
//var ParseFieldStructForDialect = func(field *StructField, dialect Dialect) (fieldValue reflect.Value, sqlType string, size int, additionalType string) {
//	// Get redirected field type
//	var (
//		reflectType = field.Struct.Type
//		dataType, _ = field.TagSettingsGet("TYPE")
//	)
//
//	for reflectType.Kind() == reflect.Ptr {
//		reflectType = reflectType.Elem()
//	}
//
//	// Get redirected field value
//	fieldValue = reflect.Indirect(reflect.New(reflectType))
//
//	if gormDataType, ok := fieldValue.Interface().(interface {
//		GormDataType(Dialect) string
//	}); ok {
//		dataType = gormDataType.GormDataType(dialect)
//	}
//
//	// Get scanner's real value
//	if dataType == "" {
//		var getScannerValue func(reflect.Value)
//		getScannerValue = func(value reflect.Value) {
//			fieldValue = value
//			if _, isScanner := reflect.New(fieldValue.Type()).Interface().(sql.Scanner); isScanner && fieldValue.Kind() == reflect.Struct {
//				getScannerValue(fieldValue.Field(0))
//			}
//		}
//		getScannerValue(fieldValue)
//	}
//
//	// Default Size
//	if num, ok := field.TagSettingsGet("SIZE"); ok {
//		size, _ = strconv.Atoi(num)
//	} else {
//		size = 255
//	}
//
//	// Default type from tag setting
//	notNull, _ := field.TagSettingsGet("NOT NULL")
//	unique, _ := field.TagSettingsGet("UNIQUE")
//	additionalType = notNull + " " + unique
//	if value, ok := field.TagSettingsGet("DEFAULT"); ok {
//		additionalType = additionalType + " DEFAULT " + value
//	}
//
//	if value, ok := field.TagSettingsGet("COMMENT"); ok {
//		additionalType = additionalType + " COMMENT " + value
//	}
//
//	return fieldValue, dataType, size, strings.TrimSpace(additionalType)
//}

func currentDatabaseAndTable(dialect Dialect, tableName string) (string, string) {
	if strings.Contains(tableName, ".") {
		splitStrings := strings.SplitN(tableName, ".", 2)
		return splitStrings[0], splitStrings[1]
	}
	return dialect.CurrentDatabase(), tableName
}
