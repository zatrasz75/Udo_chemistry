package main

import (
	"context"
	"flag"
	"fmt"
	http "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"udo_mass/config"
	"udo_mass/logger"
	"udo_mass/pkg/api"
)

type server struct {
	//	db  storage.Interface
	api *api.API
}

var cfg config.Config

// init вызывается перед main()
func init() {
	path := flag.String("config", "", "путь к конфигурационному файлу")
	flag.Parse()

	err := cfg.LoadEnv(*path)
	if err != nil {
		logger.Fatal("не удалось загрузить переменные окружения: %s", err)
	}
}

func main() {
	//formula1 := "KNO3"
	//formula2 := "KH2PO4"
	//formula3 := "K2SO4"
	//formula4 := "FeSO4*7H2O

	fmt.Println("// -------------------------------------------------------------------------")

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "продолжительность, в течение которой сервер корректно ожидает завершения существующих подключений - например, 15 секунд или 1 м")
	flag.Parse()

	// объект сервера
	var router server

	// Порт по умолчанию.
	port := cfg.AddrPort
	// Хост по умолчанию.
	host := cfg.AddrHost

	router.api = api.New(port, host)

	// Создаем HTTP сервер с заданным адресом и обработчиком.
	srv := &http.Server{
		Addr:         host + ":" + port,
		Handler:      router.api.Router(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	logger.Info("Запуск сервера на http://" + srv.Addr)

	// Запуск сервера в отдельном потоке.
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			logger.Fatal("Не удалось запустить сервер шлюза. Error:", err)
		}
	}()
	graceShutdown(*srv, wait)

}

// Выключает сервер
func graceShutdown(srv http.Server, wait time.Duration) {
	quitCH := make(chan os.Signal, 1)
	signal.Notify(quitCH, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quitCH

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	shutdownServer(srv, ctx)
}

func shutdownServer(srv http.Server, ctx context.Context) {
	err := srv.Shutdown(ctx)
	if err != nil {
		logger.Fatal("Ошибка при закрытии прослушивателей или тайм-аут контекста: %v", err)
	}

	logger.Info("Сервер успешно выключен")
}
