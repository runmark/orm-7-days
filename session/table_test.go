package session

import "testing"

type User struct {
	Id int `geeorm:"PRIMARY KEY"`
	Name string
}


func TestSession_CreateTable(t *testing.T) {

	s := New(TestDB, Dialect).Model(User{})
	defer s.Clear()

	_ = s.DropTable()
	_ = s.CreateTable()
	if !s.HasTable() {
		t.Fatal("Failed to create table User")
	}
}

