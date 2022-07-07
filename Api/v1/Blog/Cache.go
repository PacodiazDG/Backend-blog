package Blog

import (
	"github.com/PacodiazDG/Backend-blog/Extensions/RedisMercy"
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

//TokenBlackList gets if the token is blacklisted from some database
func TokenBlackList(token, idtoken string) bool {
	return RedisMercy.CheckBan(token, idtoken)
}

// Actualizar el top de los post mas vistos
func SetTop() {
	Blogs.SetCollection("Post")
	info, err := Blogs.SetTop()
	if err != nil {
		panic(err)
	}
	CacheRamPost = &info
}

//SetFastFeed
func SetFastFeed() {

	Blogs.SetCollection("Post")
	info, err := Blogs.FeedFast()
	if err != nil {
		Logs.WriteLogs(err)
		return
	}
	FastFeed = info
}

func ReflexCache() {
	SetFastFeed()
	SetTop()
}
