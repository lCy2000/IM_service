package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func (c *MySQLClient) ExecQuery(query string, args ...interface{}) error {
	_, err := c.db.Exec(query, args...)
	return err
}

func (c *MySQLClient) SelectQuery(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := c.db.Query(query, args...)
	return rows, err
}


