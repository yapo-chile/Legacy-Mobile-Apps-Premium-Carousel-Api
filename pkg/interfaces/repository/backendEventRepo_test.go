package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
)

type mockKafkaProducer struct {
	mock.Mock
}

func (m *mockKafkaProducer) SendMessage(topic string, bytes []byte) error {
	args := m.Called(topic, bytes)
	return args.Error(0)
}

func (m *mockKafkaProducer) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestPushSoldProductOK(t *testing.T) {
	mProducer := &mockKafkaProducer{}
	mProducer.On("SendMessage", mock.AnythingOfType("string"),
		mock.AnythingOfType("[]uint8")).Return(nil)
	repo := MakeBackendEventsProducer(mProducer, "")
	err := repo.PushSoldProduct(domain.Product{Type: domain.PremiumCarousel})
	assert.NoError(t, err)
	mProducer.AssertExpectations(t)
}

func TestPushSoldProductErrorProductNotSupoorted(t *testing.T) {
	mProducer := &mockKafkaProducer{}

	repo := MakeBackendEventsProducer(mProducer, "")
	err := repo.PushSoldProduct(domain.Product{Type: "arepa"})
	assert.Error(t, err)
	mProducer.AssertExpectations(t)
}
