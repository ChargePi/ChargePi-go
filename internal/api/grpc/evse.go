package grpc

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	commonv1 "github.com/ChargePi/ChargePi-go/gen/proto/common/v1"
	grpc "github.com/ChargePi/ChargePi-go/gen/proto/evse/v1"
	"github.com/ChargePi/ChargePi-go/internal/evse"
	"github.com/ChargePi/ChargePi-go/internal/evse/manager"
)

type EvseHandler struct {
	grpc.UnimplementedEvseServiceServer
	evseManager manager.Manager
}

func NewEvseHandler(manager manager.Manager) *EvseHandler {
	return &EvseHandler{
		evseManager: manager,
	}
}

func (s *EvseHandler) GetEVSEs(ctx context.Context, empty *empty.Empty) (*grpc.GetEVSEsResponse, error) {
	response := &grpc.GetEVSEsResponse{
		Evses: []*grpc.EVSE{},
	}

	for _, e := range s.evseManager.GetEVSEs() {
		evseDto := toEvse(e)
		response.Evses = append(response.Evses, evseDto)
	}

	return response, nil
}

func (s *EvseHandler) AddEVSE(ctx context.Context, request *grpc.AddEVSERequest) (*grpc.AddEVSEResponse, error) {
	response := &grpc.AddEVSEResponse{
		Status: "Failed",
	}

	// todo
	err := s.evseManager.AddEVSEFromSettings(ctx, evse.Settings{})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *EvseHandler) GetEVSE(ctx context.Context, request *grpc.GetEVSERequest) (*grpc.GetEVSEResponse, error) {
	res := &grpc.GetEVSEResponse{}

	findEVSE, err := s.evseManager.GetEVSE(int(request.GetEvseId()))
	if err != nil {
		return res, nil
	}

	res.Evse = toEvse(findEVSE)
	return res, nil
}

func (s *EvseHandler) SetEVCC(ctx context.Context, request *grpc.SetEVCCRequest) (*grpc.SetEVCCResponse, error) {
	// todo

	evse, err := s.evseManager.GetEVSE(int(request.GetEvseId()))
	if err != nil {
		return nil, err
	}

	evse.SetEvcc(nil)

	return nil, nil
}

func (s *EvseHandler) SetPowerMeter(ctx context.Context, request *grpc.SetPowerMeterRequest) (*grpc.SetPowerMeterResponse, error) {
	// todo
	return nil, nil
}

func (s *EvseHandler) GetUsageForEVSE(request *grpc.GetUsageForEVSERequest, server grpc.EvseService_GetUsageForEVSEServer) error {
	evseWithId, err := s.evseManager.GetEVSE(int(request.GetEvseId()))
	if err != nil {
		return err
	}

	ctx := server.Context()

Loop:
	for {
		select {
		case <-ctx.Done():
			break Loop
		default:

			// Sample power meter
			samples, err := evseWithId.SamplePowerMeter([]types.Measurand{types.MeasurandEnergyActiveImportRegister})
			if err != nil {
				return err
			}

			// Convert to grpc samples
			var samplesToReturn []*commonv1.Sample
			for _, sample := range samples {
				samplesToReturn = append(samplesToReturn, toSample(sample))
			}

			err = server.Send(&grpc.GetUsageForEVSEResponse{
				Samples: samplesToReturn,
			})
			if err != nil {
				return err
			}
			time.Sleep(time.Second * 10)
		}
	}

	return nil
}

func (s *EvseHandler) mustEmbedUnimplementedEvseServer() {
}

func toEvse(e evse.EVSE) *grpc.EVSE {
	return &grpc.EVSE{
		Id: int32(e.GetEvseId()),
		Evcc: &grpc.EVCC{
			Type:   e.GetEvcc().GetType(),
			Status: string(e.GetEvcc().GetState()),
		},
		PowerMeter: &grpc.PowerMeter{
			Type:    e.GetPowerMeter().GetType(),
			Enabled: false,
		},
		Status:  0,
		Session: &grpc.Session{},
	}
}

func toSample(sample types.SampledValue) *commonv1.Sample {
	consumption := &commonv1.Consumption{
		Unit: string(sample.Unit),
		// Value: sample.Value,
	}
	return &commonv1.Sample{
		Consumption: consumption,
		Measureand:  string(sample.Measurand),
		Phase:       lo.ToPtr(string(sample.Phase)),
		Location:    lo.ToPtr(string(sample.Location)),
		Context:     lo.ToPtr(string(sample.Context)),
		Timestamp:   timestamppb.New(time.Now()),
	}
}
