package blog

import (
	"github.com/PacodiazDG/Backend-blog/extensions/redisbackend"
	logs "github.com/PacodiazDG/Backend-blog/modules/logs"
)

type Top50 struct {
	Top []StoryStruct
}

var Blogs = InitControllerPost()

// Stores the feed in a variable allowing for faster access
var FastFeed []FeedStrcture

// Stores the last 50 most viewed stories allowing for quicker access
var StoryCacheVar *[]StoryStruct

// gets if the token is blacklisted from some database
func TokenBlackList(token, idtoken string) bool {
	return redisbackend.CheckBan(token, idtoken)
}

// Update the top 50 most viewed stories
func SetLastStories() {
	Blogs.SetCollection("Post")
	info, err := Blogs.SetTop()
	if err != nil {
		logs.WriteLogs(err, logs.HardError)
		panic(err)
	}
	StoryCacheVar = &info
}

// Update the fed
func SetTopFeed() {
	Blogs.SetCollection("Post")
	info, err := Blogs.FeedFast()
	if err != nil {
		logs.WriteLogs(err, logs.HardError)
		return
	}
	FastFeed = info
}

// Updates stories and fed
func ReflexCache() {
	SetTopFeed()
	SetLastStories()
}
