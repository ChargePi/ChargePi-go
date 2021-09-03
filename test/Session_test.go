package test

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/xBlaz3kx/ChargePi-go/data"
	"reflect"
	"testing"
	"time"
)

func TestSession_AddSampledValue(t *testing.T) {
	type fields struct {
		IsActive      bool
		TransactionId string
		TagId         string
		Started       string
		Consumption   []types.MeterValue
	}
	type args struct {
		samples []types.SampledValue
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []types.SampledValue
	}{
		{
			name: "ValidSampling",
			fields: fields{
				IsActive:      true,
				TransactionId: "123",
				TagId:         "123",
				Started:       time.Now().String(),
				Consumption:   []types.MeterValue{},
			},
			args: args{
				samples: []types.SampledValue{
					{
						Value:     "123.21",
						Measurand: types.MeasurandCurrentExport,
					}, {
						Value:     "123.21",
						Measurand: types.MeasurandVoltage,
					}, {
						Value:     "123.21",
						Measurand: types.MeasurandPowerActiveExport,
					},
				},
			},
			want: []types.SampledValue{
				{
					Value:     "123.21",
					Measurand: types.MeasurandCurrentExport,
				}, {
					Value:     "123.21",
					Measurand: types.MeasurandVoltage,
				}, {
					Value:     "123.21",
					Measurand: types.MeasurandPowerActiveExport,
				},
			},
		}, {
			name: "SessionNotActive",
			fields: fields{
				IsActive:      false,
				TransactionId: "",
				TagId:         "",
				Started:       "",
				Consumption:   []types.MeterValue{},
			},
			args: args{
				samples: []types.SampledValue{
					{
						Value:     "123.21",
						Measurand: types.MeasurandCurrentExport,
					}, {
						Value:     "123.21",
						Measurand: types.MeasurandVoltage,
					}, {
						Value:     "123.21",
						Measurand: types.MeasurandPowerActiveExport,
					},
				},
			},
			want: []types.SampledValue{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &data.Session{
				IsActive:      tt.fields.IsActive,
				TransactionId: tt.fields.TransactionId,
				TagId:         tt.fields.TagId,
				Started:       tt.fields.Started,
				Consumption:   tt.fields.Consumption,
			}
			session.AddSampledValue(tt.args.samples)
			var expectedResult = []types.MeterValue{}
			if tt.name == "ValidSampling" {
				expectedResult = []types.MeterValue{{Timestamp: session.Consumption[0].Timestamp, SampledValue: tt.want}}
			}
			if !reflect.DeepEqual(session.Consumption, expectedResult) {
				t.Errorf("Expected: %v, got %v", session.Consumption, tt.want)
			}
		})
	}
}

func TestSession_EndSession(t *testing.T) {
	type fields struct {
		IsActive      bool
		TransactionId string
		TagId         string
		Started       string
		Consumption   []types.MeterValue
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "StopSuccessful",
			fields: fields{
				TransactionId: "1234",
				TagId:         "1234",
			},
		}, {
			name: "Unsuccessful",
			fields: fields{
				TransactionId: "",
				TagId:         "",
				Started:       "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &data.Session{
				IsActive:      false,
				TransactionId: "",
				TagId:         "",
				Started:       "",
				Consumption:   []types.MeterValue{},
			}
			switch tt.name {
			case "StopSuccessful":
				if !session.StartSession(tt.fields.TransactionId, tt.fields.TagId) {
					t.Errorf("Could not start session")
				}
				session.EndSession()
				if session.IsActive != false && session.TagId != "" && session.Started != "" && session.TransactionId != "" {
					t.Errorf("Session not reset")
				}
				break
			case "Unsuccessful":
				session.EndSession()
				if session.IsActive != false && session.TagId != "" && session.Started != "" && session.TransactionId != "" {
					t.Errorf("Session not reset")
				}
				break
			}
		})
	}
}

