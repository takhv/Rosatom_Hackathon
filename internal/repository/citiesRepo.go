package repository

import (
	"rosatom.ru/nko/internal/models"
)

type CitiesRepo interface {
	GetAllCities() ([]models.City, error)
	GetByID(id int) (*models.City, error)
}
