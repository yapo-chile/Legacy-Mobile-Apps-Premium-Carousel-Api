package repository

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/domain"
)

type mockResult struct {
	mock.Mock
}

func (m *mockResult) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *mockResult) Scan(dest ...interface{}) {
	args := m.Called(dest)
	p := args.Get(0).([]interface{})
	for k, dv := range dest {
		reflect.ValueOf(dv).Elem().Set(reflect.ValueOf(p[k]))
	}
}

func (m *mockResult) Close() error {
	args := m.Called()
	return args.Error(0)
}

type dbHandlerMock struct {
	mock.Mock
}

func (m *dbHandlerMock) Run(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *dbHandlerMock) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *dbHandlerMock) Query(statement string, params ...interface{}) (DbResult, error) {
	args := m.Called(statement, params)
	return args.Get(0).(DbResult), args.Error(1)
}

func (m *dbHandlerMock) Insert(statement string, params ...interface{}) error {
	args := m.Called(statement, params)
	return args.Error(0)
}

func (m *dbHandlerMock) Update(statement string, params ...interface{}) error {
	args := m.Called(statement, params)
	return args.Error(0)
}

type mockProductRepoLogger struct {
	mock.Mock
}

func (m *mockProductRepoLogger) LogWarnPartialConfigNotSupported(name, value string) {
	m.Called(name, value)
}

