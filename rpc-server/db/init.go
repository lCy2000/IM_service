package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"context"
	"time"
)

type MySQLClient struct {
	db *sql.DB
}

func (c *MySQLClient) InitClient(ctx context.Context, dataSourceName string) error {
	time.Sleep(10 * time.Second)

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}

	c.db = db

	testdbErr := c.TestConnection(ctx)
	if (testdbErr != nil){
		return testdbErr
	}
	
	return nil
}

func (c* MySQLClient) TestConnection(ctx context.Context) error {
	if err := c.db.Ping(); err != nil {
		return err
	}
	return nil
}

