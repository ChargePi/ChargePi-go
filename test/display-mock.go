package test

import (
	"github.com/stretchr/testify/mock"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/notifications"
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

func (l *DisplayMock) Cleanup() {
	l.Called()
}

func (l *DisplayMock) Clear() {
	l.Called()
}

func (l *DisplayMock) GetType() string {
	return ""
}
