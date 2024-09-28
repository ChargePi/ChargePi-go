package manager

import (
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

var (
	ErrInvalidReservation = errors.New("invalid reservation")
	ErrValidationFailed   = errors.New("reservation validation failed")
)

type reservation struct {
	// EvseId is mandatory
	EvseId int `validate:"required,gte=1"`

	// Connector ID is optional
	ConnectorId *int

	// TagId used for authentication
	TagId string `validate:"gte=1,lte=20"`
}

func (r *reservation) Validate() bool {
	if r.TagId == "" {
		return false
	}

	err := validator.New().Struct(r)
	if err != nil {
		return false
	}

	return true
}

type reservations struct {
	mu           sync.Mutex
	reservations map[int]reservation
}

func newReservationsStore() *reservations {
	return &reservations{
		reservations: make(map[int]reservation),
	}
}

func (r *reservations) GetEVSEWithReservationId(reservationId int) (*reservation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// todo persist reservations?

	// Check if reservation with id exists
	reservation, isFound := r.reservations[reservationId]
	if !isFound {
		return nil, ErrReservationNotFound
	}

	return &reservation, nil
}

func (r *reservations) Reserve(reservationId int, reservationDetails reservation) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validate reservation
	if !reservationDetails.Validate() {
		return ErrValidationFailed
	}

	// Check if reservation with id already exists
	_, isFound := r.reservations[reservationId]
	if isFound {
		return ErrReservationNotFound
	}

	// Check if reservation with evse id already exists
	for _, res := range r.reservations {
		if res.EvseId == reservationDetails.EvseId {
			return ErrInvalidReservation
		}
	}

	r.reservations[reservationId] = reservationDetails
	return nil
}

func (r *reservations) RemoveReservation(reservationId int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if reservation with id exists
	_, isFound := r.reservations[reservationId]
	if !isFound {
		return ErrReservationNotFound
	}

	delete(r.reservations, reservationId)
	return nil
}

func (r *reservations) RemoveReservationsForEvse(evseId int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	hasReservation := false

	// Check if reservation with evse id exists
	for id, res := range r.reservations {
		if res.EvseId == evseId {
			hasReservation = true
			delete(r.reservations, id)
		}
	}

	if !hasReservation {
		return ErrReservationNotFound
	}

	return nil
}
