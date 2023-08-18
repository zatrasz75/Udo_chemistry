package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
	"udo_mass/pkg/calculator"
	"udo_mass/pkg/storage"
)

// API представляет собой приложение с набором обработчиков.
type API struct {
	r *mux.Router // Маршрутизатор запросов
	//db        storage.Interface // база данных
	webRoot string // Корневая директория для веб-приложения
}

// Router возвращает маршрутизатор запросов.
func (api *API) Router() *mux.Router {
	return api.r
}

// New создает новый экземпляр API и инициализирует его маршрутизатор.
func New(webRoot string) *API {
	// Создаём новый API и привязываем к нему маршрутизатор и корневую директорию для веб-приложения.
	api := &API{
		r: mux.NewRouter(),
		//	db:        db,
		webRoot: webRoot,
	}
	// Регистрируем обработчики API.
	api.endpoints()

	return api
}

// Регистрация обработчиков API.
func (api *API) endpoints() {
	api.r.HandleFunc("/", api.home).Methods(http.MethodGet)
	api.r.HandleFunc("/calculateMolarMasses", api.calculateMolarMasses).Methods(http.MethodPost)

	// веб-приложение
	api.r.PathPrefix("/web/").Handler(http.StripPrefix("/web/", http.FileServer(http.Dir("./web"))))
}

// Обработчик для статических файлов веб-приложения.
func (api *API) serveWebFiles(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	// Проверяем, что запрошенный путь начинается с "/web/".
	if !strings.HasPrefix(filePath, "/web/") {
		http.NotFound(w, r)
		return
	}

	// Проверяем, что путь после "/web/" не содержит "../" (попытка обхода пути).
	if strings.Contains(filePath, "../") {
		http.NotFound(w, r)
		return
	}

	// Строим абсолютный путь к файлу.
	absolutePath := filepath.Join(api.webRoot, filePath[5:])

	// Проверяем, что абсолютный путь находится в пределах корневой директории для веб-приложения.
	if !strings.HasPrefix(absolutePath, api.webRoot) {
		http.NotFound(w, r)
		return
	}

	// Обслуживаем статический файл.
	http.ServeFile(w, r, absolutePath)
}

func (api *API) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}

	content, err := ioutil.ReadFile("web/udo.html")
	if err != nil {
		http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
		return
	}

	// Устанавливаем правильный Content-Type для HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	w.Write(content)
}

func (api *API) calculateMolarMasses(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/calculateMolarMasses" {
		http.NotFound(w, r)
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}

	var f storage.MolarMasses
	err := json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mass := calculator.CombineChemicalFormulas(f.Nitrate, f.Phosphate, f.Potassium, f.Micro)
	log.Println(mass)

	for symbol, mass := range mass {
		fmt.Printf("%s: %.4f г/моль\n", symbol, mass)
	}

	// Здесь вы можете выполнить расчет молярных масс для каждого вещества
	// И вернуть результат в формате JSON

	//response := map[string]string{
	//	"nitrate":   f.Nitrate,
	//	"phosphate": f.Phosphate,
	//	"potassium": f.Potassium,
	//	"micro":     f.Micro,
	//}

	// Логирование данных из запроса
	slog.Info("данные из ответа формы : ", mass)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mass)

}
