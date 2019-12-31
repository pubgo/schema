package schema

import (
	"database/sql"
	"fmt"
	"github.com/pubgo/g/xerror"
)

func _names(rows *sql.Rows) (names []string) {
	n := ""
	for rows.Next() {
		xerror.PanicM(rows.Scan(&n), "rows scan error")
		names = append(names, n)
	}
	return
}

func _dialect(db *sql.DB) Dialect {
	dt := fmt.Sprintf("%T", db.Driver())
	d, ok := driverDialect[dt]
	xerror.PanicT(!ok, "unknown db driver %s\n", dt)
	return d
}

func TableNames(db *sql.DB) (_ []string, err error) {
	defer xerror.RespErr(&err)

	d := _dialect(db)
	return _names(xerror.PanicErr(db.Query(d.TableNames())).(*sql.Rows)), nil
}

func ViewNames(db *sql.DB) (_ []string, err error) {
	defer xerror.RespErr(&err)

	d := _dialect(db)
	return _names(xerror.PanicErr(db.Query(d.TableNames())).(*sql.Rows)), nil
}

func TableColumns(db *sql.DB, name string) (_ []*sql.ColumnType, err error) {
	defer xerror.RespErr(&err)

	d := _dialect(db)
	rows := xerror.PanicErr(db.Query(d.TableColumns(name))).(*sql.Rows)
	return rows.ColumnTypes()
}
func ViewColumns(db *sql.DB, name string) (_ []*sql.ColumnType, err error) {
	defer xerror.RespErr(&err)

	d := _dialect(db)
	rows := xerror.PanicErr(db.Query(d.ViewColumns(name))).(*sql.Rows)
	return rows.ColumnTypes()
}
func TableIndexs(db *sql.DB, name string) (_ []*sql.ColumnType, err error) {
	defer xerror.RespErr(&err)

	d := _dialect(db)
	rows := xerror.PanicErr(db.Query(d.TableIndexs(name))).(*sql.Rows)
	return rows.ColumnTypes()
}
