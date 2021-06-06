package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)

	generators[INSERT] = insert_
	generators[VALUES] = values_
	generators[SELECT] = select_
	generators[LIMIT] = limit_
	generators[WHERE] = where_
	generators[ORDERBY] = orderby_
}

func genBindVars(num int) string {
	vars := make([]string, num)
	for i := 0; i < num; i++ {
		vars[i] = "?"
	}
	return strings.Join(vars, ",")
}

func insert_(values ...interface{}) (string, []interface{}) {
	// INSERT INTO $tableName ($fields)
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")

	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), []interface{}{}
}

func values_(values ...interface{}) (string, []interface{}) {
	// VALUES ($v1), ($v2), ...
	var bindStr string
	var sql strings.Builder
	var vars []interface{}

	if len(values) == 0 {
		return "", nil
	}

	sql.WriteString("VALUES ")

	for i, value := range values {

		// 使用 struct 传输 values 时，需 comment 该语句
		// 使用 slice 传输 values 时，需 uncomment 该语句
		v := value.([]interface{})

		if bindStr == "" {
			// 使用 struct 传输 values 时，需 uncomment 该语句
			//bindStr = genBindVars(reflect.TypeOf(value).NumField())
			bindStr = genBindVars(len(v))
		}

		sql.WriteString(fmt.Sprintf("(%v)", bindStr))

		if i != len(values)-1 {
			sql.WriteString(",")
		}

		vars = append(vars, v...)
	}

	return sql.String(), vars
}

func select_(values ...interface{}) (string, []interface{}) {
	// SELECT $fields FROM $tableName
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %s FROM %s", fields, tableName), []interface{}{}
}

func limit_(values ...interface{}) (string, []interface{})  {
	// LIMIT $num
	return "LIMIT ?", values
}

func where_(values ...interface{}) (string, []interface{})  {
	// WHERE $conds
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

func orderby_(values ...interface{}) (string, []interface{})  {
	// OEDERBY $fields
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}



