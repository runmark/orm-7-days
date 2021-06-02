package geeorm

import (
	"database/sql"
	"example.com/mark/geeorm/log"
	"example.com/mark/geeorm/session"
)

type Engine struct {
	db *sql.DB
}

func NewEngine(driver, source string) (e *Engine, err error){
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}

	err = db.Ping()
	if err != nil {
		log.Error(err)
		return
	}

	e = &Engine{db}

	log.Info("Connect db success!")

	return
}

func (e *Engine) Close() {
	err := e.db.Close()
	if err != nil {
		log.Error(err)
	}
	log.Info("Close db success")
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db)
}