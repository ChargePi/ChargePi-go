package test

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/xBlaz3kx/ChargePi-go/chargepoint"
	"github.com/xBlaz3kx/ChargePi-go/data"
	"github.com/xBlaz3kx/ChargePi-go/hardware"
	"github.com/xBlaz3kx/ChargePi-go/hardware/power-meter"
	"reflect"
	"testing"
	"time"
)

func TestConnector_ResumeCharging(t *testing.T) {
	type fields struct {
		EvseId        int
		ConnectorId   int
		ConnectorType string
	}
	type args struct {
		session data.Session
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "ResumeSuccessful",
			fields: fields{
				EvseId:        1,
				ConnectorId:   1,
				ConnectorType: "Schuko",
			},
			args: args{session: data.Session{
				IsActive:      true,
				TransactionId: "123",
				TagId:         "123",
				Started:       time.Now().Format(time.RFC3339),
				Consumption:   []types.MeterValue{},
			}},
			wantErr: false,
		}, {
			name: "WrongStatus",
			fields: fields{
				EvseId:        1,
				ConnectorId:   1,
				ConnectorType: "Schuko",
			},
			args: args{session: data.Session{
				IsActive:      true,
				TransactionId: "123",
				TagId:         "123",
				Started:       time.Now().Format(time.RFC3339),
				Consumption:   []types.MeterValue{},
			}},
			wantErr: true,
		}, {
			name: "SessionNotActive",
			fields: fields{
				EvseId:        1,
				ConnectorId:   1,
				ConnectorType: "Schuko",
			},
			args: args{session: data.Session{
				IsActive:      false,
				TransactionId: "",
				TagId:         "",
				Started:       "",
				Consumption:   []types.MeterValue{},
			}},
			wantErr: true,
		},
	}
	connector, err := chargepoint.NewConnector(
		1,
		1,
		"Schuko",
		hardware.NewRelay(15, false),
		nil,
		false,
		5,
	)
	if err != nil {
		t.Errorf("ResumeCharging() error = %v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "ResumeSuccessful":
				connector.ConnectorStatus = core.ChargePointStatusCharging
				connector.ErrorCode = core.NoError
				break
			case "SessionNotActive":
				connector.ConnectorStatus = core.ChargePointStatusCharging
				connector.ErrorCode = core.NoError
				break
			}
			if err, _ = connector.ResumeCharging(tt.args.session); (err != nil) != tt.wantErr {
				t.Errorf("ResumeCharging() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConnector_StartCharging(t *testing.T) {
	type fields struct {
		EvseId        int
		ConnectorId   int
		ConnectorType string
	}
	type args struct {
		transactionId string
		tagId         string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "WillStartCharging",
			fields: fields{
				EvseId:        1,
				ConnectorId:   1,
				ConnectorType: "Schuko",
			},
			args: args{
				transactionId: "123",
				tagId:         "1234",
			},
			wantErr: false,
		}, {
			name: "MissingTransactionId",
			fields: fields{
				EvseId:        1,
				ConnectorId:   1,
				ConnectorType: "Schuko",
			},
			args: args{
				transactionId: "",
				tagId:         "1234",
			},
			wantErr: true,
		}, {
			name: "MissingTagId",
			fields: fields{
				EvseId:        1,
				ConnectorId:   1,
				ConnectorType: "Schuko",
			},
			args: args{
				transactionId: "1234",
				tagId:         "",
			},
			wantErr: true,
		},
	}
	relay := hardware.NewRelay(2, false)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connector, err := chargepoint.NewConnector(
				tt.fields.EvseId,
				tt.fields.ConnectorId,
				tt.fields.ConnectorType,
				relay,
				nil,
				false,
				5,
			)
			if err != nil {
				t.Errorf("StartCharging() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := connector.StartCharging(tt.args.transactionId, tt.args.tagId); (err != nil) != tt.wantErr {
				t.Errorf("StartCharging() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConnector_StopCharging(t *testing.T) {
	type fields struct {
		EvseId        int
		ConnectorId   int
		ConnectorType string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "StopCharging",
			fields: fields{
				EvseId:        1,
				ConnectorId:   1,
				ConnectorType: "Schuko",
			},
			wantErr: false,
		}, {
			name: "StopChargingFailure",
			fields: fields{
				EvseId:        1,
				ConnectorId:   1,
				ConnectorType: "Schuko",
			},
			wantErr: true,
		},
	}
	relay := hardware.NewRelay(1, false)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connector, err := chargepoint.NewConnector(
				tt.fields.EvseId,
				tt.fields.ConnectorId,
				tt.fields.ConnectorType,
				relay,
				nil,
				false,
				5,
			)
			if err != nil {
				t.Errorf("StopCharging() error at connector = %v, wantErr %v", err, tt.wantErr)
			}
			switch tt.name {
			case "StopCharging":
				err := connector.StartCharging("123", "123")
				if err != nil {
					t.Errorf("StopCharging() error while starting to charge = %v, wantErr %v", err, tt.wantErr)
				}
				break
			case "StopChargingFailure":
				break
			}
			if err := connector.StopCharging(core.ReasonLocal); (err != nil) != tt.wantErr {
				t.Errorf("StopCharging() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewConnector(t *testing.T) {
	var relay = hardware.NewRelay(21, false)
	type args struct {
		EvseId            int
		connectorId       int
		connectorType     string
		relay             *hardware.Relay
		powerMeter        *power_meter.C5460A
		powerMeterEnabled bool
		maxChargingTime   int
	}
	var tests = []struct {
		name    string
		args    args
		want    chargepoint.Connector
		wantErr bool
	}{
		{
			name: "CorrectConnector",
			args: args{
				EvseId:            1,
				connectorId:       1,
				connectorType:     "Schuko",
				relay:             relay,
				powerMeter:        nil,
				powerMeterEnabled: false,
				maxChargingTime:   5,
			},
			want: chargepoint.Connector{
				EvseId:            1,
				ConnectorId:       1,
				ConnectorType:     "Schuko",
				ConnectorStatus:   "Available",
				PowerMeterEnabled: false,
				MaxChargingTime:   5,
			},
			wantErr: false,
		}, {
			name: "InvalidEVSEID",
			args: args{
				EvseId:            0,
				connectorId:       1,
				connectorType:     "Schuko",
				relay:             relay,
				powerMeter:        nil,
				powerMeterEnabled: false,
				maxChargingTime:   5,
			},
			want:    chargepoint.Connector{},
			wantErr: true,
		}, {
			name: "InvalidConnectorId",
			args: args{
				EvseId:            1,
				connectorId:       0,
				connectorType:     "Schuko",
				relay:             relay,
				powerMeter:        nil,
				powerMeterEnabled: false,
				maxChargingTime:   5,
			},
			want:    chargepoint.Connector{},
			wantErr: true,
		}, {
			name: "RelayNotCreated",
			args: args{
				EvseId:            1,
				connectorId:       1,
				connectorType:     "Schuko",
				relay:             nil,
				powerMeter:        nil,
				powerMeterEnabled: false,
				maxChargingTime:   5,
			},
			want:    chargepoint.Connector{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := chargepoint.NewConnector(
				tt.args.EvseId,
				tt.args.connectorId,
				tt.args.connectorType,
				tt.args.relay,
				tt.args.powerMeter,
				tt.args.powerMeterEnabled,
				tt.args.maxChargingTime,
			)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewConnector() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			switch tt.name {
			case "CorrectConnector":
				if tt.args.connectorType != got.ConnectorType || tt.args.connectorId != got.ConnectorId ||
					tt.args.EvseId != got.EvseId || tt.args.powerMeterEnabled != got.PowerMeterEnabled || tt.args.maxChargingTime != got.MaxChargingTime {
					t.Errorf("NewConnector() got = %v, want %v", got, tt.want)
				}
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConnector() got = %v, want %v", got, tt.want)
			}
		})
	}
}
