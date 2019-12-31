package schema

import (
	"fmt"
)

var _ Dialect = (*mysql)(nil)

type mysql struct {
}

func (s *mysql) TableNames() string {
	return `SELECT table_name FROM information_schema.tables WHERE table_type = 'BASE TABLE' AND table_schema = database()`
}

func (s *mysql) ViewNames() string {
	return `SELECT table_name FROM information_schema.tables WHERE table_type = 'VIEW' AND table_schema = database()`
}

func (s *mysql) TableColumns(name string) string {
	return fmt.Sprintf("SELECT * FROM `%s` LIMIT 0", name)
}

func (s *mysql) ViewColumns(name string) string {
	return fmt.Sprintf("SELECT * FROM `%s` LIMIT 0", name)
}

func (s *mysql) TableIndexs(name string) string {
	return fmt.Sprintf("SHOW INDEX FROM `%s`", name)
}
