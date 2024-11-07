package main

import (
	"context"
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test_JavaCode/config"
	"test_JavaCode/internal/handler"
	"test_JavaCode/internal/repository"
	"test_JavaCode/internal/service"
	"time"
)

func main() {
	cfg := config.New()

	db, err := sql.Open("postgres", cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	repo := &repository.PostgresRepository{}
	repository.SetDB(db)

	svc := service.NewWalletService(repo)

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/wallet", func(w http.ResponseWriter, r *http.Request) {
		handler.CreateWalletHandler(w, r, svc)
	}).Methods("POST")
	router.HandleFunc("/api/v1/wallets/{walletID}", func(w http.ResponseWriter, r *http.Request) {
		handler.GetWalletBalanceHandler(w, r, svc)
	}).Methods("GET")

	server := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: router,
	}

	go func() {
		log.Printf("Сервер запущен на порту %s", cfg.AppPort)
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan
	log.Println("Получен сигнал завершения, остановка сервера...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = server.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при остановке сервера: %v", err)
	}
	log.Println("Сервер успешно остановлен.")
}
