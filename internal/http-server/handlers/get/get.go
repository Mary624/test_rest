package get

import (
	"errors"
	"fmt"
	"strconv"
	httpserver "test-example/internal/http-server"
	"test-example/internal/storage"

	"github.com/labstack/echo/v4"
)

var (
	ErrWrongParam = errors.New("wrong param")
	ErrWrongPage  = errors.New("wrong page")
)

func GetById(ctx echo.Context, getter httpserver.PeopleUpdater) (storage.Person, error) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		return storage.Person{}, fmt.Errorf("can't get int from param 'id': %s", idStr)
	}

	res, err := getter.GetPersonById(int64(id))

	if err != nil {
		return storage.Person{}, err
	}

	return res, nil
}

func GetByFilter(ctx echo.Context, getter httpserver.PeopleUpdater) ([]storage.Person, error) {

	m := make(map[string]string, 6)

	if name := ctx.QueryParam("name"); name != "" {
		m["name"] = name
	}
	if surname := ctx.QueryParam("surname"); surname != "" {
		m["surname"] = surname
	}
	if patronymic := ctx.QueryParam("patronymic"); patronymic != "" {
		m["patronymic"] = patronymic
	}
	if gender := ctx.QueryParam("gender"); gender != "" {
		if gender == "male" || gender == "female" {
			m["gender"] = gender
		} else {
			return nil, ErrWrongParam
		}
	}
	if nationality := ctx.QueryParam("nationality"); nationality != "" {
		m["nationality"] = nationality
	}
	if age := ctx.QueryParam("age"); age != "" {
		_, err := strconv.Atoi(age)
		if err != nil {
			return nil, ErrWrongParam
		}
		m["age"] = age
	}

	limit, offset := 0, 0

	c := 0

	if limitStr := ctx.QueryParam("limit"); limitStr != "" {
		l, err := strconv.Atoi(limitStr)

		if err != nil {
			return nil, fmt.Errorf("can't get int from param 'limit': %s", limitStr)
		}
		limit = l
		c++
	}

	if pageStr := ctx.QueryParam("page"); pageStr != "" {
		p, err := strconv.Atoi(pageStr)

		if err != nil {
			return nil, fmt.Errorf("can't get int from param 'offset': %s", pageStr)
		}
		if p < 1 {
			return nil, ErrWrongPage
		}
		offset = (p - 1) * limit
		c++
	}

	pL := len(m) + c

	if pL < len(ctx.QueryParams()) {
		return nil, ErrWrongParam
	}

	res, err := getter.GetPeopleByFilter(limit, offset, m)

	if err != nil {
		return nil, err
	}

	return res, nil
}
