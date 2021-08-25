package test

import (
	"os/exec"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/ws"
	goCache "github.com/patrickmn/go-cache"
	"github.com/xBlaz3kx/ChargePi-go/cache"
	"github.com/xBlaz3kx/ChargePi-go/chargepoint"
	"github.com/xBlaz3kx/ChargePi-go/hardware"
	"github.com/xBlaz3kx/ChargePi-go/settings"
	"reflect"
	"testing"
	"time"
)

func TestChargePointHandler_AddConnectors(t *testing.T) {
	type fields struct {
		connectors []chargepoint.Connector
	}
	type args struct {
		connectors []settings.Connector
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "AddConnectors",
			fields: fields{
				connectors: nil,
			},
			args: args{
				connectors: []settings.Connector{{
					EvseId:      0,
					ConnectorId: 0,
					Type:        "",
					Status:      "",
					Session: struct {
						IsActive      bool
						TransactionId string
						TagId         string
						Started       string
						Consumption   []types.MeterValue
					}{
						IsActive:      false,
						TransactionId: "",
						TagId:         "",
						Started:       "",
						Consumption:   nil,
					},
					Relay: struct {
						RelayPin     int
						InverseLogic bool
					}{},
					PowerMeter: struct {
						Enabled              bool
						PowerMeterPin        int
						SpiBus               int
						PowerUnits           string
						Consumption          float64
						ShuntOffset          float64
						VoltageDividerOffset float64
					}{},
				},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &chargepoint.ChargePointHandler{}
			handler.AddConnectors(nil)
		})
	}
}

