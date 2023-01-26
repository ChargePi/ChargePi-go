package models

type User struct {
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
	Role     string `json:"role" mapstructure:"role"`
}

type Role string

const (
	Manufacturer = Role("Manufacturer")
	Technician   = Role("Technician")
	Observer     = Role("Observer")
)
