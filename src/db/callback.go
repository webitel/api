package db

import (
	. "../config"
	"../helper"
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var COLLECTION_CALL_TRACKING = Config.Get("mongodb:callbackQueueCollection")
var COLLECTION_QUEUE = Config.Get("mongodb:callbackQueueCollection")
var COLLECTION_MEMBERS = Config.Get("mongodb:callbackMembersCollection")

func (db *DB) CallbackQueueList(r *Request, dataStructure interface{}) (err *helper.CodeError) {
	q := r.GetQueue(db, COLLECTION_QUEUE, "")
	e := q.All(dataStructure)

	if e != nil {
		err = helper.NewCodeError(500, e)
	}
	return
}

func (db *DB) CallbackQueueItem(queueId string, domainName string, dataStructure interface{}) (err *helper.CodeError) {
	if !bson.IsObjectIdHex(queueId) {
		err = helper.NewCodeError(400, errors.New("Bad queueId"))
		return
	}

	e := db.db.C(COLLECTION_QUEUE).Find(bson.M{
		"_id": bson.ObjectIdHex(queueId),
		//"domain": domainName,
	}).One(dataStructure)

	if e == mgo.ErrNotFound {
		err = helper.NewCodeError(404, e)
	} else if e != nil {
		err = helper.NewCodeError(500, e)
	}
	return
}

func (db *DB) CallbackQueueCreate(dataStructure interface{}) (err *helper.CodeError) {
	e := db.db.C(COLLECTION_QUEUE).Insert(dataStructure)

	if e == mgo.ErrNotFound {
		err = helper.NewCodeError(404, e)
	} else if e != nil {
		err = helper.NewCodeError(500, e)
	}
	return
}

func (db *DB) CallbackQueueDelete(queueId string) (err *helper.CodeError) {
	if !bson.IsObjectIdHex(queueId) {
		err = helper.NewCodeError(403, errors.New("Bad queueId"))
		return
	}

	e := db.db.C(COLLECTION_QUEUE).RemoveId(bson.ObjectIdHex(queueId))

	if e == mgo.ErrNotFound {
		err = helper.NewCodeError(404, e)
	} else if e != nil {
		err = helper.NewCodeError(500, e)
	}
	return
}

func (db *DB) CallbackQueueUpdate(queueId string, data map[string]interface{}) (err *helper.CodeError) {
	if !bson.IsObjectIdHex(queueId) {
		err = helper.NewCodeError(403, errors.New("Bad queueId"))
		return
	}

	set := bson.M{}

	for k, v := range data {
		if k == "_id" || k == "domain" || k == "queue" {
			continue
		}
		set[k] = v
	}

	e := db.db.C(COLLECTION_QUEUE).UpdateId(bson.ObjectIdHex(queueId), bson.M{
		"$set": set,
	})
	if e == mgo.ErrNotFound {
		err = helper.NewCodeError(404, e)
	} else if e != nil {
		err = helper.NewCodeError(500, e)
	}
	return
}

func (db *DB) CallbackMembersList(queueId string, r *Request, dataStructure interface{}) (err *helper.CodeError) {
	r.Filter["queue"] = queueId
	q := r.GetQueue(db, COLLECTION_MEMBERS, "")
	e := q.All(dataStructure)

	if e != nil {
		err = helper.NewCodeError(500, e)
	}
	return
}

func (db *DB) CallbackMemberCreate(dataStructure interface{}) (err *helper.CodeError) {
	e := db.db.C(COLLECTION_MEMBERS).Insert(dataStructure)

	if e == mgo.ErrNotFound {
		err = helper.NewCodeError(404, e)
	} else if e != nil {
		err = helper.NewCodeError(500, e)
	}
	return
}

func (db *DB) CallbackMemberItem(queueId string, memberId string, dataStructure interface{}) (err *helper.CodeError) {
	if !bson.IsObjectIdHex(memberId) {
		err = helper.NewCodeError(400, errors.New("Bad memberId"))
		return
	}

	e := db.db.C(COLLECTION_MEMBERS).Find(bson.M{
		"_id":   bson.ObjectIdHex(memberId),
		"queue": queueId,
		//"domain": domainName,
	}).One(dataStructure)

	if e == mgo.ErrNotFound {
		err = helper.NewCodeError(404, e)
	} else if e != nil {
		err = helper.NewCodeError(500, e)
	}
	return
}

func (db *DB) CallbackMemberDelete(queueId string, memberId string) (err *helper.CodeError) {
	if !bson.IsObjectIdHex(memberId) {
		err = helper.NewCodeError(403, errors.New("Bad memberId"))
		return
	}

	e := db.db.C(COLLECTION_MEMBERS).RemoveId(bson.ObjectIdHex(memberId))

	if e == mgo.ErrNotFound {
		err = helper.NewCodeError(404, e)
	} else if e != nil {
		err = helper.NewCodeError(500, e)
	}
	return
}

func (db *DB) CallbackMemberUpdate(queueId, memberId string, data map[string]interface{}) (err *helper.CodeError) {
	if !bson.IsObjectIdHex(queueId) {
		err = helper.NewCodeError(403, errors.New("Bad queueId"))
		return
	}
	if !bson.IsObjectIdHex(memberId) {
		err = helper.NewCodeError(403, errors.New("Bad memberId"))
		return
	}

	set := bson.M{}

	for k, v := range data {
		if k == "_id" || k == "domain" || k == "queue" || k == "comments" {
			continue
		}
		set[k] = v
	}

	e := db.db.C(COLLECTION_MEMBERS).Update(bson.M{
		"_id":   bson.ObjectIdHex(memberId),
		"queue": queueId,
	}, bson.M{
		"$set": set,
	})
	if e == mgo.ErrNotFound {
		err = helper.NewCodeError(404, e)
	} else if e != nil {
		err = helper.NewCodeError(500, e)
	}
	return
}

func (db *DB) CallbackMemberCommentAdd(queueId, memberId string, dataStructure interface{}) (err *helper.CodeError) {
	if !bson.IsObjectIdHex(queueId) {
		err = helper.NewCodeError(403, errors.New("Bad queueId"))
		return
	}
	if !bson.IsObjectIdHex(memberId) {
		err = helper.NewCodeError(403, errors.New("Bad memberId"))
		return
	}

	e := db.db.C(COLLECTION_MEMBERS).Update(bson.M{
		"_id":   bson.ObjectIdHex(memberId),
		"queue": queueId,
	}, bson.M{
		"$push": bson.M{
			"comments": dataStructure,
		},
	})

	if e == mgo.ErrNotFound {
		err = helper.NewCodeError(404, e)
	} else if e != nil {
		err = helper.NewCodeError(500, e)
	}
	return
}

func (db *DB) CallbackMemberCommentUpdate(queueId, memberId, commentId, data string) (err *helper.CodeError) {
	if !bson.IsObjectIdHex(queueId) {
		err = helper.NewCodeError(403, errors.New("Bad queueId"))
		return
	}
	if !bson.IsObjectIdHex(memberId) {
		err = helper.NewCodeError(403, errors.New("Bad memberId"))
		return
	}
	if !bson.IsObjectIdHex(commentId) {
		err = helper.NewCodeError(403, errors.New("Bad commentId"))
		return
	}

	e := db.db.C(COLLECTION_MEMBERS).Update(bson.M{
		"_id":   bson.ObjectIdHex(memberId),
		"queue": queueId,
		"comments": bson.M{
			"$elemMatch": bson.M{
				"_id": bson.ObjectIdHex(commentId),
			},
		},
	}, bson.M{
		"$set": bson.M{
			"comments.$.comment": data,
		},
	})

	if e == mgo.ErrNotFound {
		err = helper.NewCodeError(404, e)
	} else if e != nil {
		err = helper.NewCodeError(500, e)
	}
	return
}

func (db *DB) CallbackMemberCommentRemove(queueId, memberId, commentId string) (err *helper.CodeError) {
	if !bson.IsObjectIdHex(queueId) {
		err = helper.NewCodeError(403, errors.New("Bad queueId"))
		return
	}
	if !bson.IsObjectIdHex(memberId) {
		err = helper.NewCodeError(403, errors.New("Bad memberId"))
		return
	}
	if !bson.IsObjectIdHex(commentId) {
		err = helper.NewCodeError(403, errors.New("Bad commentId"))
		return
	}

	e := db.db.C(COLLECTION_MEMBERS).Update(bson.M{
		"_id":   bson.ObjectIdHex(memberId),
		"queue": queueId,
	}, bson.M{
		"$pull": bson.M{
			"comments": bson.M{
				"_id": bson.ObjectIdHex(commentId),
			},
		},
	})

	if e == mgo.ErrNotFound {
		err = helper.NewCodeError(404, e)
	} else if e != nil {
		err = helper.NewCodeError(500, e)
	}
	return
}

// region del
func (db *DB) FindCallTracking(query interface{}, dataStructure interface{}) (err error) {
	c := db.db.C(COLLECTION_CALL_TRACKING)
	err = c.Find(query).All(dataStructure)

	if err != nil {
		db.onError(err)
		return
	}

	return
}

func (db *DB) CreateGetCall(data interface{}) error {
	c := db.db.C(COLLECTION_CALL_TRACKING)
	return c.Insert(data)
}

func (db *DB) ItemGetCall(id string, dataStructure interface{}) (err error) {
	if !bson.IsObjectIdHex(id) {
		err = errors.New("Bad id") // TODO bad request
		return
	}

	c := db.db.C(COLLECTION_CALL_TRACKING)
	err = c.FindId(bson.ObjectIdHex(id)).One(dataStructure)

	if err == mgo.ErrNotFound {
		err = nil
		return
	}

	if err != nil {
		db.onError(err)
	}

	return
}

// endregion
