package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"udo_mass/pkg/calculator"
	"udo_mass/pkg/logger"
	"udo_mass/pkg/storage"
)

// Home обрабатывает GET запросы на корневой путь и отправляет содержимое HTML файла "udo.html".
func Home(w http.ResponseWriter, r *http.Request) {
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

// MultiError - это тип, из среза ошибок.
type MultiError []error

// Метод Error для MultiError формирует строку, содержащую все ошибки, разделенные точкой с запятой.
func (me MultiError) Error() string {
	var errStrs []string
	for _, err := range me {
		errStrs = append(errStrs, err.Error())
	}
	return strings.Join(errStrs, "; ")
}

// CalculateMolarMasses обрабатывает POST запросы для вычисления молекулярных масс химических соединений.
func CalculateMolarMasses(w http.ResponseWriter, r *http.Request, db storage.Database) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}

	var f storage.MolarMasses
	err := json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		http.Error(w, "Ошибка декодирования JSON", http.StatusBadRequest)
		logger.Error("Ошибка декодирования JSON", err)
		return
	}

	n := make(map[string]float64)
	p := make(map[string]float64)
	k := make(map[string]float64)
	ir := make(map[string]float64)
	// Канал ошибок из горутин
	errorCh := make(chan error)

	// Канал для передачи результатов из горутин в основной поток
	ch := make(chan map[string]float64)

	// MultiError для сбора ошибок из горутин
	var errors MultiError

	go func() {
		result := calculator.CombineChemicalFormulas(f.Nitrate, f.NitrateMass)
		ch <- result
	}()
	go func() {
		result := calculator.CombineChemicalFormulas(f.Phosphate, f.PhosphateMass)
		ch <- result
	}()
	go func() {
		result := calculator.CombineChemicalFormulas(f.Potassium, f.PotassiumMass)
		ch <- result
	}()
	go func() {
		result := calculator.CombineChemicalFormulas(f.Micro, f.MicroMass)
		ch <- result
	}()

	// Принимает результаты из горутин и собираем в соответствующие мапы
	for i := 0; i < 4; i++ {
		select {
		case result := <-ch:
			switch i {
			case 0:
				n = result
			case 1:
				p = result
			case 2:
				k = result
			case 3:
				ir = result
			}
		case err := <-errorCh:
			errors = append(errors, err)
		}
	}
	close(ch)

	response := calculator.CombineMaps(n, p, k, ir)
	fmt.Println("------------------------------------")
	// Создаем одну запись с данными о всех элементах
	for symbol, mass := range response {
		log.Printf("%s: %.4f г/литр\n", symbol, mass)
	}

	if len(response) != 0 {
		err = db.AddMolarMass(response)
		if err != nil {
			logger.Error("ошибка при вставке данных в базу данных:", err)
		}
	}

	all, err := db.AllMolarMass()
	if err != nil {
		http.Error(w, "не получилось получить данные из таблицы", http.StatusInternalServerError)
		logger.Error("не получилось получить данные из таблицы", err)
		return
	}
	// Формируем HTML-строку с данными
	var output string
	for _, v := range all {
		for id, data := range v {
			output += fmt.Sprintf("Ответ: %d<br>\n", id)
			for element, mass := range data {
				output += fmt.Sprintf("%s: %.4f г/литр<br>\n", element, mass)
			}
		}
	}

	// Устанавливаем правильный Content-Type для HTML
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(output))
}

// DelRecord обрабатывает POST запрос и удаляет запись в базе данных по её id
func DelRecord(w http.ResponseWriter, r *http.Request, db storage.Database) {
	if r.URL.Path != "/delet" {
		http.NotFound(w, r)
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}

	// Получаем id записи, которую нужно удалить, из запроса
	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат id", http.StatusBadRequest)
		logger.Error("Неверный формат id", err)
		return
	}

	// Проверяем, есть такая запись
	exists, err := db.SearchRecordById(id)
	if err != nil {
		// Обработка ошибки, если она возникла при поиске записи
		http.Error(w, "Ошибка при поиске записи", http.StatusInternalServerError)
		logger.Error("Ошибка при поиске записи", err)
		return
	}
	if !exists {
		// Запись с указанным id не существует
		http.Error(w, "Запись с указанным id не найдена", http.StatusNotFound)
		logger.Info("Запись с указанным id не найдена %s", idStr)
		return
	}

	// Удалить запись по id
	deleted, err := db.DelRecord(id)
	if err != nil {
		http.Error(w, "Ошибка при удалении записи", http.StatusInternalServerError)
		logger.Error("Ошибка при удалении записи", err)
		return
	}
	// Проверяем, была ли запись успешно удалена
	if !deleted {
		http.Error(w, "Запись с указанным id не найдена", http.StatusNotFound)
		logger.Error("Запись с указанным id не найдена", err)
		return
	} else {
		logger.Info("Запись удалена id %s", idStr)
	}

	all, err := db.AllMolarMass()
	if err != nil {
		http.Error(w, "не получилось получить данные из таблицы", http.StatusInternalServerError)
		logger.Error("не получилось получить данные из таблицы", err)
		return
	}
	// Формируем HTML-строку с данными
	var output string
	for _, v := range all {
		for id, data := range v {
			output += fmt.Sprintf("Ответ: %d<br>\n", id)
			for element, mass := range data {
				output += fmt.Sprintf("%s: %.4f г/литр<br>\n", element, mass)
			}
		}
	}

	// Устанавливаем правильный Content-Type для HTML
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(output))
}
