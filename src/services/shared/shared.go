package shared

import (
	"../../db"
	"time"
)

var DB *db.DB

func init() {
	DB = db.NewDB("")
}

func CurrentTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
