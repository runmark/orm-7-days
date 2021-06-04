package schema

import (
	"example.com/mark/geeorm/dialect"
	"reflect"
	"testing"
)

type Person struct {
	Id   int `geeorm:"PRIMARY KEY"`
	Name string
}

func TestParse(t *testing.T) {
	cases := []struct {
		input Person

		oFieldNames []string
		oFieldTypes []string
		oFieldTags  []string
	}{
		{
			Person{3, "John"},

			[]string{"Id", "Name"},
			[]string{"integer", "text"},
			[]string{"PRIMARY KEY", ""},
		},
	}

	d, _ := dialect.GetDialect("sqlite3")

	for _, c := range cases {

		s := Parse(&c.input, d)

		if !reflect.DeepEqual(s.FieldNames, c.oFieldNames) {
			t.Errorf("expect %s, but got %s\n", c.oFieldNames, s.FieldNames)
		}

		sFieldTypes, sFieldTags := []string{}, []string{}

		for _, f := range s.Fields {
			sFieldTypes = append(sFieldTypes, f.Type)
			sFieldTags = append(sFieldTags, f.Tag)
		}

		if !reflect.DeepEqual(sFieldTypes, c.oFieldTypes) || !reflect.DeepEqual(sFieldTags, c.oFieldTags) {
			t.Errorf("expect (%s, %s), but got (%s, %s)\n", c.oFieldTypes, c.oFieldTags, sFieldTypes, sFieldTags)
		}

	}

}
