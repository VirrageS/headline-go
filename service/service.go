package service

import (
	"encoding/json"
	"net/http"
	"net/url"
)

const (
	maxItemsLimit = 10
)

type HeadlineItem struct {
	Title		string `json:"title,omitempty"`
	Description	string `json:"description,omitempty"`
	Url 		string `json:"url,omitempty"`
	Points		int `json:"points,omitempty"`
}

type ByPoints []HeadlineItem

func (p ByPoints) Len() int           { return len(p) }
func (p ByPoints) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ByPoints) Less(i, j int) bool { return p[i].Points > p[j].Points }


func newRequest(method string, url *url.URL) (*http.Request, error) {
	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Headline/1.0")
	return req, nil
}

func do(req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(req)
	return resp, err
}

func decode(resp *http.Response, target interface{}) {
	body := resp.Body
	defer body.Close()

	json.NewDecoder(body).Decode(target)
}

func encode(source interface{}) (string, error) {
	bytes, err := json.Marshal(source)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
