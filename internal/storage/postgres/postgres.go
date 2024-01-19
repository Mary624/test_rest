package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"test-example/internal/config"
	"test-example/internal/storage"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(cfg config.DBConfig) (*Storage, error) {
	const op = "storage.New"

	connInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.HostDB, cfg.PortDB, cfg.UserDB, cfg.PassDB, cfg.DBName)
	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SavePerson(person storage.Person) error {
	const op = "storage.SavePerson"

	stmt, err := s.db.Prepare(`INSERT INTO people(name, surname, patronymic, age, gender, nationality)
	 VALUES($1, $2, $3, $4, $5, $6);`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(person.Name, person.Surname, person.Patronymic,
		person.Age, person.Gender, person.Nationality)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeletePerson(id int64) error {
	const op = "storage.DeletePerson"

	stmt, err := s.db.Prepare("DELETE FROM people WHERE id=$1;")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) UpdatePerson(id int64, filters map[string]string) error {
	const op = "storage.UpdatePerson"

	_, err := s.db.Exec(fmt.Sprintf(`UPDATE people 
	SET %s
	WHERE id=%d;`, getUpdateParams(filters), id))

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) GetPersonById(id int64) (storage.Person, error) {
	const op = "storage.GetPersonById"

	stmt, err := s.db.Prepare(`SELECT id, name, surname, patronymic, age, gender, nationality
	FROM people
	WHERE id=$1;`)
	if err != nil {
		return storage.Person{}, fmt.Errorf("%s: %w", op, err)
	}
	row := stmt.QueryRow(id)
	var person storage.Person
	err = row.Scan(&person.Id, &person.Name, &person.Surname, &person.Patronymic,
		&person.Age, &person.Gender, &person.Nationality)
	if errors.Is(sql.ErrNoRows, err) {
		return storage.Person{}, fmt.Errorf("%s: %w", op, storage.ErrEntryNotFound)
	}
	if err != nil {
		return storage.Person{}, fmt.Errorf("%s: %w", op, err)
	}
	return person, nil
}

func (s *Storage) GetPeopleByFilter(limit, offset int, filters map[string]string) ([]storage.Person, error) {
	const op = "storage.GetPeopleByFilter"

	limitStr := ""
	if limit != 0 {
		limitStr = fmt.Sprintf("LIMIT %d", limit)
	}
	rows, err := s.db.Query(fmt.Sprintf(`SELECT id, name, surname, patronymic, age, gender, nationality
	FROM people %s ORDER BY id ASC %s OFFSET %d;`, getQueryParams(filters), limitStr, offset))
	res := make([]storage.Person, 0, 100)
	if errors.Is(sql.ErrNoRows, err) {
		return res, fmt.Errorf("%s: %w", op, storage.ErrEntryNotFound)
	}
	if err != nil {
		return res, fmt.Errorf("%s: %w", op, err)
	}
	for rows.Next() {
		var person storage.Person
		err = rows.Scan(&person.Id, &person.Name, &person.Surname, &person.Patronymic,
			&person.Age, &person.Gender, &person.Nationality)
		if err != nil {
			return res, fmt.Errorf("%s: %w", op, err)
		}
		res = append(res, person)
	}
	return res, nil
}

func getUpdateParams(filters map[string]string) string {
	if len(filters) == 0 {
		return ""
	}
	var builder strings.Builder
	i, l := 0, len(filters)
	for key, value := range filters {
		if key != "age" {
			builder.WriteString(fmt.Sprintf("%s='%s'", key, value))
		} else {
			builder.WriteString(fmt.Sprintf("%s=%s", key, value))
		}
		if i < l-1 {
			builder.WriteString(", ")
		}
		i++
	}
	return builder.String()
}

func getQueryParams(filters map[string]string) string {
	if len(filters) == 0 {
		return ""
	}
	var builder strings.Builder
	builder.WriteString("WHERE ")
	i, l := 0, len(filters)
	for key, value := range filters {
		if key != "age" {
			builder.WriteString(fmt.Sprintf("%s='%s'", key, value))
		} else {
			builder.WriteString(fmt.Sprintf("%s=%s", key, value))
		}
		if i < l-1 {
			builder.WriteString(" AND ")
		}
		i++
	}
	return builder.String()
}
