package storage

import "errors"

var (
	ErrEntryNotFound = errors.New("person not found")
	ErrEntryExists   = errors.New("person exists")
)

const (
	DefaultAge = 1000
	DefaultStr = "none1"
)

type Person struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int64  `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}
