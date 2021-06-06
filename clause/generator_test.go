package clause

import (
	"fmt"
	"testing"
)

type User struct {
	Id int
	Name string
}

func TestValues_(t *testing.T) {

	input := []interface{}{
		User{4, "John"},
		User{5, "Amily"},
	}

	sql, vars := values_(input...)

	if sql != "VALUES (?,?),(?,?)" {
		t.Errorf("expect %v, but got %v", "VALUES (?,?),(?,?)", sql)
	}

	if fmt.Sprint(vars) != fmt.Sprint(input) {
		t.Errorf("expect %v, but got %v", input, vars)
	}

}