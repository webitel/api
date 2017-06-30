package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kataras/iris/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

type Request struct {
	Filter  map[string]interface{}
	Limit   int
	Page    int
	Sort    []string
	Columns map[string]int8
}

func MakeId() bson.ObjectId {
	return bson.NewObjectId()
}

func (r *Request) GetQueue(db *DB, collectionName string, domainName string) *mgo.Query {

	if domainName != "" {
		r.Filter["domain"] = domainName
	}

	q := db.db.C(collectionName).Find(r.Filter)

	if len(r.Columns) > 0 {
		q.Select(r.Columns)
	}

	if len(r.Sort) > 0 {
		q.Sort(r.Sort...)
	}

	if r.Limit < 1 {
		r.Limit = 40
	}

	if r.Page > 0 {
		q = q.Skip((r.Page - 1) * r.Limit)
	}

	return q.Limit(r.Limit)
}

func RequestListFromUri(ctx context.Context) *Request {
	r := &Request{}
	if ctx.FormValue("limit") != "" {
		r.Limit, _ = strconv.Atoi(ctx.FormValue("limit"))
	}
	if ctx.FormValue("page") != "" {
		r.Page, _ = strconv.Atoi(ctx.FormValue("page"))
	}

	if ctx.FormValue("filter") != "" {
		err := json.Unmarshal([]byte(ctx.FormValue("filter")), &r.Filter)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		r.Filter = make(map[string]interface{})
	}
	if ctx.FormValue("domain") != "" {
		r.Filter["domain"] = ctx.FormValue("domain")
	}

	if ctx.FormValue("columns") != "" {
		err := json.Unmarshal([]byte(ctx.FormValue("columns")), &r.Columns)
		if err != nil {
			fmt.Println(err)
		}
	}

	//fmt.Println(r)
	return r
}

func (_ *DB) ParsePGArray(array string) ([]string, error) {
	var out []string
	var arrayOpened, quoteOpened, escapeOpened bool
	item := &bytes.Buffer{}
	for _, r := range array {
		switch {
		case !arrayOpened:
			if r != '{' {
				return nil, errors.New("Doesn't appear to be a postgres array.  Doesn't start with an opening curly brace.")
			}
			arrayOpened = true
		case escapeOpened:
			item.WriteRune(r)
			escapeOpened = false
		case quoteOpened:
			switch r {
			case '\\':
				escapeOpened = true
			case '"':
				quoteOpened = false
				if item.String() == "NULL" {
					item.Reset()
				}
			default:
				item.WriteRune(r)
			}
		case r == '}':
			// done
			out = append(out, item.String())
			return out, nil
		case r == '"':
			quoteOpened = true
		case r == ',':
			// end of item
			out = append(out, item.String())
			item.Reset()
		default:
			item.WriteRune(r)
		}
	}
	return nil, errors.New("Doesn't appear to be a postgres array.  Premature end of string.")
}
