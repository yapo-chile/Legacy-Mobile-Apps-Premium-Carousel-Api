package repository

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
)

type mockSearch struct {
	mock.Mock
}

func (m *mockSearch) NewMultiMatchQuery(text interface{},
	typ string, fields ...string) Query {
	args := m.Called(text, typ, fields)
	return args.Get(0).(Query)
}

func (m *mockSearch) NewTermQuery(name string, value interface{}) Query {
	args := m.Called(name, value)
	return args.Get(0).(Query)
}

func (m *mockSearch) NewRangeQuery(name string, from, to int) Query {
	args := m.Called(name, from, to)
	return args.Get(0).(Query)
}

func (m *mockSearch) NewFunctionScoreQuery(query Query, boost float64,
	boostMode string, random bool) Query {
	args := m.Called(query, boost, boostMode, random)
	return args.Get(0).(Query)
}

func (m *mockSearch) NewBoolQuery(must []Query, mustNot []Query,
	should []Query) Query {
	args := m.Called(must, mustNot, should)
	return args.Get(0).(Query)
}

func (m *mockSearch) NewIDsQuery(ids ...string) Query {
	args := m.Called(ids)
	return args.Get(0).(Query)
}

func (m *mockSearch) NewCategoryFilter(categoryIDs ...int) Query {
	args := m.Called(categoryIDs)
	return args.Get(0).(Query)
}

func (m *mockSearch) GetDoc(index string, id string) (json.RawMessage, error) {
	args := m.Called(index, id)
	return args.Get(0).(json.RawMessage), args.Error(1)
}

func (m *mockSearch) Search(index string, query Query, from,
	size int) (SearchResult, error) {
	args := m.Called(index, query, from, size)
	return args.Get(0).(SearchResult), args.Error(1)
}

type mockSearchResult struct {
	mock.Mock
}

func (m *mockSearchResult) GetResults() (results []json.RawMessage) {
	args := m.Called()
	return args.Get(0).([]json.RawMessage)
}

func (m *mockSearchResult) TotalHits() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

type mockQuery struct {
	mock.Mock
}

func (m *mockQuery) Source() (interface{}, error) {
	args := m.Called()
	return args.Get(0), args.Error(1)
}

type mockConfig struct {
	mock.Mock
}

func (m *mockConfig) Get(s string) string {
	args := m.Called(s)
	return args.String(0)
}

func TestMakeAdRepository(t *testing.T) {
	mSearch := &mockSearch{}
	mConfig := &mockConfig{}
	expected := adRepo{
		handler:     mSearch,
		regionsConf: mConfig,
	}
	result := MakeAdRepository(mSearch, mConfig, "", "", 0)
	assert.Equal(t, &expected, result)
	mSearch.AssertExpectations(t)
	mConfig.AssertExpectations(t)
}

func TestGetUserAdsOK(t *testing.T) {
	mSearch := &mockSearch{}
	mQuery := &mockQuery{}
	mResults := &mockSearchResult{}
	mConfig := &mockConfig{}

	mSearch.On("NewTermQuery", mock.AnythingOfType("string"),
		mock.Anything).Return(mQuery)

	mSearch.On("NewCategoryFilter", mock.AnythingOfType("[]int")).Return(mQuery)

	mSearch.On("NewMultiMatchQuery",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]string"),
	).Return(mQuery)
	mSearch.On("NewIDsQuery",
		mock.AnythingOfType("[]string"),
	).Return(mQuery)
	mSearch.On("NewRangeQuery",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("int"),
		mock.AnythingOfType("int"),
	).Return(mQuery)
	mSearch.On("NewBoolQuery", mock.Anything, mock.Anything,
		mock.Anything).Return(mQuery)
	mSearch.On("NewFunctionScoreQuery",
		mock.Anything,
		mock.AnythingOfType("float64"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("bool")).Return(mQuery)

	results := []json.RawMessage{
		[]byte(`{"ListID": 123, "UserID": 2, "CategoryID": 2020, "Subject": "Autito"}`),
	}

	mSearch.On("Search", mock.AnythingOfType("string"),
		mock.Anything,
		mock.AnythingOfType("int"),
		mock.AnythingOfType("int")).Return(mResults, nil)

	mResults.On("GetResults").Return(results)
	mConfig.On("Get", mock.AnythingOfType("string")).Return("something")
	interactor := adRepo{
		handler:     mSearch,
		regionsConf: mConfig,
	}

	userAds, err := interactor.GetUserAds(0,
		domain.ProductParams{
			Categories: []int{1234, 2345},
			Exclude:    []string{"123"},
			Keywords:   []string{"key1"},
			PriceRange: 1,
		})

	expected := domain.Ads{
		{ID: "123", UserID: 2, CategoryID: 2020,
			Subject: "Autito", URL: "/something/autito_123", IsRelated: true},
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, userAds)
	mSearch.AssertExpectations(t)
	mConfig.AssertExpectations(t)
	mResults.AssertExpectations(t)
	mQuery.AssertExpectations(t)
}

