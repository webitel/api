package db

import (
	. "../config"
	"../logger"
	"gopkg.in/mgo.v2"
)

type DB struct {
	connected bool
	session   *mgo.Session
	db        *mgo.Database
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

func NewDB(uri string) *DB {
	logger.Debug("Try connect mongodb to %s", Config.Get("mongodb:uri"))

	session, err := mgo.Dial(Config.Get("mongodb:uri"))
	if err != nil {
		logger.Error("Connect to mongo error: ", err)
		return NewDB(Config.Get("mongodb:uri"))
	}
	logger.Debug("Connect to mongo success")
	return &DB{
		session: session,
		db:      session.DB("webitel"),
	}
}
