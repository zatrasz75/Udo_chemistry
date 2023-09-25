package api

import (
	"context"
	"crypto/tls"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"udo_mass/config"
	handlers "udo_mass/pkg/handlers"
	"udo_mass/pkg/logger"
	"udo_mass/pkg/storage"
	"udo_mass/pkg/storage/postgres"
)

// API представляет собой приложение с набором обработчиков.
type API struct {
	r      *mux.Router // Маршрутизатор запросов
	port   string      // Порт
	host   string      // Хост
	srv    *http.Server
	db     storage.Database // база данных
	server *handlers.Server
}

// Router возвращает маршрутизатор запросов.
func (api *API) Router() *mux.Router {
	return api.r
}

// New создает новый экземпляр API и инициализирует его маршрутизатор.
func New() *API {
	// Конфигурация
	cfg := config.New()

	db, err := postgres.New(cfg.Database.ConnStr)
	if err != nil {
		logger.Fatal("нет соединения с PostgresSQL", err)
	}
	err = db.DropMolarMassTable()
	if err != nil {
		logger.Fatal("не удалось удалить таблицу", err)
	}
	err = db.DropSessionTable()
	if err != nil {
		logger.Fatal("не удалось удалить таблицу", err)
	}
	err = db.CreatSessionTable()
	if err != nil {
		logger.Fatal("не удалось создать таблицу session_token", err)
	}
	err = db.CreatMolarMassTable()
	if err != nil {
		logger.Fatal("не удалось создать таблицу molar_mass_data", err)
	}

	// Создаём новый API и привязываем к нему маршрутизатор и корневую директорию для веб-приложения.
	api := &API{
		r:      mux.NewRouter(),
		port:   cfg.Server.AddrPort,
		host:   cfg.Server.AddrHost,
		db:     db,
		server: &handlers.Server{Db: db}, // экземпляр Server
	}
	// Регистрируем обработчики API.
	api.endpoints()

	return api
}

// Start Метод для запуска сервера
func (api *API) Start() error {
	// Конфигурация
	cfg := config.New()

	// Загрузка сертификата и ключа
	cert, err := tls.LoadX509KeyPair(cfg.Server.AddrCert, cfg.Server.AddrKey)
	if err != nil {
		return err
	}

	// Создание TLS конфигурации с возможностью ввода паролей
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	api.srv = &http.Server{
		Addr:         api.host + ":" + api.port,
		Handler:      api.r,
		TLSConfig:    tlsConfig,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
	logger.Info("Запуск сервера на https://" + api.srv.Addr)

	go func() {
		err = api.srv.ListenAndServeTLS("", "")
		if err != nil {
			logger.Error("Остановка сервера", err)
			return
		}
	}()

	return nil
}

// Stop Метод для остановки сервера
func (api *API) Stop() error {
	// Конфигурация
	cfg := config.New()

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTime)
	defer cancel()
	err := api.srv.Shutdown(ctx)
	if err != nil {
		logger.Error("Shutdown ошибка при попытке остановить сервер", err)
		return err
	}

	return nil
}

// GraceShutdown Выключает сервер при получении сигнала об остановке
func GraceShutdown(httpServer *API) {
	quitCH := make(chan os.Signal, 1)
	signal.Notify(quitCH, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quitCH

	err := shutdownServer(httpServer)
	if err != nil {
		logger.Fatal("Ошибка при остановке сервера:", err)
	}
}

func shutdownServer(httpServer *API) error {
	err := httpServer.Stop()
	if err != nil {
		logger.Error("Ошибка при закрытии прослушивателей или тайм-аут контекста: %v", err)
	}

	logger.Info("Сервер успешно выключен")

	return nil
}

// Регистрация обработчиков API.
func (api *API) endpoints() {
	api.r.HandleFunc("/", api.server.Home).Methods(http.MethodGet)
	api.r.HandleFunc("/", api.server.CalculateMolarMasses).Methods(http.MethodPost)
	api.r.HandleFunc("/delet", api.server.DelRecord).Methods(http.MethodPost)

	// веб-приложение
	api.r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))
}
