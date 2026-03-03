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

type ExploreItem struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

// Pokemon struct holding all basic info (Could be expanded further if needed)
type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
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

func (c *Cache) FetchExploreItem(url string) (ExploreItem, error) {
	// Will fetch location area dat similar to above method but for explore command.
	// Will check cache first, if not in cache then make api request
	var exploreItem ExploreItem

	// Check cache first
	if cachedData, found := c.Get(url); found {
		err := json.Unmarshal(cachedData, &exploreItem)
		if err == nil {
			return exploreItem, nil
		} else {
			return ExploreItem{}, err
		}
	} else {
		// Make api get request because it wasn't found in cache
		res, err := http.Get(url)
		if err != nil {
			return ExploreItem{}, err
		}
		defer res.Body.Close()

		err = json.NewDecoder(res.Body).Decode(&exploreItem)
		if err != nil {
			return ExploreItem{}, err
		}
		// Marshal the location area data and add it to the cache for future requests
		data, err := json.Marshal(exploreItem)
		if err != nil {
			return ExploreItem{}, err
		}
		c.Add(url, data)
	}

	return exploreItem, nil
}

func (c *Cache) FetchPokemon(url string) (Pokemon, error) {
	// Will fetch pokemon data similar to above method but on the pokemon endpoint, used for catch command.
	var pokemon Pokemon

	// Check cache first
	if cachedData, found := c.Get(url); found {
		err := json.Unmarshal(cachedData, &pokemon)
		if err == nil {
			return pokemon, nil
		} else {
			return Pokemon{}, err
		}
	} else {
		// Make api get request because it wasn't found in cache
		res, err := http.Get(url)
		if err != nil {
			return Pokemon{}, err
		}
		defer res.Body.Close()

		err = json.NewDecoder(res.Body).Decode(&pokemon)
		if err != nil {
			return Pokemon{}, err
		}
		// Marshal the location area data and add it to the cache for future requests
		data, err := json.Marshal(pokemon)
		if err != nil {
			return Pokemon{}, err
		}
		c.Add(url, data)
	}

	return pokemon, nil
}
