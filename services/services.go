package services

import (
	"time"

	"github.com/PacodiazDG/Backend-blog/api/v1/blog"
	"github.com/PacodiazDG/Backend-blog/modules/backup"
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

func ImagebackupService() {
	ticker := time.NewTicker(168 * time.Hour)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				backup.GenerateCompress()
				color.Green("backup has been successfully created :" + (time.Now()).String())
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
