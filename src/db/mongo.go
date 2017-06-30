package db

import (
	. "../config"
	"../logger"
	"database/sql"
	"fmt"
	"gopkg.in/mgo.v2"
	"strconv"
)

type DB struct {
	connected bool
	session   *mgo.Session
	db        *mgo.Database
	pg        *sql.DB
}

func (db *DB) onError(err error) {
	if err != nil && err.Error() == "EOF" {
		if db.session != nil {
			db.session.Refresh()
		}
	} else {
		logger.Warning("On error: no error!")
	}
}

func (db *DB) GetPg() *sql.DB {
	return db.pg
}

func NewDB(uri string) *DB {
	logger.Debug("Try connect mongodb to %s", Config.Get("mongodb:uri"))

	session, err := mgo.Dial(Config.Get("mongodb:uri"))
	if err != nil {
		logger.Error("Connect to mongo error: ", err)
		return NewDB(Config.Get("mongodb:uri"))
	}
	logger.Debug("Connect to mongo success")

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s",
		Config.Get("pg:user"), Config.Get("pg:password"), Config.Get("pg:dbName"), Config.Get("pg:host"))

	pg, err := sql.Open("postgres", dbinfo)
	maxConn, err := strconv.Atoi(Config.Get("pg:maxConnection"))
	if err != nil {
		panic(err)
	}

	pg.SetMaxOpenConns(maxConn)
	if err != nil {
		panic(err)
	} else {
		logger.Info("Connect to pg success")
	}
	return &DB{
		session: session,
		db:      session.DB("webitel"),
		pg:      pg,
	}
}
