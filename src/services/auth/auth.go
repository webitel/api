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

type Allows struct {
	AclRoles       []string `json:"acl/roles" bson:"acl/roles"`
	AclResource    []string `json:"resource" bson:"resource"`
	Blacklist      []string `json:"blacklist" bson:"blacklist"`
	Calendar       []string `json:"calendar" bson:"calendar"`
	RotesDefault   []string `json:"rotes/default" bson:"rotes/default"`
	RotesPublic    []string `json:"rotes/public" bson:"rotes/public"`
	RotesExtension []string `json:"rotes/extension" bson:"rotes/extension"`
	RotesDomain    []string `json:"rotes/domain" bson:"rotes/domain"`
	Channels       []string `json:"channels" bson:"channels"`
	CcTiers        []string `json:"cc/tiers" bson:"cc/tiers"`
	CcMembers      []string `json:"cc/members" bson:"cc/members"`
	CcQueue        []string `json:"cc/queue" bson:"cc/queue"`
	Dialer         []string `json:"dialer" bson:"dialer"`
	DialerMembers  []string `json:"dialer/members" bson:"dialer/members"`
	Book           []string `json:"book" bson:"book"`
	Hook           []string `json:"hook" bson:"hook"`
	Cdr            []string `json:"cdr" bson:"cdr"`
	CdrFiles       []string `json:"cdr/files" bson:"cdr/files"`
	CdrMedia       []string `json:"cdr/media" bson:"cdr/media"`
	Gateway        []string `json:"gateway" bson:"gateway"`
	GatewayProfile []string `json:"gateway/profile" bson:"gateway/profile"`
	Domain         []string `json:"domain" bson:"domain"`
	DomainItem     []string `json:"domain/item" bson:"domain/item"`
	Account        []string `json:"account" bson:"account"`
	SystemReload   []string `json:"system/reload" bson:"system/reload"`
	License        []string `json:"license" bson:"license"`
	VMail          []string `json:"vmail" bson:"vmail"`
}

type Acl struct {
	Id      bson.ObjectId       `json:"_id" bson:"_id"`
	Role    string              `json:"roles" bson:"roles"`
	Parents string              `json:"parents" bson:"parents"`
	Allows  map[string][]string `json:"allows" bson:"allows"`
}

var roles map[string]Acl

func init() {
	roles = make(map[string]Acl)

	var acl *[]Acl
	acl = &[]Acl{}
	DB.AclList(acl)

	for _, v := range *acl {
		roles[v.Role] = v
	}
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
	if roleName == "" {
		return errorForbidden
	}

	if acl, ok := roles[roleName]; ok {
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
