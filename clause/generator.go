package clause

import (
	"example.com/mark/geeorm/log"
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

	generators[UPDATE] = update_
	generators[DELETE] = delete_
	generators[COUNT] = count_
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

func limit_(values ...interface{}) (string, []interface{}) {
	// LIMIT $num
	return "LIMIT ?", values
}

func where_(values ...interface{}) (string, []interface{}) {
	// WHERE $conds
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

func orderby_(values ...interface{}) (string, []interface{}) {
	// OEDERBY $fields
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}

func update_(values ...interface{}) (string, []interface{}) {
	// UPDATE $tableName SET $field1=, $field2=...
	tableName := values[0]

	log.Infof("update_: values: %v", values)

	var keys []string
	var vars []interface{}

	for k, v := range values[1].(map[string]interface{}) {
		keys = append(keys, fmt.Sprintf("%s=?", k))
		vars = append(vars, v)
	}

	log.Infof("update_: keys: %v, vars: %v\n", keys, vars)

	sql := fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(keys, ","))

	log.Infof("update_: sql: %v, vars: %v\n", sql, vars)

	return sql, vars
}

func delete_(values ...interface{}) (string, []interface{}) {
	// DELETE FROM $tableName

	return fmt.Sprintf("DELETE FROM %s", values[0]), []interface{}{}
}

func count_(values ...interface{}) (string, []interface{}) {
	// SELECT count(*) FROM $tableName

	return select_(values[0], []string{"count (*)"})
}
