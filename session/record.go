package session

import (
	"example.com/mark/geeorm/clause"
	"example.com/mark/geeorm/log"
	"reflect"
)

func (s *Session) Insert(values ...interface{}) (nums int64, err error) {

	recordValues := make([]interface{}, 0)

	for _, value := range values {

		schema := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, schema.Name, schema.FieldNames)
		recordValues = append(recordValues, schema.RecordValues(value))

	}

	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)

	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (session *Session) Find(values interface{}) error {

	rv := reflect.Indirect(reflect.ValueOf(values))
	rte := rv.Type().Elem()

	table := session.Model(reflect.New(rte).Elem().Interface()).RefTable()

	session.clause.Set(clause.SELECT, table.Name, table.FieldNames)

	sql, vars := session.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)

	log.Infof("convert Find(%v) to SQL: %s, %s\n", values, sql, vars)

	rows, err := session.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {
		var rfields []interface{}
		rstruct := reflect.New(rte).Elem()

		for _, name := range table.FieldNames {
			rfields = append(rfields, rstruct.FieldByName(name).Addr().Interface())
		}

		err = rows.Scan(rfields...)
		if err != nil {
			return err
		}

		rv.Set(reflect.Append(rv, rstruct))
	}

	return rows.Close()
}
