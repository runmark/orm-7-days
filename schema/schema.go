package schema

type Field struct {
	Name string
	Type string
	Tag  string
}

type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	FieldMap   map[string]*Field
}

func (s *Schema) GetField(name string) *Field {
	return s.FieldMap[name]
}
