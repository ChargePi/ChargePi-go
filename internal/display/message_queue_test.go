package display

import (
	"testing"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type messageQueueTestSuite struct {
	suite.Suite
	queue *messagePriorityQueue
}

func (s *messageQueueTestSuite) SetupTest() {
	s.queue = newMessageQueue()
}

func (s *messageQueueTestSuite) TestQueueMessage() {
	tests := []struct {
		name    string
		message display.MessageInfo
	}{
		{
			name: "Queue message",
		},
		{
			name: "Queue another message",
		},
		{
			name: "Queue message with higher priority",
		},
		{
			name: "Queue message with lower priority",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s.queue.QueueMessage(tt.message)
			s.Contains(s.queue.priorityQueue.Values(), tt.message)
			// todo check priority
		})
	}
}

func (s *messageQueueTestSuite) TestGetNextMessage() {
	tests := []struct {
		name            string
		elementsToQueue []display.MessageInfo
		expected        []display.MessageInfo
		err             error
	}{
		{
			name: "No messages in queue",
		},
		{
			name: "Message dequeued",
		},
		{
			name: "Queue message with higher priority",
		},
		{
			name: "Queue message with lower priority",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			for _, message := range tt.elementsToQueue {
				s.queue.QueueMessage(message)
			}

			message, err := s.queue.GetNextMessage()
			if tt.err != nil {
				s.Error(err)
			} else {
				s.NoError(err)
				s.Contains(tt.expected, *message)
				// todo check priority
			}
		})
	}
}

func (s *messageQueueTestSuite) TestGetMessagesInQueue() {
	tests := []struct {
		name            string
		elementsToQueue []display.MessageInfo
		expected        []display.MessageInfo
		error           error
	}{
		{
			name: "Get messages in queue",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestMessageQueue(t *testing.T) {
	suite.Run(t, new(messageQueueTestSuite))
}

func Test_compareMessages(t *testing.T) {
	tests := []struct {
		name     string
		messageA display.MessageInfo
		messageB display.MessageInfo
		want     int
	}{
		{
			name:     "Message A has priority",
			messageA: display.MessageInfo{},
			messageB: display.MessageInfo{},
		},

		{
			name:     "Message B has priority",
			messageA: display.MessageInfo{},
			messageB: display.MessageInfo{},
		},
		{
			name:     "Message A and B have same priority",
			messageA: display.MessageInfo{},
			messageB: display.MessageInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := compareMessages(tt.messageA, tt.messageB)
			assert.Equal(t, tt.want, res)
		})
	}
}
