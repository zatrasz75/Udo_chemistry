package main

import (
	"fmt"
	"strconv"
	"udo_mass/pkg/calculator"
)

func main() {
	inStr := "KNO3"

	//for _, v := range inStr {
	//	// Проверяем, является ли символ буквой (элементом) или цифрой
	//	if v >= 'A' && v <= 'Z' {
	//		// Обработка предыдущего элемента и его количества (если есть)
	//		if currentElement != "" {
	//			molarMass := calculator.MolarMasses[currentElement]
	//			fmt.Printf("Молярная масса %s: %.4f\n", currentElement, molarMass)
	//
	//			if currentNumber != "" {
	//				count, _ := strconv.Atoi(currentNumber)
	//				fmt.Printf("Количество %s в соединении: %d\n", currentElement, count)
	//			}
	//		}
	//		// Сбрасываем значения для нового элемента
	//		currentElement = string(v)
	//		currentNumber = ""
	//	} else if v >= '0' && v <= '9' {
	//		currentNumber += string(v)
	//	}
	//}
	//// Обработка последнего элемента
	//if currentElement != "" {
	//	molarMass := calculator.MolarMasses[currentElement]
	//	fmt.Printf("Молярная масса %s: %.4f\n", currentElement, molarMass)
	//
	//	if currentNumber != "" {
	//		count, _ := strconv.Atoi(currentNumber)
	//		fmt.Printf("Количество %s в соединении: %d\n", currentElement, count)
	//	}
	//}

	fmt.Println("// -------------------------------------------------------------------------------")

	MolarMassCompound(inStr)

	elements := MolarMassCompoundQ(inStr)
	fmt.Println(elements)

	//// Пример: Печать информации о каждом элементе
	//for i, element := range elements {
	//	fmt.Printf("Элемент %d: %s\n", i+1, element.Symbol)
	//	fmt.Printf("Молярная масса: %.4f г/моль\n", element.Mass)
	//	fmt.Printf("Количество в соединении: %d\n", element.Count)
	//	fmt.Printf("Масса в соединении: %.4f г\n", element.MassInCompound)
	//}
	//
	//// Пример: Вычисление общей массы соединения
	//totalMolarMass := 0.0
	//for _, element := range elements {
	//	totalMolarMass += element.MassInCompound
	//}
	//fmt.Printf("Общая молярная масса: %.4f г/моль\n", totalMolarMass)

	fmt.Println("// -------------------------------------------------------------------------")
	//
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

}

func MolarMassCompound(inStr string) {
	currentElement := ""
	currentNumber := ""
	totalMolarMass := 0.0

	for _, v := range inStr {
		// Проверяем, является ли символ буквой (элементом) или цифрой
		if v >= 'A' && v <= 'Z' {
			// Обработка предыдущего элемента и его количества (если есть)
			if currentElement != "" {
				molarMass := calculator.MolarMasses[currentElement]
				fmt.Printf("Молярная масса г/моль %s: %.4f\n", currentElement, molarMass)

				if currentNumber != "" {
					count, _ := strconv.Atoi(currentNumber)
					fmt.Printf("Количество %s в соединении: %d\n", currentElement, count)

					// Учитываем количество атомов в молекуле при вычислении общей молярной массы
					massForElement := float64(count) * molarMass
					fmt.Printf("Масса %s в соединении г/моль: %.4f г\n", currentElement, massForElement)
					totalMolarMass += massForElement
				} else {
					// Учитываем один атом элемента без указания количества
					totalMolarMass += molarMass
				}
			}
			// Сбрасываем значения для нового элемента
			currentElement = string(v)
			currentNumber = ""
		} else if v >= '0' && v <= '9' {
			currentNumber += string(v)
		}
	}

	// Обработка последнего элемента
	if currentElement != "" {
		molarMass := calculator.MolarMasses[currentElement]
		fmt.Printf("Молярная масса г/моль %s: %.4f\n", currentElement, molarMass)

		if currentNumber != "" {
			count, _ := strconv.Atoi(currentNumber)
			fmt.Printf("Количество %s в соединении: %d\n", currentElement, count)

			// Учитываем количество атомов в молекуле при вычислении общей молярной массы
			massForElement := float64(count) * molarMass
			fmt.Printf("Масса %s в соединении г/моль: %.4f г\n", currentElement, massForElement)
			totalMolarMass += massForElement
		} else {
			// Учитываем один атом элемента без указания количества
			totalMolarMass += molarMass
		}
	}

	fmt.Printf("Общая молярная масса г/моль: %.4f г\n", totalMolarMass)
}

func MolarMassCompoundQ(inStr string) []calculator.Element {
	elements := []calculator.Element{}
	currentElement := calculator.Element{}
	currentNumber := ""

	for _, v := range inStr {
		if v >= 'A' && v <= 'Z' {
			if currentElement.Symbol != "" {
				elements = append(elements, currentElement)
			}
			currentElement = calculator.Element{
				Symbol: string(v),
				Mass:   calculator.MolarMasses[string(v)],
			}
			currentNumber = ""
		} else if v >= '0' && v <= '9' {
			currentNumber += string(v)
		}
	}

	if currentElement.Symbol != "" {
		elements = append(elements, currentElement)
	}

	for i, element := range elements {
		fmt.Printf("Элемент %d: %s\n", i+1, element.Symbol)
		fmt.Printf("Молярная масса: %.4f г/моль\n", element.Mass)

		if currentNumber != "" {
			count, _ := strconv.Atoi(currentNumber)
			element.Count = count
			fmt.Printf("Количество %s в соединении: %d\n", element.Symbol, element.Count)

			massInCompound := float64(element.Count) * element.Mass
			fmt.Printf("Масса %s в соединении г/моль: %.4f г\n", element.Symbol, massInCompound)
		} else {
			element.Count = 1
			massInCompound := element.Mass
			fmt.Printf("Масса %s в соединении г/моль: %.4f г\n", element.Symbol, massInCompound)
		}

		currentNumber = "" // Сбрасываем currentNumber для следующего элемента
	}

	return elements
}
