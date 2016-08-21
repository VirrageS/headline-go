package service

import (
    "net/url"

    "github.com/kataras/iris"
)

const (
    redditUrl = "https://www.reddit.com/r/programming"
)

type RedditItem struct{
    Title  string `json:"title,omitempty"`
    Url    string `json:"url,omitempty"`
    Points int `json:"score,omitempty"`
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

    items := new([]RedditItem)
    for _, c := range result.Data.Children {
        *items = append(*items, c.Item)
    }

    r.JSON(iris.StatusOK, items)
}
