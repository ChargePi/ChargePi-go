package v16

import (
	"context"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
	"github.com/xBlaz3kx/ChargePi-go/test"
	ocppVar "github.com/xBlaz3kx/ocppManager-go/configuration"
	"strings"
	"time"
)

func (s *chargePointTestSuite) TestStartStopTransaction() {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*30)
		conn        = new(test.EvseMock)
	)

	// Setup expectations
	conn.On("StartCharging", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	conn.On("ResumeCharging", mock.Anything).Return(nil, 0)
	conn.On("StopCharging", core.ReasonLocal).Return(nil)
	conn.On("RemoveReservation").Return(nil)
	conn.On("GetReservationId").Return(0)
	conn.On("GetTagId").Return(tagId)
	conn.On("GetTransactionId").Return("1")
	conn.On("GetConnectorId").Return(1)
	conn.On("GetEvseId").Return(1)
	conn.On("CalculateSessionAvgEnergyConsumption").Return(30.0)
	conn.On("GetStatus").Return(core.ChargePointStatusAvailable, core.NoError)
	conn.On("IsAvailable").Return(true).Once()
	conn.On("IsPreparing").Return(false)
	conn.On("IsCharging").Return(true).Once()
	conn.On("IsReserved").Return(false)
	conn.On("IsUnavailable").Return(false)
	conn.On("GetMaxChargingTime").Return(15)
	conn.On("SetNotificationChannel", mock.Anything).Return()

	s.manager.On("GetEVSEs").Return([]evse.EVSE{conn})
	s.manager.On("FindEVSE", 1, 1).Return(conn)
	s.manager.On("FindAvailableEVSE").Return(conn)
	s.manager.On("FindEVSEWithTagId", tagId).Return(nil).Once()
	s.manager.On("FindEVSEWithTransactionId", "1").Return(nil).Once()
	s.manager.On("StartCharging").Return()
	s.manager.On("StopCharging").Return()
	s.manager.On("StopAllEVSEs").Return()
	s.manager.On("AddEVSE", conn).Return(nil)
	s.manager.On("AddEVSEFromSettings", mock.Anything).Return(nil)
	s.manager.On("AddEVSEsFromSettings", mock.Anything).Return(nil)
	s.manager.On("RestoreEVSEStatus", mock.Anything).Return(nil)
	s.manager.On("SetNotificationChannel").Return()

	// Create and connect the Charge Point
	chargePoint := s.setupChargePoint(ctx, nil, nil, s.manager)

	time.Sleep(time.Second * 2)

	go func() {
		time.Sleep(time.Second * 3)

		//todo start charging

		// Redeclare expectations
		s.manager.On("FindEVSEWithTagId", tagId).Return(conn)
		s.manager.On("FindEVSEWithTransactionId", "1").Return(conn)

		conn.On("IsCharging").Return(true).Once()
		conn.On("IsAvailable").Return(false).Once()
		conn.On("GetStatus").Return(core.ChargePointStatusCharging, core.NoError)

		time.Sleep(time.Second * 5)

		// Stop charging

		cancel()
	}()

Loop:
	for {
		select {
		case <-ctx.Done():
			cancel()
			break Loop
		}
	}

	chargePoint.CleanUp("")
}

