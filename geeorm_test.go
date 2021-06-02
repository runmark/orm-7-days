package geeorm

import (
	_ "github.com/mattn/go-sqlite3"
	"testing"
)


func TestNewEngine(t *testing.T) {
	e, err := NewEngine("sqlite3", "cmd_test/gee.db")
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	defer e.Close()
}
