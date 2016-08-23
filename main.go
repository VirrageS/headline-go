package main

import (
	"time"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"

	"github.com/VirrageS/cache"
	"github.com/VirrageS/headline-go/service"
)

func main() {
	// this should not be in production...
	crs := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
	})
	iris.Use(crs)

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
