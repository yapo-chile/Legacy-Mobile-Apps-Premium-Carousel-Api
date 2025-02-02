package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/domain"
	"gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/usecases"
)

// adRepo implements the repository interface and gets ads from search repository
type adRepo struct {
	handler         Search
	regionsConf     Config
	imageServerLink string
	index           string
	maxAdsToDisplay int
}

// MakeAdRepository returns a fresh instance of AdRepository
func MakeAdRepository(handler Search, regionsConf Config, index,
	imageServerLink string, maxAdsToDisplay int) usecases.AdRepository {
	return &adRepo{
		handler:         handler,
		index:           index,
		imageServerLink: imageServerLink,
		regionsConf:     regionsConf,
		maxAdsToDisplay: maxAdsToDisplay,
	}
}

// GetUserAds gets user active ads from search repository using config to
// match similar ads
func (repo *adRepo) GetUserAds(userID int, productParams domain.ProductParams) (domain.Ads, error) {
	limit := repo.makeLimit(productParams)
	termQuery := repo.handler.NewTermQuery("userId", userID)
	must, mustNot := []Query{termQuery}, []Query{}

	if len(productParams.Categories) > 0 {
		must = append(must,
			repo.handler.NewCategoryFilter(productParams.Categories...))
	}

	if len(productParams.Keywords) > 0 {
		should := []Query{}
		for _, keyword := range productParams.Keywords {
			should = append(should,
				repo.handler.NewMultiMatchQuery(keyword, "cross_fields",
					"category.parentId", "category.id", "location.regionId", "location.communeId",
					"name", "body", "subject", "params.brand", "params.model",
					"params.type", "params.version"))
		}
		multiMatchBoolQuery := repo.handler.NewBoolQuery(
			[]Query{}, []Query{}, should)
		must = append(must, multiMatchBoolQuery)
	}

	if productParams.PriceRange > 0 {
		must = append(must,
			repo.handler.NewRangeQuery("price",
				productParams.PriceFrom, productParams.PriceTo))
	}
	if len(productParams.Exclude) > 0 {
		mustNot = append(mustNot,
			repo.handler.NewIDsQuery(productParams.Exclude...))
	}

	boolQuery := repo.handler.NewBoolQuery(must, mustNot, []Query{})
	scoreQuery := repo.handler.NewFunctionScoreQuery(boolQuery, 5, "multiply", true)
	result, err := repo.handler.Search(repo.index, scoreQuery, 0, limit)
	if err != nil {
		return domain.Ads{}, err
	}

	ads := repo.parseToAds(result.GetResults())
	if len(ads) < limit && productParams.FillGapsWithRandom {
		ads = repo.fillGapsWithRandom(userID, (limit - len(ads)), ads, productParams)
	}

	if len(ads) == 0 {
		return domain.Ads{}, fmt.Errorf("The specified "+
			"userID: %d don't return results elasticsearch",
			userID)
	}

	return ads, nil
}

// fillGapsWithRandom fill gaps in case of the limit is less than required ads by config.
// This method only works if config 'fillGapsWithRandom' is enabled
func (repo *adRepo) fillGapsWithRandom(userID int, delta int, ads domain.Ads,
	productParams domain.ProductParams) domain.Ads {
	exclude := []string{}
	for _, ad := range ads {
		exclude = append(exclude, ad.ID)
	}
	extraAds, _ := repo.GetUserAds(userID, domain.ProductParams{
		Exclude:            append(exclude, productParams.Exclude...),
		Categories:         productParams.Categories,
		FillGapsWithRandom: false,
		Limit:              delta,
	})
	for i, ad := range extraAds {
		ad.IsRelated = false
		extraAds[i] = ad
	}
	return repo.randomizePositions(append(ads, extraAds...))
}

// randomizePositions randomizes index in ads array
func (repo *adRepo) randomizePositions(ads domain.Ads) domain.Ads {
	for i := range ads {
		j := rand.Intn(i + 1)
		ads[i], ads[j] = ads[j], ads[i]
	}
	return ads
}

