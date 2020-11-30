package repository

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
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
	page int) (products []domain.Product, currentPage int,
	totalPages int, err error) {
	if page < 1 {
		page = 1
	}
	total := repo.GetUserProductsTotal()
	if total < 1 {
		return []domain.Product{}, page, 0, nil
	}
	totalPages = (total / repo.resultsPerPage)
	if (total % repo.resultsPerPage) > 0 {
		totalPages++
	}
	result, err := repo.makeUserProductQuery(`
		WHERE TRUE
		ORDER BY p.id DESC
		OFFSET $1 LIMIT $2`,
		(repo.resultsPerPage * (page - 1)), repo.resultsPerPage)
	if err != nil {
		return []domain.Product{}, 0, 0, err
	}
	defer result.Close()
	for result.Next() {
		product := domain.Product{}
		rawConfig := []string{}
		result.Scan(&product.ID, &product.Type, &product.UserID, &product.Email,
			&product.Status, &product.ExpiredAt, &product.CreatedAt,
			&product.Purchase.ID, &product.Purchase.Number, &product.Purchase.Type,
			&product.Purchase.Status, &product.Purchase.Price, &product.Purchase.CreatedAt,
			(*pq.StringArray)(&rawConfig))
		config, _ := repo.parseConfig(rawConfig)
		product.Config = config
		products = append(products, product)
	}
	return products, page, totalPages, nil
}

// GetReport gets sales report using interval between start date and end date.
// Returns sold products between given interval.
func (repo *productRepo) GetReport(startDate,
	endDate time.Time) (soldProducts []domain.Product, err error) {
	result, err := repo.makeUserProductQuery(`
		WHERE  p.created_at BETWEEN $1 AND $2
		ORDER BY p.id DESC`, startDate, endDate)
	if err != nil {
		return []domain.Product{}, err
	}
	defer result.Close()
	for result.Next() {
		product := domain.Product{}
		rawConfig := []string{}
		result.Scan(&product.ID, &product.Type, &product.UserID, &product.Email,
			&product.Status, &product.ExpiredAt, &product.CreatedAt,
			&product.Purchase.ID, &product.Purchase.Number, &product.Purchase.Type,
			&product.Purchase.Status, &product.Purchase.Price, &product.Purchase.CreatedAt,
			(*pq.StringArray)(&rawConfig))
		config, _ := repo.parseConfig(rawConfig)
		product.Config = config
		soldProducts = append(soldProducts, product)
	}
	return soldProducts, nil
}

func (repo *productRepo) makeUserProductQuery(conditions string,
	params ...interface{}) (DbResult, error) {
	return repo.handler.Query(`
		SELECT
			p.id, p.product_type, p.user_id, p.user_email, p.status, p.expired_at,
			p.created_at, pur.id, pur.purchase_number, pur.purchase_type,
			pur.purchase_status, pur.price, pur.created_at,
			ARRAY(
				SELECT user_product_param.name || '=' || user_product_param.value
				FROM user_product_param WHERE user_product_id = p.id
			) AS config_params
		FROM user_product as p
		JOIN purchase as pur ON (p.purchase_id = pur.id) `+
		conditions, params...,
	)
}

