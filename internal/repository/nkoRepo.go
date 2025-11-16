package repository

import (
	"rosatom.ru/nko/internal/models"
)

type NKORepo interface {
	GetAllNKO() ([]models.NKO, error)
	GetByID(id int) (*models.NKO, error)
	GetNKOName(id int) (string, error)
	SearchNKO(name string, category string) ([]models.NKO, error)
	CreateNKO(nko *models.NKO) error
}