func TestGetUserAdsWithFilledGaps(t *testing.T) {
	mSearch := &mockSearch{}
	mQuery := &mockQuery{}
	mResults := &mockSearchResult{}
	mConfig := &mockConfig{}

	mSearch.On("NewTermQuery", mock.AnythingOfType("string"),
		mock.Anything).Return(mQuery)

	mSearch.On("NewCategoryFilter", mock.AnythingOfType("[]int")).Return(mQuery)

	mSearch.On("NewMultiMatchQuery",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]string"),
	).Return(mQuery)
	mSearch.On("NewIDsQuery",
		mock.AnythingOfType("[]string"),
	).Return(mQuery)
	mSearch.On("NewRangeQuery",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("int"),
		mock.AnythingOfType("int"),
	).Return(mQuery)
	mSearch.On("NewBoolQuery", mock.Anything, mock.Anything,
		mock.Anything).Return(mQuery)
	mSearch.On("NewFunctionScoreQuery",
		mock.Anything,
		mock.AnythingOfType("float64"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("bool")).Return(mQuery)
	results1 := []json.RawMessage{
		[]byte(`{"ListID": 1234, "UserID": 2, "CategoryID": 2020, "Subject": "Autito"}`),
	}
	results2 := []json.RawMessage{
		[]byte(`{"ListID": 123, "UserID": 2, "CategoryID": 2020, "Subject": "Autito"}`),
	}

	mSearch.On("Search", mock.AnythingOfType("string"),
		mock.Anything,
		mock.AnythingOfType("int"),
		mock.AnythingOfType("int")).Return(mResults, nil)
	mResults.On("GetResults").Return(results1).Once()
	mResults.On("GetResults").Return(results2).Once()

	mConfig.On("Get", mock.AnythingOfType("string")).Return("something")
	interactor := adRepo{
		handler:         mSearch,
		regionsConf:     mConfig,
		maxAdsToDisplay: 20,
	}

	userAds, err := interactor.GetUserAds(0,
		domain.ProductParams{
			Categories:         []int{1234, 2345},
			Exclude:            []string{"123"},
			Keywords:           []string{"key1"},
			PriceRange:         1,
			FillGapsWithRandom: true,
			Limit:              2,
		})

	expected := domain.Ads{
		{ID: "1234", UserID: 2, CategoryID: 2020,
			Subject: "Autito", URL: "/something/autito_1234", IsRelated: true},
		{ID: "123", UserID: 2, CategoryID: 2020,
			Subject: "Autito", URL: "/something/autito_123", IsRelated: false},
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, userAds)
	mSearch.AssertExpectations(t)
	mConfig.AssertExpectations(t)
	mResults.AssertExpectations(t)
	mQuery.AssertExpectations(t)
}

func TestGetUserAdsZeroResults(t *testing.T) {
	mSearch := &mockSearch{}
	mQuery := &mockQuery{}
	mResults := &mockSearchResult{}
	mConfig := &mockConfig{}

	mSearch.On("NewTermQuery", mock.AnythingOfType("string"),
		mock.Anything).Return(mQuery)

	mSearch.On("NewCategoryFilter", mock.AnythingOfType("[]int")).Return(mQuery)

	mSearch.On("NewMultiMatchQuery",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]string"),
	).Return(mQuery)
	mSearch.On("NewIDsQuery",
		mock.AnythingOfType("[]string"),
	).Return(mQuery)
	mSearch.On("NewRangeQuery",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("int"),
		mock.AnythingOfType("int"),
	).Return(mQuery)
	mSearch.On("NewBoolQuery", mock.Anything, mock.Anything,
		mock.Anything).Return(mQuery)
	mSearch.On("NewFunctionScoreQuery",
		mock.Anything,
		mock.AnythingOfType("float64"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("bool")).Return(mQuery)

	results := []json.RawMessage{}

	mSearch.On("Search", mock.AnythingOfType("string"),
		mock.Anything,
		mock.AnythingOfType("int"),
		mock.AnythingOfType("int")).Return(mResults, nil)

	mResults.On("GetResults").Return(results)
	interactor := adRepo{
		handler:         mSearch,
		regionsConf:     mConfig,
		maxAdsToDisplay: 20,
	}

	_, err := interactor.GetUserAds(0,
		domain.ProductParams{
			Categories: []int{1234, 2345},
			Exclude:    []string{"123"},
			Keywords:   []string{"key1"},
			PriceRange: 1,
		})

	assert.Error(t, err)
	mSearch.AssertExpectations(t)
	mConfig.AssertExpectations(t)
	mResults.AssertExpectations(t)
	mQuery.AssertExpectations(t)
}

