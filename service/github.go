package service

import (
    "fmt"
    "net/url"
)

type Repository struct {
	Name     string `json:"name,omitempty"`
	Url      string `json:"html_url,omitempty"`
	Stars    int `json:"stargazers_count,omitempty"`
	Watchers int `json:"watchers_count,omitempty"`
	Language string `json:"language,omitempty"`
}

type RepositorySearchResult struct {
	Total             int   `json:"total_count,omitempty"`
	IncompleteResults bool  `json:"incomplete_results,omitempty"`
	Repositories      []Repository `json:"items,omitempty"`
}

type Github struct {
    Url string
}

func (g Github) Get() (string, error) {
    trending, _ := url.Parse(fmt.Sprintf("%s/search/repositories?q=created:>2016-08-13&sort=stars&order=desc", g.Url))

    request, err := NewRequest("GET", trending)
	if err != nil {
		return "", err
	}

	response, err := Do(request)
	if err != nil {
		return "", err
	}

    result := new(RepositorySearchResult)
    Decode(response, result)

    return Encode(result.Repositories)
}
