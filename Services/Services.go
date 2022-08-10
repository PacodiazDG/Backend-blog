package services

import (
	"time"

	"github.com/PacodiazDG/Backend-blog/Api/v1/Blog"
	"github.com/PacodiazDG/Backend-blog/Modules/Logs"
	"github.com/fatih/color"
)

func AutoSetCacheTop() {
	ticker := time.NewTicker(10 * time.Minute)
	quit := make(chan struct{})
	Blog.SetFastFeed()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				color.Red("Critical Error:")
				Logs.WriteLogs(r.(error))
				AutoSetCacheTop()
			}
		}()
		for {
			select {
			case <-ticker.C:
				Blog.ReflexCache()
				color.Yellow("Updated Feed Cache")
			case <-quit:
				ticker.Stop()
				return
			}
		}

	}()
}
