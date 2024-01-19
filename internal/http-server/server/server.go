package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	httpserver "test-example/internal/http-server"
	"test-example/internal/http-server/handlers/delete"
	"test-example/internal/http-server/handlers/get"
	"test-example/internal/http-server/handlers/post"
	"test-example/internal/http-server/handlers/update"
	"test-example/internal/logger"

	"github.com/labstack/echo/v4"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type Server struct {
	e   *echo.Echo
	log *slog.Logger
}

func New(env string, db httpserver.PeopleUpdater) *Server {
	e := echo.New()

	log := setupLogger(env)

	e.DELETE("/people/:id", func(ctx echo.Context) error {
		id, err := delete.Delete(ctx, db)
		if err != nil {
			log.Error("can't delete person", logger.Err(err))
			return ctx.String(http.StatusInternalServerError, err.Error())
		}
		log.Info("delete person", slog.Int("id", id))
		return ctx.String(http.StatusOK, "deleted")
	})

	e.GET("/people/:id", func(ctx echo.Context) error {
		res, err := get.GetById(ctx, db)
		if err != nil {
			log.Error("can't get person", logger.Err(err))
			return ctx.String(http.StatusInternalServerError, err.Error())
		}
		log.Info("get person")
		return ctx.JSON(http.StatusOK, res)
	})

	e.GET("/people", func(ctx echo.Context) error {
		res, err := get.GetByFilter(ctx, db)
		if err != nil {
			log.Error("can't get people", logger.Err(err))
			return ctx.String(http.StatusInternalServerError, err.Error())
		}
		log.Info("get people")
		return ctx.JSON(http.StatusOK, res)
	})

	e.PATCH("/people/:id", func(ctx echo.Context) error {
		id, err := update.Update(ctx, db)
		if err != nil {
			log.Error("can't update person", logger.Err(err))
			return ctx.String(http.StatusInternalServerError, err.Error())
		}
		log.Info("update person", slog.Int("id", id))
		return ctx.String(http.StatusOK, "updated")
	})

	e.POST("/people", func(ctx echo.Context) error {
		err := post.Save(ctx, db)
		if err != nil {
			log.Error("can't save person", logger.Err(err))
			return ctx.String(http.StatusInternalServerError, err.Error())
		}
		log.Info("save person")
		return ctx.String(http.StatusOK, "saved")
	})

	return &Server{
		e:   e,
		log: log,
	}
}

func (s *Server) Run(port int) error {
	s.log.Info("start server")
	err := s.e.Start(fmt.Sprintf(":%d", port))

	if err != nil {
		return err
	}
	return err
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
