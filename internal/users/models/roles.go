package models

import "errors"

type Role string

const (
	Manufacturer = Role("Manufacturer")
	Technician   = Role("Technician")
	Observer     = Role("Observer")
)

var ErrRoleDoesntExist = errors.New("role does not exist")

func ValidateRole(role string) error {
	switch Role(role) {
	case Manufacturer, Technician, Observer:
		return nil
	default:
		return ErrRoleDoesntExist
	}
}
