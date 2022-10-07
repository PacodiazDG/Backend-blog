package redisbackend

import (
	"encoding/json"

	"github.com/PacodiazDG/Backend-blog/database"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InsertFeedCache Get if any token is banned
func CheckBan(id, idtoken string) bool {
	_, err := database.RedisCon.Get("IDBaned" + id).Result()
	if err != redis.Nil {
		return true
	}
	_, err = database.RedisCon.Get("Tokenid" + idtoken).Result()
	return err != redis.Nil
}

// SetBan Inserta baneo a redis
func SetBan(Info User) error {
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

// SetBan Inserta baneo a redis
func RemoveBan(id primitive.ObjectID) error {
	err := database.RedisCon.Del(id.Hex()).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetFeedCache Returns the feed cache values
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

// InsertFeedCache Get if any token is banned
func InsertBanidtoken(idtoken string) error {
	return database.RedisCon.Set("Tokenid"+idtoken, "Ban", 0).Err()
}
