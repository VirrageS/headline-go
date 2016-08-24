package service

import (
	"fmt"
	"net/url"
	"sort"

	"github.com/kataras/iris"

	"github.com/VirrageS/cache"
)

const (
	hackerNewsUrl = "https://hacker-news.firebaseio.com/v0"
)

type HackerNewsItem struct {
	Title     string `json:"title,omitempty"`
	Url       string `json:"url,omitempty"`
	Points    int `json:"score,omitempty"`
}

func (h *HackerNewsItem) toHeadlineItem() *HeadlineItem {
	return &HeadlineItem{
		Title: h.Title,
		Description: "",
		Url: h.Url,
		Points: h.Points,
	}
}

type HackerNewsAPI struct {
	*iris.Context
}

func (h HackerNewsAPI) Get() {
	c := h.Context.Get("cache").(*cache.Cache)
	cached_items, ok := c.Get("hackernews")
	if ok {
		h.JSON(iris.StatusOK, cached_items)
		return
	}

	trending, _ := url.Parse(hackerNewsUrl + "/topstories.json")
	request, err := newRequest("GET", trending)
	if err != nil {
		h.JSON(iris.StatusInternalServerError, iris.Map{
			"Error": "Could not make request",
		})
		return
	}

	response, err := do(request)
	if err != nil {
		h.JSON(iris.StatusInternalServerError, iris.Map{
			"Error": "Could not do request",
		})
		return
	}

	result := new([]uint32)
	decode(response, result)

	items := make([]HackerNewsItem, 0)
	for _, id := range (*result)[:maxItemsLimit] {
		itemUrl, _ := url.Parse(fmt.Sprintf("%s/item/%v.json", hackerNewsUrl, id))
		request, err := newRequest("GET", itemUrl)
		if err != nil {
			h.JSON(iris.StatusInternalServerError, iris.Map{
				"Error": "Could not make request",
			})
			return
		}

		response, err := do(request)
		if err != nil {
			h.JSON(iris.StatusInternalServerError, iris.Map{
				"Error": "Could not do request",
			})
			return
		}

		item := new(HackerNewsItem)
		decode(response, item)

		items = append(items, *item)
	}

	headline := make([]HeadlineItem, 0)
	for _, item := range items {
		headline = append(headline, *item.toHeadlineItem())
	}

	sort.Sort(ByPoints(headline))
	c.Set("hackernews", &headline)
	h.JSON(iris.StatusOK, &headline)
}
