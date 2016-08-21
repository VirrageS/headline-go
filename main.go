package main

import (
	"fmt"

	"github.com/VirrageS/headline/service"
)

const (
	githubUrl = "https://api.github.com"
	mediumUrl = "https://api.medium.com"
	hackerRankUrl = "https://hacker-news.firebaseio.com/v0"
)

func main() {
	var hackerRank service.Service = service.HackerRank{Url: hackerRankUrl}
	items, err := hackerRank.Get()
	fmt.Print(err)
	fmt.Printf("%s\n", items)

	var github service.Service = service.Github{Url: githubUrl}
	items, err = github.Get()
	fmt.Print(err)
	fmt.Printf("%s\n", items)
}
