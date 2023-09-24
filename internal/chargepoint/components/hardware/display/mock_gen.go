// Code generated by mocktail; DO NOT EDIT.

package display

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/notifications"
)

// displayMock mock of Display.
type displayMock struct{ mock.Mock }

// NewDisplayMock creates a new displayMock.
func NewDisplayMock(tb testing.TB) *displayMock {
	tb.Helper()

	m := &displayMock{}
	m.Mock.Test(tb)

	tb.Cleanup(func() { m.AssertExpectations(tb) })

	return m
}

func (_m *displayMock) Cleanup() {
	_m.Called()
}

func (_m *displayMock) OnCleanup() *displayCleanupCall {
	return &displayCleanupCall{Call: _m.Mock.On("Cleanup"), Parent: _m}
}

func (_m *displayMock) OnCleanupRaw() *displayCleanupCall {
	return &displayCleanupCall{Call: _m.Mock.On("Cleanup"), Parent: _m}
}

type displayCleanupCall struct {
	*mock.Call
	Parent *displayMock
}

func (_c *displayCleanupCall) Panic(msg string) *displayCleanupCall {
	_c.Call = _c.Call.Panic(msg)
	return _c
}

func (_c *displayCleanupCall) Once() *displayCleanupCall {
	_c.Call = _c.Call.Once()
	return _c
}

func (_c *displayCleanupCall) Twice() *displayCleanupCall {
	_c.Call = _c.Call.Twice()
	return _c
}

func (_c *displayCleanupCall) Times(i int) *displayCleanupCall {
	_c.Call = _c.Call.Times(i)
	return _c
}

func (_c *displayCleanupCall) WaitUntil(w <-chan time.Time) *displayCleanupCall {
	_c.Call = _c.Call.WaitUntil(w)
	return _c
}

func (_c *displayCleanupCall) After(d time.Duration) *displayCleanupCall {
	_c.Call = _c.Call.After(d)
	return _c
}

func (_c *displayCleanupCall) Run(fn func(args mock.Arguments)) *displayCleanupCall {
	_c.Call = _c.Call.Run(fn)
	return _c
}

func (_c *displayCleanupCall) Maybe() *displayCleanupCall {
	_c.Call = _c.Call.Maybe()
	return _c
}

func (_c *displayCleanupCall) TypedRun(fn func()) *displayCleanupCall {
	_c.Call = _c.Call.Run(func(args mock.Arguments) {
		fn()
	})
	return _c
}

func (_c *displayCleanupCall) OnCleanup() *displayCleanupCall {
	return _c.Parent.OnCleanup()
}

func (_c *displayCleanupCall) OnClear() *displayClearCall {
	return _c.Parent.OnClear()
}

func (_c *displayCleanupCall) OnDisplayMessage(message notifications.Message) *displayDisplayMessageCall {
	return _c.Parent.OnDisplayMessage(message)
}

func (_c *displayCleanupCall) OnGetType() *displayGetTypeCall {
	return _c.Parent.OnGetType()
}

func (_c *displayCleanupCall) OnCleanupRaw() *displayCleanupCall {
	return _c.Parent.OnCleanupRaw()
}

func (_c *displayCleanupCall) OnClearRaw() *displayClearCall {
	return _c.Parent.OnClearRaw()
}

func (_c *displayCleanupCall) OnDisplayMessageRaw(message interface{}) *displayDisplayMessageCall {
	return _c.Parent.OnDisplayMessageRaw(message)
}

func (_c *displayCleanupCall) OnGetTypeRaw() *displayGetTypeCall {
	return _c.Parent.OnGetTypeRaw()
}

func (_m *displayMock) Clear() {
	_m.Called()
}

func (_m *displayMock) OnClear() *displayClearCall {
	return &displayClearCall{Call: _m.Mock.On("Clear"), Parent: _m}
}

func (_m *displayMock) OnClearRaw() *displayClearCall {
	return &displayClearCall{Call: _m.Mock.On("Clear"), Parent: _m}
}

type displayClearCall struct {
	*mock.Call
	Parent *displayMock
}

func (_c *displayClearCall) Panic(msg string) *displayClearCall {
	_c.Call = _c.Call.Panic(msg)
	return _c
}

func (_c *displayClearCall) Once() *displayClearCall {
	_c.Call = _c.Call.Once()
	return _c
}

func (_c *displayClearCall) Twice() *displayClearCall {
	_c.Call = _c.Call.Twice()
	return _c
}

func (_c *displayClearCall) Times(i int) *displayClearCall {
	_c.Call = _c.Call.Times(i)
	return _c
}