// GetUserProducts get a list of user products by email with pagination
func (repo *productRepo) GetUserProductsByEmail(email string,
	page int) (products []domain.Product, currentPage int,
	totalPages int, err error) {
	if page < 1 {
		page = 1
	}
	total := repo.GetUserProductsTotalByEmail(email)
	if total < 1 {
		return []domain.Product{}, page, 0, nil
	}
	totalPages = (total / repo.resultsPerPage)
	if (total % repo.resultsPerPage) > 0 {
		totalPages++
	}
	result, err := repo.makeUserProductQuery(
		`WHERE user_email = $1
		ORDER BY p.id DESC
		OFFSET $2 LIMIT $3`,
		email, (repo.resultsPerPage * (page - 1)),
		repo.resultsPerPage)
	if err != nil {
		return []domain.Product{}, 0, 0, err
	}
	defer result.Close()
	for result.Next() {
		product := domain.Product{}
		rawConfig := []string{}
		result.Scan(&product.ID, &product.Type, &product.UserID, &product.Email,
			&product.Status, &product.ExpiredAt, &product.CreatedAt,
			&product.Purchase.ID, &product.Purchase.Number, &product.Purchase.Type,
			&product.Purchase.Status, &product.Purchase.Price, &product.Purchase.CreatedAt,
			(*pq.StringArray)(&rawConfig))
		config, _ := repo.parseConfig(rawConfig)
		product.Config = config
		products = append(products, product)
	}
	return products, page, totalPages, nil
}

// GetUserActiveProduct gets active product for an specific userID
func (repo *productRepo) GetUserActiveProduct(userID int,
	productType domain.ProductType) (domain.Product, error) {
	result, err := repo.makeUserProductQuery(`
		WHERE  p.status = 'ACTIVE'
		AND p.user_id = $1 AND p.product_type = $2
		ORDER BY p.expired_at, p.start_at LIMIT 1`,
		userID, productType)
	if err != nil {
		return domain.Product{}, err
	}
	defer result.Close()
	product := domain.Product{}
	var configArr []string
	if result.Next() {
		result.Scan(&product.ID, &product.Type, &product.UserID, &product.Email,
			&product.Status, &product.ExpiredAt, &product.CreatedAt,
			&product.Purchase.ID, &product.Purchase.Number, &product.Purchase.Type,
			&product.Purchase.Status, &product.Purchase.Price, &product.Purchase.CreatedAt,
			(*pq.StringArray)(&configArr))
	} else {
		return domain.Product{}, usecases.ErrProductNotFound
	}
	config, err := repo.parseConfig(configArr)
	if err != nil {
		return domain.Product{}, err
	}
	product.Config = config
	return product, nil
}

// GetUserActiveProduct gets active product for an specific userProductID
func (repo *productRepo) GetUserProductByID(userProductID int) (domain.Product, error) {
	result, err := repo.makeUserProductQuery(`
		WHERE  p.id = $1`, userProductID)
	if err != nil {
		return domain.Product{}, err
	}
	defer result.Close()
	product := domain.Product{}
	var configArr []string
	if result.Next() {
		result.Scan(&product.ID, &product.Type, &product.UserID, &product.Email,
			&product.Status, &product.ExpiredAt, &product.CreatedAt,
			&product.Purchase.ID, &product.Purchase.Number, &product.Purchase.Type,
			&product.Purchase.Status, &product.Purchase.Price, &product.Purchase.CreatedAt,
			(*pq.StringArray)(&configArr))
	}

	config, err := repo.parseConfig(configArr)
	if err != nil {
		return domain.Product{}, err
	}
	product.Config = config
	return product, nil
}

// parseConfig parses rawConfiguration slice to domain.ProductParams  struct
func (repo *productRepo) parseConfig(rawConfig []string) (domain.ProductParams, error) {
	if len(rawConfig) == 0 {
		return domain.ProductParams{}, fmt.Errorf("No config found")
	}
	configs := make(map[string]string)
	for _, v := range rawConfig {
		pair := strings.Split(v, "=")
		if len(pair) >= 2 {
			configs[pair[0]] = strings.Join(pair[1:], "=")
		}
	}
	limit, _ := strconv.Atoi(configs["limit"])
	exclude, keywords := []string{}, []string{}
	if configs["exclude"] != "" {
		exclude = strings.Split(configs["exclude"], ",")
	}
	if configs["keywords"] != "" {
		keywords = strings.Split(configs["keywords"], ",")
	}
	priceRange, _ := strconv.Atoi(configs["price_range"])
	gapsWithRandom, _ := strconv.ParseBool(configs["fill_random"])
	categoriesArrayStr := strings.Split(configs["categories"], ",")
	categories := []int{}
	for _, v := range categoriesArrayStr {
		cat, _ := strconv.Atoi(v)
		if cat >= 1000 && cat <= 9999 {
			categories = append(categories, cat)
		}
	}
	return domain.ProductParams{
		Categories:         categories,
		Limit:              limit,
		Keywords:           keywords,
		Exclude:            exclude,
		PriceRange:         priceRange,
		FillGapsWithRandom: gapsWithRandom,
		Comment:            configs["comment"],
	}, nil
}

