package blog

import (
	"github.com/PacodiazDG/Backend-blog/extensions/redisbackend"
	logs "github.com/PacodiazDG/Backend-blog/modules/logs"
)

type Top50 struct {
	Top []PostStruct
}

var Blogs = InitControllerPost()

// Canche para el primer feed
var FastFeed []FeedStrcture

// Variable global para cacheRam
var CacheRamPost *[]PostStruct

// TokenBlackList gets if the token is blacklisted from some database
func TokenBlackList(token, idtoken string) bool {
	return redisbackend.CheckBan(token, idtoken)
}

// Actualizar el top de los post mas vistos
func SetTopPost() {
	Blogs.SetCollection("Post")
	info, err := Blogs.SetTop()
	if err != nil {
		logs.WriteLogs(err, logs.HardError)
		panic(err)
	}
	CacheRamPost = &info
}

// SetFastFeed
func SetTopFeed() {
	Blogs.SetCollection("Post")
	info, err := Blogs.FeedFast()
	if err != nil {
		logs.WriteLogs(err, logs.HardError)
		return
	}
	FastFeed = info
}

func ReflexCache() {
	SetTopFeed()
	SetTopPost()
}
