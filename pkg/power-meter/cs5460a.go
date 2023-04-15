package powerMeter

import (
	"context"
	"errors"
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/warthog618/gpiod"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
)

const (
	// VoltageRange Input range (+-) in mV
	VoltageRange = 0.250
	// CurrentRange Input range (+-) in mV
	CurrentRange            = 0.250
	StartSingleConvert      = 0xE0
	StartMultiConvert       = 0xE8
	SYNC0                   = 0xFE
	SYNC1                   = 0xFF
	PowerUpHaltControl      = 0xA0
	PowerDownMode0          = 0x80
	PowerDownMode1          = 0x88
	PowerDownMode2          = 0x90
	PowerDownMode3          = 0x98
	CalibrateControl        = 0xC0
	CalibrateCurrent        = 0x08
	CalibrateVoltage        = 0x10
	CalibrateCurrentVoltage = 0x18
	CalibrateGain           = 0x02
	CalibrateOffset         = 0x01
	CalibrateAll            = 0x1B
	ConfigRegister          = 0x00 << 1
	CurrentOffsetRegister   = 0x01 << 1
	CurrentGainRegister     = 0x02 << 1
	VoltageOffsetRegister   = 0x03 << 1
	VoltageGainRegister     = 0x04 << 1
	CycleCountRegister      = 0x05 << 1
	PulseRateRegister       = 0x06 << 1
	LastCurrentRegister     = 0x07 << 1
	LastVoltageRegister     = 0x08 << 1
	LastPowerRegister       = 0x09 << 1
	TotalEnergyRegister     = 0x0A << 1
	RmsCurrentRegister      = 0x0B << 1
	RmsVoltageRegister      = 0x0C << 1
	TimeBaseCaliRegister    = 0x0D << 1
	StatusRegister          = 0x0F << 1
	InterruptMaskRegister   = 0x1A << 1
	WriteRegister           = 0x40
	ReadRegister            = ^WriteRegister
	ChipReset               = 0x01 << 7
	SignBit                 = 0x01 << 23
	DataReady               = 0x01 << 23
	ConversionReady         = 0x01 << 20
)

type C5460A struct {
	pin                  int
	chipSelect           *gpiod.Line
	spiConnection        spi.Conn
	bus                  int
	voltageDividerOffset float64
	ShuntOffset          float64
	voltageMultiplier    float64
	currentMultiplier    float64
	powerMultiplier      float64
}

// NewCS5460PowerMeter creates a concrete implementation for the CS5460A chip for the charging station.
func NewCS5460PowerMeter(enablePin int, spiBus int, voltageDividerOffset float64, shuntOffset float64) (*C5460A, error) {
	if enablePin <= 0 {
		return nil, errors.New("pin cannot be negative")
	}

	if spiBus < 0 {
		return nil, errors.New("spi bus cannot be negative")
	}

	var powerMeter = C5460A{
		pin:                  enablePin,
		bus:                  spiBus,
		voltageDividerOffset: voltageDividerOffset,
		ShuntOffset:          shuntOffset,
		voltageMultiplier:    VoltageRange * voltageDividerOffset,
		currentMultiplier:    CurrentRange * shuntOffset,
		powerMultiplier:      voltageDividerOffset * shuntOffset,
	}

	return &powerMeter, nil
}

// Init prepares the CS5460A chip and SPI communication with the Chip.
func (receiver *C5460A) Init(ctx context.Context) error {
	if _, err := host.Init(); err != nil {
		return err
	}

	// Begin SPI communication
	s, err := spireg.Open(fmt.Sprintf("%d", receiver.bus))
	if err != nil {
		return err
	}

	c, err := s.Connect(500000, spi.Mode0, 8)
	if err != nil {
		return err
	}

	receiver.spiConnection = c

	// Send the sync sequence
	receiver.sendSync()
	return err
}

// sendBytes will send a payload to the Chip.
func (receiver *C5460A) sendBytes(payload []byte) {
	_ = receiver.chipSelect.SetValue(0)

	resp := make([]byte, 16)
	p := []spi.Packet{
		{
			W:      payload,
			KeepCS: true,
		},
		{
			R: resp,
		},
	}

	err := receiver.spiConnection.TxPackets(p)
	if err != nil {
		log.WithError(err).Error("Cannot transmit to the power meter")
	}

	_ = receiver.chipSelect.SetValue(1)
}

// sendSync will send a synchronisation sequence to the chip so it knows we are ready to pull data from the chip.
func (receiver *C5460A) sendSync() {
	receiver.sendBytes([]byte{SYNC1, SYNC1, SYNC1, SYNC0})
}

// sendToRegister will send data to a certain register.
func (receiver *C5460A) sendToRegister(register byte, data int32) {
	_ = receiver.chipSelect.SetValue(0)

	p := []spi.Packet{
		{
			W:      []byte{byte(int32(register) | WriteRegister)},
			KeepCS: true,
		},
		{
			W:      []byte{byte(data & 0xFF0000 >> 16), byte(data & 0xFF00 >> 8), byte(data & 0xFF)},
			KeepCS: true,
		},
	}

	err := receiver.spiConnection.TxPackets(p)
	if err != nil {
		log.WithError(err).Error("Cannot transmit to the power meter")
	}

	_ = receiver.chipSelect.SetValue(1)
}

