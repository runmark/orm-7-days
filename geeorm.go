package geeorm

import (
	"database/sql"
	"example.com/mark/geeorm/dialect"
	"example.com/mark/geeorm/log"
	"example.com/mark/geeorm/session"
	"fmt"
	"strings"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver, source string) (e *Engine, err error) {
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

	e = &Engine{db, dialect}

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
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = s.Rollback() // err is non-nil; don't change it
		} else {
			defer func() {
				if err != nil {
					_ = s.Rollback()
				}
			}()
			err = s.Commit() // err is nil; if Commit returns error update err
		}

	}()

	return txFunc(s)
}

// a - b
func difference(as []string, bs []string) (diff []string) {

	bm := map[string]struct{}{}

	for _, b := range bs {
		bm[b] = struct{}{}
	}

	for _, a := range as {
		if _, ok := bm[a]; !ok {
			diff = append(diff, a)
		}
	}

	return
}

func (e *Engine) Migrate(value interface{}) (err error) {

	_, err = e.Transaction(func(s *session.Session) (result interface{}, err error) {

		if !s.Model(value).HasTable() {
			log.Infof("table %s doesnt exists!", s.RefTable().Name)
			return nil, s.CreateTable()
		}

		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1;", s.RefTable().Name)).QueryRows()
		oldColumns, _ := rows.Columns()

		addedFields := difference(s.RefTable().FieldNames, oldColumns)
		deletedFields := difference(oldColumns, s.RefTable().FieldNames)

		log.Infof("added cols %v, deleted cols %v", addedFields, deletedFields)

		for _, field := range addedFields {
			f := s.RefTable().GetField(field)
			_, err = s.Raw(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", s.RefTable().Name, field, f.Type)).Exec()
			if err != nil {
				return
			}
		}

		if len(deletedFields) == 0 {
			return
		}

		table := s.RefTable()
		tmpTableName := table.Name + "_tmp"

		fieldsStr := strings.Join(table.FieldNames, ",")

		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s FROM %s;", tmpTableName, fieldsStr, table.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmpTableName, table.Name))

		_, err = s.Exec()

		return
	})

	return
}
