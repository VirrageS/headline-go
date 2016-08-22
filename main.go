package main

import (
	"time"

	// "github.com/iris-contrib/middleware/recovery"
	"github.com/kataras/iris"

	"github.com/VirrageS/cache"
	"github.com/VirrageS/headline-go/service"
)

func main() {
	// iris.Use(recovery.New())

	cache := cache.NewCache(time.Minute * 2)
	iris.UseFunc(func(c *iris.Context) {
		c.Set("cache", cache)
		c.Next()
	})

	iris.API("/github", service.GithubAPI{})
	iris.API("/hackerrank", service.HackerRankAPI{})
	iris.API("/reddit", service.RedditAPI{})
	iris.Listen(":8080")
}