// CreateUserProduct creates a new product for user
func (repo *productRepo) CreateUserProduct(userID int, email string,
	purchase domain.Purchase, productType domain.ProductType, expiredAt time.Time,
	config domain.ProductParams) (domain.Product, error) {
	result, err := repo.handler.Query(
		`INSERT INTO user_product(product_type, status, user_id, user_email,
			purchase_id, expired_at)
			VALUES (
				$1, 'ACTIVE', $2, $3, $4, $5
			) RETURNING id, created_at`, productType, userID, email, purchase.ID,
		expiredAt)
	if err != nil {
		return domain.Product{}, err
	}
	defer result.Close()
	var userProductID int
	var createdAt time.Time
	if result.Next() {
		result.Scan(&userProductID, &createdAt)
	} else {
		return domain.Product{},
			fmt.Errorf("next error: getting userProductID from database")
	}
	err = repo.SetConfig(userProductID, config)
	if err != nil {
		return domain.Product{}, err
	}
	return domain.Product{
		ID:        userProductID,
		Type:      productType,
		Email:     email,
		UserID:    userID,
		ExpiredAt: expiredAt,
		CreatedAt: createdAt,
		Config:    config,
		Status:    domain.ActiveProduct,
		Purchase:  purchase,
	}, nil
}

// SetConfig adds configuration to Product
func (repo *productRepo) SetConfig(userProductID int, config domain.ProductParams) error {
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
		fmt.Sprintf(`INSERT INTO user_product_param(user_product_id, name, value) VALUES %s
			ON CONFLICT (user_product_id, name) DO UPDATE set value=excluded.value`,
			strings.Join(positions, ", ")),
		insertValues...)
}

func makeConfigValues(userProductID int, config domain.ProductParams) (values [][]interface{}) {
	return [][]interface{}{
		[]interface{}{userProductID, "categories", strings.Trim(strings.Join(
			strings.Fields(fmt.Sprint(config.Categories)), ","), "[]")},
		[]interface{}{userProductID, "limit", strconv.Itoa(config.Limit)},
		[]interface{}{userProductID, "keywords", strings.Join(config.Keywords, ",")},
		[]interface{}{userProductID, "exclude", strings.Join(config.Exclude, ",")},
		[]interface{}{userProductID, "price_range", strconv.Itoa(config.PriceRange)},
		[]interface{}{userProductID, "comment", config.Comment},
		[]interface{}{userProductID, "fill_random", fmt.Sprintf("%t", config.FillGapsWithRandom)},
	}
}

// SetConfig adds configuration to Product
func (repo *productRepo) SetPartialConfig(userProductID int, configMap map[string]interface{}) error {
	for name, value := range configMap {
		switch name {
		case "status":
			if err := repo.SetStatus(userProductID,
				domain.ProductStatus(value.(string))); err != nil {
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
func (repo *productRepo) SetStatus(userProductID int, status domain.ProductStatus) error {
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

// SetExpiration sets the expiration for product
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

// ExpireProductds sets expired status for all expired products
func (repo *productRepo) ExpireProducts() error {
	return repo.handler.Update(
		`UPDATE
			user_product
		SET
			status = 'EXPIRED'
		WHERE
			expired_at < NOW()
		AND
			status = 'ACTIVE'`,
	)
}
