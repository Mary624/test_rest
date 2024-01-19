package post

import (
	"errors"
	"fmt"
	"regexp"
	"test-example/internal/client"
	httpserver "test-example/internal/http-server"
	"test-example/internal/storage"

	"github.com/labstack/echo/v4"
)

var (
	ErrNoName    = errors.New("can't get name")
	ErrNoSurname = errors.New("can't get surname")
)

func Save(ctx echo.Context, saver httpserver.PeopleUpdater) error {

	var person storage.Person

	ctx.Bind(&person)
	if person.Name == "" {
		return ErrNoName
	}
	if person.Surname == "" {
		return ErrNoSurname
	}

	r := regexp.MustCompile("^[A-z]+$")
	if !r.MatchString(person.Name) {
		return httpserver.ErrWrongValue
	}
	if !r.MatchString(person.Surname) {
		return httpserver.ErrWrongValue
	}
	if person.Patronymic != "" && !r.MatchString(person.Patronymic) {
		return httpserver.ErrWrongValue
	}

	client := client.New()
	age, err := client.InfoAge(person.Name)
	if err != nil {
		return err
	}
	gender, err := client.InfoGender(person.Name)
	if err != nil {
		return err
	}
	nationality, err := client.InfoNationality(person.Name)
	if err != nil {
		return err
	}

	person.Age = age
	person.Gender = gender
	person.Nationality = nationality

	err = saver.SavePerson(person)

	if err != nil {
		return fmt.Errorf("can't save person: %s", err.Error())
	}

	return nil
}
