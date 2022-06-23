package RedisMercy

import (
	database "github.com/PacodiazDG/Backend-blog/Database"
	"github.com/go-redis/redis"
)

// InsertFeedCache Get if any token is banned
func MidelwareBan(id, idtoken string) bool {
	_, err := database.RedisCon.Get("IDBaned" + id).Result()
	if err != redis.Nil {
		return true
	}
	_, err = database.RedisCon.Get("Tokenid" + idtoken).Result()
	return err != redis.Nil
}

//SetBan Inserta baneo a redis
func SetBan(id string) error {
	err := database.RedisCon.Set("IDBaned"+id, "True", 0).Err()
	if err != nil {
		return err
	}
	return nil
}

//SetBan Inserta baneo a redis
func RemoveBan(id string) error {
	err := database.RedisCon.Del(("IDBaned" + id)).Err()
	if err != nil {
		return err
	}
	return nil
}

//GetFeedCache Returns the feed cache values
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
