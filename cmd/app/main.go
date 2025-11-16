package main

import (
	"errors"
	"log"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	_ "modernc.org/sqlite"
	"rosatom.ru/nko/internal/handler"
	"rosatom.ru/nko/internal/middleware"
	"rosatom.ru/nko/internal/repository"
)

func main() {
	db, err := repository.NewDB()
	if err != nil {
		log.Fatal("Error with repository.NewDB() in main.go:", err)
	}
	defer db.Close()

	nkoRepo := repository.NewNKORepo(db)
	citiesRepo := repository.NewCityRepo(db)
	userRepo := repository.NewUserRepository(db)

	nkoHandler := handler.NewNKOHandler(nkoRepo)
	citiesHandler := handler.NewCitiesHandler(citiesRepo)
	userHandler := handler.NewAuthHandler(userRepo, "aaa")

	e := echo.New()

	e.Use(echomw.Logger())
	e.Use(echomw.Recover())
	e.Use(echomw.CORS())

	// Routes
	// получить все нко
	e.GET("/nko", nkoHandler.GetAllNKO)
	// получить нко по id
	e.GET("/nko/:id", nkoHandler.GetByID)
	// получить только имя нко по id
	e.GET("/nko/name:id", nkoHandler.GetNKOName)
	// поиск нко с фильтрами по имени и категории
	e.GET("/nko/search", nkoHandler.SearchNKO)
	// получение всех городов
	e.GET("/cities", citiesHandler.GetAllNKO)
	// получение города по id
	e.GET("/cities/:id", citiesHandler.GetByID)
	// регистрация
	e.POST("/register", userHandler.Register)
	// логин
	e.POST("/login", userHandler.Login)
	// просмотр профиля
	e.GET("/api/profile", userHandler.GetProfile, middleware.JWTAuth("aaa"))
	// создание нко
	e.POST("/api/nkocreate", nkoHandler.CreateNKO, middleware.JWTAuth("aaa"))

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}
