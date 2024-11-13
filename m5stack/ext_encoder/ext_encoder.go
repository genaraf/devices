package ext_encoder

import (
	"encoding/binary"
	"fmt"

	"periph.io/x/conn/v3/i2c"
)

const (
	UNIT_EXT_ENCODER_ADDR                 = 0x59
	UNIT_EXT_ENCODER_ENCODER_REG          = 0x00
	UNIT_EXT_ENCODER_METER_REG            = 0x10
	UNIT_EXT_ENCODER_METER_STRING_REG     = 0x20
	UNIT_EXT_ENCODER_RESET_REG            = 0x30
	UNIT_EXT_ENCODER_PERIMETER_REG        = 0x40
	UNIT_EXT_ENCODER_PULSE_REG            = 0x50
	UNIT_EXT_ENCODER_ZERO_PULSE_VALUE_REG = 0x60
	UNIT_EXT_ENCODER_ZERO_MODE_REG        = 0x70
	FIRMWARE_VERSION_REG                  = 0xFE
	I2C_ADDRESS_REG                       = 0xFF
)

// I2CAddr is the default I2C address for the m5stack Exth.
const I2CAddr uint16 = UNIT_EXT_ENCODER_ADDR

// Opts holds the configuration options.
type Opts struct {
	I2cAddress uint16
	PwmFreq    int16
}

// DefaultOpts are the recommended default options.
var DefaultOpts = Opts{
	I2cAddress: I2CAddr,
	PwmFreq:    1500,
}

// Dev is an handle to an M%Stack ExtEncoder units driver.
type Dev struct {
	c i2c.Dev
}

// New creates a new driver.
func New(bus i2c.Bus, opts *Opts) (*Dev, error) {
	if opts.I2cAddress < 0x01 || opts.I2cAddress > 0x70 {
		return nil, fmt.Errorf("invalid device address")
	}

	dev := &Dev{c: i2c.Dev{Bus: bus, Addr: opts.I2cAddress}}
	return dev, nil
}

func (dev *Dev) Close() {
}

func (h *Dev) readBytes(reg int, size int) ([]uint8, error) {
	r := make([]byte, size)
	err := h.c.Tx([]byte{byte(reg)}, r)
	return r, err
}

func (h *Dev) writeBytes(reg int, data []uint8) error {
	d := []byte{byte(reg)}
	d = append(d, data...)
	return h.c.Tx(d, nil)
}

func (h *Dev) getEncoderValue() (uint32, error) {
	data, err := h.readBytes(UNIT_EXT_ENCODER_ENCODER_REG, 4)
	if err != nil {
		return 0, err
	}
	value := binary.LittleEndian.Uint32(data)
	return value, err
}

func (h *Dev) getZeroPulseValue() (uint32, error) {
	data, err := h.readBytes(UNIT_EXT_ENCODER_ZERO_PULSE_VALUE_REG, 4)
	if err != nil {
		return 0, err
	}
	value := binary.LittleEndian.Uint32(data)
	return value, err
}

func (h *Dev) setZeroPulseValue(value uint32) error {
	data := make([]uint8, 4)
	binary.LittleEndian.PutUint32(data, value)
	return h.writeBytes(UNIT_EXT_ENCODER_ZERO_PULSE_VALUE_REG, data)
}

func (h *Dev) getMeterValue() (uint32, error) {
	data, err := h.readBytes(UNIT_EXT_ENCODER_METER_REG, 4)
	if err != nil {
		return 0, err
	}
	value := binary.LittleEndian.Uint32(data)
	return value, err
}

func (h *Dev) getMeterString() (string, error) {
	data, err := h.readBytes(UNIT_EXT_ENCODER_METER_STRING_REG, 9)
	if err != nil {
		return "", err
	}
	return string(data), err
}

func (h *Dev) resetEncoder() error {
	data := make([]uint8, 1)
	data[0] = 1
	return h.writeBytes(UNIT_EXT_ENCODER_RESET_REG, data)
}

func (h *Dev) setPerimeter(perimeter uint32) error {
	data := make([]uint8, 8)
	binary.LittleEndian.PutUint32(data, perimeter)
	return h.writeBytes(UNIT_EXT_ENCODER_PERIMETER_REG, data)
}

func (h *Dev) setZeroMode(mode uint8) error {
	data := make([]uint8, 1)
	data[0] = mode
	return h.writeBytes(UNIT_EXT_ENCODER_ZERO_MODE_REG, data)
}

func (h *Dev) getPerimeter() (uint32, error) {
	data, err := h.readBytes(UNIT_EXT_ENCODER_PERIMETER_REG, 4)
	if err != nil {
		return 0, err
	}
	value := binary.LittleEndian.Uint32(data)
	return value, err
}

func (h *Dev) SetPulse(pulse uint32) error {
	data := make([]uint8, 4)
	binary.LittleEndian.PutUint32(data, pulse)
	return h.writeBytes(UNIT_EXT_ENCODER_PULSE_REG, data)
}

func (h *Dev) GetPulse() (uint32, error) {
	data, err := h.readBytes(UNIT_EXT_ENCODER_PULSE_REG, 4)
	if err != nil {
		return 0, err
	}
	value := binary.LittleEndian.Uint32(data)
	return value, err
}

func (h *Dev) SetI2CAddress(addr uint8) error {
	data := make([]uint8, 1)
	data[0] = addr
	err := h.writeBytes(I2C_ADDRESS_REG, data)
	if err != nil {
		h.c.Addr = uint16(addr)
	}
	return err
}

func (h *Dev) GetI2CAddress() (uint8, error) {
	data, err := h.readBytes(I2C_ADDRESS_REG, 1)
	if err != nil {
		return 0, err
	}
	return data[0], err
}

func (h *Dev) GetFirmwareVersion() (uint8, error) {
	data, err := h.readBytes(FIRMWARE_VERSION_REG, 1)
	if err != nil {
		return 0, err
	}
	return data[0], err
}
