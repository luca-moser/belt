package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"

	mgo "gopkg.in/mgo.v2"
)

type (
	SQLConfig struct {
		Name      string
		Host      string
		User      string
		PW        string
		Port      int
		KeepAlive bool
		Pool      PoolSettings
	}

	PoolSettings struct {
		MaxOpen     int
		MaxIdle     int
		MaxLifetime int
		DialTimeout int
		Timeout     int
	}

	MongoDBConfig struct {
		Host      string
		Auth      bool
		Username  string
		Password  string
		Mechanism string
		Source    string
	}
)

func GetMSSQLConnection(config SQLConfig) (*sql.DB, error) {
	var keepAlive int
	if config.KeepAlive {
		keepAlive = 1
	}
	connectionString := fmt.Sprintf(
		"server=%s;user id=%s;password=%s;port=%d;database=%s;connection timeout=%d;dial timeout=%d;keepAlive=%d;log=1",
		config.Host, config.User, config.PW, config.Port, config.Name,
		config.Pool.Timeout, config.Pool.DialTimeout, keepAlive,
	)

	conn, err := sql.Open("mssql", connectionString)
	if err != nil {
		return nil, err
	}

	conn.SetConnMaxLifetime(time.Duration(config.Pool.MaxLifetime) * time.Second)
	conn.SetMaxIdleConns(config.Pool.MaxIdle)
	conn.SetMaxOpenConns(config.Pool.MaxOpen)
	return conn, nil
}

func GetMySQLConnection(config SQLConfig) (*sql.DB, error) {
	mysqlConfig := mysql.Config{
		User: config.User, Passwd: config.PW,
		Net: "tcp", Addr: fmt.Sprintf("%s:%d", config.Host, config.Port), DBName: config.Name,
	}
	conn, err := sql.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		return nil, err
	}
	conn.SetConnMaxLifetime(time.Duration(config.Pool.MaxLifetime) * time.Second)
	conn.SetMaxIdleConns(config.Pool.MaxIdle)
	conn.SetMaxOpenConns(config.Pool.MaxOpen)
	return conn, nil
}

func GetMongoDBConnection(config MongoDBConfig) (*mgo.Session, error) {
	var session *mgo.Session
	var err error
	if config.Auth {
		cred := &mgo.Credential{
			Username:  config.Username,
			Password:  config.Password,
			Mechanism: config.Mechanism,
			Source:    config.Source,
		}
		session, err = mgo.Dial(config.Host)
		if err = session.Login(cred); err != nil {
			return nil, err
		}
	} else {
		session, err = mgo.Dial(config.Host)
	}
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	session.SetSafe(&mgo.Safe{})
	return session, nil
}
