// Code generated by mocktail; DO NOT EDIT.

package settings

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
)

// exporterMock mock of Exporter.
type exporterMock struct{ mock.Mock }

// NewExporterMock creates a new exporterMock.
func NewExporterMock(tb testing.TB) *exporterMock {
	tb.Helper()

	m := &exporterMock{}
	m.Mock.Test(tb)

	tb.Cleanup(func() { m.AssertExpectations(tb) })

	return m
}

func (_m *exporterMock) ExportEVSESettings() []settings.EVSE {
	_ret := _m.Called()

	if _rf, ok := _ret.Get(0).(func() []settings.EVSE); ok {
		return _rf()
	}

	_ra0, _ := _ret.Get(0).([]settings.EVSE)

	return _ra0
}

func (_m *exporterMock) OnExportEVSESettings() *exporterExportEVSESettingsCall {
	return &exporterExportEVSESettingsCall{Call: _m.Mock.On("ExportEVSESettings"), Parent: _m}
}

func (_m *exporterMock) OnExportEVSESettingsRaw() *exporterExportEVSESettingsCall {
	return &exporterExportEVSESettingsCall{Call: _m.Mock.On("ExportEVSESettings"), Parent: _m}
}

type exporterExportEVSESettingsCall struct {
	*mock.Call
	Parent *exporterMock
}

func (_c *exporterExportEVSESettingsCall) Panic(msg string) *exporterExportEVSESettingsCall {
	_c.Call = _c.Call.Panic(msg)
	return _c
}

func (_c *exporterExportEVSESettingsCall) Once() *exporterExportEVSESettingsCall {
	_c.Call = _c.Call.Once()
	return _c
}

func (_c *exporterExportEVSESettingsCall) Twice() *exporterExportEVSESettingsCall {
	_c.Call = _c.Call.Twice()
	return _c
}

func (_c *exporterExportEVSESettingsCall) Times(i int) *exporterExportEVSESettingsCall {
	_c.Call = _c.Call.Times(i)
	return _c
}

func (_c *exporterExportEVSESettingsCall) WaitUntil(w <-chan time.Time) *exporterExportEVSESettingsCall {
	_c.Call = _c.Call.WaitUntil(w)
	return _c
}

func (_c *exporterExportEVSESettingsCall) After(d time.Duration) *exporterExportEVSESettingsCall {
	_c.Call = _c.Call.After(d)
	return _c
}

func (_c *exporterExportEVSESettingsCall) Run(fn func(args mock.Arguments)) *exporterExportEVSESettingsCall {
	_c.Call = _c.Call.Run(fn)
	return _c
}

func (_c *exporterExportEVSESettingsCall) Maybe() *exporterExportEVSESettingsCall {
	_c.Call = _c.Call.Maybe()
	return _c
}

func (_c *exporterExportEVSESettingsCall) TypedReturns(a []settings.EVSE) *exporterExportEVSESettingsCall {
	_c.Call = _c.Return(a)
	return _c
}

func (_c *exporterExportEVSESettingsCall) ReturnsFn(fn func() []settings.EVSE) *exporterExportEVSESettingsCall {
	_c.Call = _c.Return(fn)
	return _c
}

func (_c *exporterExportEVSESettingsCall) TypedRun(fn func()) *exporterExportEVSESettingsCall {
	_c.Call = _c.Call.Run(func(args mock.Arguments) {
		fn()
	})
	return _c
}

func (_c *exporterExportEVSESettingsCall) OnExportEVSESettings() *exporterExportEVSESettingsCall {
	return _c.Parent.OnExportEVSESettings()
}

func (_c *exporterExportEVSESettingsCall) OnExportLocalAuthList() *exporterExportLocalAuthListCall {
	return _c.Parent.OnExportLocalAuthList()
}

func (_c *exporterExportEVSESettingsCall) OnExportOcppConfiguration() *exporterExportOcppConfigurationCall {
	return _c.Parent.OnExportOcppConfiguration()
}

func (_c *exporterExportEVSESettingsCall) OnExportEVSESettingsRaw() *exporterExportEVSESettingsCall {
	return _c.Parent.OnExportEVSESettingsRaw()
}

func (_c *exporterExportEVSESettingsCall) OnExportLocalAuthListRaw() *exporterExportLocalAuthListCall {
	return _c.Parent.OnExportLocalAuthListRaw()
}

func (_c *exporterExportEVSESettingsCall) OnExportOcppConfigurationRaw() *exporterExportOcppConfigurationCall {
	return _c.Parent.OnExportOcppConfigurationRaw()
}

func (_m *exporterMock) ExportLocalAuthList() (*settings.AuthList, error) {
	_ret := _m.Called()

	if _rf, ok := _ret.Get(0).(func() (*settings.AuthList, error)); ok {
		return _rf()
	}

	_ra0, _ := _ret.Get(0).(*settings.AuthList)
	_rb1 := _ret.Error(1)

	return _ra0, _rb1
}

