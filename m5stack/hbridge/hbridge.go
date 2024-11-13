package hbridge

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"periph.io/x/conn/v3/i2c"
)

const (
	HBRIDGE_I2C_ADDR               = 0x20
	HBRIDGE_CONFIG_REG             = 0x00
	HBRIDGE_MOTOR_ADC_8BIT_REG     = 0x10
	HBRIDGE_MOTOR_ADC_12BIT_REG    = 0x20
	HBRIDGE_MOTOR_CURRENT_REG      = 0x30
	HBRIDGE_JUMP_TO_BOOTLOADER_REG = 0xFD
	HBRIDGE_FW_VERSION_REG         = 0xFE
	HBRIDGE_I2C_ADDRESS_REG        = 0xFF
)

type HbridgeDirection int

const (
	HBRIDGE_STOP HbridgeDirection = iota
	HBRIDGE_FORWARD
	HBRIDGE_BACKWARD
)

type HbridgeAnalogReadMode int

const (
	_8bit HbridgeAnalogReadMode = iota
	_12bit
)

// I2CAddr is the default I2C address for the m5stack HBrige unit.
const I2CAddr uint16 = HBRIDGE_I2C_ADDR

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

// Dev is an handle to an M5Stack HBrige Motors driver.
type Dev struct {
	c i2c.Dev
}

// New creates a new driver for M5Stack HBrige motor driver.
func New(bus i2c.Bus, opts *Opts) (*Dev, error) {
	if opts.I2cAddress < 0x01 || opts.I2cAddress > 0x70 {
		return nil, fmt.Errorf("invalid device address")
	}

	dev := &Dev{c: i2c.Dev{Bus: bus, Addr: opts.I2cAddress}}
	dev.SetDriverDirection(HBRIDGE_STOP)

	return dev, nil
}

func (dev *Dev) Close() {
	if dev != nil {
		dev.SetDriverDirection(HBRIDGE_STOP)
	}
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

func (h *Dev) GetDriverDirection() (uint8, error) {
	data, err := h.readBytes(HBRIDGE_CONFIG_REG, 1)
	if err != nil {
		return 0, err
	}
	return data[0], err
}

func (h *Dev) GetDriverSpeed8Bits() (uint8, error) {
	data, err := h.readBytes(HBRIDGE_CONFIG_REG+1, 1)
	if err != nil {
		return 0, err
	}
	return data[0], err
}

func (h *Dev) GetDriverSpeed16Bits() (uint16, error) {
	data, err := h.readBytes(HBRIDGE_CONFIG_REG+2, 2)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(data), err
}

func (h *Dev) GetDriverPWMFreq() (uint16, error) {
	data, err := h.readBytes(HBRIDGE_CONFIG_REG+4, 2)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(data), err
}

func (h *Dev) SetDriverPWMFreq(freq uint16) error {
	data := []uint8{uint8(freq & 0xff), uint8((freq >> 8) & 0xff)}
	return h.writeBytes(HBRIDGE_CONFIG_REG+4, data)
}

func (h *Dev) SetDriverDirection(dir HbridgeDirection) error {
	data := []uint8{uint8(dir)}
	return h.writeBytes(HBRIDGE_CONFIG_REG, data)
}

func (h *Dev) SetDriverSpeed8Bits(speed uint8) error {
	data := []uint8{speed}
	return h.writeBytes(HBRIDGE_CONFIG_REG+1, data)
}

func (h *Dev) SetDriverSpeed16Bits(speed uint16) error {
	data := []uint8{uint8(speed), uint8(speed >> 8)}
	return h.writeBytes(HBRIDGE_CONFIG_REG+2, data)
}

func (h *Dev) GetAnalogInput(bit HbridgeAnalogReadMode) (uint16, error) {
	if bit == _8bit {
		data, err := h.readBytes(HBRIDGE_MOTOR_ADC_8BIT_REG, 1)
		if err != nil {
			return 0, err
		}
		return uint16(data[0]), err
	} else {
		data, err := h.readBytes(HBRIDGE_MOTOR_ADC_12BIT_REG, 2)
		if err != nil {
			return 0, err
		}
		return binary.LittleEndian.Uint16(data), err
	}
}

func (h *Dev) GetMotorCurrent() (float32, error) {
	data, err := h.readBytes(HBRIDGE_MOTOR_CURRENT_REG, 4)
	if err != nil {
		return 0, err
	}
	var c float32
	binary.Read(bytes.NewReader(data), binary.LittleEndian, &c)
	return c, err
}

func (h *Dev) GetFirmwareVersion() (uint8, error) {
	data, err := h.readBytes(HBRIDGE_FW_VERSION_REG, 1)
	if err != nil {
		return 0, err
	}
	return data[0], err
}

func (h *Dev) GetI2CAddress() (uint8, error) {
	data, err := h.readBytes(HBRIDGE_I2C_ADDRESS_REG, 1)
	if err != nil {
		return 0, err
	}
	return data[0], err
}

func (h *Dev) JumpBootloader() error {
	value := uint8(1)
	return h.writeBytes(HBRIDGE_JUMP_TO_BOOTLOADER_REG, []uint8{value})
}
