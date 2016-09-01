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
	productHuntUrl = "https://www.producthunt.com/"
)

type ProductHuntAPI struct {
	*iris.Context
}

func (p ProductHuntAPI) Get() {
	c := p.Context.Get("cache").(*cache.Cache)
	cached_items, ok := c.Get("producthunt")
	if ok {
		p.JSON(iris.StatusOK, cached_items)
		return
	}

	response, err := http.Get(productHuntUrl)

	if err != nil {
		p.JSON(iris.StatusInternalServerError, iris.Map{
			"Error": "Could not make response",
		})
		return
	}

	root, err := html.Parse(response.Body)
	if err != nil {
		p.JSON(iris.StatusInternalServerError, iris.Map{
			"Error": "Could not parse body",
		})
		return
	}

	headline := make([]HeadlineItem, 0)
	products := scrape.Find(scrape.Find(root, "ul")[0], "li")

	for _, product := range products {
		// get url
		link := scrape.Find(product, "div a")[0]
		url := "https://producthunt.com" + scrape.Attr(link, "href")

		// get name
		info := scrape.Find(link, "div span")
		name := scrape.Text(info[0])

		// skip products with empty name
		if name == "" {
			continue
		}

		// get description
		description := scrape.Text(info[1])

		// get stars
		meta := scrape.Find(product, "button div")[0]

		re := regexp.MustCompile("[0-9,]+")
		number := re.FindAllString(scrape.Text(meta), -1)[0]
		points, _ := strconv.Atoi(strings.Replace(number, ",", "", -1))

		headline = append(headline, HeadlineItem{Title: name, Description: description, Url: url, Points: points})

		if len(headline) >= maxItemsLimit {
			break
		}
	}

	sort.Sort(ByPoints(headline))
	if len(headline) > 0 {
		c.Set("producthunt", &headline)
	}
	p.JSON(iris.StatusOK, &headline)
}
