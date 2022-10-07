package services

import (
	"time"

	"github.com/PacodiazDG/Backend-blog/api/v1/blog"
	"github.com/PacodiazDG/Backend-blog/modules/logs"
	"github.com/fatih/color"
)

func AutoSetCacheTop() {
	ticker := time.NewTicker(10 * time.Minute)
	quit := make(chan struct{})
	blog.SetTopFeed()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				color.Red("Critical Error:")
				logs.WriteLogs(r.(error), logs.CriticalError)
				AutoSetCacheTop()
			}
		}()
		for {
			select {
			case <-ticker.C:
				blog.ReflexCache()
				color.Yellow("Updated Feed Cache")
			case <-quit:
				ticker.Stop()
				return
			}
		}

	}()
}
