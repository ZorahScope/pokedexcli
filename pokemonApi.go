package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type apiResponse interface {
	locationAreaList
}

type locationAreaList struct {
	Count    int64    `json:"count"`
	Next     string   `json:"next"`
	Previous *string  `json:"previous"`
	Results  []Result `json:"results"`
}

type Result struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func getFromAPI[T apiResponse](url string) (T, error) {
	var result T

	res, err := http.Get(url)
	if err != nil {
		return result, fmt.Errorf("error making http request: %v", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return result, fmt.Errorf("error reading response body: %v", err)
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("error unmarshalling JSON: %v", err)
	}
	return result, nil
}
