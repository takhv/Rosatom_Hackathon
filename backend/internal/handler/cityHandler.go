package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"rosatom.ru/nko/internal/repository"
)

type CitiesHandler struct {
	repo repository.CitiesRepo
}

func NewCitiesHandler(repo repository.CitiesRepo) *CitiesHandler {
	return &CitiesHandler{repo: repo}
}

func (h *CitiesHandler) GetAllNKO(context echo.Context) error {
	cities, err := h.repo.GetAllCities()
	if err != nil {
		return context.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return context.JSON(http.StatusOK, cities)
}

func (h *CitiesHandler) GetByID(context echo.Context) error {
	idStr := context.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return context.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	city, err := h.repo.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return context.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		return context.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Internal server error",
		})
	}
	return context.JSON(http.StatusOK, city)
}
