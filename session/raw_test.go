package session

import (
	"database/sql"
	"example.com/mark/geeorm/dialect"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var err error
var TestDB *sql.DB
var Dialect dialect.Dialect

func TestMain(m *testing.M)  {

	TestDB, _ = sql.Open("sqlite3", "../cmd_test/gee.db")

	Dialect, _ = dialect.GetDialect("sqlite3")

	code := m.Run()

	_ = TestDB.Close()

	os.Exit(code)
}


func TestSession_Exec(t *testing.T) {
	s := New(TestDB, Dialect)
	defer s.Clear()

	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE user (name text);").Exec()
	r, _ := s.Raw("INSERT INTO user (`name`) VALUES (?), (?)", "Sam", "John").Exec()

	count, err := r.RowsAffected()
	if count != 2 || err != nil {
		t.Fatal("Expect 2, but got ", count)
	}
}

func TestSession_QueryRow(t *testing.T) {
	s := New(TestDB, Dialect)
	defer s.Clear()

	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE user (name text);").Exec()
	r := s.Raw("SELECT count(*) FROM User;").QueryRow()

	var count int
	err := r.Scan(&count)
	if count != 0 || err != nil {
		t.Fatal("failed to query db", err)
	}
}