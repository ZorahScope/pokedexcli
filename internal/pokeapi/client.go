package pokeapi

import (
	"encoding/json"
	"fmt"
	"github.com/zorahscope/pokedexcli/internal/pokecache"
	"io"
	"net/http"
	"time"
)

var cache = pokecache.NewCache(time.Minute * 15)

func GetFromAPI[T apiResponse](url string) (T, error) {
	var result T

	data, err := getRawData(url)
	if err != nil {
		return result, err
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("error unmarshalling JSON: %v", err)
	}
	return result, nil
}

func getRawData(url string) ([]byte, error) {
	var cachedData []byte
	var ok bool

	cachedData, ok = cache.Get(url)
	if ok {
		return cachedData, nil
	}
	res, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("error making http request: %v", err)
	}
	if res.StatusCode >= 400 {
		return []byte{}, fmt.Errorf("%v not found", res.StatusCode)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("error reading response body: %v", err)
	}
	cache.Add(url, data)
	return data, nil
}
