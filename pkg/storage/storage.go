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

type TableMolarMass struct {
	Symbol []string
	Mass   []float64
}

type Database interface {
	CreatMolarMassTable() error
	DropMolarMassTable() error
	AddMolarMass(c TableMolarMass) error
	//SearchAccount(c Account) (string, error)
	//KeysAccount(c Account) (bool, error)
	//DelAccount(c Account) (bool, error)
}
