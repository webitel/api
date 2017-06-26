package shared

import "../../db"

var DB *db.DB

func init() {
	DB = db.NewDB("")
}
