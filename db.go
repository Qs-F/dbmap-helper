// copyright 2016 de-liKeR @CreatorQsF
//
// this is Database provider for service.1.
// This pkg uses MySQL mainly, but will have Option for another rdb.
//
// design: Init provides gorp's dbmap.
// d := service.DB{sql.Open(AS-YOU-LIKE), DBTYPE(mysql, sqlite, etc...)}
// d.Init()
// d.DbMap.Select(struct{}, SQL)
// d.DbMap.Close()
package helper

import (
	"database/sql"
	"errors"
	"gopkg.in/gorp.v1"
	"reflect"
)

const (
	tagname_  = "role"
	TypeMySQL = iota
	TypeSQLite
)

type DB struct {
	DB     *sql.DB
	DBType int
	DbMap  *gorp.DbMap
}

type TableMap struct {
	Name   string
	Struct interface{}
}

func (d *DB) Init(s []TableMap) (err error) {
	if d.DBType == TypeSQLite {
		d.DbMap = &gorp.DbMap{Db: d.DB, Dialect: gorp.SqliteDialect{}}
	} else if d.DBType == TypeMySQL {
		d.DbMap = &gorp.DbMap{Db: d.DB, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	}
	for _, v := range s { // pick up 1 struct
		rt := reflect.TypeOf(v.Struct) // get struct info
		id := ""
		for j := 0; j < rt.NumField(); j++ { // to rename field name. pick up 1 struct field's tag
			field := rt.Field(j) // focus field
			if n := field.Tag.Get(tagname_); n == "id" && id == "" {
				id = field.Name
			} else if n == "id" && id != "" {
				err = errors.New("id must be only.")
				return
			}
		}
		if id == "" {
			id = "Id"
		}
		d.DbMap.AddTableWithName(v.Struct, v.Name).SetKeys(true, id) // register picked up struct as 1 table(key is Id)
	}
	err = d.DbMap.CreateTablesIfNotExists()
	return
}
