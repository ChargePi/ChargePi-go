package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_IsObserver(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		returns bool
	}{
		{
			name: "Is observer",
			user: User{
				Role: Observer,
			},
			returns: true,
		},
		{
			name: "Is manufacturer",
			user: User{
				Role: Manufacturer,
			},
			returns: false,
		},
		{
			name: "Is technician",
			user: User{
				Role: Technician,
			},
			returns: false,
		},
		{
			name: "Unknown role",
			user: User{
				Role: "Unknown",
			},
			returns: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.returns, tt.user.IsObserver())
		})
	}
}

func TestUser_IsTechnician(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		returns bool
	}{
		{
			name: "Is observer",
			user: User{
				Role: Observer,
			},
			returns: false,
		},
		{
			name: "Is manufacturer",
			user: User{
				Role: Manufacturer,
			},
			returns: false,
		},
		{
			name: "Is technician",
			user: User{
				Role: Technician,
			},
			returns: true,
		},
		{
			name: "Unknown role",
			user: User{
				Role: "Unknown",
			},
			returns: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.returns, tt.user.IsTechnician())
		})
	}
}

func TestUser_IsManufacturer(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		returns bool
	}{
		{
			name: "Is observer",
			user: User{
				Role: Observer,
			},
			returns: false,
		},
		{
			name: "Is manufacturer",
			user: User{
				Role: Manufacturer,
			},
			returns: true,
		},
		{
			name: "Is technician",
			user: User{
				Role: Technician,
			},
			returns: false,
		},
		{
			name: "Unknown role",
			user: User{
				Role: "Unknown",
			},
			returns: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.returns, tt.user.IsManufacturer())
		})
	}
}
