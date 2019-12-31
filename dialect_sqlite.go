package schema

import (
	"database/sql"
	"fmt"
)

var _ Dialect = (*sqlite3)(nil)

type sqlite3 struct {
	db *sql.DB
}

func (s *sqlite3) setDb(db *sql.DB) {
	s.db = db
}

func (s *sqlite3) TableNames() string {
	return `SELECT name FROM sqlite_master WHERE type = 'table'`
}

func (s *sqlite3) ViewNames() string {
	return `SELECT name FROM sqlite_master WHERE type = 'view'`
}

func (s *sqlite3) TableColumns(name string) string {
	return fmt.Sprintf("SELECT * FROM `%s` LIMIT 0", name)
	//return "SELECT * FROM sqlite_master WHERE tbl_name = ? and type = 'table'"
}

func (s *sqlite3) ViewColumns(name string) string {
	return fmt.Sprintf("SELECT * FROM `%s` LIMIT 0", name)
	//return "SELECT * FROM sqlite_master WHERE tbl_name = ? and type = 'view'"
}

func (s *sqlite3) TableIndexs(name string) string {
	return fmt.Sprintf("SELECT name FROM sqlite_master WHERE tbl_name = `%s` and type = 'index'", name)
}
