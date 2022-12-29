package test

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type ReaderMock struct {
	mock.Mock
}

func (r *ReaderMock) ListenForTags(ctx context.Context) {
	r.Called()
}

func (r *ReaderMock) Cleanup() {
	r.Called()
}

func (r *ReaderMock) Reset() {
	r.Called()
}

func (r *ReaderMock) GetTagChannel() <-chan string {
	return r.Called().Get(0).(chan string)
}
