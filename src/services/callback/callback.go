package callback

import (
	"../../db"
	. "../../helper"
	"../../services/auth"
	. "../../services/shared"
	"github.com/kataras/iris/core/errors"
	"gopkg.in/mgo.v2/bson"
)

type Queue struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string        `bson:"name,omitempty" json:"name,omitempty"`
	Description string        `bson:"description,omitempty" json:"description,omitempty"`
}

type location struct {
	Ip          string  `bson:"ip,omitempty" json:"ip,omitempty"`
	CountryCode string  `bson:"country_code,omitempty" json:"country_code,omitempty"`
	CountryName string  `bson:"country_name,omitempty" json:"country_name,omitempty"`
	RegionCode  string  `bson:"region_code,omitempty" json:"region_code,omitempty"`
	RegionName  string  `bson:"region_name,omitempty" json:"region_name,omitempty"`
	City        string  `bson:"city,omitempty" json:"city,omitempty"`
	ZipCode     string  `bson:"zip_code,omitempty" json:"zip_code,omitempty"`
	TimeZone    string  `bson:"time_zone,omitempty" json:"time_zone,omitempty"`
	Latitude    float32 `bson:"latitude,omitempty" json:"latitude,omitempty"`
	Longitude   float32 `bson:"longitude,omitempty" json:"longitude,omitempty"`
}

type Comment struct {
	CreatedOn int    `bson:"createdOn,omitempty" json:"createdOn,omitempty"`
	CreatedBy string `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	Data      string `bson:"data,omitempty" json:"data,omitempty"`
}

type Member struct {
	Id           bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	Queue        string        `bson:"queue,omitempty" json:"queue,omitempty"`
	WidgetId     string        `bson:"widgetId" json:"widgetId"`
	Done         bool          `bson:"done" json:"done"`
	DoneAt       int           `bson:"doneAt" json:"doneAt"`
	DoneBy       string        `bson:"doneBy" json:"doneBy"`
	Domain       string        `bson:"domain" json:"domain"`
	Number       string        `bson:"number" json:"number"`
	CreatedOn    int64         `bson:"createdOn,omitempty"`
	CallbackTime int           `bson:"callbackTime,omitempty" json:"callbackTime,omitempty"`
	Href         string        `bson:"href,omitempty" json:"href,omitempty"`
	UserAgent    string        `bson:"userAgent,omitempty" json:"userAgent,omitempty"`
	Location     location      `bson:"location,omitempty" json:"location,omitempty"`
	Comments     []Comment     `bson:"comments" json:"comments"`
}

var errorDomainIsRequired = NewCodeError(400, errors.New("Domain is required"))

// region Queue Service
func CallbackQueueList(s *auth.Session, r *db.Request) (*[]Queue, *CodeError) {
	err := auth.CheckAcl(s.RoleName, "account", "d")
	if err != nil {
		return nil, err
	}
	if s.Domain != "" {
		r.Filter["domain"] = s.Domain
	}

	if _, ok := r.Filter["domain"]; !ok {
		return nil, errorDomainIsRequired
	}

	data := &[]Queue{}
	//err = DB.CallbackQueueList(r, data)

	if err != nil {
		return data, err
	}
	return data, nil
}

func CallbackQueueItem(queueId string) (*Queue, *CodeError) {
	data := &Queue{}
	err := DB.CallbackQueueItem(queueId, "", data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func CallbackQueueCreate(q *Queue) *CodeError {

	if q.Name == "" {
		return NewCodeError(400, errors.New("Name is required"))
	}

	q.Id = db.MakeId()
	return DB.CallbackQueueCreate(q)
}

func CallbackQueueDelete(queueId string) *CodeError {
	return DB.CallbackQueueDelete(queueId)
}

// endregion

//region Members Service
func CallbackMembersList(queueId string, r *db.Request) (*[]Member, *CodeError) {
	data := &[]Member{}
	err := DB.CallbackMembersList(queueId, r, data)

	if err != nil {
		return data, err
	}
	return data, nil
}

func CallbackMemberCreate(queueId string, m *Member) *CodeError {
	if m.Number == "" {
		return NewCodeError(400, errors.New("Number is required"))
	}

	m.Id = db.MakeId()
	m.Queue = queueId
	return DB.CallbackMemberCreate(m)
}

func CallbackMemberItem(queueId string, memberId string) (*Member, *CodeError) {
	data := &Member{}
	err := DB.CallbackMemberItem(queueId, memberId, data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func CallbackMemberDelete(queueId string, memberId string) *CodeError {
	err := DB.CallbackMemberDelete(queueId, memberId)
	if err != nil {
		return err
	}

	return nil
}

// endregion
