package db

import (
	. "../config"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

var COLLECTION_ACL = Config.Get("mongodb:aclCollection")
var COLLECTION_DOMAINS = Config.Get("mongodb:domainsCollection")
var COLLECTION_AUTH = Config.Get("mongodb:authTokenCollection")

func (db *DB) AclList(dataStructure interface{}) (err error) {
	err = db.db.C(COLLECTION_ACL).Find(nil).All(dataStructure)
	return
}

type token struct {
	Name   string `bson:"roleName"`
	Expire int    `bson:"expire"`
}

type RoleDomain struct {
	Tokens []token `bson:"tokens"`
}

func (db *DB) AclFindAuth(key string, dataStructure interface{}) (err error) {
	err = db.db.C(COLLECTION_AUTH).Find(bson.M{
		"key": key,
	}).Select(bson.M{
		"roleName": 1,
		"domain":   1,
		"key":      1,
		"_id":      0,
	}).One(dataStructure)
	return
}

func (db *DB) AclFindDomain(key string, domainName string) (err error, r RoleDomain) {
	err = db.db.C(COLLECTION_DOMAINS).Find(bson.M{
		"name": domainName,
		"tokens": bson.M{
			"$elemMatch": bson.M{
				"uuid":    key,
				"enabled": true,
			},
		},
	}).Select(bson.M{
		"tokens.$": 1,
		"_id":      0,
	}).One(&r)

	return
}

func (db *DB) CheckAcl(role string, resource string, perm string) (err error) {
	var resp map[string]interface{}
	err = db.db.C(COLLECTION_ACL).Pipe([]bson.M{
		{"$match": bson.M{"roles": role}},
		{"$graphLookup": bson.M{
			"from":             COLLECTION_ACL,
			"startWith":        "$parents",
			"connectFromField": "parents",
			"connectToField":   "roles",
			"as":               "_parents",
		}},

		{"$match": bson.M{
			"$or": []bson.M{
				bson.M{"allows.blacklist": bson.M{
					"$in": []string{"r", "*"},
				}},
				bson.M{"_parents.allows.blacklist": bson.M{
					"$in": []string{"r", "*"},
				}},
			},
		}},

		{"$project": bson.M{
			"roles": 1,
		}},
	}).One(&resp)

	fmt.Println(resp)

	//err = db.db.C("aclPermissions_view").Find(bson.M{
	//	"roles": "admin",
	//	"allows.blacklist": bson.M{
	//		"$in": []string{"*", "r"},
	//	},
	//}).Select(bson.M{
	//	"_id":                       1,
	//	"allows.blacklist":          1,
	//	"_parents.allows.blacklist": 1,
	//}).One(&resp)
	//fmt.Println(err, resp)
	return
}
