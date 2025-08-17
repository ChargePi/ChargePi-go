package display

import (
	"github.com/emirpasic/gods/queues/priorityqueue"
	"github.com/emirpasic/gods/utils"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
	"github.com/pkg/errors"
)

var (
	ErrNoMessages = errors.New("no messages in queue")
	ErrorDequeue  = errors.New("cannot dequeue message")
)

type messagePriorityQueue struct {
	priorityQueue *priorityqueue.Queue
}

// newMessageQueue creates a new MessageStore instance.
func newMessageQueue() *messagePriorityQueue {
	return &messagePriorityQueue{
		priorityQueue: priorityqueue.NewWith(compareMessages),
	}
}

func compareMessages(a, b interface{}) int {
	messageA := a.(display.MessageInfo)
	messageB := b.(display.MessageInfo)

	return utils.IntComparator(messageA.Priority, messageB.Priority)
}

func (ms *messagePriorityQueue) QueueMessage(message display.MessageInfo) {
	ms.priorityQueue.Enqueue(message)
}

func (ms *messagePriorityQueue) GetNextMessage() (*display.MessageInfo, error) {
	if ms.priorityQueue.Size() == 0 {
		return nil, ErrNoMessages
	}

	val, err := ms.priorityQueue.Dequeue()
	if !err {
		return nil, ErrorDequeue
	}

	cast := val.(display.MessageInfo)
	return &cast, nil
}

func (ms *messagePriorityQueue) GetMessagesInQueue() ([]display.MessageInfo, error) {
	var messages []display.MessageInfo

	// Iterate over all messages
	for _, message := range ms.priorityQueue.Values() {
		cast, canCast := message.(display.MessageInfo)
		if !canCast {
			return nil, errors.New("error casting message")
		}

		messages = append(messages, cast)
	}

	return messages, nil
}