func (_m *exporterMock) OnExportLocalAuthList() *exporterExportLocalAuthListCall {
	return &exporterExportLocalAuthListCall{Call: _m.Mock.On("ExportLocalAuthList"), Parent: _m}
}

func (_m *exporterMock) OnExportLocalAuthListRaw() *exporterExportLocalAuthListCall {
	return &exporterExportLocalAuthListCall{Call: _m.Mock.On("ExportLocalAuthList"), Parent: _m}
}

type exporterExportLocalAuthListCall struct {
	*mock.Call
	Parent *exporterMock
}

func (_c *exporterExportLocalAuthListCall) Panic(msg string) *exporterExportLocalAuthListCall {
	_c.Call = _c.Call.Panic(msg)
	return _c
}

func (_c *exporterExportLocalAuthListCall) Once() *exporterExportLocalAuthListCall {
	_c.Call = _c.Call.Once()
	return _c
}

func (_c *exporterExportLocalAuthListCall) Twice() *exporterExportLocalAuthListCall {
	_c.Call = _c.Call.Twice()
	return _c
}

func (_c *exporterExportLocalAuthListCall) Times(i int) *exporterExportLocalAuthListCall {
	_c.Call = _c.Call.Times(i)
	return _c
}

func (_c *exporterExportLocalAuthListCall) WaitUntil(w <-chan time.Time) *exporterExportLocalAuthListCall {
	_c.Call = _c.Call.WaitUntil(w)
	return _c
}

func (_c *exporterExportLocalAuthListCall) After(d time.Duration) *exporterExportLocalAuthListCall {
	_c.Call = _c.Call.After(d)
	return _c
}

func (_c *exporterExportLocalAuthListCall) Run(fn func(args mock.Arguments)) *exporterExportLocalAuthListCall {
	_c.Call = _c.Call.Run(fn)
	return _c
}

func (_c *exporterExportLocalAuthListCall) Maybe() *exporterExportLocalAuthListCall {
	_c.Call = _c.Call.Maybe()
	return _c
}

func (_c *exporterExportLocalAuthListCall) TypedReturns(a *settings.AuthList, b error) *exporterExportLocalAuthListCall {
	_c.Call = _c.Return(a, b)
	return _c
}

func (_c *exporterExportLocalAuthListCall) ReturnsFn(fn func() (*settings.AuthList, error)) *exporterExportLocalAuthListCall {
	_c.Call = _c.Return(fn)
	return _c
}

func (_c *exporterExportLocalAuthListCall) TypedRun(fn func()) *exporterExportLocalAuthListCall {
	_c.Call = _c.Call.Run(func(args mock.Arguments) {
		fn()
	})
	return _c
}

func (_c *exporterExportLocalAuthListCall) OnExportEVSESettings() *exporterExportEVSESettingsCall {
	return _c.Parent.OnExportEVSESettings()
}

func (_c *exporterExportLocalAuthListCall) OnExportLocalAuthList() *exporterExportLocalAuthListCall {
	return _c.Parent.OnExportLocalAuthList()
}

func (_c *exporterExportLocalAuthListCall) OnExportOcppConfiguration() *exporterExportOcppConfigurationCall {
	return _c.Parent.OnExportOcppConfiguration()
}

func (_c *exporterExportLocalAuthListCall) OnExportEVSESettingsRaw() *exporterExportEVSESettingsCall {
	return _c.Parent.OnExportEVSESettingsRaw()
}

func (_c *exporterExportLocalAuthListCall) OnExportLocalAuthListRaw() *exporterExportLocalAuthListCall {
	return _c.Parent.OnExportLocalAuthListRaw()
}

func (_c *exporterExportLocalAuthListCall) OnExportOcppConfigurationRaw() *exporterExportOcppConfigurationCall {
	return _c.Parent.OnExportOcppConfigurationRaw()
}

func (_m *exporterMock) ExportOcppConfiguration() configuration.Config {
	_ret := _m.Called()

	if _rf, ok := _ret.Get(0).(func() configuration.Config); ok {
		return _rf()
	}

	_ra0, _ := _ret.Get(0).(configuration.Config)

	return _ra0
}

func (_m *exporterMock) OnExportOcppConfiguration() *exporterExportOcppConfigurationCall {
	return &exporterExportOcppConfigurationCall{Call: _m.Mock.On("ExportOcppConfiguration"), Parent: _m}
}

func (_m *exporterMock) OnExportOcppConfigurationRaw() *exporterExportOcppConfigurationCall {
	return &exporterExportOcppConfigurationCall{Call: _m.Mock.On("ExportOcppConfiguration"), Parent: _m}
}

type exporterExportOcppConfigurationCall struct {
	*mock.Call
	Parent *exporterMock
}

func (_c *exporterExportOcppConfigurationCall) Panic(msg string) *exporterExportOcppConfigurationCall {
	_c.Call = _c.Call.Panic(msg)
	return _c
}

func (_c *exporterExportOcppConfigurationCall) Once() *exporterExportOcppConfigurationCall {
	_c.Call = _c.Call.Once()
	return _c
}

