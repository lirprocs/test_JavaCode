package service

import (
	"errors"
	"github.com/google/uuid"
	"test_JavaCode/internal/repository"
)

type WalletService interface {
	HandleWalletOperation(walletID uuid.UUID, operationType string, amount int64) error
	GetWalletBalance(walletID uuid.UUID) (int64, error)
}

type WalletServiceImpl struct {
	repo repository.WalletRepository
}

func NewWalletService(repo repository.WalletRepository) WalletService {
	return &WalletServiceImpl{repo: repo}
}

func (s *WalletServiceImpl) HandleWalletOperation(walletID uuid.UUID, operationType string, amount int64) error {
	switch operationType {
	case "DEPOSIT":
		return s.repo.UpdateBalance(walletID, amount)
	case "WITHDRAW":
		return s.repo.UpdateBalance(walletID, -amount)
	default:
		return errors.New("invalid operation type")
	}
}

func (s *WalletServiceImpl) GetWalletBalance(walletID uuid.UUID) (int64, error) {
	return s.repo.GetBalance(walletID)
}
