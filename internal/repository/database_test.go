//go:build docker
// +build docker

package repository

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var testDB *sql.DB

func initTestDB(t *testing.T) {
	var err error
	godotenv.Load("config_test.env")
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("POSTGRES_DB"))
	testDB, err = sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatal("Не удалось подключиться к базе данных:", err)
	}
	SetDB(testDB)

	_, err = testDB.Exec("DELETE FROM wallets")
	if err != nil {
		t.Fatal("Не удалось очистить таблицу wallets:", err)
	}
}

func TestUpdateBalance(t *testing.T) {
	initTestDB(t)
	defer tearDown()

	repo := &PostgresRepository{}

	walletID := uuid.New()

	// Тест 1: Увеличение баланса
	err := repo.UpdateBalance(walletID, 100)
	assert.NoError(t, err)

	balance, err := repo.GetBalance(walletID)
	assert.NoError(t, err)
	assert.Equal(t, int64(100), balance)

	// Тест 2: Уменьшение баланса
	err = repo.UpdateBalance(walletID, -50)
	assert.NoError(t, err)

	balance, err = repo.GetBalance(walletID)
	assert.NoError(t, err)
	assert.Equal(t, int64(50), balance)

	// Тест 3: Проверка недостатка средств
	err = repo.UpdateBalance(walletID, -100)
	assert.Error(t, err)
	assert.Equal(t, "insufficient funds", err.Error())

	// Тест 4: Проверка создания нового кошелька
	newWalletID := uuid.New()
	balance, err = repo.GetBalance(newWalletID)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), balance)

	// Тест 5: Проверка обновления баланса нового кошелька
	err = repo.UpdateBalance(newWalletID, 200)
	assert.NoError(t, err)

	balance, err = repo.GetBalance(newWalletID)
	assert.NoError(t, err)
	assert.Equal(t, int64(200), balance)
}

func TestGetBalance(t *testing.T) {
	initTestDB(t)
	defer tearDown()

	repo := &PostgresRepository{}

	walletID := uuid.New()

	balance, err := repo.GetBalance(walletID)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), balance)

	err = repo.UpdateBalance(walletID, 50)
	assert.NoError(t, err)

	balance, err = repo.GetBalance(walletID)
	assert.NoError(t, err)
	assert.Equal(t, int64(50), balance)
}

func TestConcurrentUpdateBalance(t *testing.T) {
	initTestDB(t)
	defer tearDown()

	repo := &PostgresRepository{}
	walletID := uuid.New()
	initialBalance := int64(1000)
	err := repo.UpdateBalance(walletID, initialBalance)
	assert.NoError(t, err)

	numGoroutines := 100
	amountPerGoroutine := int64(10)
	expectedBalance := initialBalance + int64(numGoroutines)*amountPerGoroutine

	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			err := repo.UpdateBalance(walletID, amountPerGoroutine)
			assert.NoError(t, err)
			done <- true
		}()
	}

	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	finalBalance, err := repo.GetBalance(walletID)
	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, finalBalance)
}

func tearDown() {
	if testDB != nil {
		testDB.Close()
	}
}
