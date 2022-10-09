package test

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/charge-point"
)

type (
	DisplayMock struct {
		mock.Mock
	}
)

/*------------------ Display mock ------------------*/

func (l *DisplayMock) DisplayMessage(message chargePoint.Message) {
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

func (l *DisplayMock) GetLcdChannel() chan<- chargePoint.Message {
	return l.Called().Get(0).(chan chargePoint.Message)
}

func (l *DisplayMock) GetType() string {
	return ""
}
