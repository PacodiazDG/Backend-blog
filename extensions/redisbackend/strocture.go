package redisbackend

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User-specific structure used in redis
type UserRedisJson struct {
	ID           primitive.ObjectID `bson:"_id" `
	Blocked      bool               `bson:"Blocked"`
	Reason       string             `bson:"Reason"`
	LoginAttempt int                `bson:"LoginAttempt"`
	Details      string             `bson:"Details"`
}

type TokenBan struct {
	Reason  string `bson:"Reason"`
	Details string `bson:"Details"`
}
