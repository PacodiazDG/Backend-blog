package Blog

import (
	"github.com/PacodiazDG/Backend-blog/Extensions/RedisBackend"
	"github.com/PacodiazDG/Backend-blog/Modules/Logs"
)

type Top50 struct {
	Top []PostSimpleStruct
}

var Blogs = InitControllerPost()

// Canche para el primer feed
var FastFeed []FeedStrcture

// Variable global para cacheRam
var CacheRamPost *[]PostSimpleStruct

// TokenBlackList gets if the token is blacklisted from some database
func TokenBlackList(token, idtoken string) bool {
	return RedisBackend.CheckBan(token, idtoken)
}

// Actualizar el top de los post mas vistos
func SetTopPost() {
	Blogs.SetCollection("Post")
	info, err := Blogs.SetTop()
	if err != nil {
		Logs.WriteLogs(err, Logs.HardError)
		panic(err)
	}
	CacheRamPost = &info
}

// SetFastFeed
func SetTopFeed() {
	Blogs.SetCollection("Post")
	info, err := Blogs.FeedFast()
	if err != nil {
		Logs.WriteLogs(err, Logs.HardError)
		return
	}
	FastFeed = info
}

func ReflexCache() {
	SetTopFeed()
	SetTopPost()
}
