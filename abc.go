package schema

import (
	"reflect"
)

// Dialect
// interface contains behaviors that differ across SQL database
type Dialect interface {
	TableNames() string
	ViewNames() string
	TableColumns(name string) string
	ViewColumns(name string) string
	TableIndexs(name string) string
}

type ColumnType struct {
	Name             string
	Length           int64
	DatabaseTypeName string
	DecimalSize      int64
	Nullable         bool
	ScanType         reflect.Type
}
