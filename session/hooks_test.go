package session

import (
	"example.com/mark/geeorm/log"
	"testing"
)

type Account struct {
	ID int `geeorm:"PRIMARY KEY"`
	Password string
}

func (account *Account) BeforeInsert(s *Session) error {
	log.Info("before insert", account)
	account.ID += 1000
	return nil
}

func (a *Account) AfterQuery(s *Session) error  {
	log.Info("after query", a)
	a.Password = "******"
	return nil
}

func TestSession_CallMethod(t *testing.T) {
	s:= NewSession().Model(&Account{})

	_ = s.DropTable()
	_ = s.CreateTable()
	_, _ = s.Insert(&Account{1, "123456"}, &Account{2, "wert"})

	u := &Account{}

	err := s.First(u)

	if err != nil || u.ID != 1001 || u.Password != "******" {
		t.Fatal("failed to call hooks after query, got ", u)
	}
}