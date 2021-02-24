package jdatabase

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ijidan/jgo/jgo/jlogger"
)

//数据库连接
type Connection struct {
	host     string
	port     int64
	dbName   string
	username string
	password string
	connect  *sql.DB
}

//构造函数
func (c *Connection) construct(host string, port int64, dbName string, username string, password string) {
	c.host = host
	c.port = port
	c.dbName = dbName
	c.username = username
	c.password = password
}


//获取连接
func (c *Connection) GetConnect(host string, port int64, dbName string, username string, password string) (*sql.DB, error) {
	c.construct(host, port, dbName, username, password)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.username, c.password, c.host, port, c.dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	c.connect = db
	return c.connect, nil
}

//关闭连接
func CloseConnection(db *sql.DB) {
	if db != nil {
		jlogger.Info("db closed!")
		_ = db.Close()
	}
}