// readFromRegister reads from the Chip registers.
func (receiver *C5460A) readFromRegister(register int32) int32 {
	_ = receiver.chipSelect.SetValue(0)

	received := make([]byte, 32)
	p := []spi.Packet{
		{
			W:      []byte{byte((ReadRegister) | register)},
			KeepCS: true,
		},
		{
			R: received,
		},
	}

	err := receiver.spiConnection.TxPackets(p)
	if err != nil {
		log.WithError(err).Error("Cannot transmit to the power meter")
	}

	var value int32 = 0x0
	for i := 0; i < len(received); i++ {
		value <<= 8
		value |= int32(received[i])
	}

	err = receiver.chipSelect.SetValue(1)
	if err != nil {
		return 0
	}

	return value
}

func (receiver *C5460A) getStatus() int32 {
	return receiver.readFromRegister(StatusRegister)
}

func (receiver *C5460A) clearStatus(status int32) {
	receiver.sendToRegister(StatusRegister, status)
}

func (receiver *C5460A) startConverting() {
	receiver.clearStatus(ConversionReady)
	receiver.sendBytes([]byte{StartMultiConvert})
	for {
		if receiver.getStatus()&ConversionReady == 0 {
			break
		}
	}
}

func (receiver *C5460A) stopConverting() {
	receiver.sendBytes([]byte{PowerUpHaltControl})
}

func (receiver *C5460A) Reset() {
	receiver.sendToRegister(ConfigRegister, ChipReset)
	receiver.sendSync()
}

func (receiver *C5460A) calibrate(command byte) {
	receiver.stopConverting()
	receiver.clearStatus(DataReady)
	receiver.sendBytes([]byte{CalibrateControl | (command & CalibrateAll)})
	for {
		if receiver.getStatus()&DataReady == 0 {
			break
		}
	}
	receiver.clearStatus(DataReady)
	receiver.startConverting()
}

func (receiver *C5460A) calibrateVoltage() int32 {
	receiver.calibrate(CalibrateVoltage | CalibrateOffset)
	return receiver.readFromRegister(VoltageOffsetRegister)
}

func (receiver *C5460A) setVoltageOffset(value int32) {
	receiver.stopConverting()
	receiver.sendToRegister(VoltageOffsetRegister, value)
	receiver.startConverting()
}

func (receiver *C5460A) calibrateCurrentOffset() int32 {
	receiver.calibrate(CalibrateCurrent | CalibrateOffset)
	return receiver.readFromRegister(CurrentOffsetRegister)
}

func (receiver *C5460A) setCurrentOffset(value int32) {
	receiver.stopConverting()
	receiver.sendToRegister(CurrentOffsetRegister, value)
	receiver.startConverting()
}

func (receiver *C5460A) GetEnergy() (*types.SampledValue, error) {
	value := float64(receiver.readFromRegister(TotalEnergyRegister) << 16)

	return &types.SampledValue{
		Value:     fmt.Sprintf("%.1f", value),
		Measurand: types.MeasurandEnergyActiveImportRegister,
	}, nil
}

func (receiver *C5460A) GetPower(phase int) (*types.SampledValue, error) {
	value := float64(receiver.readFromRegister(LastPowerRegister)) * receiver.powerMultiplier

	return &types.SampledValue{
		Value:     fmt.Sprintf("%.1f", value),
		Measurand: types.MeasurandPowerActiveImport,
	}, nil
}

func (receiver *C5460A) GetApparentPower(phase int) (*types.SampledValue, error) {
	value := float64(receiver.readFromRegister(LastPowerRegister)) * receiver.powerMultiplier

	return &types.SampledValue{
		Value:     fmt.Sprintf("%.1f", value),
		Measurand: types.MeasurandPowerActiveImport,
	}, nil
}

func (receiver *C5460A) GetReactivePower(phase int) (*types.SampledValue, error) {
	value := float64(receiver.readFromRegister(LastPowerRegister)) * receiver.powerMultiplier

	return &types.SampledValue{
		Value:     fmt.Sprintf("%.1f", value),
		Measurand: types.MeasurandPowerActiveImport,
	}, nil
}

func (receiver *C5460A) GetCurrent(phase int) (*types.SampledValue, error) {
	value := float64(receiver.readFromRegister(LastCurrentRegister)) * receiver.currentMultiplier

	return &types.SampledValue{
		Value:     fmt.Sprintf("%.1f", value),
		Phase:     getPhase(phase),
		Measurand: types.MeasurandCurrentImport,
	}, nil
}

func (receiver *C5460A) GetVoltage(phase int) (*types.SampledValue, error) {
	value := float64(receiver.readFromRegister(LastVoltageRegister)) * receiver.voltageMultiplier

	return &types.SampledValue{
		Value:     fmt.Sprintf("%.1f", value),
		Phase:     getPhase(phase),
		Measurand: types.MeasurandVoltage,
	}, nil
}

func (receiver *C5460A) getRMSCurrent(phase int) (*types.SampledValue, error) {
	value := float64(receiver.readFromRegister(RmsCurrentRegister)) * receiver.currentMultiplier

	return &types.SampledValue{
		Value:     fmt.Sprintf("%.1f", value),
		Phase:     getPhase(phase),
		Measurand: types.MeasurandCurrentImport,
	}, nil
}

func (receiver *C5460A) getRMSVoltage(phase int) (*types.SampledValue, error) {
	value := float64(receiver.readFromRegister(RmsVoltageRegister)) * receiver.voltageMultiplier
	return &types.SampledValue{
		Value:     fmt.Sprintf("%.1f", value),
		Phase:     getPhase(phase),
		Measurand: types.MeasurandVoltage,
	}, nil
}

func (receiver *C5460A) GetType() string {
	return TypeC5460A
}

func getPhase(phase int) types.Phase {
	switch phase {
	case 1:
		return types.PhaseL1
	case 2:
		return types.PhaseL2
	case 3:
		return types.PhaseL3
	}
	return ""
}
