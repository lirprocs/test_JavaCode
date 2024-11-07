package handler

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"test_JavaCode/internal/service"
)

type WalletRequest struct {
	WalletID      uuid.UUID `json:"walletId"`
	OperationType string    `json:"operationType"`
	Amount        int64     `json:"amount"`
}

func CreateWalletHandler(w http.ResponseWriter, r *http.Request, svc service.WalletService) {
	var req WalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err := svc.HandleWalletOperation(req.WalletID, req.OperationType, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetWalletBalanceHandler(w http.ResponseWriter, r *http.Request, svc service.WalletService) {
	walletID, _ := uuid.Parse(r.URL.Path[len("/api/v1/wallets/"):])

	balance, err := svc.GetWalletBalance(walletID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]int64{"balance": balance})
}
