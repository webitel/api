package acl

import (
	. "../../helper"
	"../../services/auth"
	. "../../services/shared"
	"database/sql"
	"github.com/kataras/iris/core/errors"
	"github.com/lib/pq"
)

type Group struct {
	Id       int64  `json:"id"`
	ParentId int64  `json:"parentId,omitempty"`
	Name     string `json:"name"`
}
type Permission struct {
	Id         int64    `json:"id"`
	ObjectType string   `json:"objectType,omitempty"`
	ObjectId   string   `json:"objectId,omitempty"`
	GroupId    int64    `json:"groupId"`
	Allows     []string `json:"allows"`
}

var sqlList = `
	SELECT
		id,
		name,
		parent_id
	FROM acl_group
`

var sqlPermList = `
	SELECT
		id,
		COALESCE(object_type, '') as object_type,
		COALESCE(object_id, '') as object_id,
		group_id,
		perm
	FROM acl_permission
	WHERE group_id = (select acl_group.id from acl_group where acl_group.name = $1)
`

var sqlUpdatePerm = `
	with upsert as (
	  update acl_permission
	  set (perm) = ($1)
	  where group_id = (select acl_group.id from acl_group where acl_group.name = $2 LIMIT 1)  AND object_type = $3
	  returning *
	)
	INSERT INTO acl_permission (object_type, group_id, perm)
	select $3, (select acl_group.id from acl_group where acl_group.name = $2 LIMIT 1), $1
	WHERE NOT EXISTS (SELECT * FROM upsert);
`

var sqlRemoveGroupByName = `
	DELETE FROM acl_group WHERE name = $1
`

var sqlCreateGroup = `
	INSERT INTO acl_group (name, parent_id)
	VALUES ($1, (select id from acl_group where name = $2))
`

func List(s *auth.Session) (err *CodeError, data []*Group) {
	err = auth.CheckAcl(s, "acl/roles", "r")
	if err != nil {
		return
	}

	rows, e := DB.GetPg().Query(sqlList)
	defer rows.Close()
	if e != nil {
		err = NewCodeError(500, e)
		return
	}

	var p sql.NullInt64

	for rows.Next() {
		g := new(Group)
		e = rows.Scan(&g.Id, &g.Name, &p)

		if p.Valid {
			g.ParentId = p.Int64
		}

		if e != nil {
			err = NewCodeError(500, e)
			return
		}

		data = append(data, g)
	}

	if e = rows.Err(); e != nil {
		err = NewCodeError(500, e)
		return
	}

	return
}

func RemoveGroup(s *auth.Session, groupName string) (err *CodeError) {
	err = auth.CheckAcl(s, "acl/roles", "d")
	if err != nil {
		return
	}

	_, e := DB.GetPg().Exec(sqlRemoveGroupByName, groupName)

	if e != nil {
		err = NewCodeError(500, e)
	}
	return
}

var errBadRequestNameRequired = NewCodeError(400, errors.New("Name is required"))

func CreateGroup(s *auth.Session, groupName, parentName string) (err *CodeError) {

	if groupName == "" {
		err = errBadRequestNameRequired
		return
	}

	err = auth.CheckAcl(s, "acl/roles", "c")
	if err != nil {
		return
	}

	_, e := DB.GetPg().Exec(sqlCreateGroup, groupName, parentName)

	if e != nil {
		err = NewCodeError(500, e)
	}
	return
}

func ListPerms(s *auth.Session, groupName string) (err *CodeError, data []*Permission) {
	err = auth.CheckAcl(s, "acl/resource", "r")
	if err != nil {
		return
	}

	rows, e := DB.GetPg().Query(sqlPermList, groupName)
	defer rows.Close()
	if e != nil {
		err = NewCodeError(500, e)
		return
	}

	var tmp string
	for rows.Next() {
		g := new(Permission)

		e = rows.Scan(&g.Id, &g.ObjectType, &g.ObjectId, &g.GroupId, &tmp)

		if tmp != "" {
			g.Allows, e = DB.ParsePGArray(tmp)
		}

		if e != nil {
			err = NewCodeError(500, e)
			return
		}

		data = append(data, g)
	}

	if e = rows.Err(); e != nil {
		err = NewCodeError(500, e)
		return
	}

	return
}

func UpdatePermission(s *auth.Session, groupName string, data map[string][]string) (err *CodeError) {
	err = auth.CheckAcl(s, "acl/resource", "u")
	if err != nil {
		return
	}

	for objectType, perm := range data {
		_, e := DB.GetPg().Exec(sqlUpdatePerm, pq.Array(perm), groupName, objectType)
		if e != nil {
			err = NewCodeError(500, e)
		}
		return
	}
	return
}
