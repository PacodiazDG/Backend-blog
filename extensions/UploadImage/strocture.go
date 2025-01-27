package uploadimage

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// /
type FileUserRegister struct {
	UserId bson.ObjectId
	Time   time.Time
	Size   int
	Ref    bson.ObjectId
}
