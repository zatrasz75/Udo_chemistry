package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lib/pq"
	"time"
	"udo_mass/pkg/logger"
	Database "udo_mass/pkg/storage"
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
    symbols TEXT[] NOT NULL,
    masses FLOAT[] NOT NULL
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

// AddMolarMass Добавляет данные в базу Postgres
func (s *Store) AddMolarMass(c Database.TableMolarMass) error {
	symbols := pq.StringArray(c.Symbol)
	masses := pq.Float64Array(c.Mass)

	_, err := s.db.Exec(context.Background(),
		"INSERT INTO molar_mass_data(symbols, masses) VALUES ($1, $2);", symbols, masses)
	if err != nil {
		logger.Error("ошибка при вставке данных в базу данных:", err)
		return err
	}

	return nil
}
