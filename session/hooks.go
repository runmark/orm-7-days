package session

import (
	"example.com/mark/geeorm/log"
	"reflect"
)

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

func (s *Session) CallMethod(method string, value interface{}) {
	rm := reflect.ValueOf(s.RefTable().Model).MethodByName(method)

	if value != nil {
		rm = reflect.ValueOf(value).MethodByName(method)
	}

	params := []reflect.Value{}
	params = append(params, reflect.ValueOf(s))

	results := rm.Call(params)

	if rm.IsValid() {

		if len(results) > 0 {

			if err, ok := results[0].Interface().(error); ok {
				log.Error(err)
			}

		}
	}

}
