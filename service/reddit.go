package service

import (
	"net/url"
	"sort"

	"github.com/kataras/iris"

	"github.com/VirrageS/cache"
)

const (
	redditUrl = "https://www.reddit.com/r/programming"
)

type RedditItem struct{
	Title	string `json:"title,omitempty"`
	Url		string `json:"url,omitempty"`
	Points	int `json:"score,omitempty"`
}

func (r *RedditItem) toHeadlineItem() *HeadlineItem {
	return &HeadlineItem{
		Title: r.Title,
		Description: "",
		Url: r.Url,
		Points: r.Points,
	}
}

type RedditResultChildren struct {
	Kind string `json:"kind,omitempty"`
	Item RedditItem `json:"data,omitempty"`
}

type RedditResult struct {
	Kind string `json:"kind,omitempty"`
	Data struct {
		Modhash string `json:"modhash,omitempty"`
		Children []RedditResultChildren `json:"children,omitempty"`
	} `json:"data,omitempty"`
}

type RedditAPI struct {
	*iris.Context
}

func (r RedditAPI) Get() {
	c := r.Context.Get("cache").(*cache.Cache)
	cached_items, ok := c.Get("reddit")
	if ok {
		r.JSON(iris.StatusOK, cached_items)
		return
	}

	trending, _ := url.Parse(redditUrl + "/hot.json")
	request, err := newRequest("GET", trending)
	if err != nil {
		r.JSON(iris.StatusInternalServerError, iris.Map{
			"Error": "Could not make request",
		})
		return
	}

	response, err := do(request)
	if err != nil {
		r.JSON(iris.StatusInternalServerError, iris.Map{
			"Error": "Could not do request",
		})
		return
	}

	result := new(RedditResult)
	decode(response, result)

	limit := len(result.Data.Children)
	if limit > maxItemsLimit {
		limit = maxItemsLimit
	}

	headline := make([]HeadlineItem, 0)
	for _, c := range result.Data.Children[:limit] {
		headline = append(headline, *c.Item.toHeadlineItem())
	}

	sort.Sort(ByPoints(headline))
	if len(headline) > 0 {
		c.Set("reddit", &headline)
	}

	r.JSON(iris.StatusOK, &headline)
}
