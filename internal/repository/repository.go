package repository

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type WalletRepository interface {
	UpdateBalance(walletID uuid.UUID, amount int64) error
	GetBalance(walletID uuid.UUID) (int64, error)
}

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

type PostgresRepository struct{}

func (r *PostgresRepository) UpdateBalance(walletID uuid.UUID, amount int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentBalance int64
	err = tx.QueryRow("SELECT balance FROM wallets WHERE id=$1 FOR UPDATE", walletID).Scan(&currentBalance)

	if err == sql.ErrNoRows {
		_, err = tx.Exec("INSERT INTO wallets (id, balance) VALUES ($1, $2)", walletID, 0)
		if err != nil {
			return err
		}
		currentBalance = 0
	} else if err != nil {
		return err
	}

	newBalance := currentBalance + amount
	if newBalance < 0 {
		return errors.New("insufficient funds")
	}

	_, err = tx.Exec("UPDATE wallets SET balance=$1 WHERE id=$2", newBalance, walletID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *PostgresRepository) GetBalance(walletID uuid.UUID) (int64, error) {
	var balance int64

	err := db.QueryRow("SELECT balance FROM wallets WHERE id=$1", walletID).Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = db.Exec("INSERT INTO wallets (id, balance) VALUES ($1, $2)", walletID, 0)
			if err != nil {
				return 0, err
			}
			return 0, nil
		}
		return 0, err
	}
	return balance, nil
}
