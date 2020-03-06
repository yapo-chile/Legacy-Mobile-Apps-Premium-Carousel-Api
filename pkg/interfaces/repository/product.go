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
	logger         ProductRepositoryLogger
}

// ProductRepositoryLogger logs product repository events
type ProductRepositoryLogger interface {
	LogWarnPartialConfigNotSupported(name, value string)
}

// MakeProductRepository creates a new instance of ProductRepository
func MakeProductRepository(handler DbHandler, resultsPerPage int,
	logger ProductRepositoryLogger) usecases.ProductRepository {
	return &productRepo{
		handler:        handler,
		resultsPerPage: resultsPerPage,
		logger:         logger,
	}
}

// GetUserProductsTotal get the total of user products
func (repo *productRepo) GetUserProductsTotal() (total int) {
	result, err := repo.handler.Query(`SELECT COUNT(*) as total FROM
		user_product`)
	if err != nil {
		return 0
	}
	defer result.Close()
	if result.Next() {
		result.Scan(&total)
	}
	return total
}

// GetUserProductsTotal get the total of user products
func (repo *productRepo) GetUserProductsTotalByEmail(email string) (total int) {
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
func (repo *productRepo) GetUserProducts(
	page int) (products []usecases.Product, currentPage int,
	totalPages int, err error) {
	if page < 1 {
		page = 1
	}
	total := repo.GetUserProductsTotal()
	if total < 1 {
		return []usecases.Product{}, page, 0, nil
	}
	totalPages = (total / repo.resultsPerPage)
	if (total % repo.resultsPerPage) > 0 {
		totalPages++
	}
	result, err := repo.handler.Query(
		`SELECT
		p.id, p.product_type, p.user_id, p.user_email, p.status, p.expired_at,
		p.created_at,
		ARRAY(
			SELECT user_product_config.name || '=' || user_product_config.value
			FROM user_product_config WHERE user_product_id = p.id
		) AS config_params,
		p.comment
		FROM user_product as p
		WHERE TRUE
		ORDER BY p.id DESC
		OFFSET $1 LIMIT $2`,
		(repo.resultsPerPage * (page - 1)), repo.resultsPerPage)
	if err != nil {
		return []usecases.Product{}, 0, 0, err
	}
	defer result.Close()
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

// GetUserProducts get a list of user products by email with pagination
func (repo *productRepo) GetUserProductsByEmail(email string,
	page int) (products []usecases.Product, currentPage int,
	totalPages int, err error) {
	if page < 1 {
		page = 1
	}
	total := repo.GetUserProductsTotalByEmail(email)
	if total < 1 {
		return []usecases.Product{}, page, 0, nil
	}
	totalPages = (total / repo.resultsPerPage)
	if (total % repo.resultsPerPage) > 0 {
		totalPages++
	}
	result, err := repo.handler.Query(
		`SELECT
		p.id, p.product_type, p.user_id, p.user_email, p.status, p.expired_at,
		p.created_at,
		ARRAY(
			SELECT user_product_config.name || '=' || user_product_config.value
			FROM user_product_config WHERE user_product_id = p.id
		) AS config_params,
		p.comment
		FROM user_product as p
		WHERE user_email = $1
		ORDER BY p.id DESC
		OFFSET $2 LIMIT $3`,
		email, (repo.resultsPerPage * (page - 1)), repo.resultsPerPage)
	if err != nil {
		return []usecases.Product{}, 0, 0, err
	}
	defer result.Close()
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

// GetUserActiveProduct gets active product for an specific userProductID
func (repo *productRepo) GetUserProductByID(userProductID int) (usecases.Product, error) {
	result, err := repo.handler.Query(`SELECT
	p.id, p.product_type, p.user_id, p.user_email, p.status, p.expired_at,
	p.created_at,
	ARRAY(
		SELECT user_product_config.name || '=' || user_product_config.value
		FROM user_product_config WHERE user_product_id = id
	) AS config_params,
	p.comment
	FROM user_product as p
	WHERE  p.id = $1`,
		userProductID)
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
	var userProductID int
	var createdAt time.Time
	if result.Next() {
		result.Scan(&userProductID, &createdAt)
	} else {
		return usecases.Product{},
			fmt.Errorf("next error: getting userProductID from database")
	}
	err = repo.SetConfig(userProductID, config)
	if err != nil {
		return usecases.Product{}, err
	}
	return usecases.Product{
		ID:        userProductID,
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

// SetConfig adds configuration to Product
func (repo *productRepo) SetConfig(userProductID int, config usecases.CpConfig) error {
	values := makeConfigValues(userProductID, config)
	insertValues, positions := []interface{}{}, []string{}
	counter := 0
	for _, v := range values {
		temp := []string{}
		for range v {
			counter++
			temp = append(temp, fmt.Sprintf("$%d", counter))
		}
		positions = append(positions, "("+strings.Join(temp, ",")+")")
		insertValues = append(insertValues, []interface{}{v[0], v[1], v[2]}...)
	}
	return repo.handler.Insert(
		fmt.Sprintf(`INSERT INTO user_product_config(user_product_id, name, value) VALUES %s
			ON CONFLICT (user_product_id, name) DO UPDATE set value=excluded.value`,
			strings.Join(positions, ", ")),
		insertValues...)
}

func makeConfigValues(userProductID int, config usecases.CpConfig) (values [][]interface{}) {
	return [][]interface{}{
		[]interface{}{userProductID, "categories", strings.Trim(strings.Join(
			strings.Fields(fmt.Sprint(config.Categories)), ","), "[]")},
		[]interface{}{userProductID, "limit", strconv.Itoa(config.Limit)},
		[]interface{}{userProductID, "custom_query", config.CustomQuery},
		[]interface{}{userProductID, "exclude", strings.Join(config.Exclude, ",")},
		[]interface{}{userProductID, "price_range", strconv.Itoa(config.PriceRange)},
		[]interface{}{userProductID, "gaps_with_random", fmt.Sprintf("%t", config.FillGapsWithRandom)},
	}
}

// SetConfig adds configuration to Product
func (repo *productRepo) SetPartialConfig(userProductID int, configMap map[string]interface{}) error {
	for name, value := range configMap {
		switch name {
		case "status":
			if err := repo.SetStatus(userProductID,
				usecases.ProductStatus(value.(string))); err != nil {
				return err
			}
		default:
			repo.logger.LogWarnPartialConfigNotSupported(name,
				fmt.Sprintf("%+v", value))
		}
	}
	return nil
}

// SetStatus sets the user product status
func (repo *productRepo) SetStatus(userProductID int, status usecases.ProductStatus) error {
	result, err := repo.handler.
		Query(
			`UPDATE user_product SET status=$1 WHERE id=$2`,
			status,
			userProductID,
		)
	if err != nil {
		return err
	}
	return result.Close()
}

// SetStatus sets the user product status
func (repo *productRepo) SetExpiration(userProductID int, expiredAt time.Time) error {
	result, err := repo.handler.
		Query(
			`UPDATE user_product SET expired_at=$1 WHERE id=$2`,
			expiredAt,
			userProductID,
		)
	if err != nil {
		return err
	}
	return result.Close()
}
