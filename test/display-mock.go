package test

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/notifications"
)

type (
	DisplayMock struct {
		mock.Mock
	}
)

/*------------------ Display mock ------------------*/

func (l *DisplayMock) DisplayMessage(message notifications.Message) {
	l.Called(message)
}

func (l *DisplayMock) ListenForMessages(ctx context.Context) {
	l.Called()
}

func (l *DisplayMock) Cleanup() {
	l.Called()
}

func (l *DisplayMock) Clear() {
	l.Called()
}

func (l *DisplayMock) GetLcdChannel() chan<- notifications.Message {
	return l.Called().Get(0).(chan notifications.Message)
}

func (l *DisplayMock) GetType() string {
	return ""
}
