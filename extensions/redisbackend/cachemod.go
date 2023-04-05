package redisbackend

import (
	"encoding/json"

	"github.com/PacodiazDG/Backend-blog/database"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Get if any token is banned
func CheckBan(id, idtoken string) bool {
	_, err := database.RedisCon.Get(id).Result()
	if err != redis.Nil {
		return true
	}
	_, err = database.RedisCon.Get(idtoken).Result()
	return err != redis.Nil
}

// Insert ban to user in redisdb
func SetBan(Info UserRedisJson) error {
	JsonUser, err := json.Marshal(Info)
	if err != nil {
		return err
	}
	err = database.RedisCon.Set(Info.ID.Hex(), JsonUser, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func SetBanToken(Token, details string) error {
	err := database.RedisCon.Set(Token, details, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// remove a user's ban on redis
func RemoveBan(id primitive.ObjectID) error {
	err := database.RedisCon.Del(id.Hex()).Err()
	if err != nil {
		return err
	}
	return nil
}

// Returns the feed cache values
func GetFeedCache(IDfeed string) string {
	val, err := database.RedisCon.Get("feed" + IDfeed).Result()
	if err != nil {
		panic(err)
	}
	return val
}

// InsertFeedCache
func InsertFeedCache(data, value string) {
	err := database.RedisCon.Set("feed"+(value), data, 0).Err()
	if err != nil {
		panic(err)
	}
}

// Get if any token is banned
func InsertBanidtoken(idtoken string) error {
	return database.RedisCon.Set("Tokenid"+idtoken, "Ban", 0).Err()
}