func (s *chargePointTestSuite) TestStartStopTransactionWithReader() {
	var (
		ctx, cancel   = context.WithTimeout(context.Background(), time.Second*30)
		conn          = new(test.EvseMock)
		readerChannel = make(chan string)
	)

	// Setup expectations
	conn.On("StartCharging", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	conn.On("ResumeCharging", mock.Anything).Return(nil, 0)
	conn.On("StopCharging", mock.Anything).Return(nil)
	conn.On("RemoveReservation").Return(nil)
	conn.On("GetReservationId").Return(0)
	conn.On("GetTransactionId").Return("1")
	conn.On("GetTagId").Return(strings.ToUpper(tagId))
	conn.On("GetConnectorId").Return(1)
	conn.On("GetEvseId").Return(1)
	conn.On("CalculateSessionAvgEnergyConsumption").Return(30.0)
	conn.On("GetStatus").Return(core.ChargePointStatusAvailable, core.NoError)
	conn.On("IsAvailable").Return(true).Once()
	conn.On("IsPreparing").Return(false)
	conn.On("IsCharging").Return(true).Once()
	conn.On("IsReserved").Return(false)
	conn.On("IsUnavailable").Return(false)
	conn.On("GetMaxChargingTime").Return(15)
	conn.On("SetNotificationChannel", mock.Anything).Return()

	s.manager.On("GetEVSEs").Return([]evse.EVSE{conn})
	s.manager.On("FindEVSE", 1, 1).Return(conn)
	s.manager.On("FindAvailableEVSE").Return(conn)
	s.manager.On("FindEVSEWithTagId", strings.ToUpper(tagId)).Return(nil).Once()
	s.manager.On("FindEVSEWithTransactionId", "1").Return(nil).Once()
	s.manager.On("StartCharging").Return()
	s.manager.On("StopCharging").Return()
	s.manager.On("StopAllEVSEs").Return()
	s.manager.On("AddEVSE", conn).Return(nil)
	s.manager.On("AddEVSEFromSettings", mock.Anything).Return(nil)
	s.manager.On("AddEVSEsFromSettings", mock.Anything).Return(nil)
	s.manager.On("RestoreEVSEStatus", mock.Anything).Return(nil)
	s.manager.On("SetNotificationChannel").Return()

	// Mock tagReader
	s.tagReader.On("ListenForTags").Return()
	s.tagReader.On("Cleanup").Return()
	s.tagReader.On("Reset").Return()
	s.tagReader.On("GetTagChannel").Return(readerChannel)

	// Create and connect the Charge Point
	cp := s.setupChargePoint(ctx, nil, s.tagReader, s.manager)

	// Simulate reading a card
	go func() {
		time.Sleep(time.Second * 3)
		// Card read once - start charging
		log.Debug("Sending tag to reader")
		readerChannel <- tagId

		time.Sleep(time.Second * 10)

		s.manager.On("FindEVSEWithTagId", strings.ToUpper(tagId)).Return(conn)
		conn.On("IsCharging").Return(true).Once()
		conn.On("IsAvailable").Return(false).Once()
		conn.On("GetStatus").Return(core.ChargePointStatusCharging, core.NoError)

		// Card read second time - stop charging
		log.Debug("Sending tag to reader")
		readerChannel <- tagId
		time.Sleep(time.Second * 4)

		cancel()
	}()

Loop:
	for {
		select {
		case <-ctx.Done():
			cancel()
			break Loop
		}
	}

	cp.CleanUp("")
	//s.tagReader.AssertCalled(s.T(), "ListenForTags")
}

func (s *chargePointTestSuite) TestRemoteStartStopTransaction() {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*30)
		conn        = new(test.EvseMock)
	)

	// Setup expectations
	conn.On("StartCharging", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	conn.On("ResumeCharging", mock.Anything).Return(nil, 0)
	conn.On("StopCharging", mock.Anything).Return(nil)
	conn.On("RemoveReservation").Return(nil)
	conn.On("GetReservationId").Return(0)
	conn.On("GetTransactionId").Return("1")
	conn.On("GetTagId").Return(strings.ToUpper(tagId))
	conn.On("GetConnectorId").Return(1)
	conn.On("GetEvseId").Return(1)
	conn.On("CalculateSessionAvgEnergyConsumption").Return(30.0)
	conn.On("GetStatus").Return(core.ChargePointStatusAvailable, core.NoError)
	conn.On("IsAvailable").Return(true).Once()
	conn.On("IsPreparing").Return(false)
	conn.On("IsCharging").Return(true).Once()
	conn.On("IsReserved").Return(false)
	conn.On("IsUnavailable").Return(false)
	conn.On("GetMaxChargingTime").Return(15)
	conn.On("SetNotificationChannel", mock.Anything).Return()

	s.manager.On("GetEVSEs").Return([]evse.EVSE{conn})
	s.manager.On("FindEVSE", 1, 1).Return(conn)
	s.manager.On("FindAvailableEVSE").Return(conn)
	s.manager.On("FindEVSEWithTagId", strings.ToUpper(tagId)).Return(nil).Once()
	s.manager.On("FindEVSEWithTransactionId", "1").Return(nil).Once()
	s.manager.On("StartCharging").Return()
	s.manager.On("StopCharging").Return()
	s.manager.On("StopAllEVSEs").Return()
	s.manager.On("AddEVSE", conn).Return(nil)
	s.manager.On("AddEVSEFromSettings", mock.Anything).Return(nil)
	s.manager.On("AddEVSEsFromSettings", mock.Anything).Return(nil)
	s.manager.On("RestoreEVSEStatus", mock.Anything).Return(nil)
	s.manager.On("SetNotificationChannel").Return()

	// Create and connect the Charge Point
	cp := s.setupChargePoint(ctx, nil, nil, s.manager)

	// Simulate reading a card
	go func() {
		time.Sleep(time.Second * 3)

		// Request remote start transaction
		err := s.centralSystem.RemoteStartTransaction("/"+chargePointId, func(confirmation *core.RemoteStartTransactionConfirmation, err error) {
			s.Require().NoError(err)
			s.Assert().Equal(types.RemoteStartStopStatusAccepted, confirmation.Status)
		}, tagId)
		s.Require().NoError(err)

		time.Sleep(time.Second * 10)

		s.manager.On("FindEVSEWithTagId", strings.ToUpper(tagId)).Return(conn)
		conn.On("IsCharging").Return(true).Once()
		conn.On("IsAvailable").Return(false).Once()
		conn.On("GetStatus").Return(core.ChargePointStatusCharging, core.NoError)

		// Request remote stop transaction
		err = s.centralSystem.RemoteStopTransaction("/"+chargePointId, func(confirmation *core.RemoteStopTransactionConfirmation, err error) {}, 1)
		s.Require().NoError(err)

		cancel()
	}()

Loop:
	for {
		select {
		case <-ctx.Done():
			cancel()
			break Loop
		}
	}

	cp.CleanUp("")
	//s.tagReader.AssertCalled(s.T(), "ListenForTags")
}

