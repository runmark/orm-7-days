package session

import (
	"example.com/mark/geeorm/clause"
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
