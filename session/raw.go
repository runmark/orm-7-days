package session

import (
	"database/sql"
	"example.com/mark/geeorm/clause"
	"example.com/mark/geeorm/dialect"
	"example.com/mark/geeorm/log"
	"example.com/mark/geeorm/schema"
	"strings"
)

type Session struct {
	db       *sql.DB
	tx       *sql.Tx
	dialect  dialect.Dialect
	refTable *schema.Schema
	clause   clause.Clause
	sql      strings.Builder
	sqlVars  []interface{}
}

type CommonDB interface {
	Exec(query string, args ...interface{}) (result sql.Result, err error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (rows *sql.Rows, err error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{db: db, dialect: dialect}
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	result, err = s.DB().Exec(s.sql.String(), s.sqlVars...)
	if err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	rows, err = s.DB().Query(s.sql.String(), s.sqlVars...)
	if err != nil {
		log.Error(err)
	}
	return
}