func (s *chargePointTestSuite) TestGetAndChangeConfiguration() {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*30)
	)

	s.manager.On("GetEVSEs").Return([]evse.EVSE{})
	s.manager.On("AddEVSE", mock.Anything).Return(nil)
	s.manager.On("AddEVSEFromSettings", mock.Anything).Return(nil)
	s.manager.On("AddEVSEsFromSettings", mock.Anything).Return(nil)
	s.manager.On("RestoreEVSEStatus", mock.Anything).Return(nil)
	s.manager.On("SetNotificationChannel").Return()

	// Create and connect the Charge Point
	cp := s.setupChargePoint(ctx, nil, nil, s.manager)

	go func() {
		time.Sleep(time.Second * 3)

		// Get the configuration
		err := s.centralSystem.GetConfiguration("/"+chargePointId, func(confirmation *core.GetConfigurationConfirmation, err error) {
			s.Require().NoError(err)
			s.Assert().NotEmpty(confirmation.ConfigurationKey)
		}, []string{})
		s.Require().NoError(err)

		time.Sleep(time.Second * 2)

		// Change the configuration
		err = s.centralSystem.ChangeConfiguration("/"+chargePointId, func(confirmation *core.ChangeConfigurationConfirmation, err error) {
			s.Require().NoError(err)
			s.Assert().EqualValues(core.ConfigurationStatusAccepted, confirmation.Status)
			cancel()
		}, ocppVar.AuthorizationCacheEnabled.String(), "false")
		s.Require().NoError(err)
	}()

Loop:
	for {
		select {
		case <-ctx.Done():
			cancel()
			break Loop
		}
	}

	cp.CleanUp("")
}
