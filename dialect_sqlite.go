package schema

import (
	"database/sql"
	"fmt"
)

var sqlite = dialect{
	queries: [3]string{
		// columnTypes query.
		"SELECT * FROM `%s` LIMIT 0",
		// tableNames query.
		pack(`
			SELECT name
			FROM
				sqlite_master
			WHERE
				type = 'table'
		`),
		// viewNames query.
		pack(`
			SELECT name
			FROM
				sqlite_master
			WHERE
				type = 'view'
		`),
	},
}

var _ Dialect = (*sqlite3)(nil)

type sqlite3 struct {
}

func (s sqlite3) TableNames(db *sql.DB) ([]string, error) {
	panic("implement me")
}

func (s sqlite3) ViewNames(db *sql.DB) ([]string, error) {
	panic("implement me")
}

func (s sqlite3) TableColumns(db *sql.DB, name ...string) (map[string][]*sql.ColumnType, error) {
	panic("implement me")
}

func (s sqlite3) ViewColumns(db *sql.DB, name ...string) (map[string][]*sql.ColumnType, error) {
	panic("implement me")
}

func (s sqlite3) TableIndexs(db *sql.DB, name ...string) (map[string][]*sql.ColumnType, error) {
	panic("implement me")
}

func (sqlite3) GetName() string {
	return "sqlite3"
}

func (s sqlite3) HasIndex(tableName string, indexName string) bool {
	var count int
	s.db.QueryRow(fmt.Sprintf("SELECT count(*) FROM sqlite_master WHERE tbl_name = ? AND sql LIKE '%%INDEX %v ON%%'", indexName), tableName).Scan(&count)
	return count > 0
}

func (s sqlite3) HasTable(tableName string) bool {
	var count int
	s.db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?", tableName).Scan(&count)
	return count > 0
}

func (s sqlite3) HasColumn(tableName string, columnName string) bool {
	var count int
	s.db.QueryRow(fmt.Sprintf("SELECT count(*) FROM sqlite_master WHERE tbl_name = ? AND (sql LIKE '%%\"%v\" %%' OR sql LIKE '%%%v %%');\n", columnName, columnName), tableName).Scan(&count)
	return count > 0
}

func (s sqlite3) CurrentDatabase() (name string) {
	var (
		ifaces   = make([]interface{}, 3)
		pointers = make([]*string, 3)
		i        int
	)
	for i = 0; i < 3; i++ {
		ifaces[i] = &pointers[i]
	}
	if err := s.db.QueryRow("PRAGMA database_list").Scan(ifaces...); err != nil {
		return
	}
	if pointers[1] != nil {
		name = *pointers[1]
	}
	return
}
