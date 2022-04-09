package powerMeter

import (
	"errors"
	"github.com/warthog618/gpiod"
)

const (
	// VoltageRange Input range (+-) in mV
	VoltageRange float64 = 0.250
	// CurrentRange Input range (+-) in mV
	CurrentRange            float64 = 0.250
	StartSingleConvert              = 0xE0
	StartMultiConvert               = 0xE8
	SYNC0                   byte    = 0xFE
	SYNC1                   byte    = 0xFF
	PowerUpHaltControl              = 0xA0
	PowerDownMode0                  = 0x80
	PowerDownMode1                  = 0x88
	PowerDownMode2                  = 0x90
	PowerDownMode3                  = 0x98
	CalibrateControl                = 0xC0
	CalibrateCurrent                = 0x08
	CalibrateVoltage                = 0x10
	CalibrateCurrentVoltage         = 0x18
	CalibrateGain                   = 0x02
	CalibrateOffset                 = 0x01
	CalibrateAll                    = 0x1B
	ConfigRegister                  = 0x00 << 1
	CurrentOffsetRegister           = 0x01 << 1
	CurrentGainRegister             = 0x02 << 1
	VoltageOffsetRegister           = 0x03 << 1
	VoltageGainRegister             = 0x04 << 1
	CycleCountRegister              = 0x05 << 1
	PulseRateRegister               = 0x06 << 1
	LastCurrentRegister             = 0x07 << 1
	LastVoltageRegister             = 0x08 << 1
	LastPowerRegister               = 0x09 << 1
	TotalEnergyRegister             = 0x0A << 1
	RmsCurrentRegister              = 0x0B << 1
	RmsVoltageRegister              = 0x0C << 1
	TimeBaseCaliRegister            = 0x0D << 1
	StatusRegister                  = 0x0F << 1
	InterruptMaskRegister           = 0x1A << 1
	WriteRegister           int32   = 0x40
	ReadRegister            int32   = ^WriteRegister
	ChipReset                       = 0x01 << 7
	SignBit                         = 0x01 << 23
	DataReady                       = 0x01 << 23
	ConversionReady                 = 0x01 << 20
)

type C5460A struct {
	EnablePin            int
	chipSelect           *gpiod.Line
	SpiBus               int
	VoltageDividerOffset float64
	ShuntOffset          float64
	VoltageMultiplier    float64
	CurrentMultiplier    float64
	PowerMultiplier      float64
}

func NewCS5460PowerMeter(enablePin int, spiBus int, voltageDividerOffset float64, shuntOffset float64) (*C5460A, error) {
	if enablePin <= 0 {
		return nil, errors.New("pin cannot be negative")
	}

	if spiBus < 0 {
		return nil, errors.New("spi bus cannot be negative")
	}

	var powerMeter = C5460A{EnablePin: enablePin, SpiBus: spiBus, VoltageDividerOffset: voltageDividerOffset, ShuntOffset: shuntOffset}
	err := powerMeter.init()
	if err != nil {
		return nil, err
	}

	/*var spiDev rpio.SpiDev
	if spiDev = rpio.Spi0; powerMeter.SpiBus == 0 {
		spiDev = rpio.Spi1
	}
	rpio.SpiSpeed(500000)
	err := rpio.SpiBegin(spiDev)*/
	if err != nil {
		return nil, err
	}
	powerMeter.sendSync()
	return &powerMeter, nil
}

func (receiver *C5460A) init() error {
	receiver.VoltageMultiplier = VoltageRange * receiver.VoltageDividerOffset
	receiver.CurrentMultiplier = CurrentRange * receiver.ShuntOffset
	receiver.PowerMultiplier = receiver.VoltageDividerOffset * receiver.ShuntOffset
	// Refer to gpiod docs
	c, err := gpiod.NewChip("gpiochip0")
	if err != nil {
		return err
	}
	receiver.chipSelect, err = c.RequestLine(receiver.EnablePin)
	return err
}

func (receiver *C5460A) sendBytes(payload []byte) {
	receiver.chipSelect.SetValue(0)
	//rpio.SpiTransmit(payload...)
	receiver.chipSelect.SetValue(1)
}

func (receiver *C5460A) sendSync() {
	receiver.sendBytes([]byte{SYNC1, SYNC1, SYNC1, SYNC0})
}

func (receiver *C5460A) sendToRegister(register byte, data int32) {
	receiver.chipSelect.SetValue(0)
	//rpio.SpiTransmit(byte(int32(register) | WriteRegister))
	//rpio.SpiTransmit([]byte{byte(data & 0xFF0000 >> 16), byte(data & 0xFF00 >> 8), byte(data & 0xFF)}...)
	receiver.chipSelect.SetValue(1)
}

func (receiver *C5460A) readFromRegister(register int32) int32 {
	receiver.chipSelect.SetValue(0)
	//rpio.SpiTransmit(byte((ReadRegister) | register))
	var received = []byte{0x0} //rpio.SpiReceive(3)
	var value int32 = 0x0
	for i := 0; i < len(received); i++ {
		value <<= 8
		value |= int32(received[i])
	}
	receiver.chipSelect.SetValue(1)
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

func (receiver *C5460A) GetEnergy() float64 {
	return float64(receiver.readFromRegister(TotalEnergyRegister) << 16)
}
func (receiver *C5460A) GetPower() float64 {
	return float64(receiver.readFromRegister(LastPowerRegister)) * receiver.PowerMultiplier
}
func (receiver *C5460A) GetCurrent() float64 {
	return float64(receiver.readFromRegister(LastCurrentRegister)) * receiver.CurrentMultiplier
}
func (receiver *C5460A) GetVoltage() float64 {
	return float64(receiver.readFromRegister(LastVoltageRegister)) * receiver.VoltageMultiplier
}
func (receiver *C5460A) GetRMSCurrent() float64 {
	return float64(receiver.readFromRegister(RmsCurrentRegister)) * receiver.CurrentMultiplier
}
func (receiver *C5460A) GetRMSVoltage() float64 {
	return float64(receiver.readFromRegister(RmsVoltageRegister)) * receiver.VoltageMultiplier
}