func TestGetUserAdsSearchError(t *testing.T) {
	mSearch := &mockSearch{}
	mQuery := &mockQuery{}
	mResults := &mockSearchResult{}
	mConfig := &mockConfig{}

	mSearch.On("NewTermQuery", mock.AnythingOfType("string"),
		mock.Anything).Return(mQuery)

	mSearch.On("NewCategoryFilter", mock.AnythingOfType("[]int")).Return(mQuery)

	mSearch.On("NewMultiMatchQuery",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]string"),
	).Return(mQuery)
	mSearch.On("NewIDsQuery",
		mock.AnythingOfType("[]string"),
	).Return(mQuery)
	mSearch.On("NewRangeQuery",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("int"),
		mock.AnythingOfType("int"),
	).Return(mQuery)
	mSearch.On("NewBoolQuery", mock.Anything, mock.Anything,
		mock.Anything).Return(mQuery)
	mSearch.On("NewFunctionScoreQuery",
		mock.Anything,
		mock.AnythingOfType("float64"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("bool")).Return(mQuery)

	mSearch.On("Search", mock.AnythingOfType("string"),
		mock.Anything,
		mock.AnythingOfType("int"),
		mock.AnythingOfType("int")).Return(mResults, fmt.Errorf("e"))

	interactor := adRepo{
		handler:     mSearch,
		regionsConf: mConfig,
	}

	_, err := interactor.GetUserAds(0,
		domain.ProductParams{
			Categories: []int{1234, 2345},
			Exclude:    []string{"123"},
			Keywords:   []string{"key1"},
			PriceRange: 1,
		})

	assert.Error(t, err)
	mSearch.AssertExpectations(t)
	mConfig.AssertExpectations(t)
	mResults.AssertExpectations(t)
	mQuery.AssertExpectations(t)
}

func TestGetMainImage(t *testing.T) {
	interactor := adRepo{
		imageServerLink: "http://img.yapo.cl/%s/%s/%s.jpg",
	}

	image := interactor.getMainImage([]Media{{SeqNo: 0, ID: 1}})
	expected := domain.Image{
		Full:   "http://img.yapo.cl/images/00/0000000001.jpg",
		Medium: "http://img.yapo.cl/thumbsli/00/0000000001.jpg",
		Small:  "http://img.yapo.cl/thumbs/00/0000000001.jpg",
	}
	assert.Equal(t, expected, image)

	image = interactor.getMainImage([]Media{{SeqNo: 1, ID: 100}})
	expected = domain.Image{
		Full:   "http://img.yapo.cl/images/00/0000000100.jpg",
		Medium: "http://img.yapo.cl/thumbsli/00/0000000100.jpg",
		Small:  "http://img.yapo.cl/thumbs/00/0000000100.jpg",
	}
	assert.Equal(t, expected, image)

	image = interactor.getMainImage([]Media{})
	expected = domain.Image{}
	assert.Equal(t, expected, image)
}

func TestGetAdOK(t *testing.T) {
	mSearch := &mockSearch{}
	mConfig := &mockConfig{}
	interactor := adRepo{
		handler:         mSearch,
		regionsConf:     mConfig,
		maxAdsToDisplay: 20,
	}
	var result json.RawMessage
	result = []byte(`{"ListID": 123,
	 "UserID": 2, "CategoryID": 2020, "Subject": "Autito"}`)
	mSearch.On("GetDoc",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
	).Return(result, nil)
	mConfig.On("Get", mock.AnythingOfType("string")).Return("something")

	userAds, err := interactor.GetAd("123")

	expected := domain.Ad{ID: "123", UserID: 2, CategoryID: 2020,
		Subject: "Autito", URL: "/something/autito_123", IsRelated: true}
	assert.NoError(t, err)
	assert.Equal(t, expected, userAds)
	mConfig.AssertExpectations(t)
	mSearch.AssertExpectations(t)
}

func TestGetAdGetDocError(t *testing.T) {
	mSearch := &mockSearch{}
	mConfig := &mockConfig{}
	interactor := adRepo{
		handler:     mSearch,
		regionsConf: mConfig,
	}
	var result json.RawMessage
	result = []byte(`{"ListID": 123,
	 "UserID": 2, "CategoryID": 2020, "Subject": "Autito"}`)
	mSearch.On("GetDoc",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
	).Return(result, fmt.Errorf("e"))

	_, err := interactor.GetAd("123")

	assert.Error(t, err)
	mConfig.AssertExpectations(t)
	mSearch.AssertExpectations(t)
}

func TestGetAdUnmarshalError(t *testing.T) {
	mSearch := &mockSearch{}
	mConfig := &mockConfig{}
	interactor := adRepo{
		handler:     mSearch,
		regionsConf: mConfig,
	}
	var result json.RawMessage
	result = []byte(`{"ListID": 123,
	 "UserID": 2, "CategoryID": 2020, "Subject": "`)
	mSearch.On("GetDoc",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
	).Return(result, nil)

	_, err := interactor.GetAd("123")

	assert.Error(t, err)
	mConfig.AssertExpectations(t)
	mSearch.AssertExpectations(t)
}