func (_c *exporterExportOcppConfigurationCall) Twice() *exporterExportOcppConfigurationCall {
	_c.Call = _c.Call.Twice()
	return _c
}

func (_c *exporterExportOcppConfigurationCall) Times(i int) *exporterExportOcppConfigurationCall {
	_c.Call = _c.Call.Times(i)
	return _c
}

func (_c *exporterExportOcppConfigurationCall) WaitUntil(w <-chan time.Time) *exporterExportOcppConfigurationCall {
	_c.Call = _c.Call.WaitUntil(w)
	return _c
}

func (_c *exporterExportOcppConfigurationCall) After(d time.Duration) *exporterExportOcppConfigurationCall {
	_c.Call = _c.Call.After(d)
	return _c
}

func (_c *exporterExportOcppConfigurationCall) Run(fn func(args mock.Arguments)) *exporterExportOcppConfigurationCall {
	_c.Call = _c.Call.Run(fn)
	return _c
}

func (_c *exporterExportOcppConfigurationCall) Maybe() *exporterExportOcppConfigurationCall {
	_c.Call = _c.Call.Maybe()
	return _c
}

func (_c *exporterExportOcppConfigurationCall) TypedReturns(a configuration.Config) *exporterExportOcppConfigurationCall {
	_c.Call = _c.Return(a)
	return _c
}

func (_c *exporterExportOcppConfigurationCall) ReturnsFn(fn func() configuration.Config) *exporterExportOcppConfigurationCall {
	_c.Call = _c.Return(fn)
	return _c
}

func (_c *exporterExportOcppConfigurationCall) TypedRun(fn func()) *exporterExportOcppConfigurationCall {
	_c.Call = _c.Call.Run(func(args mock.Arguments) {
		fn()
	})
	return _c
}

func (_c *exporterExportOcppConfigurationCall) OnExportEVSESettings() *exporterExportEVSESettingsCall {
	return _c.Parent.OnExportEVSESettings()
}

func (_c *exporterExportOcppConfigurationCall) OnExportLocalAuthList() *exporterExportLocalAuthListCall {
	return _c.Parent.OnExportLocalAuthList()
}

func (_c *exporterExportOcppConfigurationCall) OnExportOcppConfiguration() *exporterExportOcppConfigurationCall {
	return _c.Parent.OnExportOcppConfiguration()
}

func (_c *exporterExportOcppConfigurationCall) OnExportEVSESettingsRaw() *exporterExportEVSESettingsCall {
	return _c.Parent.OnExportEVSESettingsRaw()
}

func (_c *exporterExportOcppConfigurationCall) OnExportLocalAuthListRaw() *exporterExportLocalAuthListCall {
	return _c.Parent.OnExportLocalAuthListRaw()
}

func (_c *exporterExportOcppConfigurationCall) OnExportOcppConfigurationRaw() *exporterExportOcppConfigurationCall {
	return _c.Parent.OnExportOcppConfigurationRaw()
}

// importerMock mock of Importer.
type importerMock struct{ mock.Mock }

// NewImporterMock creates a new importerMock.
func NewImporterMock(tb testing.TB) *importerMock {
	tb.Helper()

	m := &importerMock{}
	m.Mock.Test(tb)

	tb.Cleanup(func() { m.AssertExpectations(tb) })

	return m
}

func (_m *importerMock) ImportEVSESettings(evseSettings []settings.EVSE) error {
	_ret := _m.Called(evseSettings)

	if _rf, ok := _ret.Get(0).(func([]settings.EVSE) error); ok {
		return _rf(evseSettings)
	}

	_ra0 := _ret.Error(0)

	return _ra0
}

func (_m *importerMock) OnImportEVSESettings(settings []settings.EVSE) *importerImportEVSESettingsCall {
	return &importerImportEVSESettingsCall{Call: _m.Mock.On("ImportEVSESettings", settings), Parent: _m}
}

func (_m *importerMock) OnImportEVSESettingsRaw(settings interface{}) *importerImportEVSESettingsCall {
	return &importerImportEVSESettingsCall{Call: _m.Mock.On("ImportEVSESettings", settings), Parent: _m}
}

type importerImportEVSESettingsCall struct {
	*mock.Call
	Parent *importerMock
}

func (_c *importerImportEVSESettingsCall) Panic(msg string) *importerImportEVSESettingsCall {
	_c.Call = _c.Call.Panic(msg)
	return _c
}

func (_c *importerImportEVSESettingsCall) Once() *importerImportEVSESettingsCall {
	_c.Call = _c.Call.Once()
	return _c
}

func (_c *importerImportEVSESettingsCall) Twice() *importerImportEVSESettingsCall {
	_c.Call = _c.Call.Twice()
	return _c
}

func (_c *importerImportEVSESettingsCall) Times(i int) *importerImportEVSESettingsCall {
	_c.Call = _c.Call.Times(i)
	return _c
}

