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

var errorForbidden = helper.NewCodeError(403, errors.New("Forbidden"))

func CheckAcl(s *Session, resource string, operation string) *helper.CodeError {
	if s.Id == "" || operation == "" {
		return errorForbidden
	}

	err, _, _, domainName := DB.CheckAcl(s.Id, resource, operation)
	if err != nil {
		return err
	}
	s.Domain = domainName
	return nil
}

func GetSessionFromContext(ctx context.Context) *Session {
	return ctx.Values().Get("user").(*Session)
}