func TestSession_StartSession(t *testing.T) {
	type fields struct {
		IsActive      bool
		TransactionId string
		TagId         string
		Started       string
		Consumption   []types.MeterValue
	}
	type args struct {
		transactionId string
		tagId         string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "StartSuccessful",
			fields: fields{
				IsActive:      false,
				TransactionId: "",
				TagId:         "",
				Started:       "",
				Consumption:   []types.MeterValue{},
			},
			args: args{
				transactionId: "123",
				tagId:         "1234",
			},
			want: true,
		},
		{
			name: "MissingTransactionID",
			fields: fields{
				IsActive:      false,
				TransactionId: "",
				TagId:         "",
				Started:       "",
				Consumption:   []types.MeterValue{},
			},
			args: args{
				transactionId: "",
				tagId:         "1234",
			},
			want: false,
		}, {
			name: "MissingTagID",
			fields: fields{
				IsActive:      false,
				TransactionId: "",
				TagId:         "",
				Started:       "",
				Consumption:   []types.MeterValue{},
			},
			args: args{
				transactionId: "1234",
				tagId:         "",
			},
			want: false,
		}, {
			name: "TransactionAlreadyActive",
			fields: fields{
				IsActive:      true,
				TransactionId: "1234",
				TagId:         "1234",
				Started:       time.Now().String(),
				Consumption:   []types.MeterValue{},
			},
			args: args{
				transactionId: "1234",
				tagId:         "1234",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &data.Session{
				IsActive:      tt.fields.IsActive,
				TransactionId: tt.fields.TransactionId,
				TagId:         tt.fields.TagId,
				Started:       tt.fields.Started,
				Consumption:   tt.fields.Consumption,
			}
			if got := session.StartSession(tt.args.transactionId, tt.args.tagId); got != tt.want {
				t.Errorf("StartSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_CalculateAvgPower(t *testing.T) {
	type fields struct {
		Started     string
		Consumption []types.MeterValue
	}
	tests := []struct {
		name   string
		fields fields
		want   float32
	}{
		{
			name: "Avg10",
			fields: fields{
				Started: time.Now().Format(time.RFC3339),
				Consumption: []types.MeterValue{
					{
						Timestamp: &types.DateTime{Time: time.Now()},
						SampledValue: []types.SampledValue{
							{
								Value:     "10",
								Measurand: types.MeasurandPowerActiveExport,
							},
						},
					}, {
						Timestamp: &types.DateTime{Time: time.Now()},
						SampledValue: []types.SampledValue{
							{
								Value:     "10",
								Measurand: types.MeasurandPowerActiveExport,
							},
						},
					},
				},
			},
			want: 10,
		}, {
			name: "AvgMixed",
			fields: fields{
				Started: time.Now().Format(time.RFC3339),
				Consumption: []types.MeterValue{
					{
						Timestamp: &types.DateTime{Time: time.Now()},
						SampledValue: []types.SampledValue{
							{
								Value:     "10",
								Measurand: types.MeasurandPowerActiveExport,
							},
						},
					}, {
						Timestamp: &types.DateTime{Time: time.Now()},
						SampledValue: []types.SampledValue{
							{
								Value:     "2",
								Measurand: types.MeasurandCurrentExport,
							}, {
								Value:     "15",
								Measurand: types.MeasurandVoltage,
							},
						},
					},
				},
			},
			want: 20,
		}, {
			name: "AvgMixedFaulty",
			fields: fields{
				Started: time.Now().Format(time.RFC3339),
				Consumption: []types.MeterValue{
					{
						Timestamp: &types.DateTime{Time: time.Now()},
						SampledValue: []types.SampledValue{
							{
								Value:     "10",
								Measurand: types.MeasurandPowerActiveExport,
							},
						},
					}, {
						Timestamp: &types.DateTime{Time: time.Now()},
						SampledValue: []types.SampledValue{
							{
								Value:     "2",
								Measurand: types.MeasurandCurrentExport,
							},
						},
					},
				},
			},
			want: 10,
		}, {
			name: "3Measurands",
			fields: fields{
				Started: time.Now().Format(time.RFC3339),
				Consumption: []types.MeterValue{
					{
						Timestamp: &types.DateTime{Time: time.Now()},
						SampledValue: []types.SampledValue{
							{
								Value:     "10",
								Measurand: types.MeasurandPowerActiveExport,
							},
							{
								Value:     "2",
								Measurand: types.MeasurandCurrentExport,
							},
							{
								Value:     "15",
								Measurand: types.MeasurandVoltage,
							},
						},
					}, {
						Timestamp: &types.DateTime{Time: time.Now()},
						SampledValue: []types.SampledValue{
							{
								Value:     "2",
								Measurand: types.MeasurandCurrentExport,
							},
							{
								Value:     "15",
								Measurand: types.MeasurandVoltage,
							},
							{
								Value:     "10",
								Measurand: types.MeasurandPowerActiveExport,
							},
						},
					},
				},
			},
			want: 10,
		}, {
			name: "ZeroSamples",
			fields: fields{
				Started: time.Now().Format(time.RFC3339),
				Consumption: []types.MeterValue{
					{
						Timestamp:    &types.DateTime{Time: time.Now()},
						SampledValue: []types.SampledValue{},
					}, {
						Timestamp:    &types.DateTime{Time: time.Now()},
						SampledValue: []types.SampledValue{},
					},
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &data.Session{
				Started:     tt.fields.Started,
				Consumption: tt.fields.Consumption,
			}
			if got := session.CalculateAvgPower(); got != tt.want {
				t.Errorf("CalculateAvgPower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_CalculateEnergyConsumptionWithAvgPower(t *testing.T) {
	type fields struct {
		Started     string
		Consumption []types.MeterValue
	}
	tests := []struct {
		name   string
		fields fields
		want   float32
	}{
		{
			name: "ApproxConsumption",
			fields: fields{
				Started: time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
				Consumption: []types.MeterValue{
					{
						Timestamp: &types.DateTime{Time: time.Now()},
						SampledValue: []types.SampledValue{
							{
								Value:     "10",
								Measurand: types.MeasurandPowerActiveExport,
							},
						},
					},
				},
			},
			want: float32(300 * 10),
		}, {
			name: "StartTimeNow",
			fields: fields{
				Started: time.Now().Format(time.RFC3339),
				Consumption: []types.MeterValue{
					{
						Timestamp: &types.DateTime{Time: time.Now()},
						SampledValue: []types.SampledValue{
							{
								Value:     "10",
								Measurand: types.MeasurandPowerActiveExport,
							},
						},
					},
				},
			},
			want: 0.0,
		}, {
			name: "NoMeasurements",
			fields: fields{
				Started: time.Now().Format(time.RFC3339),
				Consumption: []types.MeterValue{
					{
						Timestamp:    &types.DateTime{Time: time.Now()},
						SampledValue: []types.SampledValue{},
					},
				},
			},
			want: 0.0,
		}, {
			name: "NotStarted",
			fields: fields{
				Started: "",
				Consumption: []types.MeterValue{
					{
						Timestamp:    &types.DateTime{Time: time.Now()},
						SampledValue: []types.SampledValue{},
					},
				},
			},
			want: 0.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &data.Session{
				Started:     tt.fields.Started,
				Consumption: tt.fields.Consumption,
			}
			if got := session.CalculateEnergyConsumptionWithAvgPower(); got != tt.want && tt.want < got-(got/0.005) {
				t.Errorf("CalculateEnergyConsumptionWithAvgPower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_CalculateEnergyConsumption(t *testing.T) {
	type fields struct {
		IsActive      bool
		TransactionId string
		TagId         string
		Started       string
		Consumption   []types.MeterValue
	}
	tests := []struct {
		name   string
		fields fields
		want   float32
	}{
		{
			name: "EnergyOk",
			fields: fields{
				Started: "",
				Consumption: []types.MeterValue{
					{
						Timestamp: &types.DateTime{Time: time.Now()},
						SampledValue: []types.SampledValue{
							{
								Value:     "1",
								Measurand: types.MeasurandEnergyActiveExportInterval,
							},
						},
					},
				},
			},
			want: 1,
		}, {
			name: "NoEnergySampled",
			fields: fields{
				Started: "",
				Consumption: []types.MeterValue{
					{
						Timestamp: &types.DateTime{Time: time.Now()},
						SampledValue: []types.SampledValue{
							{
								Value:     "10",
								Measurand: types.MeasurandVoltage,
							}, {
								Value:     "10",
								Measurand: types.MeasurandCurrentExport,
							}, {
								Value:     "10",
								Measurand: types.MeasurandPowerActiveExport,
							},
						},
					},
				},
			},
			want: 0.0,
		}, {
			name: "MultipleEnergySamples",
			fields: fields{
				Started: "",
				Consumption: []types.MeterValue{
					{
						Timestamp: &types.DateTime{Time: time.Now()},
						SampledValue: []types.SampledValue{
							{
								Value:     "1",
								Measurand: types.MeasurandEnergyActiveExportInterval,
							},
						},
					}, {
						Timestamp: &types.DateTime{Time: time.Now()},
						SampledValue: []types.SampledValue{
							{
								Value:     "10",
								Measurand: types.MeasurandEnergyActiveExportInterval,
							},
						},
					}, {
						Timestamp: &types.DateTime{Time: time.Now()},
						SampledValue: []types.SampledValue{
							{
								Value:     "11",
								Measurand: types.MeasurandEnergyActiveExportInterval,
							},
						},
					},
				},
			},
			want: 22.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &data.Session{
				Started:     tt.fields.Started,
				Consumption: tt.fields.Consumption,
			}
			if got := session.CalculateEnergyConsumption(); got != tt.want {
				t.Errorf("CalculateEnergyConsumption() = %v, want %v", got, tt.want)
			}
		})
	}
}
