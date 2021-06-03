package dialect

import (
	"reflect"
	"testing"
	"time"
)

func TestSqlite3_DataTypeOf(t *testing.T) {
	d := sqlite3{}

	cases := []struct{
		Value interface{}
		Type string
	}{
		{"Tom", "text"},
		{true, "bool"},
		{123, "integer"},
		{12.3, "real"},
		{[]int{1,2,3}, "blob"},
		{time.Now(), "datetime"},
	}

	for _,c := range cases {
		typ := d.DataTypeOf(reflect.ValueOf(c.Value))
		if typ != c.Type {
			t.Errorf("expect %s, but got %s", c.Type, typ)
		}
	}
}