func (_c *displayClearCall) WaitUntil(w <-chan time.Time) *displayClearCall {
	_c.Call = _c.Call.WaitUntil(w)
	return _c
}

func (_c *displayClearCall) After(d time.Duration) *displayClearCall {
	_c.Call = _c.Call.After(d)
	return _c
}

func (_c *displayClearCall) Run(fn func(args mock.Arguments)) *displayClearCall {
	_c.Call = _c.Call.Run(fn)
	return _c
}

func (_c *displayClearCall) Maybe() *displayClearCall {
	_c.Call = _c.Call.Maybe()
	return _c
}

func (_c *displayClearCall) TypedRun(fn func()) *displayClearCall {
	_c.Call = _c.Call.Run(func(args mock.Arguments) {
		fn()
	})
	return _c
}

func (_c *displayClearCall) OnCleanup() *displayCleanupCall {
	return _c.Parent.OnCleanup()
}

func (_c *displayClearCall) OnClear() *displayClearCall {
	return _c.Parent.OnClear()
}

func (_c *displayClearCall) OnDisplayMessage(message notifications.Message) *displayDisplayMessageCall {
	return _c.Parent.OnDisplayMessage(message)
}

func (_c *displayClearCall) OnGetType() *displayGetTypeCall {
	return _c.Parent.OnGetType()
}

func (_c *displayClearCall) OnCleanupRaw() *displayCleanupCall {
	return _c.Parent.OnCleanupRaw()
}

func (_c *displayClearCall) OnClearRaw() *displayClearCall {
	return _c.Parent.OnClearRaw()
}

func (_c *displayClearCall) OnDisplayMessageRaw(message interface{}) *displayDisplayMessageCall {
	return _c.Parent.OnDisplayMessageRaw(message)
}

func (_c *displayClearCall) OnGetTypeRaw() *displayGetTypeCall {
	return _c.Parent.OnGetTypeRaw()
}

func (_m *displayMock) DisplayMessage(message notifications.Message) {
	_m.Called(message)
}

func (_m *displayMock) OnDisplayMessage(message notifications.Message) *displayDisplayMessageCall {
	return &displayDisplayMessageCall{Call: _m.Mock.On("DisplayMessage", message), Parent: _m}
}

func (_m *displayMock) OnDisplayMessageRaw(message interface{}) *displayDisplayMessageCall {
	return &displayDisplayMessageCall{Call: _m.Mock.On("DisplayMessage", message), Parent: _m}
}

type displayDisplayMessageCall struct {
	*mock.Call
	Parent *displayMock
}

func (_c *displayDisplayMessageCall) Panic(msg string) *displayDisplayMessageCall {
	_c.Call = _c.Call.Panic(msg)
	return _c
}

func (_c *displayDisplayMessageCall) Once() *displayDisplayMessageCall {
	_c.Call = _c.Call.Once()
	return _c
}

func (_c *displayDisplayMessageCall) Twice() *displayDisplayMessageCall {
	_c.Call = _c.Call.Twice()
	return _c
}

func (_c *displayDisplayMessageCall) Times(i int) *displayDisplayMessageCall {
	_c.Call = _c.Call.Times(i)
	return _c
}

func (_c *displayDisplayMessageCall) WaitUntil(w <-chan time.Time) *displayDisplayMessageCall {
	_c.Call = _c.Call.WaitUntil(w)
	return _c
}

func (_c *displayDisplayMessageCall) After(d time.Duration) *displayDisplayMessageCall {
	_c.Call = _c.Call.After(d)
	return _c
}

func (_c *displayDisplayMessageCall) Run(fn func(args mock.Arguments)) *displayDisplayMessageCall {
	_c.Call = _c.Call.Run(fn)
	return _c
}

func (_c *displayDisplayMessageCall) Maybe() *displayDisplayMessageCall {
	_c.Call = _c.Call.Maybe()
	return _c
}

func (_c *displayDisplayMessageCall) TypedRun(fn func(notifications.Message)) *displayDisplayMessageCall {
	_c.Call = _c.Call.Run(func(args mock.Arguments) {
		_message, _ := args.Get(0).(notifications.Message)
		fn(_message)
	})
	return _c
}

func (_c *displayDisplayMessageCall) OnCleanup() *displayCleanupCall {
	return _c.Parent.OnCleanup()
}

func (_c *displayDisplayMessageCall) OnClear() *displayClearCall {
	return _c.Parent.OnClear()
}

