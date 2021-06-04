package session

import (
	"example.com/mark/geeorm/log"
	"example.com/mark/geeorm/schema"
	"fmt"
	"reflect"
	"strings"
)

func (s *Session) Model(origin interface{}) *Session {

	if s.refTable == nil || reflect.TypeOf(s.refTable.Model) != reflect.TypeOf(origin) {
		s.refTable = schema.Parse(origin, s.dialect)
	}

	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set!")
	}
	return s.refTable
}

func (s *Session) CreateTable() (err error) {

	fieldStrs := make([]string, 0)

	for _, field := range s.RefTable().Fields {
		fieldStr := fmt.Sprintf("%v %v ", field.Name, field.Type)
		if field.Tag != "" {
			fieldStr += field.Tag
		}
		fieldStrs = append(fieldStrs, fieldStr)
	}

	createSql := fmt.Sprintf("CREATE TABLE %v (%v);", s.refTable.Name, strings.Join(fieldStrs, ","))

	_, err = s.Raw(createSql).Exec()

	return
}

func (s *Session) DropTable() error {

	dropSql := fmt.Sprintf("DROP TABLE IF EXISTS %s;", s.refTable.Name)

	_, err := s.Raw(dropSql).Exec()

	return err

}

func (s *Session) HasTable() (ok bool) {
	sql, args := s.dialect.TableExistSQL(s.RefTable().Name)

	row := s.Raw(sql, args...).QueryRow()

	var name string
	err := row.Scan(&name)
	if err != nil {
		log.Error(err)
	}

	return name == s.RefTable().Name
}
