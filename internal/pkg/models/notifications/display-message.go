package notifications

import "time"

type (
	// Message Object representing the message that will be displayed on the Display.
	// Each array element in Messages represents a line being displayed on the 16x2 screen.
	Message struct {
		Messages        []string
		MessageDuration time.Duration
	}
)

// NewMessage creates a new message for the Display.
func NewMessage(duration time.Duration, messages []string) Message {
	return Message{
		Messages:        messages,
		MessageDuration: duration,
	}
}
