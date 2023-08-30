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
func CalculateMolarMasses(w http.ResponseWriter, r *http.Request) {
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
	for symbol, mass := range response {
		log.Printf("%s: %.4f г/литр\n", symbol, mass)
	}

	// Устанавливаем правильный Content-Type для HTML
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Ошибка кодирования JSON", http.StatusInternalServerError)
		logger.Error("Ошибка кодирования JSON", err)
	}
}
