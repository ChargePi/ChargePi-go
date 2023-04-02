package database

type SessionRepository interface {
	AddSession()
	GetSession()
	UpdateSession()
	DeleteSession()
}
