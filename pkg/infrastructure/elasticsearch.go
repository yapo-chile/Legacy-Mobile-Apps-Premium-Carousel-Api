package infrastructure

import (
	"context"
	"encoding/json"

	"github.com/olivere/elastic/v7"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/interfaces/loggers"
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/interfaces/repository"
)

type elasticsearch struct {
	client *elastic.Client
	logger loggers.Logger
}

func NewElasticsearch(host, port string, logger loggers.Logger) *elasticsearch {
	client, _ := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(host+":"+port),
	)
	esversion, err := client.ElasticsearchVersion(host + ":" + port)
	if err != nil {
		logger.Error("Error connecting to elasticsearch")
		return nil
	}
	logger.Info("Connected to elasticsearch version: %s", esversion)
	return &elasticsearch{
		client: client,
		logger: logger,
	}
}

func (e *elasticsearch) Search(index string,
	query repository.Query, from,
	size int) (repository.SearchResult, error) {
	res, err := e.client.Search().
		Index(index).
		Query(query).
		From(from).Size(size).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	result := searchResult(*res)
	return &result, nil
}

func (e *elasticsearch) NewMultiMatchQuery(text interface{}, typ string, fields ...string) repository.Query {
	return elastic.NewMultiMatchQuery(text, fields...).
		Type(typ)
}

func (e *elasticsearch) NewTermQuery(name string, value interface{}) repository.Query {
	return elastic.NewTermQuery(name, value)
}

func (e *elasticsearch) NewBoolQueryMust(query ...repository.Query) repository.Query {
	input := []elastic.Query{}
	for _, q := range query {
		input = append(input, q.(elastic.Query))
	}
	return elastic.NewBoolQuery().Must(input...)
}

func (e *elasticsearch) NewFunctionScoreQuery(query repository.Query, boost float64,
	boostMode string, random bool) repository.Query {
	q := elastic.NewFunctionScoreQuery().
		Query(query).Boost(boost).BoostMode(boostMode)
	if random {
		randomFunction := elastic.NewRandomFunction()
		q = q.AddScoreFunc(randomFunction)
	}
	return q
}

type searchResult elastic.SearchResult

func (r *searchResult) GetResults() (results []json.RawMessage) {
	for _, hit := range r.Hits.Hits {
		results = append(results, hit.Source)
	}
	return
}

func (r *searchResult) TotalHits() int64 {
	return r.Hits.TotalHits.Value
}
