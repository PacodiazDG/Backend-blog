package RedisBackend

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id" `
	Blocked      bool               `bson:"Blocked"`
	Reason       string             `bson:"Reason"`
	Visible      bool               `bson:"Visible"`
	Date         time.Time          `bson:"Date"`
	LoginAttempt int                `bson:"LoginAttempt"`
	Details      string             `bson:"Details"`
}
