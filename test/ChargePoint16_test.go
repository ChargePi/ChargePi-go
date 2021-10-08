package test

import (
	"math/rand"
	"os/exec"
	"time"

	"github.com/lorenzodonini/ocpp-go/ws"
	"github.com/xBlaz3kx/ChargePi-go/chargepoint"
	"github.com/xBlaz3kx/ChargePi-go/hardware"
	"reflect"
	"testing"
)

func TestChargePointHandler_AddConnectors(t *testing.T) {
	type fields struct {
		connectors []chargepoint.Connector
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "AddConnectors",
			fields: fields{
				connectors: nil,
			},
		}, {
			name: "NoConnectors",
			fields: fields{
				connectors: nil,
			},
		}, {
			name: "InvalidConnectorSettings",
			fields: fields{
				connectors: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &chargepoint.ChargePointHandler{}
			switch tt.name {
			case "AddConnectors":
				break
			}
			handler.AddConnectors()
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
	relay := hardware.NewRelay(10, false)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &chargepoint.ChargePointHandler{}
			for i, connector := range tt.fields.connectors {
				newConnector, err := chargepoint.NewConnector(
					connector.ConnectorId,
					connector.ConnectorId,
					connector.ConnectorType,
					relay,
					nil,
					false,
					connector.MaxChargingTime,
				)
				if err != nil {
					t.Errorf("error adding a connector for findAvailableConnector(); %v", err)
					continue
				}
				handler.Connectors = append(handler.Connectors, newConnector)
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
	relay := hardware.NewRelay(15, false)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &chargepoint.ChargePointHandler{}
			for _, connector := range tt.fields.connectors {
				newConnector, err := chargepoint.NewConnector(
					connector.EvseId,
					connector.ConnectorId,
					connector.ConnectorType,
					relay,
					nil,
					false,
					connector.MaxChargingTime,
				)
				if err != nil {
					t.Errorf("error adding a connector for findConnectorWithId(); %v", err)
					continue
				}
				handler.Connectors = append(handler.Connectors, newConnector)
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
	relay := hardware.NewRelay(2, false)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &chargepoint.ChargePointHandler{}
			for i, connector := range tt.fields.connectors {
				rand.Seed(time.Now().UnixNano())
				newConnector, err := chargepoint.NewConnector(
					connector.ConnectorId,
					connector.ConnectorId,
					connector.ConnectorType,
					relay,
					nil,
					false,
					connector.MaxChargingTime,
				)
				if err != nil {
					t.Errorf("error adding a connector for findConnectorWithId(); %v", err)
					continue
				}
				handler.Connectors = append(handler.Connectors, newConnector)
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
	relay := hardware.NewRelay(1, false)
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &chargepoint.ChargePointHandler{}
			for _, connector := range tt.fields.connectors {
				newConnector, err := chargepoint.NewConnector(
					connector.ConnectorId,
					connector.ConnectorId,
					connector.ConnectorType,
					relay,
					nil,
					false,
					connector.MaxChargingTime,
				)
				if err != nil {
					t.Errorf("error adding a connector for findConnectorWithId(); %v", err)
					continue
				}
				handler.Connectors = append(handler.Connectors, newConnector)
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
