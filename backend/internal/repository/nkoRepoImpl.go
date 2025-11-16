package repository

import (
	"database/sql"
	"time"

	"rosatom.ru/nko/internal/models"
)

type nkoRepo struct {
	db *sql.DB
}

func NewNKORepo(db *sql.DB) NKORepo {
	return &nkoRepo{
		db: db,
	}
}

func (r *nkoRepo) GetAllNKO() ([]models.NKO, error) {
	rows, err := r.db.Query("SELECT * FROM ngos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nkoList []models.NKO
	for rows.Next() {
		var nko models.NKO
		var desc, volDesc, phone, addr, logo, website, social sql.NullString
		err := rows.Scan(
			&nko.ID, &nko.Name, &nko.Category,
			&desc, &volDesc, &phone, &addr, &logo, &website, &social,
			&nko.City_id, &nko.Status, &nko.Created_at, &nko.Updated_at,
		)
		if err != nil {
			return nil, err
		}

		nko.Description = desc.String
		nko.Volunteer_description = volDesc.String
		nko.Phone = phone.String
		nko.Address = addr.String
		nko.Logo_url = logo.String
		nko.Website_url = website.String
		nko.Social_links = social.String

		nkoList = append(nkoList, nko)
	}
	return nkoList, nil
}

func (r *nkoRepo) GetByID(id int) (*models.NKO, error) {
	var nko models.NKO
	var desc, volDesc, phone, addr, logo, website, social sql.NullString
	err := r.db.QueryRow("SELECT * FROM ngos WHERE id = ?", id).Scan(
		&nko.ID, &nko.Name, &nko.Category,
		&desc, &volDesc, &phone, &addr, &logo, &website, &social,
		&nko.City_id, &nko.Status, &nko.Created_at, &nko.Updated_at,
	)
	if err != nil {
		return nil, err
	}

	nko.Description = desc.String
	nko.Volunteer_description = volDesc.String
	nko.Phone = phone.String
	nko.Address = addr.String
	nko.Logo_url = logo.String
	nko.Website_url = website.String
	nko.Social_links = social.String

	return &nko, nil
}

func (r *nkoRepo) GetNKOName(id int) (string, error) {
	var nkoName string
	err := r.db.QueryRow("SELECT name FROM ngos WHERE id = ?", id).Scan(&nkoName)
	if err != nil {
		return "noname", err
	}

	return nkoName, nil
}

func (r *nkoRepo) SearchNKO(name string, category string) ([]models.NKO, error) {
	query := "SELECT * FROM ngos WHERE 1=1"
	args := []interface{}{}

	if name != "" {
		query += " AND name LIKE ?"
		args = append(args, "%"+name+"%")
	}

	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ngos []models.NKO
	for rows.Next() {
		var nko models.NKO
		var desc, volDesc, phone, addr, logo, website, social sql.NullString

		err := rows.Scan(
			&nko.ID, &nko.Name, &nko.Category,
			&desc, &volDesc, &phone, &addr, &logo, &website, &social,
			&nko.City_id, &nko.Status, &nko.Created_at, &nko.Updated_at,
		)
		if err != nil {
			return nil, err
		}

		nko.Description = desc.String
		nko.Volunteer_description = volDesc.String
		nko.Phone = phone.String
		nko.Address = addr.String
		nko.Logo_url = logo.String
		nko.Website_url = website.String
		nko.Social_links = social.String

		ngos = append(ngos, nko)
	}

	if len(ngos) == 0 {
		return nil, err
	}

	return ngos, nil
}

func (r *nkoRepo) CreateNKO(nko *models.NKO) error {
	query := `
        INSERT INTO ngos (
            name, category, description, volunteer_description, 
            phone, address, logo_url, website_url, social_links, city_id
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

	result, err := r.db.Exec(
		query,
		nko.Name,
		nko.Category,
		toNullString(nko.Description),
		toNullString(nko.Volunteer_description),
		toNullString(nko.Phone),
		toNullString(nko.Address),
		toNullString(nko.Logo_url),
		toNullString(nko.Website_url),
		toNullString(nko.Social_links),
		nko.City_id,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	nko.ID = int(id)
	nko.Status = "pending"
	nko.Created_at = time.Now().Format("2006-01-02 15:04:05")
	nko.Updated_at = nko.Created_at

	return nil
}

func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
