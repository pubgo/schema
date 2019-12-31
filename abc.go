package schema

// Dialect
// interface contains behaviors that differ across SQL database
type Dialect interface {
	TableNames() string
	ViewNames() string
	TableColumns(name string) string
	ViewColumns(name string) string
	TableIndexs(name string) string
}
