package repository

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"
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

func TestMakeProductRepositoryOK(t *testing.T) {
	mockDB := &dbHandlerMock{}
	repo := MakeProductRepository(mockDB, 10)
	assert.Equal(t, &productRepo{
		handler:        mockDB,
		resultsPerPage: 10,
	}, repo)
	mockDB.AssertExpectations(t)
}

func TestGetUserProductsTotalOk(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]interface {}"),
	).Return(mResult, nil)
	mResult.On("Close").Return(nil)
	mResult.On("Next").Return(true).Once()
	mResult.On("Scan", mock.Anything).Return([]interface{}{123})
	repo := MakeProductRepository(mockDB, 10)
	result := repo.GetUserProductsTotal("123")
	assert.Equal(t, 123, result)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
}

func TestGetUserProductsTotalError(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]interface {}"),
	).Return(mResult, fmt.Errorf("e"))
	repo := MakeProductRepository(mockDB, 10)
	result := repo.GetUserProductsTotal("123")
	assert.Equal(t, 0, result)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
}

func TestGetUserProductsOk(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
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
		11, usecases.PremiumCarousel, "1", "test@mail.com", usecases.ActiveProduct,
		testTime, testTime, []string{"categories=2020,1020"}, "comentario"}).Once()
	repo := MakeProductRepository(mockDB, 10)
	result, currentPage,
		totalPages, err := repo.GetUserProducts("test@email.com", 0)
	expected := []usecases.Product{
		{
			ID:        11,
			Type:      usecases.PremiumCarousel,
			UserID:    "1",
			Email:     "test@mail.com",
			Status:    usecases.ActiveProduct,
			ExpiredAt: testTime,
			CreatedAt: testTime,
			Config: usecases.CpConfig{
				Categories: []int{2020, 1020},
				Exclude:    []string{},
			},
			Comment: "comentario",
		},
	}
	assert.Equal(t, 1, currentPage)
	assert.Equal(t, 2, totalPages)
	assert.Equal(t, expected, result)
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
}

func TestGetUserProductsZeroResults(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]interface {}"),
	).Return(mResult, fmt.Errorf("e")).Once()
	repo := MakeProductRepository(mockDB, 10)
	result, currentPage,
		totalPages, err := repo.GetUserProducts("test@email.com", 0)
	expected := []usecases.Product{}
	assert.Equal(t, 1, currentPage)
	assert.Equal(t, 0, totalPages)
	assert.Equal(t, expected, result)
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
}

func TestGetUserProductsQueryError(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
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
	repo := MakeProductRepository(mockDB, 10)
	_, _, _, err := repo.GetUserProducts("test@email.com", 0)
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
}

func TestGetUserActiveProductOk(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mResult.On("Close").Return(nil)
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, nil).Once()
	mResult.On("Next").Return(true).Once()
	testTime := time.Now()
	mResult.On("Scan", mock.Anything).Return([]interface{}{
		11, usecases.PremiumCarousel, "1",
		"test@mail.com", usecases.ActiveProduct, testTime, testTime,
		[]string{"categories=2020,1020",
			"exclude=11111,22222"}, "comentario"}).Once()
	repo := MakeProductRepository(mockDB, 10)
	result, err := repo.GetUserActiveProduct("test@email.com",
		usecases.PremiumCarousel)
	expected := usecases.Product{
		ID:        11,
		Type:      usecases.PremiumCarousel,
		UserID:    "1",
		Email:     "test@mail.com",
		Status:    usecases.ActiveProduct,
		ExpiredAt: testTime,
		CreatedAt: testTime,
		Config: usecases.CpConfig{
			Categories: []int{2020, 1020},
			Exclude:    []string{"11111", "22222"},
		},
		Comment: "comentario",
	}
	assert.Equal(t, expected, result)
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
}

func TestGetUserActiveProductError(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, fmt.Errorf("err")).Once()
	repo := MakeProductRepository(mockDB, 10)
	_, err := repo.GetUserActiveProduct("test@email.com",
		usecases.PremiumCarousel)
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
}

func TestGetUserActiveProductErrorNoConfig(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mResult.On("Close").Return(nil)
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, nil).Once()
	mResult.On("Next").Return(true).Once()
	testTime := time.Now()
	mResult.On("Scan", mock.Anything).Return([]interface{}{
		11, usecases.PremiumCarousel, "1",
		"test@mail.com", usecases.ActiveProduct,
		testTime, testTime, []string{}, "comentario"}).Once()
	repo := MakeProductRepository(mockDB, 10)
	_, err := repo.GetUserActiveProduct("test@email.com",
		usecases.PremiumCarousel)
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
}

func TestAddUserProductOk(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
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
	repo := MakeProductRepository(mockDB, 10)
	result, err := repo.AddUserProduct("1", "test@mail.com", "comentario",
		usecases.PremiumCarousel, testTime, usecases.CpConfig{
			Categories: []int{2020, 1020},
			Exclude:    []string{"11111", "22222"},
		})
	expected := usecases.Product{
		ID:        11,
		Type:      usecases.PremiumCarousel,
		UserID:    "1",
		Email:     "test@mail.com",
		Status:    usecases.ActiveProduct,
		ExpiredAt: testTime,
		CreatedAt: testTime,
		Config: usecases.CpConfig{
			Categories: []int{2020, 1020},
			Exclude:    []string{"11111", "22222"},
		},
		Comment: "comentario",
	}
	assert.Equal(t, expected, result)
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
}

func TestAddUserProductBadUserId(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	testTime := time.Now()
	repo := MakeProductRepository(mockDB, 10)
	_, err := repo.AddUserProduct("aaaaaa", "test@mail.com", "comentario",
		usecases.PremiumCarousel, testTime, usecases.CpConfig{
			Categories: []int{2020, 1020},
			Exclude:    []string{"11111", "22222"},
		})
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
}

func TestAddUserProductQueryError(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, fmt.Errorf("err")).Once()
	testTime := time.Now()
	repo := MakeProductRepository(mockDB, 10)
	_, err := repo.AddUserProduct("1", "test@mail.com", "comentario",
		usecases.PremiumCarousel, testTime, usecases.CpConfig{
			Categories: []int{2020, 1020},
			Exclude:    []string{"11111", "22222"},
		})
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
}

func TestAddUserProductNextError(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mResult, nil).Once()
	mResult.On("Close").Return(nil).Once()
	mResult.On("Next").Return(false).Once()
	testTime := time.Now()
	repo := MakeProductRepository(mockDB, 10)
	_, err := repo.AddUserProduct("1", "test@mail.com", "comentario",
		usecases.PremiumCarousel, testTime, usecases.CpConfig{
			Categories: []int{2020, 1020},
			Exclude:    []string{"11111", "22222"},
		})
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
}

func TestAddUserProductAddConfigError(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
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
	repo := MakeProductRepository(mockDB, 10)
	_, err := repo.AddUserProduct("1", "test@mail.com", "comentario",
		usecases.PremiumCarousel, testTime, usecases.CpConfig{
			Categories: []int{2020, 1020},
			Exclude:    []string{"11111", "22222"},
		})
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
}
