//go:build !docker
// +build !docker

package repository

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) UpdateBalance(walletID uuid.UUID, amount int64) error {
	args := m.Called(walletID, amount)
	return args.Error(0)
}

func (m *MockWalletRepository) GetBalance(walletID uuid.UUID) (int64, error) {
	args := m.Called(walletID)
	return args.Get(0).(int64), args.Error(1)
}

func TestUpdateBalance(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	walletID := uuid.New()

	mockRepo.On("UpdateBalance", walletID, int64(100)).Return(nil)

	err := mockRepo.UpdateBalance(walletID, 100)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetBalance(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	walletID := uuid.New()

	mockRepo.On("GetBalance", walletID).Return(int64(100), nil)

	balance, err := mockRepo.GetBalance(walletID)

	assert.NoError(t, err)
	assert.Equal(t, int64(100), balance)
	mockRepo.AssertExpectations(t)
}

func TestGetBalanceWithError(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	walletID := uuid.New()

	mockRepo.On("GetBalance", walletID).Return(int64(0), errors.New("db error"))

	balance, err := mockRepo.GetBalance(walletID)

	assert.Error(t, err)
	assert.Equal(t, int64(0), balance)
	mockRepo.AssertExpectations(t)
}
