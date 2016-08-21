package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Service interface {
	Get() (string, error)
}

// TODO: create some kind of global struct to manage results
// with fields something like in hackerrank

func NewRequest(method string, url *url.URL) (*http.Request, error) {
	fmt.Printf("%s\n", url.String())
	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	return req, nil
}

func Do(req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(req)
	return resp, err
}

func Decode(resp *http.Response, target interface{}) {
	body := resp.Body
	defer body.Close()

	json.NewDecoder(body).Decode(target)
}

func Encode(source interface{}) (string, error) {
	bytes, err := json.Marshal(source)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
