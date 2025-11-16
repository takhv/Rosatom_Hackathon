package repository

import (
	"database/sql"

	"rosatom.ru/nko/internal/models"
)

type citiesRepo struct {
	db *sql.DB
}

func NewCityRepo(db *sql.DB) CitiesRepo {
	return &citiesRepo{
		db: db,
	}
}

func (r *citiesRepo) GetAllCities() ([]models.City, error) {
	rows, err := r.db.Query("SELECT * FROM cities")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var citiesList []models.City
	for rows.Next() {
		var city models.City
		err := rows.Scan(&city.ID, &city.Name, &city.Region)
		if err != nil {
			return nil, err
		}

		citiesList = append(citiesList, city)
	}
	return citiesList, nil
}

func (r *citiesRepo) GetByID(id int) (*models.City, error) {
	var city models.City
	err := r.db.QueryRow("SELECT * FROM ngos WHERE id = ?", id).Scan(&city.ID, &city.Name, &city.Region)
	if err != nil {
		return nil, err
	}

	return &city, nil
}
