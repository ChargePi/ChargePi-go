package models

import (
	"testing"
)

func TestValidateRole(t *testing.T) {
	tests := []struct {
		name    string
		role    string
		wantErr bool
	}{
		{
			name:    "Manufacturer",
			role:    "Manufacturer",
			wantErr: false,
		},
		{
			name:    "Technician",
			role:    "Technician",
			wantErr: false,
		},
		{
			name:    "Observer",
			role:    "Observer",
			wantErr: false,
		},
		{
			name:    "Invalid",
			role:    "Invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateRole(tt.role); (err != nil) != tt.wantErr {
				t.Errorf("ValidateRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
