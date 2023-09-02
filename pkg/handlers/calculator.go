package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

	n := calculator.CombineChemicalFormulas(f.Nitrate, f.NitrateMass)
	p := calculator.CombineChemicalFormulas(f.Phosphate, f.PhosphateMass)
	k := calculator.CombineChemicalFormulas(f.Potassium, f.PotassiumMass)
	ir := calculator.CombineChemicalFormulas(f.Micro, f.MicroMass)

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
		logger.Error("не получилось получить данные из таблицы", err)
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
