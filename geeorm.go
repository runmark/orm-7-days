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
		log.Errorf("dialect %s doesnt exists!", driver)
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

// 将所有的操作放到一个回调函数中，作为入参传递给 engine.Transaction()，发生任何错误，自动回滚，如果没有错误发生，则提交。
type TxFunc func(s *session.Session) (interface{}, error)

func (e *Engine) Transaction(txFunc TxFunc) (result interface{}, err error) {

	s := e.NewSession()

	err = s.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {

		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p)  // re-throw panic after Rollback
		} else if err != nil  {
			_ = s.Rollback()  // err is non-nil; don't change it
		} else {
			defer func() {
				if err != nil {
					_ = s.Rollback()
				}
			}()
			err = s.Commit()  // err is nil; if Commit returns error update err
		}

	}()

	return txFunc(s)
}