func (_c *displayDisplayMessageCall) OnDisplayMessage(message notifications.Message) *displayDisplayMessageCall {
	return _c.Parent.OnDisplayMessage(message)
}

func (_c *displayDisplayMessageCall) OnGetType() *displayGetTypeCall {
	return _c.Parent.OnGetType()
}

func (_c *displayDisplayMessageCall) OnCleanupRaw() *displayCleanupCall {
	return _c.Parent.OnCleanupRaw()
}

func (_c *displayDisplayMessageCall) OnClearRaw() *displayClearCall {
	return _c.Parent.OnClearRaw()
}

func (_c *displayDisplayMessageCall) OnDisplayMessageRaw(message interface{}) *displayDisplayMessageCall {
	return _c.Parent.OnDisplayMessageRaw(message)
}

func (_c *displayDisplayMessageCall) OnGetTypeRaw() *displayGetTypeCall {
	return _c.Parent.OnGetTypeRaw()
}

func (_m *displayMock) GetType() string {
	_ret := _m.Called()

	if _rf, ok := _ret.Get(0).(func() string); ok {
		return _rf()
	}

	_ra0 := _ret.String(0)

	return _ra0
}

func (_m *displayMock) OnGetType() *displayGetTypeCall {
	return &displayGetTypeCall{Call: _m.Mock.On("GetType"), Parent: _m}
}

func (_m *displayMock) OnGetTypeRaw() *displayGetTypeCall {
	return &displayGetTypeCall{Call: _m.Mock.On("GetType"), Parent: _m}
}

type displayGetTypeCall struct {
	*mock.Call
	Parent *displayMock
}

func (_c *displayGetTypeCall) Panic(msg string) *displayGetTypeCall {
	_c.Call = _c.Call.Panic(msg)
	return _c
}

func (_c *displayGetTypeCall) Once() *displayGetTypeCall {
	_c.Call = _c.Call.Once()
	return _c
}

func (_c *displayGetTypeCall) Twice() *displayGetTypeCall {
	_c.Call = _c.Call.Twice()
	return _c
}

func (_c *displayGetTypeCall) Times(i int) *displayGetTypeCall {
	_c.Call = _c.Call.Times(i)
	return _c
}

func (_c *displayGetTypeCall) WaitUntil(w <-chan time.Time) *displayGetTypeCall {
	_c.Call = _c.Call.WaitUntil(w)
	return _c
}

func (_c *displayGetTypeCall) After(d time.Duration) *displayGetTypeCall {
	_c.Call = _c.Call.After(d)
	return _c
}

func (_c *displayGetTypeCall) Run(fn func(args mock.Arguments)) *displayGetTypeCall {
	_c.Call = _c.Call.Run(fn)
	return _c
}

func (_c *displayGetTypeCall) Maybe() *displayGetTypeCall {
	_c.Call = _c.Call.Maybe()
	return _c
}

func (_c *displayGetTypeCall) TypedReturns(a string) *displayGetTypeCall {
	_c.Call = _c.Return(a)
	return _c
}

func (_c *displayGetTypeCall) ReturnsFn(fn func() string) *displayGetTypeCall {
	_c.Call = _c.Return(fn)
	return _c
}

func (_c *displayGetTypeCall) TypedRun(fn func()) *displayGetTypeCall {
	_c.Call = _c.Call.Run(func(args mock.Arguments) {
		fn()
	})
	return _c
}

func (_c *displayGetTypeCall) OnCleanup() *displayCleanupCall {
	return _c.Parent.OnCleanup()
}

func (_c *displayGetTypeCall) OnClear() *displayClearCall {
	return _c.Parent.OnClear()
}

func (_c *displayGetTypeCall) OnDisplayMessage(message notifications.Message) *displayDisplayMessageCall {
	return _c.Parent.OnDisplayMessage(message)
}

func (_c *displayGetTypeCall) OnGetType() *displayGetTypeCall {
	return _c.Parent.OnGetType()
}

func (_c *displayGetTypeCall) OnCleanupRaw() *displayCleanupCall {
	return _c.Parent.OnCleanupRaw()
}

func (_c *displayGetTypeCall) OnClearRaw() *displayClearCall {
	return _c.Parent.OnClearRaw()
}

func (_c *displayGetTypeCall) OnDisplayMessageRaw(message interface{}) *displayDisplayMessageCall {
	return _c.Parent.OnDisplayMessageRaw(message)
}

func (_c *displayGetTypeCall) OnGetTypeRaw() *displayGetTypeCall {
	return _c.Parent.OnGetTypeRaw()
}