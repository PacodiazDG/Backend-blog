package services

import (
	"time"

	"github.com/PacodiazDG/Backend-blog/api/v1/blog"
	"github.com/PacodiazDG/Backend-blog/components/backup"
	"github.com/PacodiazDG/Backend-blog/components/logs"
	"github.com/PacodiazDG/Backend-blog/database"
	"github.com/fatih/color"
)

func AutoSetCacheTop() {
	ticker := time.NewTicker(10 * time.Minute)
	quit := make(chan struct{})
	blog.SetTopFeed()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				color.Red("[Svc] Critical error:  please check the error logs")
				logs.WriteLogs(r.(error), logs.CriticalError)
				AutoSetCacheTop()
			}
		}()
		for {
			select {
			case <-ticker.C:
				blog.ReflexCache()
				color.Yellow("[Svc] Updated Feed Cache")
			case <-quit:
				ticker.Stop()
				return
			}
		}

	}()
}

// Service to generate a backup at a certain time.
func ImagebackupService() {
	ticker := time.NewTicker(168 * time.Hour)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				backup.GenerateBackup("./Serverfiles")
				color.Green("[Svc] backup has been successfully created :" + (time.Now()).String())
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func RedistestConection() {
	color.Green("[Svc] Redis Service Registed :" + (time.Now()).String())
	ticker := time.NewTicker(10 * time.Minute)
	quit := make(chan struct{})
	go func() {
		defer func() {
			if r := recover(); r != nil {
				color.Red("[Svc] Critical error:  please check the error logs. RedistestConection")
				logs.WriteLogs(r.(error), logs.CriticalError)
				RedistestConection()
			}
		}()
		for {
			select {
			case <-ticker.C:
				err := database.RedisCon.Set("key", time.Now().Format(time.RFC850), 0).Err()
				if err != nil {
					color.Red("[Svc] Redis Conection test failded:" + (time.Now()).String())
					logs.WriteLogs(err, logs.CriticalError)
					panic(err)
				}
				color.Green("[Svc] Redis Conection test :" + (time.Now()).String())
			case <-quit:
				ticker.Stop()
				return
			}
		}

	}()
}
