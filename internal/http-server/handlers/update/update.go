package update

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	httpserver "test-example/internal/http-server"
	"test-example/internal/storage"

	"github.com/labstack/echo/v4"
)

var (
	ErrNoUpdatedFields = errors.New("no fields to update")
)

func Update(ctx echo.Context, updater httpserver.PeopleUpdater) (int, error) {

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("can't get int from param 'id': %s", idStr)
	}

	var person storage.Person
	person.Name = storage.DefaultStr
	person.Surname = storage.DefaultStr
	person.Patronymic = storage.DefaultStr
	person.Age = storage.DefaultAge
	person.Gender = storage.DefaultStr
	person.Nationality = storage.DefaultStr
	ctx.Bind(&person)

	m := make(map[string]string, 6)
	if person.Name != storage.DefaultStr {
		m["name"] = person.Name
	}
	if person.Surname != storage.DefaultStr {
		m["surname"] = person.Surname
	}
	if person.Patronymic != storage.DefaultStr {
		m["patronymic"] = person.Patronymic
	}
	if person.Age != storage.DefaultAge {
		m["age"] = strconv.Itoa(int(person.Age))
	}
	if person.Nationality != storage.DefaultStr {
		m["nationality"] = person.Nationality
	}

	if len(m) == 0 {
		return id, ErrNoUpdatedFields
	}

	r := regexp.MustCompile("^[A-z]+$")
	for key, value := range m {
		if key == "patronymic" && value == "" {
			continue
		}
		if key != "age" && !r.MatchString(value) {
			return 0, httpserver.ErrWrongValue
		}
	}
	if person.Gender != storage.DefaultStr {
		if person.Gender == "male" || person.Gender == "female" {
			m["gender"] = person.Gender
		} else {
			return 0, httpserver.ErrWrongValue
		}
	}

	err = updater.UpdatePerson(int64(id), m)

	if err != nil {
		return 0, fmt.Errorf("can't update person: %s", err.Error())
	}

	return id, nil
}
