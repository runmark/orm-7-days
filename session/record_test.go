package session

import "testing"

var (
	user1 = &User{3, "Tom"}
	user2 = &User{4, "Sam"}
	user3 = &User{5, "John"}
)

func testRecordInit(t *testing.T) *Session {
	t.Helper()
	s := NewSession().Model(&User{})

	err1 := s.DropTable()
	err2 := s.CreateTable()
	_, err3 := s.Insert(user1, user2)

	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatal("failed init test records")
	}

	return s
}

func TestSession_Insert(t *testing.T) {

	s := testRecordInit(t)

	affected, err := s.Insert(user3)
	if err != nil || affected != 1 {
		t.Fatal("failed to create record")
	}

}

func TestSession_Find(t *testing.T) {
	s := testRecordInit(t)

	var users []User
	err := s.Find(&users)

	if err != nil || len(users) != 2 {
		t.Fatal("failed to query all")
	}

}
