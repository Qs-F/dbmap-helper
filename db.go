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
package service

import (
	"database/sql"
	"errors"
	"gopkg.in/gorp.v1"
	"reflect"
)

const (
	TAGNAME    = "role"
	TypeMySQL  = "mysql"
	TypeSQLite = "sqlite"
)

type DB struct {
	DB     *sql.DB
	DBType string
	DbMap  *gorp.DbMap
}

type TableMap struct {
	Name   string
	Struct interface{}
}

// gorp contains some dialects. how we use:
// switch Dialect by DBtype
// sql.Open must be written main.go by user. this is perfect !!!!!!
func (d *DB) Init(s []TableMap) (err error) {
	if d.DBType == "sqlite" {
		d.DbMap = &gorp.DbMap{Db: d.DB, Dialect: gorp.SqliteDialect{}}
	} else if d.DBType == "mysql" {
		d.DbMap = &gorp.DbMap{Db: d.DB, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	}
	for _, v := range s { // pick up 1 struct
		rt := reflect.TypeOf(v.Struct) // get struct info
		id := ""
		for j := 0; j < rt.NumField(); j++ { // to rename field name. pick up 1 struct field's tag
			field := rt.Field(j) // focus field
			if n := field.Tag.Get(TAGNAME); n == "id" && id == "" {
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
