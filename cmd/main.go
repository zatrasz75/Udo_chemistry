package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
	"udo_mass/pkg/api"
	"udo_mass/pkg/middl"
)

type server struct {
	//	db  storage.Interface
	api *api.API
}

// init вызывается перед main()
func init() {
	// загружает значения из файла .env в систему
	if err := godotenv.Load(); err != nil {
		log.Print("Файл .env не найден.")
	}
}

func main() {
	//inStr := "KH2PO4"

	fmt.Println("// -------------------------------------------------------------------------")

	//// Создаем и заполняем карту молярных масс
	//molarMassMap := make(map[string]float64)
	//massMap := calculator.MolarMassCompound(molarMassMap)
	//
	//for symbol, mass := range massMap {
	//	fmt.Printf("%s: %.4f г/моль\n", symbol, mass)
	//}
	//
	//fmt.Println(massMap)

	fmt.Println("// -------------------------------------------------------------------------")

	//// Получение молярных масс из пакета calculator
	//k := calculator.MolarMasses["K"]
	//n := calculator.MolarMasses["N"]
	//o := calculator.MolarMasses["O"]
	//
	//fmt.Println(k, n, o)
	//
	//// Вычисление массовой доли азота (N) в нитрате калия (KNO3)
	//totalMolarMass := k + n + (o * 3) // Общая молярная масса нитрата калия
	//fmt.Printf("Общая молярная масса %.4f г/моль\n", totalMolarMass)
	//
	//nitrogenFraction := (n / totalMolarMass) * 100 // Массовая доля азота в процентах
	//fmt.Printf("Массовая доля азота в процентах: %.4f\n", nitrogenFraction)
	//
	//// Вычисление массы азота в граммах на моль (г/моль) в нитрате калия.
	//nitrogenGramsPerMole := nitrogenFraction * (n / 100) // Массовая доля азота в долях от 1 моля азота
	//fmt.Printf("Масса азота в нитрате калия: %.4f г/моль\n", nitrogenGramsPerMole)

	fmt.Println("// ----------------------------------------------------------------")

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "продолжительность, в течение которой сервер корректно ожидает завершения существующих подключений - например, 15 секунд или 1 м")
	flag.Parse()

	port := os.Getenv("APP_PORT")
	if port == "" {
		log.Fatal("переменная окружения APP_PORT не задана")
	}
	host := os.Getenv("APP_HOST")
	if host == "" {
		log.Fatal("Переменная окружения APP_HOST не задана")
	}

	// объект сервера
	var router server

	// Получаем текущий путь к main.go
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Не удалось получить текущий каталог:", err)
	}
	// Получаем абсолютный путь к каталогу web/
	webRoot := filepath.Join(currentDir, "../web")

	router.api = api.New(webRoot)

	// Логирования запросов.
	router.api.Router().Use(middl.Middle)

	log.Println("Запуск сервера на ", "http://"+host+":"+port)

	// Создаем HTTP сервер с заданным адресом и обработчиком.
	srv := &http.Server{
		Addr:         host + ":" + port,
		Handler:      router.api.Router(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запуск сервера в отдельном потоке.
	go func() {
		err = srv.ListenAndServe()
		if err != nil {
			log.Fatal("Не удалось запустить сервер шлюза. Error:", err)
		}
	}()
	graceShutdown(*srv, wait)

}

// Выключает сервер
func graceShutdown(srv http.Server, wait time.Duration) {
	quitCH := make(chan os.Signal, 1)
	signal.Notify(quitCH, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quitCH

	// Создаем крайний срок для ожидания.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Останавливаем сервер с таймаутом ожидания.
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Printf("Ошибка при закрытии прослушивателей или тайм-аут контекста %v", err)
		return
	}
	log.Printf("Выключение сервера")
	os.Exit(0)
}
