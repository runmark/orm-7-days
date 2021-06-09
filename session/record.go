package session

import (
	"errors"
	"example.com/mark/geeorm/clause"
	"example.com/mark/geeorm/log"
	"reflect"
)

func (s *Session) Insert(values ...interface{}) (nums int64, err error) {

	recordValues := make([]interface{}, 0)

	for _, value := range values {

		//s.CallMethod(BeforeInsert, value)

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

	//session.CallMethod(BeforeQuery, nil)

	rv := reflect.Indirect(reflect.ValueOf(values))
	rte := rv.Type().Elem()

	// 为什么直接使用 rv 不可以？
	// 因为 values 底层类型是 []struct，而不是 struct.
	//table := session.Model(rv.Interface()).RefTable()
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

		//session.CallMethod(AfterQuery, rstruct.Addr().Interface())

		rv.Set(reflect.Append(rv, rstruct))


	}

	return rows.Close()
}

// support map[string]interface{}
// also support kv list: "Name", "Tom", "Age", 18, ....
func (session *Session) Update(kv ...interface{}) (int64, error) {

	m, ok := kv[0].(map[string]interface{})

	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}

	log.Infof("Update: kv: %v, m: %v", kv, m)

	session.clause.Set(clause.UPDATE, session.RefTable().Name, m)

	log.Infof("Update: tableName: %v, args: %v", session.RefTable().Name, m)

	sql, vars := session.clause.Build(clause.UPDATE, clause.WHERE)

	log.Infof("Update: sql: %v, vars: %v\n", sql, vars)

	result, err := session.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (s *Session) Delete() (int64, error) {
	table := s.RefTable()
	s.clause.Set(clause.DELETE, table.Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)

	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (s *Session) Count() (int64, error) {
	table := s.RefTable()
	s.clause.Set(clause.COUNT, table.Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)

	row := s.Raw(sql, vars...).QueryRow()

	var count int64

	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	vars = append(vars, desc)
	vars = append(vars, args...)
	s.clause.Set(clause.WHERE, vars...)
	return s
}

func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY,  desc)
	return s
}

func (s *Session) First(value interface{}) error {
	// Find() + OrderBy() + Limit(1)
	rv := reflect.Indirect(reflect.ValueOf(value))
	rt := rv.Type()

	rts := reflect.SliceOf(rt)

	rvs := reflect.New(rts).Elem()

	//rvs.Set(reflect.Append(rvs, rv))

	// 注意因为 Find() 的入参为 切片指针，
	//因此需要将 rvs 改为指针形式传入，即 rvs.Addr()
	err := s.Limit(1).Find(rvs.Addr().Interface())
	if err != nil {
		return err
	}

	log.Infof("rvs is %v\n", rvs)

	if rvs.Len() == 0 {
		return errors.New("NOT FOUND")
	}

	rv.Set(rvs.Index(0))

	return nil
}