func TestMakeProductRepositoryOK(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)
	assert.Equal(t, &productRepo{
		handler:        mockDB,
		resultsPerPage: 10,
		logger:         mLogger,
	}, repo)
	mockDB.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserProductsTotalOk(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]interface {}"),
	).Return(mResult, nil)
	mResult.On("Close").Return(nil)
	mResult.On("Next").Return(true).Once()
	mResult.On("Scan", mock.Anything).Return([]interface{}{123})
	repo := MakeProductRepository(mockDB, 10, mLogger)
	result := repo.GetUserProductsTotal()
	assert.Equal(t, 123, result)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserProductsTotalError(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]interface {}"),
	).Return(mResult, fmt.Errorf("e"))
	repo := MakeProductRepository(mockDB, 10, mLogger)
	result := repo.GetUserProductsTotal()
	assert.Equal(t, 0, result)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserProductsTotalByEmailOk(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]interface {}"),
	).Return(mResult, nil)
	mResult.On("Close").Return(nil)
	mResult.On("Next").Return(true).Once()
	mResult.On("Scan", mock.Anything).Return([]interface{}{123})
	repo := MakeProductRepository(mockDB, 10, mLogger)
	result := repo.GetUserProductsTotalByEmail("123@123.cl")
	assert.Equal(t, 123, result)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserProductsTotalByEmailError(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]interface {}"),
	).Return(mResult, fmt.Errorf("e"))
	repo := MakeProductRepository(mockDB, 10, mLogger)
	result := repo.GetUserProductsTotalByEmail("123@test.cl")
	assert.Equal(t, 0, result)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserProductsByEmailOk(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)

	mResult.On("Close").Return(nil)
	// Get Products Total mocks
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]interface {}"),
	).Return(mResult, nil).Once()
	mResult.On("Next").Return(true).Once()
	mResult.On("Scan", mock.Anything).Return([]interface{}{11}).Once()
	// get products query
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, nil).Once()
	mResult.On("Next").Return(true).Once()
	mResult.On("Next").Return(false).Once()
	testTime := time.Now()
	mResult.On("Scan", mock.Anything).Return([]interface{}{
		11, domain.PremiumCarousel, 1, "test@mail.com", domain.ActiveProduct,
		testTime, testTime, 0, 0, domain.AdminPurchase, domain.AcceptedPurchase,
		100, testTime, []string{"categories=2020,1020",
			"keywords=a,b,c", "comment=comentario"}}).Once()

	result, currentPage,
		totalPages, err := repo.GetUserProductsByEmail("test@email.com", 0)
	expected := []domain.Product{
		{
			ID:        11,
			Type:      domain.PremiumCarousel,
			UserID:    1,
			Email:     "test@mail.com",
			Status:    domain.ActiveProduct,
			ExpiredAt: testTime,
			CreatedAt: testTime,
			Purchase: domain.Purchase{
				ID:        0,
				Number:    0,
				Price:     100,
				Type:      domain.AdminPurchase,
				Status:    domain.AcceptedPurchase,
				CreatedAt: testTime,
			},
			Config: domain.ProductParams{
				Categories: []int{2020, 1020},
				Exclude:    []string{},
				Keywords:   []string{"a", "b", "c"},
				Comment:    "comentario",
			},
		},
	}
	assert.Equal(t, 1, currentPage)
	assert.Equal(t, 2, totalPages)
	assert.Equal(t, expected, result)
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserProductsOk(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)

	mResult.On("Close").Return(nil)
	// Get Products Total mocks
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]interface {}"),
	).Return(mResult, nil).Once()
	mResult.On("Next").Return(true).Once()
	mResult.On("Scan", mock.Anything).Return([]interface{}{11}).Once()
	// get products query
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, nil).Once()
	mResult.On("Next").Return(true).Once()
	mResult.On("Next").Return(false).Once()
	testTime := time.Now()
	mResult.On("Scan", mock.Anything).Return([]interface{}{
		11, domain.PremiumCarousel, 1, "test@mail.com", domain.ActiveProduct,
		testTime, testTime, 0, 0, domain.AdminPurchase, domain.AcceptedPurchase,
		100, testTime, []string{"categories=2020,1020",
			"keywords=a,b,c", "comment=comentario"}}).Once()
	result, currentPage,
		totalPages, err := repo.GetUserProducts(0)
	expected := []domain.Product{
		{
			ID:        11,
			Type:      domain.PremiumCarousel,
			UserID:    1,
			Email:     "test@mail.com",
			Status:    domain.ActiveProduct,
			ExpiredAt: testTime,
			CreatedAt: testTime,
			Purchase: domain.Purchase{
				ID:        0,
				Number:    0,
				Price:     100,
				Type:      domain.AdminPurchase,
				Status:    domain.AcceptedPurchase,
				CreatedAt: testTime,
			},
			Config: domain.ProductParams{
				Categories: []int{2020, 1020},
				Exclude:    []string{},
				Keywords:   []string{"a", "b", "c"},
				Comment:    "comentario",
			},
		},
	}
	assert.Equal(t, 1, currentPage)
	assert.Equal(t, 2, totalPages)
	assert.Equal(t, expected, result)
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetReportOk(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)
	mResult.On("Close").Return(nil)

	// get products query
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, nil).Once()
	mResult.On("Next").Return(true).Once()
	mResult.On("Next").Return(false).Once()
	testTime := time.Now()
	mResult.On("Scan", mock.Anything).Return([]interface{}{
		11, domain.PremiumCarousel, 1, "test@mail.com", domain.ActiveProduct,
		testTime, testTime, 0, 0, domain.AdminPurchase, domain.AcceptedPurchase,
		100, testTime, []string{"categories=2020,1020",
			"keywords=a,b,c", "comment=comentario"}}).Once()
	result, err := repo.GetReport(testTime, testTime)
	expected := []domain.Product{
		{
			ID:        11,
			Type:      domain.PremiumCarousel,
			UserID:    1,
			Email:     "test@mail.com",
			Status:    domain.ActiveProduct,
			ExpiredAt: testTime,
			CreatedAt: testTime,
			Purchase: domain.Purchase{
				ID:        0,
				Number:    0,
				Price:     100,
				Type:      domain.AdminPurchase,
				Status:    domain.AcceptedPurchase,
				CreatedAt: testTime,
			},
			Config: domain.ProductParams{
				Categories: []int{2020, 1020},
				Exclude:    []string{},
				Keywords:   []string{"a", "b", "c"},
				Comment:    "comentario",
			},
		},
	}

	assert.Equal(t, expected, result)
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetReportQueryError(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, fmt.Errorf("err")).Once()
	_, err := repo.GetReport(time.Now(), time.Now())
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserProductsZeroResults(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]interface {}"),
	).Return(mResult, fmt.Errorf("e")).Once()
	result, currentPage,
		totalPages, err := repo.GetUserProducts(0)
	expected := []domain.Product{}
	assert.Equal(t, 1, currentPage)
	assert.Equal(t, 0, totalPages)
	assert.Equal(t, expected, result)
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserProductsByEmailZeroResults(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]interface {}"),
	).Return(mResult, fmt.Errorf("e")).Once()
	result, currentPage,
		totalPages, err := repo.GetUserProductsByEmail("test@email.com", 0)
	expected := []domain.Product{}
	assert.Equal(t, 1, currentPage)
	assert.Equal(t, 0, totalPages)
	assert.Equal(t, expected, result)
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserProductsByEmailQueryError(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)
	mResult.On("Close").Return(nil)
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]interface {}"),
	).Return(mResult, nil).Once()
	mResult.On("Next").Return(true).Once()
	mResult.On("Scan", mock.Anything).Return([]interface{}{11}).Once()
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, fmt.Errorf("err")).Once()
	_, _, _, err := repo.GetUserProductsByEmail("test@email.com", 0)
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserProductsQueryError(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)
	mResult.On("Close").Return(nil)
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]interface {}"),
	).Return(mResult, nil).Once()
	mResult.On("Next").Return(true).Once()
	mResult.On("Scan", mock.Anything).Return([]interface{}{11}).Once()
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, fmt.Errorf("err")).Once()
	_, _, _, err := repo.GetUserProducts(0)
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserActiveProductOk(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)
	mResult.On("Close").Return(nil)
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, nil).Once()
	mResult.On("Next").Return(true).Once()
	testTime := time.Now()
	mResult.On("Scan", mock.Anything).Return([]interface{}{
		11, domain.PremiumCarousel, 1, "test@mail.com", domain.ActiveProduct,
		testTime, testTime, 0, 0, domain.AdminPurchase, domain.AcceptedPurchase,
		100, testTime, []string{"categories=2020,1020",
			"keywords=a,b,c", "comment=comentario"}}).Once()
	result, err := repo.GetUserActiveProduct(1,
		domain.PremiumCarousel)
	expected := domain.Product{
		ID:        11,
		Type:      domain.PremiumCarousel,
		UserID:    1,
		Email:     "test@mail.com",
		Status:    domain.ActiveProduct,
		ExpiredAt: testTime,
		CreatedAt: testTime,
		Purchase: domain.Purchase{
			ID:        0,
			Number:    0,
			Price:     100,
			Type:      domain.AdminPurchase,
			Status:    domain.AcceptedPurchase,
			CreatedAt: testTime,
		},
		Config: domain.ProductParams{
			Categories: []int{2020, 1020},
			Exclude:    []string{},
			Keywords:   []string{"a", "b", "c"},
			Comment:    "comentario",
		},
	}
	assert.Equal(t, expected, result)
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserActiveProductErrorProductNotFound(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)
	mResult.On("Close").Return(nil)
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, nil).Once()
	mResult.On("Next").Return(false).Once()
	_, err := repo.GetUserActiveProduct(1,
		domain.PremiumCarousel)

	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserActiveProductError(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, fmt.Errorf("err")).Once()
	_, err := repo.GetUserActiveProduct(1,
		domain.PremiumCarousel)
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserActiveProductErrorNoConfig(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)
	mResult.On("Close").Return(nil)
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, nil).Once()
	mResult.On("Next").Return(true).Once()
	testTime := time.Now()
	mResult.On("Scan", mock.Anything).Return([]interface{}{
		11, domain.PremiumCarousel, 1, "test@mail.com", domain.ActiveProduct,
		testTime, testTime, 0, 0, domain.AdminPurchase, domain.AcceptedPurchase,
		100, testTime, []string{}}).Once()
	_, err := repo.GetUserActiveProduct(1,
		domain.PremiumCarousel)
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestCreateUserProductOk(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)
	mResult.On("Close").Return(nil).Once()
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, nil).Once()
	mResult.On("Next").Return(true).Once()
	mockDB.On("Insert",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(nil).Once()
	testTime := time.Now()
	mResult.On("Scan", mock.Anything).Return([]interface{}{
		11, testTime}).Once()
	result, err := repo.CreateUserProduct(1, "test@mail.com", domain.Purchase{},
		domain.PremiumCarousel, testTime, domain.ProductParams{
			Categories: []int{2020, 1020},
			Exclude:    []string{"11111", "22222"},
		})
	expected := domain.Product{
		ID:        11,
		Type:      domain.PremiumCarousel,
		UserID:    1,
		Email:     "test@mail.com",
		Status:    domain.ActiveProduct,
		ExpiredAt: testTime,
		CreatedAt: testTime,
		Purchase:  domain.Purchase{},
		Config: domain.ProductParams{
			Categories: []int{2020, 1020},
			Exclude:    []string{"11111", "22222"},
		},
	}
	assert.Equal(t, expected, result)
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestCreateUserProductQueryError(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, fmt.Errorf("err")).Once()
	testTime := time.Now()
	_, err := repo.CreateUserProduct(1, "test@mail.com", domain.Purchase{},
		domain.PremiumCarousel, testTime, domain.ProductParams{
			Categories: []int{2020, 1020},
			Exclude:    []string{"11111", "22222"},
		})
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestCreateUserProductNextError(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, nil).Once()
	mResult.On("Close").Return(nil).Once()
	mResult.On("Next").Return(false).Once()
	testTime := time.Now()
	_, err := repo.CreateUserProduct(1, "test@mail.com", domain.Purchase{},
		domain.PremiumCarousel, testTime, domain.ProductParams{
			Categories: []int{2020, 1020},
			Exclude:    []string{"11111", "22222"},
		})
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestCreateUserProductAddConfigError(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)
	mResult.On("Close").Return(nil).Once()
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, nil).Once()
	mResult.On("Next").Return(true).Once()
	mockDB.On("Insert",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(fmt.Errorf("e")).Once()
	testTime := time.Now()
	mResult.On("Scan", mock.Anything).Return([]interface{}{
		11, testTime}).Once()
	_, err := repo.CreateUserProduct(1, "test@mail.com", domain.Purchase{},
		domain.PremiumCarousel, testTime, domain.ProductParams{
			Categories: []int{2020, 1020},
			Exclude:    []string{"11111", "22222"},
		})
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserProductByIDOK(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)
	mResult.On("Close").Return(nil)
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, nil).Once()
	mResult.On("Next").Return(true).Once()
	testTime := time.Now()
	mResult.On("Scan", mock.Anything).Return([]interface{}{
		11, domain.PremiumCarousel, 1, "test@mail.com", domain.ActiveProduct,
		testTime, testTime, 0, 0, domain.AdminPurchase, domain.AcceptedPurchase,
		100, testTime, []string{"categories=2020,1020",
			"keywords=a,b,c", "exclude=1,2,3", "comment=comentario"}}).Once()
	result, err := repo.GetUserProductByID(11)
	expected := domain.Product{
		ID:        11,
		Type:      domain.PremiumCarousel,
		UserID:    1,
		Email:     "test@mail.com",
		Status:    domain.ActiveProduct,
		ExpiredAt: testTime,
		CreatedAt: testTime,
		Purchase: domain.Purchase{
			ID:        0,
			Number:    0,
			Price:     100,
			Type:      domain.AdminPurchase,
			Status:    domain.AcceptedPurchase,
			CreatedAt: testTime,
		},
		Config: domain.ProductParams{
			Categories: []int{2020, 1020},
			Exclude:    []string{"1", "2", "3"},
			Keywords:   []string{"a", "b", "c"},
			Comment:    "comentario",
		},
	}
	assert.Equal(t, expected, result)
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserProductByIDQueryError(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, fmt.Errorf("err")).Once()
	_, err := repo.GetUserProductByID(11)
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserProductByIDParseConfigError(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)
	mResult.On("Close").Return(nil)
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, nil).Once()
	mResult.On("Next").Return(true).Once()
	testTime := time.Now()
	mResult.On("Scan", mock.Anything).Return([]interface{}{
		11, domain.PremiumCarousel, 1, "test@mail.com", domain.ActiveProduct,
		testTime, testTime, 0, 0, domain.AdminPurchase, domain.AcceptedPurchase,
		100, testTime, []string{}}).Once()
	_, err := repo.GetUserProductByID(11)

	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestSetPartialConfigOK(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)
	mResult.On("Close").Return(nil)
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, nil).Once()
	mLogger.On("LogWarnPartialConfigNotSupported",
		mock.Anything, mock.Anything)
	err := repo.SetPartialConfig(11, map[string]interface{}{
		"status": "ACTIVE",
		"other":  "not supported",
	})
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestSetPartialConfigError(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, fmt.Errorf("err")).Once()

	err := repo.SetPartialConfig(11, map[string]interface{}{
		"status": "ACTIVE",
	})
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestSetExpirationOK(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)
	mResult.On("Close").Return(nil)
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, nil).Once()
	err := repo.SetExpiration(11, time.Now())
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestSetExpirationError(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, fmt.Errorf("err")).Once()
	err := repo.SetExpiration(11, time.Now())
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestSetExpirateProductsOK(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)
	mockDB.On("Update",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(nil)
	err := repo.ExpireProducts()
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestSetExpirateProductsError(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mLogger := &mockProductRepoLogger{}
	repo := MakeProductRepository(mockDB, 10, mLogger)
	mockDB.On("Update",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(fmt.Errorf("err"))
	err := repo.ExpireProducts()
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}
