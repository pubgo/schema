package schema

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"regexp"
	"unicode/utf8"
)

var mysql1 = dialect{
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
				table_schema = database()
		`),
		// viewNames query.
		pack(`
			SELECT table_name
			FROM
				information_schema.tables
			WHERE
				table_type = 'VIEW' AND
				table_schema = database()
		`),
	},
}

var _ Dialect = (*mysql)(nil)

var mysqlIndexRegex = regexp.MustCompile(`^(.+)\((\d+)\)$`)

type mysql struct {
}

func (s mysql) TableNames(db *sql.DB) ([]string, error) {
	panic("implement me")
}

func (s mysql) ViewNames(db *sql.DB) ([]string, error) {
	panic("implement me")
}

func (s mysql) TableColumns(db *sql.DB, name ...string) (map[string][]*sql.ColumnType, error) {
	panic("implement me")
}

func (s mysql) ViewColumns(db *sql.DB, name ...string) (map[string][]*sql.ColumnType, error) {
	panic("implement me")
}

func (s mysql) TableIndexs(db *sql.DB, name ...string) (map[string][]*sql.ColumnType, error) {
	panic("implement me")
}

func (mysql) GetName() string {
	return "mysql"
}

func (mysql) Quote(key string) string {
	return fmt.Sprintf("`%s`", key)
}

func (s mysql) HasForeignKey(tableName string, foreignKeyName string) bool {
	var count int
	currentDatabase, tableName := currentDatabaseAndTable(&s, tableName)
	s.db.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS WHERE CONSTRAINT_SCHEMA=? AND TABLE_NAME=? AND CONSTRAINT_NAME=? AND CONSTRAINT_TYPE='FOREIGN KEY'", currentDatabase, tableName, foreignKeyName).Scan(&count)
	return count > 0
}

func (s mysql) CurrentDatabase() (name string) {
	s.db.QueryRow("SELECT DATABASE()").Scan(&name)
	return
}

func (mysql) SelectFromDummyTable() string {
	return "FROM DUAL"
}

func (s mysql) BuildKeyName(kind, tableName string, fields ...string) string {
	keyName := s.commonDialect.BuildKeyName(kind, tableName, fields...)
	if utf8.RuneCountInString(keyName) <= 64 {
		return keyName
	}
	h := sha1.New()
	h.Write([]byte(keyName))
	bs := h.Sum(nil)

	// sha1 is 40 characters, keep first 24 characters of destination
	destRunes := []rune(keyNameRegex.ReplaceAllString(fields[0], "_"))
	if len(destRunes) > 24 {
		destRunes = destRunes[:24]
	}

	return fmt.Sprintf("%s%x", string(destRunes), bs)
}

func (s commonDialect) HasIndex(tableName string, indexName string) bool {
	var count int
	currentDatabase, tableName := currentDatabaseAndTable(&s, tableName)
	s.db.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.STATISTICS WHERE table_schema = ? AND table_name = ? AND index_name = ?", currentDatabase, tableName, indexName).Scan(&count)
	return count > 0
}

// NormalizeIndexAndColumn returns index name and column name for specify an index prefix length if needed
func (mysql) NormalizeIndexAndColumn(indexName, columnName string) (string, string) {
	submatch := mysqlIndexRegex.FindStringSubmatch(indexName)
	if len(submatch) != 3 {
		return indexName, columnName
	}
	indexName = submatch[1]
	columnName = fmt.Sprintf("%s(%s)", columnName, submatch[2])
	return indexName, columnName
}

func (mysql) DefaultValueStr() string {
	return "VALUES()"
}
