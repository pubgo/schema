package schema

import (
	"database/sql"
	"fmt"
)

// TableNames
// returns a list of all table names in the current schema
// (not including system tables).
func TableNames(db *sql.DB) ([]string, error) {
	dt := fmt.Sprintf("%T", db.Driver())
	d, ok := driverDialect[dt]
	if !ok {
		return nil, fmt.Errorf("unknown db driver %s\n", dt)
	}

	d.setDb(db)
	return d.TableNames()
}

// ViewNames returns a list of all view names in the current schema
// (not including system views).
func ViewNames(db *sql.DB) ([]string, error) {
	dt := fmt.Sprintf("%T", db.Driver())
	d, ok := driverDialect[dt]
	if !ok {
		return nil, fmt.Errorf("unknown db driver %s\n", dt)
	}

	d.setDb(db)
	return d.ViewNames()
}

func TableColumns(db *sql.DB, name ...string) (map[string][]*sql.ColumnType, error) {
	dt := fmt.Sprintf("%T", db.Driver())
	d, ok := driverDialect[dt]
	if !ok {
		return nil, fmt.Errorf("unknown db driver %s\n", dt)
	}

	d.setDb(db)
	return d.TableColumns(name...)
}
func ViewColumns(db *sql.DB, name ...string) (map[string][]*sql.ColumnType, error) {
	dt := fmt.Sprintf("%T", db.Driver())
	d, ok := driverDialect[dt]
	if !ok {
		return nil, fmt.Errorf("unknown db driver %s\n", dt)
	}

	d.setDb(db)
	return d.ViewColumns(name...)
}
func TableIndexs(db *sql.DB, name ...string) (map[string][]*sql.ColumnType, error) {
	dt := fmt.Sprintf("%T", db.Driver())
	d, ok := driverDialect[dt]
	if !ok {
		return nil, fmt.Errorf("unknown db driver %s\n", dt)
	}

	d.setDb(db)
	return d.TableIndexs(name...)
}
