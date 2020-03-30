package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"

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

func (e *elasticsearch) GetDoc(index string, id string) (json.RawMessage, error) {
	res, err := e.client.Get().
		Index(index).
		Id(id).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	if !res.Found {
		return nil, fmt.Errorf("%s not found in elasticsearch", id)
	}
	return res.Source, nil
}

func (e *elasticsearch) NewMultiMatchQuery(text interface{}, typ string,
	fields ...string) repository.Query {
	return elastic.NewMultiMatchQuery(text, fields...).
		Type(typ)
}

func (e *elasticsearch) NewTermQuery(name string, value interface{}) repository.Query {
	return elastic.NewTermQuery(name, value)
}

func (e *elasticsearch) NewRangeQuery(name string, from, to int) repository.Query {
	q := elastic.NewRangeQuery(name)
	if from > 0 {
		q = q.Gte(from)
	}
	if to > 0 {
		q = q.Lte(to)
	}
	return q
}

func (e *elasticsearch) NewBoolQuery(must []repository.Query,
	mustNot []repository.Query, should []repository.Query) repository.Query {
	inputMust := []elastic.Query{}
	for _, q := range must {
		inputMust = append(inputMust, q.(elastic.Query))
	}
	inputMustNot := []elastic.Query{}
	for _, q := range mustNot {
		inputMustNot = append(inputMustNot, q.(elastic.Query))
	}
	inputShould := []elastic.Query{}
	for _, q := range should {
		inputShould = append(inputShould, q.(elastic.Query))
	}
	return elastic.NewBoolQuery().Must(inputMust...).MustNot(inputMustNot...).
		Should(inputShould...)
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

func (e *elasticsearch) NewIDsQuery(ids ...string) repository.Query {
	return elastic.NewIdsQuery().Ids(ids...)
}

func (e *elasticsearch) NewCategoryFilter(categoryIDs ...int) repository.Query {
	inputShould := []elastic.Query{}
	for _, cat := range categoryIDs {
		if cat > 9999 || cat < 1000 {
			continue
		}
		if (cat % 1000) == 0 {
			inputShould = append(inputShould,
				elastic.NewRangeQuery("CategoryID").Gte(cat).Lt(cat+1000))
		} else {
			inputShould = append(inputShould,
				elastic.NewTermQuery("CategoryID", cat))

		}
	}
	return elastic.NewBoolQuery().Should(inputShould...)
}

type searchResult elastic.SearchResult

func (r *searchResult) GetResults() (results []json.RawMessage) {
	for _, hit := range r.Hits.Hits {
		results = append(results, hit.Source)
	}
	return
}

func (r *searchResult) TotalHits() int64 {
	if r.Hits != nil && r.Hits.TotalHits != nil {
		return r.Hits.TotalHits.Value
	}
	return 0
}