func (_c *importerImportEVSESettingsCall) WaitUntil(w <-chan time.Time) *importerImportEVSESettingsCall {
	_c.Call = _c.Call.WaitUntil(w)
	return _c
}

func (_c *importerImportEVSESettingsCall) After(d time.Duration) *importerImportEVSESettingsCall {
	_c.Call = _c.Call.After(d)
	return _c
}

func (_c *importerImportEVSESettingsCall) Run(fn func(args mock.Arguments)) *importerImportEVSESettingsCall {
	_c.Call = _c.Call.Run(fn)
	return _c
}

func (_c *importerImportEVSESettingsCall) Maybe() *importerImportEVSESettingsCall {
	_c.Call = _c.Call.Maybe()
	return _c
}

func (_c *importerImportEVSESettingsCall) TypedReturns(a error) *importerImportEVSESettingsCall {
	_c.Call = _c.Return(a)
	return _c
}

func (_c *importerImportEVSESettingsCall) ReturnsFn(fn func([]settings.EVSE) error) *importerImportEVSESettingsCall {
	_c.Call = _c.Return(fn)
	return _c
}

func (_c *importerImportEVSESettingsCall) TypedRun(fn func([]settings.EVSE)) *importerImportEVSESettingsCall {
	_c.Call = _c.Call.Run(func(args mock.Arguments) {
		_settings, _ := args.Get(0).([]settings.EVSE)
		fn(_settings)
	})
	return _c
}

func (_c *importerImportEVSESettingsCall) OnImportEVSESettings(settings []settings.EVSE) *importerImportEVSESettingsCall {
	return _c.Parent.OnImportEVSESettings(settings)
}

func (_c *importerImportEVSESettingsCall) OnImportLocalAuthList(list settings.AuthList) *importerImportLocalAuthListCall {
	return _c.Parent.OnImportLocalAuthList(list)
}

func (_c *importerImportEVSESettingsCall) OnImportOcppConfiguration(config configuration.Config) *importerImportOcppConfigurationCall {
	return _c.Parent.OnImportOcppConfiguration(config)
}

func (_c *importerImportEVSESettingsCall) OnImportEVSESettingsRaw(settings interface{}) *importerImportEVSESettingsCall {
	return _c.Parent.OnImportEVSESettingsRaw(settings)
}

func (_c *importerImportEVSESettingsCall) OnImportLocalAuthListRaw(list interface{}) *importerImportLocalAuthListCall {
	return _c.Parent.OnImportLocalAuthListRaw(list)
}

func (_c *importerImportEVSESettingsCall) OnImportOcppConfigurationRaw(config interface{}) *importerImportOcppConfigurationCall {
	return _c.Parent.OnImportOcppConfigurationRaw(config)
}

func (_m *importerMock) ImportLocalAuthList(list settings.AuthList) error {
	_ret := _m.Called(list)

	if _rf, ok := _ret.Get(0).(func(settings.AuthList) error); ok {
		return _rf(list)
	}

	_ra0 := _ret.Error(0)

	return _ra0
}

func (_m *importerMock) OnImportLocalAuthList(list settings.AuthList) *importerImportLocalAuthListCall {
	return &importerImportLocalAuthListCall{Call: _m.Mock.On("ImportLocalAuthList", list), Parent: _m}
}

func (_m *importerMock) OnImportLocalAuthListRaw(list interface{}) *importerImportLocalAuthListCall {
	return &importerImportLocalAuthListCall{Call: _m.Mock.On("ImportLocalAuthList", list), Parent: _m}
}

type importerImportLocalAuthListCall struct {
	*mock.Call
	Parent *importerMock
}

func (_c *importerImportLocalAuthListCall) Panic(msg string) *importerImportLocalAuthListCall {
	_c.Call = _c.Call.Panic(msg)
	return _c
}

func (_c *importerImportLocalAuthListCall) Once() *importerImportLocalAuthListCall {
	_c.Call = _c.Call.Once()
	return _c
}

func (_c *importerImportLocalAuthListCall) Twice() *importerImportLocalAuthListCall {
	_c.Call = _c.Call.Twice()
	return _c
}

func (_c *importerImportLocalAuthListCall) Times(i int) *importerImportLocalAuthListCall {
	_c.Call = _c.Call.Times(i)
	return _c
}

func (_c *importerImportLocalAuthListCall) WaitUntil(w <-chan time.Time) *importerImportLocalAuthListCall {
	_c.Call = _c.Call.WaitUntil(w)
	return _c
}

func (_c *importerImportLocalAuthListCall) After(d time.Duration) *importerImportLocalAuthListCall {
	_c.Call = _c.Call.After(d)
	return _c
}

func (_c *importerImportLocalAuthListCall) Run(fn func(args mock.Arguments)) *importerImportLocalAuthListCall {
	_c.Call = _c.Call.Run(fn)
	return _c
}

func (_c *importerImportLocalAuthListCall) Maybe() *importerImportLocalAuthListCall {
	_c.Call = _c.Call.Maybe()
	return _c
}

