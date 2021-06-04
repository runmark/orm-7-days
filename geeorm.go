package geeorm

import (
	"database/sql"
	"example.com/mark/geeorm/dialect"
	"example.com/mark/geeorm/log"
	"example.com/mark/geeorm/session"
)

type Engine struct {
	db *sql.DB
	dialect dialect.Dialect
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

	dialect, ok := dialect.GetDialect(driver)
	if !ok {
		log.Error("dialect %v doesnt exists!", driver)
		return
	}

	e = &Engine{db,dialect}

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
	return session.New(e.db, e.dialect)
}