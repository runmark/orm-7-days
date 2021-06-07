package schema

import (
	"example.com/mark/geeorm/dialect"
	"go/ast"
	"reflect"
)

type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	FieldMap   map[string]*Field
}

type Field struct {
	Name string
	Type string
	Tag  string
}

func (s *Schema) GetField(name string) *Field {
	return s.FieldMap[name]
}

func (s *Schema) RecordValues(origin interface{}) (fieldValues []interface{}) {

	rov := reflect.Indirect(reflect.ValueOf(origin))

	for _, name := range s.FieldNames {
		fieldValues = append(fieldValues, rov.FieldByName(name).Interface())
	}

	return
}

func Parse(origin interface{}, dialect dialect.Dialect) (schema *Schema) {
	rv := reflect.Indirect(reflect.ValueOf(origin))
	rt := rv.Type()

	schema = &Schema{
		Model:    origin,
		Name:     rt.Name(),
		FieldMap: make(map[string]*Field),
	}

	for i := 0; i < rt.NumField(); i++ {

		rfield := rt.Field(i)

		if !rfield.Anonymous && ast.IsExported(rfield.Name) {

			schema.FieldNames = append(schema.FieldNames, rfield.Name)

			field := &Field{}
			field.Name = rfield.Name
			//field.Type = dialect.DataTypeOf(reflect.New(rfield.Type).Elem())
			field.Type = dialect.DataTypeOf(reflect.Indirect(reflect.New(rfield.Type)))

			rtag := rfield.Tag
			t, ok := rtag.Lookup("geeorm")
			if ok {
				field.Tag = t
			}

			schema.Fields = append(schema.Fields, field)

			schema.FieldMap[field.Name] = field
		}

	}

	return
}
