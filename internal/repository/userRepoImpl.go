package repository

import (
	"database/sql"
	"fmt"
	"time"

	"rosatom.ru/nko/internal/models"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	result, err := r.db.Exec(`
        INSERT INTO users (email, password_hash, full_name, ngo_id, role) 
        VALUES (?, ?, ?, ?, ?)`,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.NgoID,
		user.Role,
	)

	if err != nil {
		return err
	}

	user.Created_at = time.Now().Format("2006-01-02 15:04:05")

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	var ngoID sql.NullInt64

	err := r.db.QueryRow(`
        SELECT id, email, password_hash, full_name, ngo_id, role, created_at 
        FROM users WHERE email = ?
    `, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&ngoID,
		&user.Role,
		&user.Created_at,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	if ngoID.Valid {
		ngoIDInt := int(ngoID.Int64)
		user.NgoID = &ngoIDInt
	}

	return &user, nil
}

func (r *userRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	var ngoID sql.NullInt64

	err := r.db.QueryRow(`
        SELECT id, email, password_hash, full_name, ngo_id, role, created_at 
        FROM users WHERE id = ?
    `, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&ngoID,
		&user.Role,
		&user.Created_at,
	)

	if err != nil {
		return nil, err
	}

	if ngoID.Valid {
		ngoIDInt := int(ngoID.Int64)
		user.NgoID = &ngoIDInt
	}

	return &user, nil
}
