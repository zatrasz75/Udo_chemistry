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
    session_id INT REFERENCES "session_token"(id),
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

// CreatSessionTable Создает таблицу session_token
func (s *Store) CreatSessionTable() error {
	qwery := `CREATE TABLE IF NOT EXISTS "session_token" (
    id SERIAL PRIMARY KEY,
    token TEXT CHECK (LENGTH(token) = 88),
    created_at TIMESTAMP DEFAULT NOW()
);`

	_, err := s.db.Exec(context.Background(), qwery)
	if err != nil {
		logger.Error("не удалось создать таблицу session_token", err)
		return err
	}

	return nil
}

// DropSessionTable Удаляет таблицу session_token
func (s *Store) DropSessionTable() error {
	drop := `DROP TABLE IF EXISTS "session_token";`

	_, err := s.db.Exec(context.Background(), drop)
	if err != nil {
		return err
	}

	return nil
}

// GetSessionTokenID проверяет есть ли такая запись и возвращает ее id
func (s *Store) GetSessionTokenID(token string) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(),
		"SELECT id FROM session_token WHERE token = $1;", token).Scan(&id)

	return id, err
}

// AddSessionToken Добавляет запись о сессии
func (s *Store) AddSessionToken(token string) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(),
		"INSERT INTO session_token(token) VALUES ($1) RETURNING id;", token).Scan(&id)
	if err != nil {
		logger.Error("ошибка при вставке данных в базу данных:", err)
		return 0, err
	}
	return id, err
}

// AddMolarMass Добавляет данные в таблицу
func (s *Store) AddMolarMass(c map[string]float64, id int) error {
	jsonData, err := json.Marshal(c)
	if err != nil {
		logger.Error("ошибка при преобразовании данных в JSON:", err)
	}

	_, err = s.db.Exec(context.Background(),
		"INSERT INTO molar_mass_data(data, session_id) VALUES ($1, $2);", jsonData, id)
	if err != nil {
		logger.Error("ошибка при вставке данных в базу данных:", err)
		return err
	}

	return nil
}

// AllMolarMass Выводит все данные из таблицы molar_mass_data
func (s *Store) AllMolarMass(sessionID int) ([]map[int]map[string]float64, error) {
	var data []map[int]map[string]float64

	rows, err := s.db.Query(context.Background(), "SELECT id, data FROM molar_mass_data WHERE session_id = $1;", sessionID)
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

// SearchRecordById Проверяет, существует запись с таким id.
func (s *Store) SearchRecordById(id int, sessionID int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM molar_mass_data WHERE id = $1 AND session_id = $2)"

	var exists bool
	err := s.db.QueryRow(context.Background(), query, id, sessionID).Scan(&exists)
	if err != nil {
		logger.Error("нет такой записи", err)
		return false, err
	}
	return exists, err
}

// DelRecord Удаляет запись с выбранным id
func (s *Store) DelRecord(id int, sessionID int) (bool, error) {
	delet := "DELETE FROM molar_mass_data WHERE id = $1 AND session_id = $2"

	_, err := s.db.Exec(context.Background(), delet, id, sessionID)
	if err != nil {
		logger.Error("не удалось удалить запись", err)
		return false, err
	}
	return true, nil
}
