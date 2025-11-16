package models

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	FullName     string `json:"full_name"`
	NgoID        *int   `json:"ngo_id,omitempty"`
	Role         string `json:"role"`
	Created_at   string `json:"created_at"`
}
