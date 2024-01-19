package delete

import (
	"fmt"
	"strconv"
	httpserver "test-example/internal/http-server"

	"github.com/labstack/echo/v4"
)

func Delete(ctx echo.Context, deleter httpserver.PeopleUpdater) (int, error) {

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return 0, fmt.Errorf("can't get int from param 'id': %s", idStr)
	}

	err = deleter.DeletePerson(int64(id))

	if err != nil {
		return 0, fmt.Errorf("can't delete person: %s", err.Error())
	}

	return id, nil
}
