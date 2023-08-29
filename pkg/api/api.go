package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"udo_mass/config"
	"udo_mass/logger"
	"udo_mass/pkg/calculator"
	"udo_mass/pkg/storage"
)

// API представляет собой приложение с набором обработчиков.
type API struct {
	r    *mux.Router    // Маршрутизатор запросов
	cfg  *config.Config // Конфигурация
	port string         // Порт
	host string         // Хост
	//db        storage.Interface // база данных
}

// Router возвращает маршрутизатор запросов.
func (api *API) Router() *mux.Router {
	return api.r
}

// New создает новый экземпляр API и инициализирует его маршрутизатор.
func New(cfg *config.Config, port string, host string) *API {
	// Создаём новый API и привязываем к нему маршрутизатор и корневую директорию для веб-приложения.
	api := &API{
		r:    mux.NewRouter(),
		cfg:  cfg,
		port: port,
		host: host,
		//	db:        db,
	}
	// Регистрируем обработчики API.
	api.endpoints()

	return api
}

// Регистрация обработчиков API.
func (api *API) endpoints() {
	api.r.HandleFunc("/", api.home).Methods(http.MethodGet)
	api.r.HandleFunc("/", api.calculateMolarMasses).Methods(http.MethodPost)

	// веб-приложение
	api.r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))
}

func (api *API) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}

	content, err := ioutil.ReadFile("web/html/udo.html")
	if err != nil {
		http.Error(w, "Ошибка чтения файла", http.StatusNotFound)
		logger.Error("Ошибка чтения файла", err)
		return
	}

	// Устанавливаем правильный Content-Type для HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	w.Write(content)
}

func (api *API) calculateMolarMasses(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}

	var f storage.MolarMasses
	err := json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Error("Ошибка декодирования JSON", err)
		return
	}

	n := calculator.CombineChemicalFormulas(f.Nitrate, f.NitrateMass)
	p := calculator.CombineChemicalFormulas(f.Phosphate, f.PhosphateMass)
	k := calculator.CombineChemicalFormulas(f.Potassium, f.PotassiumMass)
	ir := calculator.CombineChemicalFormulas(f.Micro, f.MicroMass)

	response := calculator.CombineMaps(n, p, k, ir)
	fmt.Println("------------------------------------")
	for symbol, mass := range response {
		log.Printf("%s: %.4f г/литр\n", symbol, mass)
	}

	// Устанавливаем правильный Content-Type для HTML
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
