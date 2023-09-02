package postgres

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
	"udo_mass/pkg/logger"
)

// Store Хранилище данных
type Store struct {
	db *pgxpool.Pool
}

// New Конструктор
func New(constr string) (*Store, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var db, err = pgxpool.Connect(ctx, constr)
	if err != nil {
		return nil, err
	}
	s := Store{
		db: db,
	}
	return &s, nil
}

// CreatMolarMassTable Создает таблицу molar_mass_data
func (s *Store) CreatMolarMassTable() error {
	qwery := `CREATE TABLE IF NOT EXISTS "molar_mass_data" (
    id SERIAL PRIMARY KEY,
    data JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);`

	_, err := s.db.Exec(context.Background(), qwery)
	if err != nil {
		logger.Error("не удалось создать таблицу molar_mass_data", err)
		return err
	}

	return nil
}

// DropMolarMassTable Удаляет таблицу molar_mass_data
func (s *Store) DropMolarMassTable() error {
	drop := `DROP TABLE IF EXISTS "molar_mass_data";`

	_, err := s.db.Exec(context.Background(), drop)
	if err != nil {
		return err
	}

	return nil
}

// AddMolarMass Добавляет данные в таблицу molar_mass_data
func (s *Store) AddMolarMass(response map[string]float64) error {
	jsonData, err := json.Marshal(response)
	if err != nil {
		logger.Error("ошибка при преобразовании данных в JSON:", err)
	}

	_, err = s.db.Exec(context.Background(),
		"INSERT INTO molar_mass_data(data) VALUES ($1);", jsonData)
	if err != nil {
		logger.Error("ошибка при вставке данных в базу данных:", err)
		return err
	}

	return nil
}

// AllMolarMass Выводит все данные из таблицы molar_mass_data
func (s *Store) AllMolarMass() ([]map[int]map[string]float64, error) {
	var data []map[int]map[string]float64

	rows, err := s.db.Query(context.Background(), "SELECT id, data FROM molar_mass_data;")
	if err != nil {
		logger.Error("ошибка при выполнении запроса:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var jsonData []byte
		err = rows.Scan(&id, &jsonData)
		if err != nil {
			logger.Error("ошибка при сканировании строки:", err)
			return nil, err
		}

		var result map[string]float64
		err = json.Unmarshal(jsonData, &result)
		if err != nil {
			logger.Error("ошибка при разборе JSON:", err)
			return nil, err
		}
		dataPoint := map[int]map[string]float64{
			id: result,
		}

		data = append(data, dataPoint)
	}

	if err = rows.Err(); err != nil {
		logger.Error("ошибка при обработке результатов запроса:", err)
		return nil, err
	}

	return data, nil
}

// DelRecord Удаляет запись с выбранным id
func (s *Store) DelRecord(id int) (bool, error) {
	delet := "DELETE FROM molar_mass_data WHERE id = $1"

	_, err := s.db.Exec(context.Background(), delet, id)
	if err != nil {
		logger.Error("не удалось удалить запись", err)
		return false, err
	}
	return true, nil
}
