package service

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// TODO: create some kind of global struct to manage results
// with fields something like in hackerrank

func newRequest(method string, url *url.URL) (*http.Request, error) {
	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
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
