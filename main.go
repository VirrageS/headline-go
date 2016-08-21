package main

import (
	"github.com/iris-contrib/middleware/recovery"
	"github.com/kataras/iris"

	"github.com/VirrageS/headline-go/service"
)

func main() {
	iris.Use(recovery.New())

	iris.API("/github", service.GithubAPI{})
	iris.API("/hackerrank", service.HackerRankAPI{})
	iris.API("/reddit", service.RedditAPI{})
	iris.Listen(":8080")
}
