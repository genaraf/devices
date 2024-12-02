package servo_unit

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"periph.io/x/conn/v3/i2c"
)

const (
	M5_UNIT_8SERVO_DEFAULT_ADDR         = 0x25
	M5_UNIT_8SERVO_MODE_REG             = 0x00
	M5_UNIT_8SERVO_OUTPUT_CTL_REG       = 0x10
	M5_UNIT_8SERVO_DIGITAL_INPUT_REG    = 0x20
	M5_UNIT_8SERVO_ANALOG_INPUT_8B_REG  = 0x30
	M5_UNIT_8SERVO_ANALOG_INPUT_12B_REG = 0x40
	M5_UNIT_8SERVO_SERVO_ANGLE_8B_REG   = 0x50
	M5_UNIT_8SERVO_SERVO_PULSE_16B_REG  = 0x60
	M5_UNIT_8SERVO_RGB_24B_REG          = 0x70
	M5_UNIT_8SERVO_PWM_8B_REG           = 0x90
	M5_UNIT_8SERVO_CURRENT_REG          = 0xA0
	JUMP_TO_BOOTLOADER_REG              = 0xFD
	FIRMWARE_VERSION_REG                = 0xFE
	I2C_ADDRESS_REG                     = 0xFF

	M5_UNIT_8SERVO_FW_VERSION_REG = 0xFE
	M5_UNIT_8SERVO_ADDRESS_REG    = 0xFF
)

type ExtIOMode int

const (
	DIGITAL_INPUT_MODE ExtIOMode = iota
	DIGITAL_OUTPUT_MODE
	ADC_INPUT_MODE
	SERVO_CTL_MODE
	RGB_LED_MODE
	PWM_MODE
)

type AnalogReadMode int

const (
	A8bit AnalogReadMode = iota
	A12bit
)

// Dev is an handle to an M%Stack 8Servo unit driver.
type Dev struct {
	c i2c.Dev
}

// I2CAddr is the default I2C address for the m5stack 8Servo unit.
const I2CAddr uint16 = M5_UNIT_8SERVO_DEFAULT_ADDR

// Opts holds the configuration options.
type Opts struct {
	I2cAddress uint16
}

// DefaultOpts are the recommended default options.
var DefaultOpts = Opts{
	I2cAddress: I2CAddr,
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

func (h *Dev) SetAllPinMode(mode ExtIOMode) error {
	data := make([]uint8, 8)
	for i := range data {
		data[i] = uint8(mode)
	}

	err := h.writeBytes(M5_UNIT_8SERVO_MODE_REG, data)
	if err != nil {
		return err
	}
	return nil
}

func (h *Dev) SetOnePinMode(pin uint8, mode ExtIOMode) error {
	if pin > 8 {
		return fmt.Errorf("wrong pin number")
	}
	return h.writeBytes(M5_UNIT_8SERVO_MODE_REG+int(pin), []uint8{uint8(mode)})
}

func (h *Dev) GetOnePinMode(pin uint8) (ExtIOMode, error) {
	if pin > 8 {
		return 0, fmt.Errorf("wrong pin number")
	}
	data, err := h.readBytes(M5_UNIT_8SERVO_MODE_REG+int(pin), 1)
	if err != nil {
		return 0, err
	}
	return ExtIOMode(data[0]), nil
}

func (h *Dev) SetDigitalOutput(pin uint8, state uint8) error {
	if pin > 7 {
		return fmt.Errorf("wrong pin number")
	}
	reg := M5_UNIT_8SERVO_OUTPUT_CTL_REG + pin
	return h.writeBytes(int(reg), []uint8{state})
}

func (h *Dev) SetLEDColor(pin uint8, color uint32) error {
	if pin > 7 {
		return fmt.Errorf("wrong pin number")
	}
	data := []uint8{
		uint8((color >> 16) & 0xff),
		uint8((color >> 8) & 0xff),
		uint8(color & 0xff),
	}
	reg := pin*3 + M5_UNIT_8SERVO_RGB_24B_REG
	return h.writeBytes(int(reg), data)
}

func (h *Dev) SetServoAngle(pin uint8, angle uint8) error {
	reg := pin + M5_UNIT_8SERVO_SERVO_ANGLE_8B_REG
	return h.writeBytes(int(reg), []uint8{angle})
}

func (h *Dev) SetPWM(pin uint8, angle uint8) error {
	reg := pin + M5_UNIT_8SERVO_PWM_8B_REG
	return h.writeBytes(int(reg), []uint8{angle})
}

func (h *Dev) SetServoPulse(pin uint8, pulse uint16) error {
	data := make([]uint8, 2)
	reg := pin*2 + M5_UNIT_8SERVO_SERVO_PULSE_16B_REG
	data[1] = uint8((pulse >> 8) & 0xff)
	data[0] = uint8(pulse & 0xff)
	return h.writeBytes(int(reg), data)
}

func (h *Dev) GetDigitalInput(pin uint8) (bool, error) {
	reg := pin + M5_UNIT_8SERVO_DIGITAL_INPUT_REG
	data, err := h.readBytes(int(reg), 1)
	if err != nil {
		return data[0] != 0, nil
	}
	return false, err
}

func (h *Dev) GetAnalogInput(pin uint8, bit AnalogReadMode) (uint16, error) {
	var err error
	if bit == A8bit {
		reg := pin + M5_UNIT_8SERVO_ANALOG_INPUT_8B_REG
		data, err := h.readBytes(int(reg), 1)
		if err == nil {
			return uint16(data[0]), nil
		}
	} else {
		reg := pin*2 + M5_UNIT_8SERVO_ANALOG_INPUT_12B_REG
		data, err := h.readBytes(int(reg), 2)
		if err == nil {
			return (uint16(data[1]) << 8) | uint16(data[0]), nil
		}
	}
	return 0, err
}

func (h *Dev) GetServoCurrent() (float32, error) {
	data := make([]uint8, 4)
	data, err := h.readBytes(M5_UNIT_8SERVO_CURRENT_REG, 4)
	if err != nil {
		return 0, err
	}
	var c float32
	binary.Read(bytes.NewReader(data), binary.LittleEndian, &c)
	return c, nil
}

func (h *Dev) SetI2CAddress(addr uint8) error {
	err := h.writeBytes(JUMP_TO_BOOTLOADER_REG, []uint8{addr})
	if err == nil {
		h.c.Addr = uint16(addr)
	}
	return err
}

func (h *Dev) GetFirmwareVersion() (uint8, error) {
	data, err := h.readBytes(FIRMWARE_VERSION_REG, 1)
	if err != nil {
		return 0, err
	}
	return data[0], err
}

func (h *Dev) GetI2CAddress() (uint8, error) {
	data, err := h.readBytes(I2C_ADDRESS_REG, 1)
	if err != nil {
		return 0, err
	}
	return data[0], err
}

func (h *Dev) JumpBootloader() error {
	value := uint8(1)
	return h.writeBytes(JUMP_TO_BOOTLOADER_REG, []uint8{value})
}