func (_c *importerImportLocalAuthListCall) TypedReturns(a error) *importerImportLocalAuthListCall {
	_c.Call = _c.Return(a)
	return _c
}

func (_c *importerImportLocalAuthListCall) ReturnsFn(fn func(settings.AuthList) error) *importerImportLocalAuthListCall {
	_c.Call = _c.Return(fn)
	return _c
}

func (_c *importerImportLocalAuthListCall) TypedRun(fn func(settings.AuthList)) *importerImportLocalAuthListCall {
	_c.Call = _c.Call.Run(func(args mock.Arguments) {
		_list, _ := args.Get(0).(settings.AuthList)
		fn(_list)
	})
	return _c
}

func (_c *importerImportLocalAuthListCall) OnImportEVSESettings(settings []settings.EVSE) *importerImportEVSESettingsCall {
	return _c.Parent.OnImportEVSESettings(settings)
}

func (_c *importerImportLocalAuthListCall) OnImportLocalAuthList(list settings.AuthList) *importerImportLocalAuthListCall {
	return _c.Parent.OnImportLocalAuthList(list)
}

func (_c *importerImportLocalAuthListCall) OnImportOcppConfiguration(config configuration.Config) *importerImportOcppConfigurationCall {
	return _c.Parent.OnImportOcppConfiguration(config)
}

func (_c *importerImportLocalAuthListCall) OnImportEVSESettingsRaw(settings interface{}) *importerImportEVSESettingsCall {
	return _c.Parent.OnImportEVSESettingsRaw(settings)
}

func (_c *importerImportLocalAuthListCall) OnImportLocalAuthListRaw(list interface{}) *importerImportLocalAuthListCall {
	return _c.Parent.OnImportLocalAuthListRaw(list)
}

func (_c *importerImportLocalAuthListCall) OnImportOcppConfigurationRaw(config interface{}) *importerImportOcppConfigurationCall {
	return _c.Parent.OnImportOcppConfigurationRaw(config)
}

func (_m *importerMock) ImportOcppConfiguration(config configuration.Config) error {
	_ret := _m.Called(config)

	if _rf, ok := _ret.Get(0).(func(configuration.Config) error); ok {
		return _rf(config)
	}

	_ra0 := _ret.Error(0)

	return _ra0
}

func (_m *importerMock) OnImportOcppConfiguration(config configuration.Config) *importerImportOcppConfigurationCall {
	return &importerImportOcppConfigurationCall{Call: _m.Mock.On("ImportOcppConfiguration", config), Parent: _m}
}

func (_m *importerMock) OnImportOcppConfigurationRaw(config interface{}) *importerImportOcppConfigurationCall {
	return &importerImportOcppConfigurationCall{Call: _m.Mock.On("ImportOcppConfiguration", config), Parent: _m}
}

type importerImportOcppConfigurationCall struct {
	*mock.Call
	Parent *importerMock
}

func (_c *importerImportOcppConfigurationCall) Panic(msg string) *importerImportOcppConfigurationCall {
	_c.Call = _c.Call.Panic(msg)
	return _c
}

func (_c *importerImportOcppConfigurationCall) Once() *importerImportOcppConfigurationCall {
	_c.Call = _c.Call.Once()
	return _c
}

func (_c *importerImportOcppConfigurationCall) Twice() *importerImportOcppConfigurationCall {
	_c.Call = _c.Call.Twice()
	return _c
}

func (_c *importerImportOcppConfigurationCall) Times(i int) *importerImportOcppConfigurationCall {
	_c.Call = _c.Call.Times(i)
	return _c
}

func (_c *importerImportOcppConfigurationCall) WaitUntil(w <-chan time.Time) *importerImportOcppConfigurationCall {
	_c.Call = _c.Call.WaitUntil(w)
	return _c
}

func (_c *importerImportOcppConfigurationCall) After(d time.Duration) *importerImportOcppConfigurationCall {
	_c.Call = _c.Call.After(d)
	return _c
}

func (_c *importerImportOcppConfigurationCall) Run(fn func(args mock.Arguments)) *importerImportOcppConfigurationCall {
	_c.Call = _c.Call.Run(fn)
	return _c
}

func (_c *importerImportOcppConfigurationCall) Maybe() *importerImportOcppConfigurationCall {
	_c.Call = _c.Call.Maybe()
	return _c
}

func (_c *importerImportOcppConfigurationCall) TypedReturns(a error) *importerImportOcppConfigurationCall {
	_c.Call = _c.Return(a)
	return _c
}

func (_c *importerImportOcppConfigurationCall) ReturnsFn(fn func(configuration.Config) error) *importerImportOcppConfigurationCall {
	_c.Call = _c.Return(fn)
	return _c
}

func (_c *importerImportOcppConfigurationCall) TypedRun(fn func(configuration.Config)) *importerImportOcppConfigurationCall {
	_c.Call = _c.Call.Run(func(args mock.Arguments) {
		_config, _ := args.Get(0).(configuration.Config)
		fn(_config)
	})
	return _c
}

