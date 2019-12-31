package schema

// driverDialect,
// mapping database/sql driver names to database dialects.
// This is somewhat fragile.
var driverDialect = map[string]Dialect{
	"*sqlite3.SQLiteDriver":       &sqlite3{},  // github.com/mattn/go-sqlite3
	"*sqlite.impl":                &sqlite3{},  // github.com/gwenn/gosqlite
	"sqlite3.Driver":              &sqlite3{},  // github.com/mxk/go-sqlite - TODO(js) No datatypes.
	"*pq.Driver":                  &postgres{}, // github.com/lib/pq
	"*stdlib.Driver":              &postgres{}, // github.com/jackc/pgx
	"*pgsqldriver.postgresDriver": &postgres{}, // github.com/jbarham/gopgsqldriver - TODO(js) No datatypes.
	"*mysql.MySQLDriver":          &mysql{},    // github.com/go-sql-driver/mysql
	"*godrv.Driver":               &mysql{},    // github.com/ziutek/mymysql - TODO(js) No datatypes.
}
