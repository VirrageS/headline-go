package service

import (
	"fmt"
	"net/url"

	"github.com/kataras/iris"

	"github.com/VirrageS/cache"
)

const (
	hackerRankUrl = "https://hacker-news.firebaseio.com/v0"
	hackerRankLimit = 10
)

type HackerRankItem struct {
	Title     string `json:"title,omitempty"`
	Url       string `json:"url,omitempty"`
	Points    uint32 `json:"score,omitempty"`
}

type HackerRankAPI struct {
	*iris.Context
}

func (h HackerRankAPI) Get() {
	c := h.Context.Get("cache").(*cache.Cache)
	cached_items, ok := c.Get("hackerrank")
	if ok {
		h.JSON(iris.StatusOK, cached_items)
		return
	}

	trending, _ := url.Parse(hackerRankUrl + "/topstories.json")
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

	items := make([]HackerRankItem, 0)
	for _, id := range (*result)[:hackerRankLimit] {
		itemUrl, _ := url.Parse(fmt.Sprintf("%s/item/%v.json", hackerRankUrl, id))
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

		item := new(HackerRankItem)
		decode(response, item)

		items = append(items, *item)
	}

	c.Set("hackerrank", &items)
	h.JSON(iris.StatusOK, &items)
}