func (_c *importerImportOcppConfigurationCall) OnImportEVSESettings(settings []settings.EVSE) *importerImportEVSESettingsCall {
	return _c.Parent.OnImportEVSESettings(settings)
}

func (_c *importerImportOcppConfigurationCall) OnImportLocalAuthList(list settings.AuthList) *importerImportLocalAuthListCall {
	return _c.Parent.OnImportLocalAuthList(list)
}

func (_c *importerImportOcppConfigurationCall) OnImportOcppConfiguration(config configuration.Config) *importerImportOcppConfigurationCall {
	return _c.Parent.OnImportOcppConfiguration(config)
}

func (_c *importerImportOcppConfigurationCall) OnImportEVSESettingsRaw(settings interface{}) *importerImportEVSESettingsCall {
	return _c.Parent.OnImportEVSESettingsRaw(settings)
}

func (_c *importerImportOcppConfigurationCall) OnImportLocalAuthListRaw(list interface{}) *importerImportLocalAuthListCall {
	return _c.Parent.OnImportLocalAuthListRaw(list)
}

func (_c *importerImportOcppConfigurationCall) OnImportOcppConfigurationRaw(config interface{}) *importerImportOcppConfigurationCall {
	return _c.Parent.OnImportOcppConfigurationRaw(config)
}

// managerMock mock of Manager.
type managerMock struct{ mock.Mock }

// NewManagerMock creates a new managerMock.
func NewManagerMock(tb testing.TB) *managerMock {
	tb.Helper()

	m := &managerMock{}
	m.Mock.Test(tb)

	tb.Cleanup(func() { m.AssertExpectations(tb) })

	return m
}

func (_m *managerMock) GetChargePointSettings() settings.ChargePoint {
	_ret := _m.Called()

	if _rf, ok := _ret.Get(0).(func() settings.ChargePoint); ok {
		return _rf()
	}

	_ra0, _ := _ret.Get(0).(settings.ChargePoint)

	return _ra0
}

func (_m *managerMock) OnGetChargePointSettings() *managerGetChargePointSettingsCall {
	return &managerGetChargePointSettingsCall{Call: _m.Mock.On("GetChargePointSettings"), Parent: _m}
}

func (_m *managerMock) OnGetChargePointSettingsRaw() *managerGetChargePointSettingsCall {
	return &managerGetChargePointSettingsCall{Call: _m.Mock.On("GetChargePointSettings"), Parent: _m}
}

type managerGetChargePointSettingsCall struct {
	*mock.Call
	Parent *managerMock
}

func (_c *managerGetChargePointSettingsCall) Panic(msg string) *managerGetChargePointSettingsCall {
	_c.Call = _c.Call.Panic(msg)
	return _c
}

func (_c *managerGetChargePointSettingsCall) Once() *managerGetChargePointSettingsCall {
	_c.Call = _c.Call.Once()
	return _c
}

func (_c *managerGetChargePointSettingsCall) Twice() *managerGetChargePointSettingsCall {
	_c.Call = _c.Call.Twice()
	return _c
}

func (_c *managerGetChargePointSettingsCall) Times(i int) *managerGetChargePointSettingsCall {
	_c.Call = _c.Call.Times(i)
	return _c
}

func (_c *managerGetChargePointSettingsCall) WaitUntil(w <-chan time.Time) *managerGetChargePointSettingsCall {
	_c.Call = _c.Call.WaitUntil(w)
	return _c
}

func (_c *managerGetChargePointSettingsCall) After(d time.Duration) *managerGetChargePointSettingsCall {
	_c.Call = _c.Call.After(d)
	return _c
}

func (_c *managerGetChargePointSettingsCall) Run(fn func(args mock.Arguments)) *managerGetChargePointSettingsCall {
	_c.Call = _c.Call.Run(fn)
	return _c
}

func (_c *managerGetChargePointSettingsCall) Maybe() *managerGetChargePointSettingsCall {
	_c.Call = _c.Call.Maybe()
	return _c
}

func (_c *managerGetChargePointSettingsCall) TypedReturns(a settings.ChargePoint) *managerGetChargePointSettingsCall {
	_c.Call = _c.Return(a)
	return _c
}

func (_c *managerGetChargePointSettingsCall) ReturnsFn(fn func() settings.ChargePoint) *managerGetChargePointSettingsCall {
	_c.Call = _c.Return(fn)
	return _c
}

func (_c *managerGetChargePointSettingsCall) TypedRun(fn func()) *managerGetChargePointSettingsCall {
	_c.Call = _c.Call.Run(func(args mock.Arguments) {
		fn()
	})
	return _c
}

func (_c *managerGetChargePointSettingsCall) OnGetChargePointSettings() *managerGetChargePointSettingsCall {
	return _c.Parent.OnGetChargePointSettings()
}

