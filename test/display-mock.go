package test

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/xBlaz3kx/ChargePi-go/internal/models"
)

type (
	DisplayMock struct {
		mock.Mock
	}
)

/*------------------ Display mock ------------------*/

func (l *DisplayMock) DisplayMessage(message models.Message) {
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

func (l *DisplayMock) GetLcdChannel() chan<- models.Message {
	return l.Called().Get(0).(chan models.Message)
}
