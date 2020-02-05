package repository

import (
	"fmt"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"
)

// AdRepo implements the repository interface and gets ads from search-ms
type AdRepo struct {
	handler         Elasticsearch
	regionsConf     Config
	path            string
	maxAdsToDisplay int
}

// MakeAdRepository returns a fresh instance of AdRepo
func MakeAdRepository(handler Elasticsearch, regionsConf Config, path string, maxAdsToDisplay int) usecases.AdRepository {
	return &AdRepo{
		handler:         handler,
		path:            path,
		regionsConf:     regionsConf,
		maxAdsToDisplay: maxAdsToDisplay,
	}
}

// SearchInput object to recieve search-ms input data type
type SearchInput map[string]interface{}

// GetUserAds gets user active ads list from search-ms. The pagination starts
// from page 1, also page 0 means page 1
func (repo *AdRepo) GetUserAds(userID string, cpConfig usecases.CpConfig) (domain.Ads, error) {
	termQuery := repo.handler.NewTermQuery("UserID", userID)
	multiMatchQuery := repo.handler.NewMultiMatchQuery("hola", "cross_fields", "Category^0.5",
		"SubCategory^0.5", "Region^0.5", "Commune", "name", "Body",
		"Subject^2", "Params.Brand", "Params.Model",
		"Params.Type", "Params.Version")

	boolQuery := repo.handler.NewBoolQueryMust(termQuery, multiMatchQuery)
	scoreQuery := repo.handler.NewFunctionScoreQuery(boolQuery, 5, "multiply", true)
	result, err := repo.handler.Search("ads", scoreQuery, 0, 10)
	if err != nil {
		panic(err)
	}

	if result.TotalHits() > 0 {
		fmt.Printf("Found a total of %d ads\n", result.TotalHits())
		for _, hit := range result.GetResults() {
			fmt.Printf("Result: %+v\n", string(hit))
		}
	} else {
		// No hits
		fmt.Print("Ads not found\n")
	}
	return domain.Ads{}, nil
}
