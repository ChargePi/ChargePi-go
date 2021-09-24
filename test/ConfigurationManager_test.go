package test

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/patrickmn/go-cache"
	cache2 "github.com/xBlaz3kx/ChargePi-go/cache"
	"github.com/xBlaz3kx/ChargePi-go/settings"
	"reflect"
	"testing"
	"time"
)

func TestGetConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		want    settings.OCPPConfig
		wantErr bool
	}{
		{
			name: "ConfigurationFound",
			want: settings.OCPPConfig{
				Version: 1,
				Keys: []core.ConfigurationKey{
					{
						Key:      "Test1",
						Readonly: false,
						Value:    "60",
					}, {
						Key:      "Test2",
						Readonly: false,
						Value:    "ABCD",
					},
				},
			},
			wantErr: false,
		}, {
			name:    "ConfigurationFound",
			want:    settings.OCPPConfig{},
			wantErr: true,
		},
	}
	err := settings.WriteToFile("test/configuration.json", settings.OCPPConfig{
		Version: 1,
		Keys: []core.ConfigurationKey{
			{
				Key:      "Test1",
				Readonly: false,
				Value:    "60",
			}, {
				Key:      "Test2",
				Readonly: false,
				Value:    "ABCD",
			},
		},
	})
	if err != nil {
		return
	}
	var config = settings.OCPPConfig{}
	//settings.DecodeFile("test/configuration.json", &config)
	if cache2.Cache == nil {
		cache2.Cache = cache.New(time.Minute*10, time.Minute*10)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "ConfigurationFound":
				cache2.Cache.Set("OCPPConfiguration", &config, cache.DefaultExpiration)
				break
			case "ConfigurationNotInCache":
				break
			}
			got, err := settings.GetConfiguration()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfiguration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfiguration() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetConfigurationValue(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := settings.GetConfigurationValue(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfigurationValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetConfigurationValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOCPPConfig_GetConfig(t *testing.T) {
	type fields struct {
		Version int
		Keys    []core.ConfigurationKey
	}
	tests := []struct {
		name   string
		fields fields
		want   []core.ConfigurationKey
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &settings.OCPPConfig{
				Version: tt.fields.Version,
				Keys:    tt.fields.Keys,
			}
			if got := config.GetConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOCPPConfig_GetConfigurationValue(t *testing.T) {
	type fields struct {
		Version int
		Keys    []core.ConfigurationKey
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &settings.OCPPConfig{
				Version: tt.fields.Version,
				Keys:    tt.fields.Keys,
			}
			got, err := config.GetConfigurationValue(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfigurationValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetConfigurationValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOCPPConfig_UpdateKey(t *testing.T) {
	type fields struct {
		Version int
		Keys    []core.ConfigurationKey
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &settings.OCPPConfig{
				Version: tt.fields.Version,
				Keys:    tt.fields.Keys,
			}
			if err := config.UpdateKey(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("UpdateKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateConfigurationFile(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := settings.UpdateConfigurationFile(); (err != nil) != tt.wantErr {
				t.Errorf("UpdateConfigurationFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateConnectorSessionInfo(t *testing.T) {
	type args struct {
		evseId      int
		connectorId int
		session     *settings.Session
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestUpdateConnectorStatus(t *testing.T) {
	type args struct {
		evseId      int
		connectorId int
		status      core.ChargePointStatus
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestUpdateKey(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := settings.UpdateKey(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("UpdateKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWriteToFile(t *testing.T) {
	type args struct {
		filename  string
		structure interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := settings.WriteToFile(tt.args.filename, tt.args.structure); (err != nil) != tt.wantErr {
				t.Errorf("WriteToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
