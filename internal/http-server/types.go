package httpserver

import (
	"errors"
	"test-example/internal/storage"
)

var (
	ErrWrongValue = errors.New("wrong value")
)

type PeopleUpdater interface {
	DeletePerson(int64) error
	SavePerson(storage.Person) error
	UpdatePerson(int64, map[string]string) error
	GetPersonById(int64) (storage.Person, error)
	GetPeopleByFilter(int, int, map[string]string) ([]storage.Person, error)
}
