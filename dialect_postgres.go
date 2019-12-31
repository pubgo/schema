package schema

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

var postgres1 = dialect{
	queries: [3]string{
		// columnTypes query.
		"SELECT * FROM `%s` LIMIT 0",
		// tableNames query.
		pack(`
			SELECT table_name
			FROM
				information_schema.tables
			WHERE
				table_type = 'BASE TABLE' AND
				table_schema = current_schema()
		`),
		// viewNames query.
		pack(`
			SELECT table_name
			FROM
				information_schema.tables
			WHERE
				table_type = 'VIEW' AND
				table_schema = current_schema()
		`),
	},
}

var _ Dialect = (*postgres)(nil)

type postgres struct {
}

func (s postgres) TableNames(db *sql.DB) ([]string, error) {
	panic("implement me")
}

func (s postgres) ViewNames(db *sql.DB) ([]string, error) {
	panic("implement me")
}

func (s postgres) TableColumns(db *sql.DB, name ...string) (map[string][]*sql.ColumnType, error) {
	panic("implement me")
}

func (s postgres) ViewColumns(db *sql.DB, name ...string) (map[string][]*sql.ColumnType, error) {
	panic("implement me")
}

func (s postgres) TableIndexs(db *sql.DB, name ...string) (map[string][]*sql.ColumnType, error) {
	panic("implement me")
}

func (postgres) GetName() string {
	return "postgres"
}

func (postgres) BindVar(i int) string {
	return fmt.Sprintf("$%v", i)
}

func (s postgres) HasIndex(tableName string, indexName string) bool {
	var count int
	s.db.QueryRow("SELECT count(*) FROM pg_indexes WHERE tablename = $1 AND indexname = $2 AND schemaname = CURRENT_SCHEMA()", tableName, indexName).Scan(&count)
	return count > 0
}

func (s postgres) HasForeignKey(tableName string, foreignKeyName string) bool {
	var count int
	s.db.QueryRow("SELECT count(con.conname) FROM pg_constraint con WHERE $1::regclass::oid = con.conrelid AND con.conname = $2 AND con.contype='f'", tableName, foreignKeyName).Scan(&count)
	return count > 0
}

func (s postgres) HasTable(tableName string) bool {
	var count int
	s.db.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.tables WHERE table_name = $1 AND table_type = 'BASE TABLE' AND table_schema = CURRENT_SCHEMA()", tableName).Scan(&count)
	return count > 0
}

func (s postgres) HasColumn(tableName string, columnName string) bool {
	var count int
	s.db.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.columns WHERE table_name = $1 AND column_name = $2 AND table_schema = CURRENT_SCHEMA()", tableName, columnName).Scan(&count)
	return count > 0
}

func (s postgres) CurrentDatabase() (name string) {
	s.db.QueryRow("SELECT CURRENT_DATABASE()").Scan(&name)
	return
}

func (s postgres) LastInsertIDReturningSuffix(tableName, key string) string {
	return fmt.Sprintf("RETURNING %v.%v", tableName, key)
}

func (postgres) SupportLastInsertID() bool {
	return false
}

func isUUID(value reflect.Value) bool {
	if value.Kind() != reflect.Array || value.Type().Len() != 16 {
		return false
	}
	typename := value.Type().Name()
	lower := strings.ToLower(typename)
	return "uuid" == lower || "guid" == lower
}

func isJSON(value reflect.Value) bool {
	_, ok := value.Interface().(json.RawMessage)
	return ok
}
