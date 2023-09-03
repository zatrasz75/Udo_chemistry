package storage

type MolarMasses struct {
	NitrateMass   string `json:"nitrate_mass"`
	Nitrate       string `json:"nitrate"`
	PhosphateMass string `json:"phosphate_mass"`
	Phosphate     string `json:"phosphate"`
	PotassiumMass string `json:"potassium_mass"`
	Potassium     string `json:"potassium"`
	MicroMass     string `json:"micro_mass"`
	Micro         string `json:"micro"`
}

type Database interface {
	CreatMolarMassTable() error
	DropMolarMassTable() error
	AddMolarMass(c map[string]float64) error
	AllMolarMass() ([]map[int]map[string]float64, error)
	DelRecord(id int) (bool, error)
	SearchRecordById(id int) (bool, error)
}