func (_c *managerGetChargePointSettingsCall) OnSetChargePointSettings(settings settings.ChargePoint) *managerSetChargePointSettingsCall {
	return _c.Parent.OnSetChargePointSettings(settings)
}

func (_c *managerGetChargePointSettingsCall) OnSetupOcppConfiguration(version configuration.ProtocolVersion, supportedProfiles []string) *managerSetupOcppConfigurationCall {
	return _c.Parent.OnSetupOcppConfiguration(version, supportedProfiles...)
}

func (_c *managerGetChargePointSettingsCall) OnGetChargePointSettingsRaw() *managerGetChargePointSettingsCall {
	return _c.Parent.OnGetChargePointSettingsRaw()
}

func (_c *managerGetChargePointSettingsCall) OnSetChargePointSettingsRaw(settings interface{}) *managerSetChargePointSettingsCall {
	return _c.Parent.OnSetChargePointSettingsRaw(settings)
}

func (_c *managerGetChargePointSettingsCall) OnSetupOcppConfigurationRaw(version interface{}, supportedProfiles interface{}) *managerSetupOcppConfigurationCall {
	return _c.Parent.OnSetupOcppConfigurationRaw(version, supportedProfiles)
}

func (_m *managerMock) SetChargePointSettings(chargePointSettings settings.ChargePoint) error {
	_ret := _m.Called(chargePointSettings)

	if _rf, ok := _ret.Get(0).(func(settings.ChargePoint) error); ok {
		return _rf(chargePointSettings)
	}

	_ra0 := _ret.Error(0)

	return _ra0
}

func (_m *managerMock) OnSetChargePointSettings(settings settings.ChargePoint) *managerSetChargePointSettingsCall {
	return &managerSetChargePointSettingsCall{Call: _m.Mock.On("SetChargePointSettings", settings), Parent: _m}
}

func (_m *managerMock) OnSetChargePointSettingsRaw(settings interface{}) *managerSetChargePointSettingsCall {
	return &managerSetChargePointSettingsCall{Call: _m.Mock.On("SetChargePointSettings", settings), Parent: _m}
}

type managerSetChargePointSettingsCall struct {
	*mock.Call
	Parent *managerMock
}

func (_c *managerSetChargePointSettingsCall) Panic(msg string) *managerSetChargePointSettingsCall {
	_c.Call = _c.Call.Panic(msg)
	return _c
}

func (_c *managerSetChargePointSettingsCall) Once() *managerSetChargePointSettingsCall {
	_c.Call = _c.Call.Once()
	return _c
}

func (_c *managerSetChargePointSettingsCall) Twice() *managerSetChargePointSettingsCall {
	_c.Call = _c.Call.Twice()
	return _c
}

func (_c *managerSetChargePointSettingsCall) Times(i int) *managerSetChargePointSettingsCall {
	_c.Call = _c.Call.Times(i)
	return _c
}

func (_c *managerSetChargePointSettingsCall) WaitUntil(w <-chan time.Time) *managerSetChargePointSettingsCall {
	_c.Call = _c.Call.WaitUntil(w)
	return _c
}

func (_c *managerSetChargePointSettingsCall) After(d time.Duration) *managerSetChargePointSettingsCall {
	_c.Call = _c.Call.After(d)
	return _c
}

func (_c *managerSetChargePointSettingsCall) Run(fn func(args mock.Arguments)) *managerSetChargePointSettingsCall {
	_c.Call = _c.Call.Run(fn)
	return _c
}

func (_c *managerSetChargePointSettingsCall) Maybe() *managerSetChargePointSettingsCall {
	_c.Call = _c.Call.Maybe()
	return _c
}

func (_c *managerSetChargePointSettingsCall) TypedReturns(a error) *managerSetChargePointSettingsCall {
	_c.Call = _c.Return(a)
	return _c
}

func (_c *managerSetChargePointSettingsCall) ReturnsFn(fn func(settings.ChargePoint) error) *managerSetChargePointSettingsCall {
	_c.Call = _c.Return(fn)
	return _c
}

func (_c *managerSetChargePointSettingsCall) TypedRun(fn func(settings.ChargePoint)) *managerSetChargePointSettingsCall {
	_c.Call = _c.Call.Run(func(args mock.Arguments) {
		_settings, _ := args.Get(0).(settings.ChargePoint)
		fn(_settings)
	})
	return _c
}

func (_c *managerSetChargePointSettingsCall) OnGetChargePointSettings() *managerGetChargePointSettingsCall {
	return _c.Parent.OnGetChargePointSettings()
}

func (_c *managerSetChargePointSettingsCall) OnSetChargePointSettings(settings settings.ChargePoint) *managerSetChargePointSettingsCall {
	return _c.Parent.OnSetChargePointSettings(settings)
}

func (_c *managerSetChargePointSettingsCall) OnSetupOcppConfiguration(version configuration.ProtocolVersion, supportedProfiles []string) *managerSetupOcppConfigurationCall {
	return _c.Parent.OnSetupOcppConfiguration(version, supportedProfiles...)
}

