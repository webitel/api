package auth

import (
	"../../helper"
	. "../shared"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/errors"
	"gopkg.in/mgo.v2/bson"
)

type Session struct {
	jwt.StandardClaims
	Id       string `json:"id,omitempty" bson:"key"`
	RoleName string `bson:"roleName"`
	Version  int8   `json:"v,omitempty"`
	Domain   string `json:"d,omitempty" bson:"domain"`
	Type     string `json:"t,omitempty"`
}

type Acl struct {
	Id      bson.ObjectId       `json:"_id" bson:"_id"`
	Role    string              `json:"roles" bson:"roles"`
	Parents string              `json:"parents" bson:"parents"`
	Allows  map[string][]string `json:"allows" bson:"allows"`
}

type roles struct {
	values map[string]Acl
}

func (r *roles) Reload() {
	var acl *[]Acl
	acl = &[]Acl{}
	DB.AclList(acl)

	for _, v := range *acl {
		r.values[v.Role] = v
	}
}

func NewRoles() (r *roles) {
	r = &roles{}
	r.values = make(map[string]Acl)
	r.Reload()
	return
}

var _roles *roles

func init() {
	_roles = NewRoles()
}

func AclFindDomain(key, domainName string) (error, string) {
	err, r := DB.AclFindDomain(key, domainName)
	if err != nil {
		return errorForbidden, ""
	}

	if len(r.Tokens) == 0 {
		return errorForbidden, ""
	}

	if r.Tokens[0].Name == "" {
		return errorForbidden, ""
	}

	return nil, r.Tokens[0].Name
}

func AclFindAuth(key string, s *Session) (err error) {
	DB.AclFindAuth(key, s)
	if s.RoleName == "" {
		return errorForbidden
	}
	return
}

var errorForbidden = helper.NewCodeError(403, errors.New("Forbidden"))

func CheckAcl(roleName string, resource string, operation string) *helper.CodeError {
	if roleName == "" || operation == "" {
		return errorForbidden
	}

	if acl, ok := _roles.values[roleName]; ok {
		if allow, ok := acl.Allows[resource]; ok {
			for _, p := range allow {
				if p == "*" || p == operation {
					return nil
				}
			}
		}

		if acl.Parents != "" {
			return CheckAcl(acl.Parents, resource, operation)
		}
	}
	return errorForbidden
}

func GetSessionFromContext(ctx context.Context) *Session {
	return ctx.Values().Get("user").(*Session)
}
