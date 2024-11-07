//go:build !docker
// +build !docker

package handler

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) HandleWalletOperation(walletID uuid.UUID, operationType string, amount int64) error {
	args := m.Called(walletID, operationType, amount)
	return args.Error(0)
}

func (m *MockWalletService) GetWalletBalance(walletID uuid.UUID) (int64, error) {
	args := m.Called(walletID)
	return args.Get(0).(int64), args.Error(1)
}

func TestCreateWalletHandler(t *testing.T) {
	mockSvc := new(MockWalletService)
	walletID := uuid.New()

	mockSvc.On("HandleWalletOperation", walletID, "DEPOSIT", int64(100)).Return(nil)

	req := &WalletRequest{
		WalletID:      walletID,
		OperationType: "DEPOSIT",
		Amount:        100,
	}
	body, _ := json.Marshal(req)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewReader(body))

	rr := httptest.NewRecorder()

	CreateWalletHandler(rr, request, mockSvc)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetWalletBalanceHandler(t *testing.T) {
	mockSvc := new(MockWalletService)
	walletID := uuid.New()

	mockSvc.On("GetWalletBalance", walletID).Return(int64(100), nil)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/"+walletID.String(), nil)

	rr := httptest.NewRecorder()

	GetWalletBalanceHandler(rr, request, mockSvc)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]int64
	_ = json.NewDecoder(rr.Body).Decode(&response)
	assert.Equal(t, int64(100), response["balance"])

	mockSvc.AssertExpectations(t)
}
