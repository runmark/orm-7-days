package dialect

import "reflect"

var dialectsMap = make(map[string]Dialect)

type Dialect interface {
	DataTypeOf(typ reflect.Value) string
	TableExistSQL(tableName string) (string, []interface{})
}

func RegisterDialect(name string, d Dialect) {
	dialectsMap[name] = d
}

func GetDialect(name string) (d Dialect, ok bool) {
	d, ok = dialectsMap[name]
	return
}
