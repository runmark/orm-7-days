package geeorm

import (
	"errors"
	"example.com/mark/geeorm/session"
	_ "github.com/mattn/go-sqlite3"
	"reflect"
	"testing"
)

func TestNewEngine(t *testing.T) {
	e, err := NewEngine("sqlite3", "cmd_test/gee.db")
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	defer e.Close()
}

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	e, err := NewEngine("sqlite3", "./cmd_test/gee.db")
	if err != nil {
		t.Fatal("failed to connect", err)
	}

	return e
}

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func TestEngine_Transaction(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		transactionRollback(t)
	})

	t.Run("commit", func(t *testing.T) {
		transactionCommit(t)
	})
}

func transactionRollback(t *testing.T) {
	e := OpenDB(t)
	defer e.Close()

	s := e.NewSession()
	_ = s.Model(&User{}).DropTable()

	_, err := e.Transaction(func(s *session.Session) (interface{}, error) {
		_ = s.Model(&User{}).CreateTable()
		_, _ = s.Insert(&User{"Tom", 18})
		return nil, errors.New("rollback")
	})

	if err == nil || s.HasTable() {
		t.Fatal("failed to rollback")
	}

}

func transactionCommit(t *testing.T) {
	e := OpenDB(t)
	defer e.Close()

	s := e.NewSession()
	s.Model(&User{}).DropTable()

	_, err := e.Transaction(func(s *session.Session) (interface{}, error) {
		_ = s.Model(&User{}).CreateTable()
		_, _ = s.Insert(&User{"Tom", 18})
		return nil, nil
	})

	user := User{}
	err = s.First(&user)
	if err != nil || user.Age != 18 || user.Name != "Tom" {
		t.Fatal("commit failed", err)
	}

}


func TestEngine_Migrate(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text PRIMARY KEY, XXX integer);").Exec()
	_, _ = s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	engine.Migrate(&User{})

	rows, _ := s.Raw("SELECT * FROM User").QueryRows()
	columns, _ := rows.Columns()
	if !reflect.DeepEqual(columns, []string{"Name", "Age"}) {
		t.Fatal("Failed to migrate table User, got columns", columns)
	}
}