var notAlphaNumbericRegex, _ = regexp.Compile("[^a-zA-Z0-9]+")
var specialCases = strings.NewReplacer("á", "a", "é", "e", "í", "i", "ó", "o",
	"ú", "u", "'", "", "ñ", "n")

// parseToAds parses raw searchRepository response to domain object
func (repo *adRepo) parseToAds(results []json.RawMessage) (ads domain.Ads) {
	fmt.Printf("results %+v\n", results)
	log.Printf("results %+v\n", results)
	for _, hit := range results {
		result := usecases.Ad{}
		err := json.Unmarshal(hit, &result) // nolint
		log.Printf("result Unmarshal err %+v\n", err)
		fmt.Printf("result Unmarshal err %+v\n", err)
		ads = append(ads, repo.fillAd(result))
	}
	return
}

// fillAd parses search's document to domain.Ad struct
func (repo *adRepo) fillAd(result usecases.Ad) domain.Ad {
	regionKey := fmt.Sprintf("region.%d.link", result.Location.RegionID)
	regionName := repo.regionsConf.Get(regionKey)
	currency := "peso"
	if val, ok := result.Params["currency"]; ok {
		currency, ok = val.Value.(string)
		if !ok {
			currency = "peso"
		}
	}
	return domain.Ad{
		ID:         strconv.FormatInt(result.ListID, 10),
		UserID:     int(result.UserID),
		CategoryID: int(result.Category.ParentID),
		Subject:    result.Subject,
		Price:      result.Price,
		Currency:   currency,
		URL: "/" + strings.Join(
			[]string{
				notAlphaNumbericRegex.ReplaceAllString(
					specialCases.Replace(strings.ToLower(regionName)), "_"),
				notAlphaNumbericRegex.ReplaceAllString(
					specialCases.Replace(strings.ToLower(result.Subject)), "_") +
					"_" + strconv.FormatInt(result.ListID, 10),
			},
			"/",
		),
		IsRelated: true,
		Image:     repo.getMainImage(result.Media),
	}
}

// getMainImage gets the main image for required ad using media struct
func (repo *adRepo) getMainImage(imgs []usecases.AdMedia) domain.Image {
	if len(imgs) == 0 {
		return domain.Image{}
	}
	for _, img := range imgs {
		if img.SeqNo == 0 {
			return repo.fillImage(img.ID)
		}
	}
	return repo.fillImage(imgs[0].ID)
}

// fillImage parses the image ID to domain Image struct
func (repo *adRepo) fillImage(ID int) domain.Image {
	IDstr := fmt.Sprintf("%010d", ID)
	return domain.Image{
		Full:   fmt.Sprintf(repo.imageServerLink, "images", IDstr[:2], IDstr),
		Medium: fmt.Sprintf(repo.imageServerLink, "thumbsli", IDstr[:2], IDstr),
		Small:  fmt.Sprintf(repo.imageServerLink, "thumbs", IDstr[:2], IDstr),
	}
}

// makeLimit determines the real limit based on configuration
func (repo *adRepo) makeLimit(productParams domain.ProductParams) int {
	if productParams.Limit > 0 && productParams.Limit < repo.maxAdsToDisplay {
		return productParams.Limit
	}
	return repo.maxAdsToDisplay
}

// GetAd gets ad in search Repository using listID
func (repo *adRepo) GetAd(listID string) (domain.Ad, error) {
	termQuery := repo.handler.NewTermQuery("listId", listID)
	log.Printf("termQuery %s\n", termQuery)
	fmt.Printf("termQuery %s\n", termQuery)
	res, err := repo.handler.Search(repo.index, termQuery, 0, 10)
	log.Printf("Search res:%+v err:%+v\n", res, err)
	fmt.Printf("Search res:%+v err:%+v\n", res, err)
	if err != nil {
		return domain.Ad{}, err
	}
	ads := repo.parseToAds(res.GetResults())
	return ads[0], nil
}
