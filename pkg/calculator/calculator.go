package calculator

import (
	"fmt"
	"strconv"
	"unicode"
)

type Element struct {
	Symbol         string
	Mass           float64
	Count          int
	MassInCompound float64
}

// MolarMasses Создание карты с химическими элементами и их молярными массами
var molarMasses = map[string]float64{
	"H":  1.0081,
	"He": 4.0026,
	"Li": 6.94,
	"Be": 9.0122,
	"B":  10.811,
	"C":  12.01,
	"N":  14.0111,
	"O":  16.00,
	"F":  18.99,
	"Ne": 20.18,
	"Na": 22.99,
	"Mg": 24.31,
	"Al": 26.98,
	"Si": 28.09,
	"P":  30.97,
	"S":  32.07,
	"Cl": 35.45,
	"K":  39.10,
	"Ar": 39.95,
	"Ca": 40.08,
	"Sc": 44.96,
	"Ti": 47.87,
	"V":  50.94,
	"Cr": 52.00,
	"Mn": 54.94,
	"Fe": 55.85,
	"Ni": 58.69,
	"Co": 58.93,
	"Cu": 63.55,
	"Zn": 65.38,
	"Ga": 69.72,
	"Ge": 72.63,
	"As": 74.92,
	"Se": 78.97,
	"Br": 79.90,
	"Kr": 83.80,
	"Rb": 85.47,
	"Sr": 87.62,
	"Y":  88.91,
	"Zr": 91.22,
	"Nb": 92.91,
	"Mo": 95.94,
	"Tc": 98.00,
	"Ru": 101.07,
	"Rh": 102.91,
	"Pd": 106.42,
	"Ag": 107.87,
	"Cd": 112.41,
	"In": 114.82,
	"Sn": 118.71,
	"Sb": 121.76,
	"I":  126.90,
	"Te": 127.60,
	"Xe": 131.29,
	"Cs": 132.91,
	"Ba": 137.33,
	"La": 138.91,
	"Ce": 140.12,
	"Pr": 140.91,
	"Nd": 144.24,
	"Pm": 145.00,
	"Sm": 150.36,
	"Eu": 152.00,
	"Gd": 157.25,
	"Tb": 158.93,
	"Dy": 162.50,
	"Ho": 164.93,
	"Er": 167.26,
	"Tm": 168.93,
	"Yb": 173.05,
	"Lu": 175.00,
	"Hf": 178.49,
	"Ta": 180.95,
	"W":  183.84,
	"Re": 186.21,
	"Os": 190.23,
	"Ir": 192.22,
	"Pt": 195.08,
	"Au": 196.97,
	"Hg": 200.59,
	"Tl": 204.38,
	"Pb": 207.2,
	"Bi": 208.98,
	"Th": 232.04,
	"Pa": 231.04,
	"U":  238.03,
	"Np": 237.00,
	"Pu": 244.00,
	"Am": 243.00,
	"Cm": 247.00,
	"Bk": 247.00,
	"Cf": 251.00,
	"Es": 252.00,
	"Fm": 257.00,
	"Md": 258.00,
	"No": 259.00,
	"Lr": 262.00,
}

// CombineChemicalFormulas Объединяет строковые значения элементов и атомов в список
func CombineChemicalFormulas(formula, gram string) map[string]float64 {
	charCount := make(map[string]int)

	for i, char := range formula {
		if !unicode.IsUpper(char) && !unicode.IsDigit(char) {
			el := string(formula[i-1]) + string(char)
			charCount[el]++

			delete(charCount, string(formula[i-1]))

		} else if !unicode.IsDigit(char) && unicode.IsUpper(char) {
			charCount[string(char)]++

		} else {
			digit, _ := strconv.Atoi(string(char))
			el := formula[i-1]
			charCount[string(el)] += digit - 1
		}
	}

	malarMass := molarMassCompound(charCount)
	totalSubstance := totalMassOfTheSubstance(malarMass, gram)

	return totalSubstance
}

// MolarMassCompound Посчитывает малярную массу каждого элемента
func molarMassCompound(charCount map[string]int) map[string]float64 {
	compoundMasses := make(map[string]float64)
	totalMass := 0.0

	for element, count := range charCount {
		if mass, found := molarMasses[element]; found {
			massForElement := mass * float64(count)
			compoundMasses[element] = massForElement
			totalMass += massForElement
		} else {
			fmt.Printf("Не удалось найти молекулярную массу для элемента %s\n", element)
		}
	}

	return compoundMasses
}

// TotalMassOfTheSubstance Подсчитывает массу в зависимости от внесенного вещества в граммах
func totalMassOfTheSubstance(substance map[string]float64, mass string) map[string]float64 {
	m, _ := strconv.Atoi(mass)
	var total float64
	result := make(map[string]float64)

	for _, v := range substance {
		total += v
	}

	for element, count := range substance {
		result[element] = float64(m) * count / total
	}

	return result
}

// CombineMaps Объединяет полученные результаты в один список
func CombineMaps(maps ...map[string]float64) map[string]float64 {
	combined := make(map[string]float64)

	for _, m := range maps {
		for element, value := range m {
			combined[element] += value
		}
	}

	return combined
}
