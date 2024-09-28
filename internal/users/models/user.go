package models

type User struct {
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
	Role     Role   `json:"role" mapstructure:"role"`
}

func (u User) IsManufacturer() bool {
	return u.Role == Manufacturer
}

func (u User) IsTechnician() bool {
	return u.Role == Technician
}

func (u User) IsObserver() bool {
	return u.Role == Observer
}
