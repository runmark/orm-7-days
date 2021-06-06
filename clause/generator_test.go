package clause

import (
	"fmt"
	"testing"
)

type User struct {
	Id   int
	Name string
}

type cov interface {}

func TestValues_(t *testing.T) {

	input := []interface{}{
		User{4, "John"},
		User{5, "Amily"},
	}
	//input := [][](interface{}){
	//	{User{4, "John"}},
	//		{User{5, "Amily"}},
	//}

	//input := []User{
	//	{4, "John"},{5, "Amily"},
	//}

	//inputConvt := []interface{}{}
	//
	//for _, inp := range input {
	//	inputConvt = append(inputConvt, cov(inp))
	//}

	sql, vars := values_(input...)

	if sql != "VALUES (?,?),(?,?)" {
		t.Errorf("expect %v, but got %v", "VALUES (?,?),(?,?)", sql)
	}

	if fmt.Sprint(vars) != fmt.Sprint(input) {
		t.Errorf("expect %v, but got %v", input, vars)
	}

}
