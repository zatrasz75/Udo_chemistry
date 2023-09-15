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
	CreatSessionTable() error
	DropMolarMassTable() error
	DropSessionTable() error
	AddMolarMass(c map[string]float64, id int) error
	AllMolarMass(int) ([]map[int]map[string]float64, error)
	DelRecord(id int, sessionID int) (bool, error)
	SearchRecordById(id int, sessionID int) (bool, error)
	GetSessionTokenID(token string) (int, error)
	AddSessionToken(token string) (int, error)
}
