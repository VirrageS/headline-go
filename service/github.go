package service

import (
	"net/http"
	"regexp"
	"strconv"

	"golang.org/x/net/html"

	"github.com/kataras/iris"

	"github.com/VirrageS/scrape"
)

const (
	githubUrl = "https://github.com/trending"
)

type Repository struct {
	Title	   string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Url		 string `json:"url,omitempty"`
	Points	  int `json:"points,omitempty"`
}

type GithubAPI struct {
	*iris.Context
}

func (g GithubAPI) Get() {
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

	repositories := new([]Repository)
	for _, repo := range repos {
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

		re := regexp.MustCompile("[0-9]+")
		numbers := re.FindAllString(scrape.Text(meta), -1)
		stars, _ := strconv.Atoi(numbers[0])

		*repositories = append(*repositories, Repository{Title: name, Description: description, Url: url, Points: stars})
	}

	g.JSON(iris.StatusOK, repositories)
}
