package schema

import (
	"fmt"
)

var _ Dialect = (*postgres)(nil)

type postgres struct {
}

func (s postgres) TableNames() string {
	return `SELECT table_name FROM information_schema.tables WHERE table_type = 'BASE TABLE' AND table_schema = current_schema()`
}

func (s postgres) ViewNames() string {
	return `SELECT table_name FROM information_schema.tables WHERE table_type = 'VIEW' AND table_schema = current_schema()`
}

func (s postgres) TableColumns(name string) string {
	return fmt.Sprintf("SELECT * FROM `%s` LIMIT 0", name)
}

func (s postgres) ViewColumns(name string) string {
	return fmt.Sprintf("SELECT * FROM `%s` LIMIT 0", name)
}

func (s postgres) TableIndexs(name string) string {
	return fmt.Sprintf("SELECT indexname FROM pg_indexes WHERE tablename = `%s` AND schemaname = CURRENT_SCHEMA()", name)
}
