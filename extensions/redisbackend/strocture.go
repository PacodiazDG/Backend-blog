package redisbackend

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User-specific structure used in redis
type UserRedisJson struct {
	ID           primitive.ObjectID `bson:"_id" `
	Blocked      bool               `bson:"Blocked"`
	Reason       string             `bson:"Reason"`
	Visible      bool               `bson:"Visible"`
	Date         time.Time          `bson:"Date"`
	LoginAttempt int                `bson:"LoginAttempt"`
	Details      string             `bson:"Details"`
}