func (_c *managerSetChargePointSettingsCall) OnGetChargePointSettingsRaw() *managerGetChargePointSettingsCall {
	return _c.Parent.OnGetChargePointSettingsRaw()
}

func (_c *managerSetChargePointSettingsCall) OnSetChargePointSettingsRaw(settings interface{}) *managerSetChargePointSettingsCall {
	return _c.Parent.OnSetChargePointSettingsRaw(settings)
}

func (_c *managerSetChargePointSettingsCall) OnSetupOcppConfigurationRaw(version interface{}, supportedProfiles interface{}) *managerSetupOcppConfigurationCall {
	return _c.Parent.OnSetupOcppConfigurationRaw(version, supportedProfiles)
}

func (_m *managerMock) SetupOcppConfiguration(version configuration.ProtocolVersion, supportedProfiles ...string) {
	_m.Called(version, supportedProfiles)
}

func (_m *managerMock) OnSetupOcppConfiguration(version configuration.ProtocolVersion, supportedProfiles ...string) *managerSetupOcppConfigurationCall {
	return &managerSetupOcppConfigurationCall{Call: _m.Mock.On("SetupOcppConfiguration", version, supportedProfiles), Parent: _m}
}

func (_m *managerMock) OnSetupOcppConfigurationRaw(version interface{}, supportedProfiles interface{}) *managerSetupOcppConfigurationCall {
	return &managerSetupOcppConfigurationCall{Call: _m.Mock.On("SetupOcppConfiguration", version, supportedProfiles), Parent: _m}
}

type managerSetupOcppConfigurationCall struct {
	*mock.Call
	Parent *managerMock
}

func (_c *managerSetupOcppConfigurationCall) Panic(msg string) *managerSetupOcppConfigurationCall {
	_c.Call = _c.Call.Panic(msg)
	return _c
}

func (_c *managerSetupOcppConfigurationCall) Once() *managerSetupOcppConfigurationCall {
	_c.Call = _c.Call.Once()
	return _c
}

func (_c *managerSetupOcppConfigurationCall) Twice() *managerSetupOcppConfigurationCall {
	_c.Call = _c.Call.Twice()
	return _c
}

func (_c *managerSetupOcppConfigurationCall) Times(i int) *managerSetupOcppConfigurationCall {
	_c.Call = _c.Call.Times(i)
	return _c
}

func (_c *managerSetupOcppConfigurationCall) WaitUntil(w <-chan time.Time) *managerSetupOcppConfigurationCall {
	_c.Call = _c.Call.WaitUntil(w)
	return _c
}

func (_c *managerSetupOcppConfigurationCall) After(d time.Duration) *managerSetupOcppConfigurationCall {
	_c.Call = _c.Call.After(d)
	return _c
}

func (_c *managerSetupOcppConfigurationCall) Run(fn func(args mock.Arguments)) *managerSetupOcppConfigurationCall {
	_c.Call = _c.Call.Run(fn)
	return _c
}

func (_c *managerSetupOcppConfigurationCall) Maybe() *managerSetupOcppConfigurationCall {
	_c.Call = _c.Call.Maybe()
	return _c
}

func (_c *managerSetupOcppConfigurationCall) TypedRun(fn func(configuration.ProtocolVersion, ...string)) *managerSetupOcppConfigurationCall {
	_c.Call = _c.Call.Run(func(args mock.Arguments) {
		_version, _ := args.Get(0).(configuration.ProtocolVersion)
		_supportedProfiles, _ := args.Get(1).([]string)
		fn(_version, _supportedProfiles...)
	})
	return _c
}

func (_c *managerSetupOcppConfigurationCall) OnGetChargePointSettings() *managerGetChargePointSettingsCall {
	return _c.Parent.OnGetChargePointSettings()
}

func (_c *managerSetupOcppConfigurationCall) OnSetChargePointSettings(settings settings.ChargePoint) *managerSetChargePointSettingsCall {
	return _c.Parent.OnSetChargePointSettings(settings)
}

func (_c *managerSetupOcppConfigurationCall) OnSetupOcppConfiguration(version configuration.ProtocolVersion, supportedProfiles ...string) *managerSetupOcppConfigurationCall {
	return _c.Parent.OnSetupOcppConfiguration(version, supportedProfiles...)
}

func (_c *managerSetupOcppConfigurationCall) OnGetChargePointSettingsRaw() *managerGetChargePointSettingsCall {
	return _c.Parent.OnGetChargePointSettingsRaw()
}

func (_c *managerSetupOcppConfigurationCall) OnSetChargePointSettingsRaw(settings interface{}) *managerSetChargePointSettingsCall {
	return _c.Parent.OnSetChargePointSettingsRaw(settings)
}

func (_c *managerSetupOcppConfigurationCall) OnSetupOcppConfigurationRaw(version interface{}, supportedProfiles interface{}) *managerSetupOcppConfigurationCall {
	return _c.Parent.OnSetupOcppConfigurationRaw(version, supportedProfiles)
}
