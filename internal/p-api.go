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

func (c *Cache) FetchLocationAreas(url string) (LocationArea, error) {
	// Starts by checking cache for the requested URL. If the data is found in the cache and is still valid (not expired), it returns the cached data.
	// If the data is not in the cache or has expired, it makes an HTTP GET request to the specified URL to fetch the location area data.
	// After receiving the response, it decodes the JSON data into a LocationArea struct and returns it.

	var locationArea LocationArea

	// Check cache first
	if cachedData, found := c.Get(url); found {

		err := json.Unmarshal(cachedData, &locationArea)
		if err == nil {
			return locationArea, nil
		}

	} else {
		// Make api get request if not in cache
		res, err := http.Get(url)

		if err != nil {
			return LocationArea{}, err
		}
		defer res.Body.Close()

		err = json.NewDecoder(res.Body).Decode(&locationArea)
		if err != nil {
			return LocationArea{}, err
		}
		// Marshal the location area data and add it to the cache for future requests
		data, err := json.Marshal(locationArea)
		if err != nil {
			return LocationArea{}, err
		}
		c.Add(url, data)

	}

	return locationArea, nil
}
