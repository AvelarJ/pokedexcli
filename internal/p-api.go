package internal

import (
	"encoding/json"
	"net/http"
)

type LocationArea struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func FetchLocationAreas(url string) (LocationArea, error) {
	res, err := http.Get(url)
	if err != nil {
		return LocationArea{}, err
	}
	defer res.Body.Close()

	locationArea := LocationArea{}
	err = json.NewDecoder(res.Body).Decode(&locationArea)
	if err != nil {
		return LocationArea{}, err
	}
	return locationArea, nil
}
