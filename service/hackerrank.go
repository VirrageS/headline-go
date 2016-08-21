package service

import (
	"fmt"
	"net/url"
)

type HackerRankItem struct {
	Id        uint32 `json:"id,omitempty"`
	Title     string `json:"title,omitempty"`
	Url       string `json:"url,omitempty"`
	Points    uint32 `json:"score,omitempty"`
	CreatedAt uint32 `json:"time,omitempty"`
}

type HackerRank struct {
	Url string
}

func (h HackerRank) Get() (string, error) {
	trending, err := url.Parse(h.Url + "/topstories.json")
	if err != nil {
		return "", err
	}

	request, err := NewRequest("GET", trending)
	if err != nil {
		return "", err
	}

	response, err := Do(request)
	if err != nil {
		return "", err
	}

	result := new([]uint32)
	Decode(response, result)

	items := new([]HackerRankItem)
	for _, id := range (*result)[:10] {
		itemUrl, _ := url.Parse(fmt.Sprintf("%s/item/%v.json", h.Url, id))
		request, err := NewRequest("GET", itemUrl)
		if err != nil {
			return "", err
		}

		response, err := Do(request)
		if err != nil {
			return "", err
		}

		item := new(HackerRankItem)
		Decode(response, item)

		*items = append(*items, *item)
	}

	return Encode(items)
}
