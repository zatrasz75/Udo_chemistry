package main

import (
	"fmt"
	"udo_mass/pkg/calculator"
)

func main() {
	inStr := "KNO3"

	fmt.Println("// -------------------------------------------------------------------------")

	// Создаем и заполняем карту молярных масс
	molarMassMap := make(map[string]float64)
	massMap := calculator.MolarMassCompound(inStr, molarMassMap)

	for symbol, mass := range massMap {
		fmt.Printf("%s: %.4f г/моль\n", symbol, mass)
	}

	fmt.Println(massMap)

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

}
