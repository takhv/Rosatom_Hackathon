package models

type CreateNKO struct {
	Name                 string `json:"name"`
	Category             string `json:"category"`
	Description          string `json:"description"`
	VolunteerDescription string `json:"volunteer_description"`
	Phone                string `json:"phone"`
	Address              string `json:"address"`
	LogoURL              string `json:"logo_url"`
	WebsiteURL           string `json:"website_url"`
	SocialLinks          string `json:"social_links"`
	CityID               int    `json:"city_id"`
}
