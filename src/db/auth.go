package db

import (
	"../helper"
	"github.com/kataras/iris/core/errors"
	_ "github.com/lib/pq"
)

//var pg *sql.DB
//
//func init() {
//	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s",
//		"webitel", "webitel", "webitel", "10.10.10.200")
//	var err error
//	pg, err = sql.Open("postgres", dbinfo)
//	pg.SetMaxOpenConns(100)
//	if err != nil {
//		panic(err)
//	}
//
//}
//
//var COLLECTION_ACL = Config.Get("mongodb:aclCollection")
//
//func (db *DB) AclList(dataStructure interface{}) (err error) {
//	err = db.db.C(COLLECTION_ACL).Find(nil).All(dataStructure)
//	return
//}
//
//type token struct {
//	Name   string `bson:"roleName"`
//	Expire int    `bson:"expire"`
//}
//
//type RoleDomain struct {
//	Tokens []token `bson:"tokens"`
//}

var errorPerm = helper.NewCodeError(403, errors.New("Forbidden"))

func (db *DB) CheckAcl(key string, resource string, perm string) (err *helper.CodeError, groupId int, userName string, domainName string) {
	// (group_id, userName, domainName)
	rows, e := db.pg.Query(`select group_id, username, domain_name from check_permission($1, $2, $3)`, key, resource, perm)
	defer rows.Close()
	if e != nil {
		err = helper.NewCodeError(500, e)
		return
	}

	rows.Next()
	rows.Scan(&groupId, &userName, &domainName)
	if groupId == 0 {
		err = errorPerm
	}
	return
}
