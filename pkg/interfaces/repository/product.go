package repository

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"
)

// productRepo holds connections to get user products
type productRepo struct {
	handler        DbHandler
	resultsPerPage int
}

// MakeProductRepository creates a new instance of ProductRepository
func MakeProductRepository(handler DbHandler, resultsPerPage int) usecases.ProductRepository {
	return &productRepo{
		handler:        handler,
		resultsPerPage: resultsPerPage,
	}
}

// GetUserProductsTotal get the total of user products
func (repo *productRepo) GetUserProductsTotal(email string) (total int) {
	result, err := repo.handler.Query(`SELECT COUNT(*) as total FROM
	user_product WHERE user_email=$1`, email)
	if err != nil {
		return 0
	}
	defer result.Close()
	if result.Next() {
		result.Scan(&total)
	}
	return total
}

// GetUserProducts get a list of user products with pagination
func (repo *productRepo) GetUserProducts(email string,
	page int) ([]usecases.Product, int, int, error) {
	if page < 1 {
		page = 1
	}
	total := repo.GetUserProductsTotal(email)
	if total < 1 {
		return []usecases.Product{}, page, 0, nil
	}
	totalPages := (total / repo.resultsPerPage)
	if (total % repo.resultsPerPage) > 0 {
		totalPages++
	}
	result, err := repo.handler.Query(`
	SELECT
		p.id, p.product_type, p.user_id, p.user_email, p.status, p.expired_at,
		p.created_at,
		ARRAY(
			SELECT user_product_config.name || '=' || user_product_config.value
			FROM user_product_config WHERE user_product_id = p.id
		) AS config_params,
		p.comment
	FROM user_product as p
	WHERE  user_email = $1 order by p.expired_at desc OFFSET $2 LIMIT $3`,
		email, ((repo.resultsPerPage - 1) * page), repo.resultsPerPage)
	if err != nil {
		return []usecases.Product{}, 0, 0, err
	}
	defer result.Close()
	products := []usecases.Product{}
	for result.Next() {
		product := usecases.Product{}
		rawConfig := []string{}
		result.Scan(&product.ID, &product.Type, &product.UserID, &product.Email,
			&product.Status, &product.ExpiredAt, &product.CreatedAt,
			(*pq.StringArray)(&rawConfig), &product.Comment)
		config, _ := repo.parseConfig(rawConfig)
		product.Config = config
		products = append(products, product)
	}
	return products, page, totalPages, nil
}

// GetUserActiveProduct gets active product for an specific userID
func (repo *productRepo) GetUserActiveProduct(userID string,
	productType usecases.ProductType) (usecases.Product, error) {
	result, err := repo.handler.Query(`SELECT
	p.id, p.product_type, p.user_id, p.user_email, p.status, p.expired_at,
	p.created_at,
	ARRAY(
		SELECT user_product_config.name || '=' || user_product_config.value
		FROM user_product_config WHERE user_product_id = id
	) AS config_params,
	p.comment
	FROM user_product as p
	WHERE  p.status = 'ACTIVE' AND p.user_id = $1 and p.product_type = $2`,
		userID, productType)
	if err != nil {
		return usecases.Product{}, err
	}
	defer result.Close()
	product := usecases.Product{}
	var configArr []string
	if result.Next() {
		result.Scan(&product.ID, &product.Type, &product.UserID, &product.Email,
			&product.Status, &product.ExpiredAt, &product.CreatedAt,
			(*pq.StringArray)(&configArr), &product.Comment)
	}

	config, err := repo.parseConfig(configArr)
	if err != nil {
		return usecases.Product{}, err
	}
	product.Config = config
	return product, nil
}

// parseConfig parses rawConfiguration slice to usecases.CpConfig  struct
func (repo *productRepo) parseConfig(rawConfig []string) (usecases.CpConfig, error) {
	if len(rawConfig) == 0 {
		return usecases.CpConfig{}, fmt.Errorf("No config found")
	}
	configs := make(map[string]string)
	for _, v := range rawConfig {
		pair := strings.Split(v, "=")
		if len(pair) >= 2 {
			configs[pair[0]] = strings.Join(pair[1:], "=")
		}
	}
	limit, _ := strconv.Atoi(configs["limit"])
	customQuery := configs["custom_query"]
	exclude := []string{}
	if configs["exclude"] != "" {
		exclude = strings.Split(configs["exclude"], ",")
	}
	priceRange, _ := strconv.Atoi(configs["price_range"])
	gapsWithRandom, _ := strconv.ParseBool(configs["gaps_with_random"])
	categoriesArrayStr := strings.Split(configs["categories"], ",")
	categories := []int{}
	for _, v := range categoriesArrayStr {
		cat, _ := strconv.Atoi(v)
		if cat >= 1000 && cat <= 9999 {
			categories = append(categories, cat)
		}
	}
	return usecases.CpConfig{
		Categories:         categories,
		Limit:              limit,
		CustomQuery:        customQuery,
		Exclude:            exclude,
		PriceRange:         priceRange,
		FillGapsWithRandom: gapsWithRandom,
	}, nil
}

// AddUserProduct add a new product to user
func (repo *productRepo) AddUserProduct(userID, email, comment string,
	productType usecases.ProductType, expiredAt time.Time,
	config usecases.CpConfig) (usecases.Product, error) {
	userIDint, err := strconv.Atoi(userID)
	if err != nil {
		return usecases.Product{}, err
	}
	result, err := repo.handler.Query(
		`INSERT INTO user_product(product_type, status, user_id, user_email,
			expired_at, comment)
			VALUES (
				$1, 'ACTIVE', $2, $3, $4, $5
			) RETURNING id, created_at`, productType, userIDint, email, expiredAt, comment)
	if err != nil {
		return usecases.Product{}, err
	}
	defer result.Close()
	var productID int
	var createdAt time.Time
	if result.Next() {
		result.Scan(&productID, &createdAt)
	} else {
		return usecases.Product{},
			fmt.Errorf("next error: getting productID from database")
	}
	err = repo.addConfig(productID, config)
	if err != nil {
		return usecases.Product{}, err
	}
	return usecases.Product{
		ID:        productID,
		Type:      productType,
		Email:     email,
		UserID:    userID,
		ExpiredAt: expiredAt,
		CreatedAt: createdAt,
		Comment:   comment,
		Config:    config,
		Status:    usecases.ActiveProduct,
	}, nil
}

// addConfig adds configuration to Product
func (repo *productRepo) addConfig(productID int, config usecases.CpConfig) error {
	params := [][]string{
		[]string{"($1, $2, $3)", "categories", strings.Trim(strings.Join(
			strings.Fields(fmt.Sprint(config.Categories)), ","), "[]")},
		[]string{"($4, $5, $6)", "limit", strconv.Itoa(config.Limit)},
		[]string{"($7, $8, $9)", "custom_query", config.CustomQuery},
		[]string{"($10, $11, $12)", "exclude", strings.Join(config.Exclude, ",")},
		[]string{"($13, $14, $15)", "price_range", strconv.Itoa(config.PriceRange)},
		[]string{"($16, $17, $18)", "gaps_with_random", fmt.Sprintf("%t", config.FillGapsWithRandom)},
	}
	values, positions := []interface{}{}, []string{}
	for _, v := range params {
		positions = append(positions, v[0])
		values = append(values, []interface{}{productID, v[1], v[2]}...)
	}

	return repo.handler.Insert(
		fmt.Sprintf(`INSERT INTO user_product_config(user_product_id, name, value) VALUES %s`,
			strings.Join(positions, ", ")),
		values...)
}
