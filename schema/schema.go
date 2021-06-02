package schema

type Field struct {
	name string
	typ string
	tag string
}

type Schema struct {
	name string
	fields []*Field
}
