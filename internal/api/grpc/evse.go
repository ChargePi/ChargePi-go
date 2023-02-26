package grpc

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	"github.com/xBlaz3kx/ChargePi-go/pkg/grpc"
)

type Service struct {
	grpc.UnimplementedEvseServer
	evseManager evse.Manager
}

func NewEvseService(manager evse.Manager) *Service {
	return &Service{
		evseManager: manager,
	}
}

func (s *Service) GetEVSEs(ctx context.Context, empty *empty.Empty) (*grpc.GetEvsesResponse, error) {
	response := &grpc.GetEvsesResponse{
		EVSEs: []*grpc.EVSE{},
	}

	for _, e := range s.evseManager.GetEVSEs() {
		evSe := toEvse(e)
		response.EVSEs = append(response.EVSEs, evSe)
	}

	return response, nil
}

func (s *Service) AddEVSE(ctx context.Context, request *grpc.SetEVCCRequest) (*grpc.AddEvseResponse, error) {
	response := &grpc.AddEvseResponse{
		Status: "Failed",
	}

	// todo
	// s.evseManager.AddEVSE()

	return response, nil
}

func (s *Service) GetEVSE(ctx context.Context, request *grpc.GetEvseRequest) (*grpc.GetEvseResponse, error) {
	res := &grpc.GetEvseResponse{}

	findEVSE, err := s.evseManager.FindEVSE(int(request.EvseId))
	if err != nil {
		return res, nil
	}

	res.EVSE = toEvse(findEVSE)
	return res, nil
}

func (s *Service) SetEVCC(ctx context.Context, request *grpc.SetEVCCRequest) (*grpc.SetEvccResponse, error) {
	// todo
	return nil, nil
}

func (s *Service) SetPowerMeter(ctx context.Context, request *grpc.SetPowerMeterRequest) (*grpc.SetPowerMeterDetails, error) {
	// todo
	return nil, nil
}

func (s *Service) GetUsageForEVSE(request *grpc.GetUsageForEVSERequest, server grpc.Evse_GetUsageForEVSEServer) error {
	evseWithId, err := s.evseManager.FindEVSE(int(request.EvseId))
	if err != nil {
		return err
	}

	powerMeter := evseWithId.GetPowerMeter()
	if util.IsNilInterfaceOrPointer(powerMeter) {
		return nil
	}

	ctx := server.Context()

Loop:
	for {
		select {
		case <-ctx.Done():
			break Loop
		default:
			// todo
			// Get measurements from the power meter
			powerMeter.GetPower()
			powerMeter.GetEnergy()
			powerMeter.GetCurrent(1)
			powerMeter.GetVoltage(1)

			powerMeter.GetCurrent(2)
			powerMeter.GetVoltage(2)

			powerMeter.GetCurrent(3)
			powerMeter.GetVoltage(3)
			time.Sleep(time.Second * 10)
		}
	}

	return nil
}

func (s *Service) mustEmbedUnimplementedEvseServer() {
}

func toEvse(e evse.EVSE) *grpc.EVSE {
	return &grpc.EVSE{
		Id: int32(e.GetEvseId()),
		EVCC: &grpc.EVCC{
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
