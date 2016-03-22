package helper

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

type TestS struct {
	Field1 string `db:"field_1"`
	Gl     int64  `db:"field_2"`
	GH     string
	Id     int64
}
type Test2 struct {
	Name  string
	Id    int64
	TmpId int64 `role:"id"`
}

func TestInit(t *testing.T) {
	db, err := sql.Open("sqlite3", "./test-edg32yd.db")
	if err != nil {
		t.Error(err.Error())
	}
	d := &DB{db, TypeSQLite, nil}
	o := []TableMap{{"table1", TestS{}}, {"table2", Test2{}}}
	err = d.Init(o)
	if err != nil {
		t.Error(err.Error())
	}
	err = d.DbMap.Insert(&Test2{"fumi", 0, 0})
	if err != nil {
		t.Error(err.Error())
	}
	l, err := d.DbMap.Select(Test2{}, "select * from table2 where name = ?", "fumi")
	if err != nil {
		t.Error(err.Error())
	}
	for _, v := range l {
		p := v.(*Test2)
		t.Log(p.Name)
	}
}
