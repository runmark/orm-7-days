package session

import "testing"

var (
	user1 = &User{3, "Tom", 20}
	user2 = &User{4, "Sam", 30}
	user3 = &User{5, "John", 50}
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


func TestSession_Limit(t *testing.T) {
	s := testRecordInit(t)
	var users []User
	err := s.Limit(1).Find(&users)
	if err != nil || len(users) != 1 {
		t.Fatal("failed to query with limit condition")
	}
}

func TestSession_Update(t *testing.T) {
	s := testRecordInit(t)
	affected, _ := s.Where("Name = ?", "Tom").Update("Age", 100)
	u := &User{}
	_ = s.OrderBy("Age DESC").First(u)

	if affected != 1 || u.Age != 100 {
		t.Errorf("expect (1, 100), but got (%v, %v)", affected, u.Id)
	}
}

func TestSession_DeleteAndCount(t *testing.T) {
	s := testRecordInit(t)
	affected, _ := s.Where("Name = ?", "Tom").Delete()
	count, _ := s.Count()

	if affected != 1 || count != 1 {
		t.Fatal("failed to delete or count")
	}
}
