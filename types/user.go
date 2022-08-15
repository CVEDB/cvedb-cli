package types

import (
	"github.com/google/uuid"
	"time"
)

type Config struct {
	User struct {
		Token   string
		VaultId uuid.UUID
	}
	BaseUrl string
}

type User struct {
	ID            int     `json:"id"`
	IsActive      bool    `json:"is_active"`
	Email         string  `json:"email"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	Onboarding    bool    `json:"onboarding"`
	Profile       Profile `json:"profile"`
	InitialCredit int     `json:"initial_credit"`
}

type Profile struct {
	VaultInfo  VaultInfo `json:"vault_info"`
	Bio        string    `json:"bio"`
	Type       int       `json:"type"`
	Username   string    `json:"username"`
	EntityType string    `json:"entity_type"`
}

type VaultInfo struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Type         int       `json:"type"`
	Metadata     string    `json:"metadata"`
	CreatedDate  time.Time `json:"created_date"`
	ModifiedDate time.Time `json:"modified_date"`
}

type Hives struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Page     int    `json:"page"`
	Last     int    `json:"last"`
	Count    int    `json:"count"`
	Results  []Hive `json:"results"`
}

type Hive struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	HiveType     string    `json:"hive_type"`
	VaultInfo    uuid.UUID `json:"vault_info"`
	Cluster      string    `json:"cluster"`
	State        string    `json:"state"`
	CreatedDate  time.Time `json:"created_date"`
	ModifiedDate time.Time `json:"modified_date"`
	Bees         []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Mem         string `json:"mem"`
		CPU         string `json:"cpu"`
		Total       int    `json:"total"`
		Running     int    `json:"running"`
		Up          int    `json:"up"`
		Down        int    `json:"down"`
		Error       int    `json:"error"`
	} `json:"bees"`
}