func TestChargePointHandler_findAvailableConnector(t *testing.T) {
	type fields struct {
		connectors             []chargepoint.Connector
		chargingConnectorIndex int
	}
	tests := []struct {
		name   string
		fields fields
		want   chargepoint.Connector
	}{
		{
			name: "FirstConnectorAvailable",
			fields: fields{
				connectors: []chargepoint.Connector{
					{
						EvseId:          1,
						ConnectorId:     1,
						ConnectorType:   "Schuko",
						MaxChargingTime: 0,
					}, {
						EvseId:          1,
						ConnectorId:     2,
						ConnectorType:   "Schuko",
						MaxChargingTime: 0,
					}, {
						EvseId:          1,
						ConnectorId:     3,
						ConnectorType:   "Schuko",
						MaxChargingTime: 0,
					},
				},
				chargingConnectorIndex: -1,
			},
			want: chargepoint.Connector{
				EvseId:          0,
				ConnectorId:     0,
				ConnectorType:   "Schuko",
				MaxChargingTime: 0,
			},
		}, {
			name: "SecondConnectorAvailable",
			fields: fields{
				connectors: []chargepoint.Connector{
					{
						EvseId:          1,
						ConnectorId:     1,
						ConnectorType:   "Schuko",
						MaxChargingTime: 0,
					}, {
						EvseId:          1,
						ConnectorId:     2,
						ConnectorType:   "Schuko",
						MaxChargingTime: 0,
					}, {
						EvseId:          1,
						ConnectorId:     3,
						ConnectorType:   "Schuko",
						MaxChargingTime: 0,
					},
				},
				chargingConnectorIndex: 0,
			},
			want: chargepoint.Connector{
				EvseId:          1,
				ConnectorId:     2,
				ConnectorType:   "Schuko",
				MaxChargingTime: 0,
			},
		},
		{
			name: "NoConnectorAvailable",
			fields: fields{
				connectors: []chargepoint.Connector{
					{
						EvseId:          1,
						ConnectorId:     1,
						ConnectorType:   "Schuko",
						MaxChargingTime: 0,
					},
				},
				chargingConnectorIndex: 0,
			},
			want: chargepoint.Connector{
				EvseId:          0,
				ConnectorId:     0,
				ConnectorType:   "Schuko",
				MaxChargingTime: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &chargepoint.ChargePointHandler{}
			for i, connector := range tt.fields.connectors {
				newConnector, err := chargepoint.NewConnector(
					connector.ConnectorId,
					connector.ConnectorId,
					connector.ConnectorType,
					hardware.NewRelay(1, false),
					nil,
					false,
					connector.MaxChargingTime,
				)
				if err != nil {
					t.Errorf("error adding a connector for findAvailableConnector(); %v", err)
					continue
				}
				handler.Connectors = append(handler.Connectors, &newConnector)
				if i == tt.fields.chargingConnectorIndex {
					connector2 := handler.FindConnectorWithId(connector.ConnectorId)
					err = connector2.StartCharging("123", "123")
					if err != nil {
						t.Errorf("Cannot start charging connector %d: %v", connector2.ConnectorId, err)
					}
				}
			}
			got := handler.FindAvailableConnector()
			if got == nil && tt.want.ConnectorId == 0 && tt.want.EvseId == 0 {
				return
			} else if got == nil {
				t.Errorf("findAvailableConnector() = %v, want %v", got, tt.want)
			} else if !got.IsAvailable() && got.ConnectorId != tt.want.ConnectorId && got.EvseId != tt.want.EvseId {
				t.Errorf("findAvailableConnector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChargePointHandler_findConnectorWithId(t *testing.T) {
	type fields struct {
		connectors []chargepoint.Connector
	}
	type args struct {
		connectorID int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   chargepoint.Connector
	}{
		{
			name: "ID1",
			fields: fields{
				connectors: []chargepoint.Connector{
					{
						EvseId:          1,
						ConnectorId:     1,
						ConnectorType:   "Schuko",
						MaxChargingTime: 0,
					}, {
						EvseId:          1,
						ConnectorId:     2,
						ConnectorType:   "Schuko",
						MaxChargingTime: 0,
					}, {
						EvseId:          1,
						ConnectorId:     3,
						ConnectorType:   "Schuko",
						MaxChargingTime: 0,
					},
				}},
			args: args{connectorID: 1},
			want: chargepoint.Connector{
				EvseId:          1,
				ConnectorId:     1,
				ConnectorType:   "Schuko",
				MaxChargingTime: 0,
			},
		},
		{
			name: "ID2",
			fields: fields{
				connectors: []chargepoint.Connector{
					{
						EvseId:          1,
						ConnectorId:     1,
						ConnectorType:   "Schuko",
						MaxChargingTime: 0,
					}, {
						EvseId:          1,
						ConnectorId:     2,
						ConnectorType:   "Schuko",
						MaxChargingTime: 0,
					}, {
						EvseId:          1,
						ConnectorId:     3,
						ConnectorType:   "Schuko",
						MaxChargingTime: 0,
					},
				}},
			args: args{connectorID: 2},
			want: chargepoint.Connector{
				EvseId:          1,
				ConnectorId:     2,
				ConnectorType:   "Schuko",
				MaxChargingTime: 0,
			},
		}, {
			name: "NoConnector",
			fields: fields{
				connectors: []chargepoint.Connector{
					{
						EvseId:          1,
						ConnectorId:     1,
						ConnectorType:   "Schuko",
						MaxChargingTime: 0,
					}, {
						EvseId:          1,
						ConnectorId:     2,
						ConnectorType:   "Schuko",
						MaxChargingTime: 0,
					}, {
						EvseId:          1,
						ConnectorId:     3,
						ConnectorType:   "Schuko",
						MaxChargingTime: 0,
					},
				}},
			args: args{connectorID: 5},
			want: chargepoint.Connector{
				EvseId:          0,
				ConnectorId:     0,
				ConnectorType:   "Schuko",
				MaxChargingTime: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &chargepoint.ChargePointHandler{}
			for _, connector := range tt.fields.connectors {
				newConnector, err := chargepoint.NewConnector(
					connector.EvseId,
					connector.ConnectorId,
					connector.ConnectorType,
					hardware.NewRelay(1, false),
					nil,
					false,
					connector.MaxChargingTime,
				)
				if err != nil {
					t.Errorf("error adding a connector for findConnectorWithId(); %v", err)
					continue
				}
				handler.Connectors = append(handler.Connectors, &newConnector)
			}
			got := handler.FindConnectorWithId(tt.args.connectorID)
			if got == nil && tt.want.ConnectorId == 0 && tt.want.EvseId == 0 {
				return
			} else if got == nil {
				t.Errorf("findConnectorWithId() = %v, want %v", got, tt.want)
			} else if !got.IsAvailable() && got.ConnectorId != tt.want.ConnectorId && got.EvseId != tt.want.EvseId {
				t.Errorf("findConnectorWithId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChargePointHandler_findConnectorWithTagId(t *testing.T) {
	type fields struct {
		connectors []chargepoint.Connector
	}
	type args struct {
		tagId                  string
		chargingConnectorIndex int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   chargepoint.Connector
	}{
		{
			name: "NoConnector",
			fields: fields{
				connectors: []chargepoint.Connector{
					{
						EvseId:          1,
						ConnectorId:     1,
						ConnectorType:   "Schuko",
						MaxChargingTime: 0,
					}, {
						EvseId:          1,
						ConnectorId:     2,
						ConnectorType:   "Schuko",
						MaxChargingTime: 0,
					}, {
						EvseId:          1,
						ConnectorId:     3,
						ConnectorType:   "Schuko",
						MaxChargingTime: 0,
					},
				}},
			args: args{tagId: "1234",
				chargingConnectorIndex: -1,
			},
			want: chargepoint.Connector{
				EvseId:          0,
				ConnectorId:     0,
				ConnectorType:   "Schuko",
				MaxChargingTime: 0,
			},
		}, {
			name: "FoundConnector",
			fields: fields{
				connectors: []chargepoint.Connector{
					{
						EvseId:          1,
						ConnectorId:     1,
						ConnectorType:   "Schuko",
						MaxChargingTime: 0,
					}, {
						EvseId:          1,
						ConnectorId:     2,
						ConnectorType:   "Schuko",
						MaxChargingTime: 0,
					}, {
						EvseId:          1,
						ConnectorId:     3,
						ConnectorType:   "Schuko",
						MaxChargingTime: 0,
					},
				}},
			args: args{tagId: "1234",
				chargingConnectorIndex: 0,
			},
			want: chargepoint.Connector{
				EvseId:          1,
				ConnectorId:     1,
				ConnectorType:   "Schuko",
				MaxChargingTime: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &chargepoint.ChargePointHandler{}
			for i, connector := range tt.fields.connectors {
				newConnector, err := chargepoint.NewConnector(
					connector.ConnectorId,
					connector.ConnectorId,
					connector.ConnectorType,
					hardware.NewRelay(1, false),
					nil,
					false,
					connector.MaxChargingTime,
				)
				if err != nil {
					t.Errorf("error adding a connector for findConnectorWithId(); %v", err)
					continue
				}
				handler.Connectors = append(handler.Connectors, &newConnector)
				if i == tt.args.chargingConnectorIndex {
					connector2 := handler.FindConnectorWithId(newConnector.ConnectorId)
					err = connector2.StartCharging("1234", tt.args.tagId)
					if err != nil {
						t.Errorf("Cannot start charging connector %d: %v", connector.ConnectorId, err)
					}
				}
			}
			got := handler.FindConnectorWithTagId(tt.args.tagId)
			if got == nil && tt.want.ConnectorId == 0 && tt.want.EvseId == 0 {
				return
			} else if got == nil {
				t.Errorf("findConnectorWithTagId() = %v, want %v", got, tt.want)
			} else if !got.IsAvailable() && got.ConnectorId != tt.want.ConnectorId && got.EvseId != tt.want.EvseId {
				t.Errorf("findConnectorWithTagId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChargePointHandler_findConnectorWithTransactionId(t *testing.T) {
	type fields struct {
		connectors []chargepoint.Connector
	}
	type args struct {
		transactionId          string
		chargingConnectorIndex int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   chargepoint.Connector
	}{
		{
			name: "OneConnectorWithTransactionId",
			fields: fields{connectors: []chargepoint.Connector{
				{
					EvseId:          1,
					ConnectorId:     1,
					ConnectorType:   "Schuko",
					MaxChargingTime: 0,
				}, {
					EvseId:          1,
					ConnectorId:     2,
					ConnectorType:   "Schuko",
					MaxChargingTime: 0,
				}, {
					EvseId:          1,
					ConnectorId:     3,
					ConnectorType:   "Schuko",
					MaxChargingTime: 0,
				},
			}},
			args: args{transactionId: "1234", chargingConnectorIndex: 0},
			want: chargepoint.Connector{
				EvseId:          1,
				ConnectorId:     1,
				ConnectorType:   "Schuko",
				MaxChargingTime: 0,
			},
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &chargepoint.ChargePointHandler{}
			for _, connector := range tt.fields.connectors {
				newConnector, err := chargepoint.NewConnector(
					connector.ConnectorId,
					connector.ConnectorId,
					connector.ConnectorType,
					hardware.NewRelay(1, false),
					nil,
					false,
					connector.MaxChargingTime,
				)
				if err != nil {
					t.Errorf("error adding a connector for findConnectorWithId(); %v", err)
					continue
				}
				handler.Connectors = append(handler.Connectors, &newConnector)
				if i == tt.args.chargingConnectorIndex {
					connector2 := handler.FindConnectorWithId(newConnector.ConnectorId)
					err = connector2.StartCharging(tt.args.transactionId, "1234")
					if err != nil {
						t.Errorf("Cannot start charging connector %d: %v", connector.ConnectorId, err)
					}
				}
			}
			got := handler.FindConnectorWithTransactionId(tt.args.transactionId)
			if got == nil && tt.want.ConnectorId == 0 && tt.want.EvseId == 0 {
				return
			} else if got == nil {
				t.Errorf("findConnectorWithTransactionId() = %v, want %v", got, tt.want)
			} else if !got.IsAvailable() && got.ConnectorId != tt.want.ConnectorId && got.EvseId != tt.want.EvseId {
				t.Errorf("findConnectorWithTransactionId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTLSClient(t *testing.T) {
	type args struct {
		CACertificatePath     string
		ClientCertificatePath string
		ClientKeyPath         string
	}
	tests := []struct {
		name string
		args args
		want *ws.Client
	}{
		{
			name: "GoodCertificates",
			args: args{
				CACertificatePath:     "certs/ca.crt",
				ClientCertificatePath: "certs/cp/charge-point.crt",
				ClientKeyPath:         "certs/cp/charge-point.key",
			},
			want: ws.NewClient(),
		}, {
			name: "NoCertificates",
			args: args{
				CACertificatePath:     "",
				ClientCertificatePath: "",
				ClientKeyPath:         "",
			},
			want: nil,
		}, {
			name: "MissingCACertificate",
			args: args{
				CACertificatePath:     "",
				ClientCertificatePath: "certs/cp/central-system.crt",
				ClientKeyPath:         "certs/cp/central-system.key",
			},
			want: nil,
		}, {
			name: "MissingClientCert",
			args: args{
				CACertificatePath:     "certs/cs/ca.crt",
				ClientCertificatePath: "",
				ClientKeyPath:         "certs/cp/central-system.key",
			},
			want: nil,
		}, {
			name: "MissingClientKey",
			args: args{
				CACertificatePath:     "certs/ca.crt",
				ClientCertificatePath: "certs/cp/central-system.crt",
				ClientKeyPath:         "certs/cp/central-system.key",
			},
			want: nil,
		}, {
			name: "InvalidFilePath",
			args: args{
				CACertificatePath:     "certs/cs/ca123.crt",
				ClientCertificatePath: "certs/cp/central-system.crt",
				ClientKeyPath:         "certs/cp/central-system.key",
			},
			want: nil,
		},
	}
	exec.Command("sudo ./create-test-certs.sh")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := chargepoint.GetTLSClient(tt.args.CACertificatePath, tt.args.ClientCertificatePath, tt.args.ClientKeyPath); !reflect.DeepEqual(got, tt.want) {
				if tt.name == "GoodCertificates" && got != nil {
					return
				}
				t.Errorf("getTLSClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_retrySendingRequest(t *testing.T) {
	type args struct {
		cacheRetryKey string
		maxRetries    int
		interval      int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "RetryTwice",
			args: args{
				cacheRetryKey: "RetryTwice",
				maxRetries:    2,
				interval:      5,
			},
			wantErr: false,
		}, {
			name: "ZeroRetries",
			args: args{
				cacheRetryKey: "ZeroRetries",
				maxRetries:    0,
				interval:      10,
			},
			wantErr: false,
		},
		{
			name: "ZeroRetries2",
			args: args{
				cacheRetryKey: "ZeroRetries2",
				maxRetries:    -1,
				interval:      10,
			},
			wantErr: false,
		}, {
			name: "FiveBackToBackRetries",
			args: args{
				cacheRetryKey: "FiveBackToBackRetries",
				maxRetries:    5,
				interval:      0,
			},
			wantErr: false,
		},
	}
	cache.Cache = goCache.New(time.Minute*10, time.Minute*10)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			/*if err := chargepoint.RetrySendingRequest(tt.args.cacheRetryKey, tt.args.maxRetries, tt.args.interval, fmt.Println, "Job executed"); (err != nil) != tt.wantErr {
				t.Errorf("retrySendingRequest() error = %v, wantErr %v", err, tt.wantErr)
			}*/
			time.Sleep(time.Duration(tt.args.maxRetries*tt.args.interval+5) * time.Second)

		})
	}
}
