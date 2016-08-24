package service

import (
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html"

	"github.com/kataras/iris"

	"github.com/VirrageS/cache"
	"github.com/VirrageS/scrape"
)

const (
	githubUrl = "https://github.com/trending"
)

type GithubAPI struct {
	*iris.Context
}

func (g GithubAPI) Get() {
	c := g.Context.Get("cache").(*cache.Cache)
	cached_items, ok := c.Get("github")
	if ok {
		g.JSON(iris.StatusOK, cached_items)
		return
	}

	response, err := http.Get(githubUrl)

	if err != nil {
		g.JSON(iris.StatusInternalServerError, iris.Map{
			"Error": "Could not make response",
		})
		return
	}

	root, err := html.Parse(response.Body)
	if err != nil {
		g.JSON(iris.StatusInternalServerError, iris.Map{
			"Error": "Could not parse body",
		})
		return
	}

	repos := scrape.Find(root, ".repo-list-item")

	headline := make([]HeadlineItem, 0)
	for _, repo := range repos[:maxItemsLimit] {
		// get url
		link := scrape.Find(repo, ".repo-list-name a")[0]
		url := "https://github.com" + scrape.Attr(link, "href")

		// get name
		name := scrape.Text(link)

		// get description
		desc := scrape.Find(repo, ".repo-list-description")
		description := ""
		if len(desc) > 0 {
			description = scrape.Text(desc[0])
		}

		// get stars
		meta := scrape.Find(repo, ".repo-list-meta")[0]

		re := regexp.MustCompile("[0-9,]+")
		points := re.FindAllString(scrape.Text(meta), -1)[0]
		stars, _ := strconv.Atoi(strings.Replace(points, ",", "", -1))

		headline = append(headline, HeadlineItem{Title: name, Description: description, Url: url, Points: stars})
	}

	sort.Sort(ByPoints(headline))
	c.Set("github", &headline)
	g.JSON(iris.StatusOK, &headline)
}
