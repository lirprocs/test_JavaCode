//go:build !docker
// +build !docker

package service

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

func TestHandleWalletOperationDeposit(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	walletID := uuid.New()

	mockRepo.On("UpdateBalance", walletID, int64(100)).Return(nil)

	svc := NewWalletService(mockRepo)

	err := svc.HandleWalletOperation(walletID, "DEPOSIT", 100)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestHandleWalletOperationWithdraw(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	walletID := uuid.New()

	mockRepo.On("UpdateBalance", walletID, int64(-50)).Return(nil)

	svc := NewWalletService(mockRepo)

	err := svc.HandleWalletOperation(walletID, "WITHDRAW", 50)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestHandleWalletOperationInvalid(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	walletID := uuid.New()

	svc := NewWalletService(mockRepo)

	err := svc.HandleWalletOperation(walletID, "INVALID", 100)

	assert.Error(t, err)
	assert.Equal(t, "invalid operation type", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetWalletBalance(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	walletID := uuid.New()

	mockRepo.On("GetBalance", walletID).Return(int64(200), nil)

	svc := NewWalletService(mockRepo)

	balance, err := svc.GetWalletBalance(walletID)

	assert.NoError(t, err)
	assert.Equal(t, int64(200), balance)
	mockRepo.AssertExpectations(t)
}

func TestGetWalletBalanceWithError(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	walletID := uuid.New()

	mockRepo.On("GetBalance", walletID).Return(int64(0), errors.New("db error"))

	svc := NewWalletService(mockRepo)

	balance, err := svc.GetWalletBalance(walletID)

	assert.Error(t, err)
	assert.Equal(t, int64(0), balance)
	mockRepo.AssertExpectations(t)
}
