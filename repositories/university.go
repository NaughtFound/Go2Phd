package repositories

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"naughtfound.github.io/go2phd/models"
)

type APIUniversity struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

type UniversityRepository struct {
	httpClient *http.Client
}

func NewUniversityRepository() *UniversityRepository {
	return &UniversityRepository{
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (r *UniversityRepository) GetByName(name string) (models.University, error) {
	apiURL := fmt.Sprintf("http://universities.hipolabs.com/search?name=%s", url.QueryEscape(name))

	resp, err := r.httpClient.Get(apiURL)
	if err != nil {
		return models.University{Name: name, Country: "Unknown"}, err
	}
	defer resp.Body.Close()

	var results []APIUniversity
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return models.University{Name: name, Country: "Unknown"}, err
	}

	if len(results) > 0 {
		return models.University{
			Name:    results[0].Name,
			Country: results[0].Country,
		}, nil
	}

	return models.University{
		Name:    name,
		Country: "Unknown",
	}, nil
}
