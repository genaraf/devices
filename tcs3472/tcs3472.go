package tcs3472

import (
	"fmt"
	"time"

	"periph.io/x/conn/v3/i2c"
)

const (
	TCS3472_ADDRESS      = 0x29
	TCS34725_COMMAND_BIT = 0x80 /**< Command bit **/
	TCS3472_ENABLE       = 0x00
	TCS3472_ATIME        = 0x01
	TCS3472_WTIME        = 0x03
	TCS3472_CONFIG       = 0x0D
	TCS3472_CONTROL      = 0x0F
	TCS3472_ID           = 0x12
	TCS3472_STATUS       = 0x13
	TCS3472_CLEAR_LOW    = 0x14
	TCS3472_CLEAR_HIGH   = 0x15
	TCS3472_RED_LOW      = 0x16
	TCS3472_RED_HIGH     = 0x17
	TCS3472_GREEN_LOW    = 0x18
	TCS3472_GREEN_HIGH   = 0x19
	TCS3472_BLUE_LOW     = 0x1A
	TCS3472_BLUE_HIGH    = 0x1B

	TCS3472_POWER_ON  = 0x01
	TCS3472_POWER_OFF = 0x00
	TCS3472_AEN       = 0x02
)

/*
 * 60-Hz period: 16.67ms, 50-Hz period: 20ms
 * 100ms is evenly divisible by 50Hz periods and by 60Hz periods
 */
type IntegrationTime int

const (
	TCS34725_INTEGRATIONTIME_2_4MS IntegrationTime = 0xFF // 2.4ms - 1 cycle - Max Count: 1024
	TCS34725_INTEGRATIONTIME_24MS  IntegrationTime = 0xF6 // 24.0ms - 10 cycles - Max Count: 10240
	TCS34725_INTEGRATIONTIME_50MS  IntegrationTime = 0xEB // 50.4ms - 21 cycles - Max Count: 21504
	TCS34725_INTEGRATIONTIME_60MS  IntegrationTime = 0xE7 // 60.0ms - 25 cycles - Max Count: 25700
	TCS34725_INTEGRATIONTIME_101MS IntegrationTime = 0xD6 // 100.8ms - 42 cycles - Max Count: 43008
	TCS34725_INTEGRATIONTIME_120MS IntegrationTime = 0xCE // 120.0ms - 50 cycles - Max Count: 51200
	TCS34725_INTEGRATIONTIME_154MS IntegrationTime = 0xC0 // 153.6ms - 64 cycles - Max Count: 65535
	TCS34725_INTEGRATIONTIME_180MS IntegrationTime = 0xB5 // 180.0ms - 75 cycles - Max Count: 65535
	TCS34725_INTEGRATIONTIME_199MS IntegrationTime = 0xAD // 199.2ms - 83 cycles - Max Count: 65535
	TCS34725_INTEGRATIONTIME_240MS IntegrationTime = 0x9C // 240.0ms - 100 cycles - Max Count: 65535
	TCS34725_INTEGRATIONTIME_300MS IntegrationTime = 0x83 // 300.0ms - 125 cycles - Max Count: 65535
	TCS34725_INTEGRATIONTIME_360MS IntegrationTime = 0x6A // 360.0ms - 150 cycles - Max Count: 65535
	TCS34725_INTEGRATIONTIME_401MS IntegrationTime = 0x59 // 400.8ms - 167 cycles - Max Count: 65535
	TCS34725_INTEGRATIONTIME_420MS IntegrationTime = 0x51 // 420.0ms - 175 cycles - Max Count: 65535
	TCS34725_INTEGRATIONTIME_480MS IntegrationTime = 0x38 // 480.0ms - 200 cycles - Max Count: 65535
	TCS34725_INTEGRATIONTIME_499MS IntegrationTime = 0x30 // 499.2ms - 208 cycles - Max Count: 65535
	TCS34725_INTEGRATIONTIME_540MS IntegrationTime = 0x1F // 540.0ms - 225 cycles - Max Count: 65535
	TCS34725_INTEGRATIONTIME_600MS IntegrationTime = 0x06 // 600.0ms - 250 cycles - Max Count: 65535
	TCS34725_INTEGRATIONTIME_614MS IntegrationTime = 0x00 // 614.4ms - 256 cycles - Max Count: 65535
)

type TCS34725Gain byte

const (
	TCS34725Gain1X  TCS34725Gain = 0x00 // No gain
	TCS34725Gain4X  TCS34725Gain = 0x01 // 4x gain
	TCS34725Gain16X TCS34725Gain = 0x02 // 16x gain
	TCS34725Gain60X TCS34725Gain = 0x03 // 60x gain
)

type Color struct {
	Red   uint16
	Green uint16
	Blue  uint16
	Clear uint16
}

// Dev is an handle to an M%Stack 8Servo unit driver.
type Dev struct {
	c     i2c.Dev
	Gain  TCS34725Gain
	ITime IntegrationTime
}

// I2CAddr is the default I2C address for the m5stack 8Servo unit.
const I2CAddr uint16 = TCS3472_ADDRESS

// Opts holds the configuration options.
type Opts struct {
	I2cAddress uint16
	Gain       TCS34725Gain
	ITime      IntegrationTime
}

// DefaultOpts are the recommended default options.
var DefaultOpts = Opts{
	I2cAddress: I2CAddr,
	Gain:       TCS34725Gain1X,
	ITime:      TCS34725_INTEGRATIONTIME_154MS,
}

// New creates a new driver.
func New(bus i2c.Bus, opts *Opts) (*Dev, error) {
	if opts.I2cAddress < 0x01 || opts.I2cAddress > 0x70 {
		return nil, fmt.Errorf("invalid device address")
	}

	dev := &Dev{c: i2c.Dev{Bus: bus, Addr: opts.I2cAddress}}
	err := dev.SetIntegrationTime(opts.ITime)
	if err != nil {
		return nil, err
	}
	err = dev.SetGain(opts.Gain)
	if err != nil {
		return nil, err
	}
	return dev, nil
}

func (h *Dev) Close() {
	h.PowerOff()
}

func (h *Dev) readBytes(reg int, size int) ([]uint8, error) {
	r := make([]byte, size)
	err := h.c.Tx([]byte{byte(TCS34725_COMMAND_BIT | reg)}, r)
	return r, err
}

func (h *Dev) writeBytes(reg int, data []uint8) error {
	d := []byte{byte(TCS34725_COMMAND_BIT | reg)}
	d = append(d, data...)
	return h.c.Tx(d, nil)
}

func (h *Dev) SetIntegrationTime(time IntegrationTime) error {
	buf := []byte{byte(time)}
	err := h.writeBytes(TCS3472_ATIME, buf)
	if err != nil {
		return err
	}
	h.ITime = time
	return nil
}
func (h *Dev) SetGain(gain TCS34725Gain) error {
	buf := []byte{byte(gain)}
	err := h.writeBytes(TCS3472_ATIME, buf)
	if err != nil {
		return err
	}
	h.Gain = gain
	return nil
}

func (h *Dev) Status() (byte, error) {
	data, err := h.readBytes(TCS3472_STATUS, 1)
	if err != nil {
		return 0, err
	}

	return data[0], nil
}

func (h *Dev) PowerOn() error {
	buf := []byte{byte(TCS3472_POWER_ON)}
	err := h.writeBytes(TCS3472_ENABLE, buf)
	if err != nil {
		return err
	}
	time.Sleep(3 * time.Millisecond)

	buf[0] = byte(TCS3472_POWER_ON | TCS3472_AEN)
	err = h.writeBytes(TCS3472_ENABLE, buf)
	if err != nil {
		return err
	}

	time.Sleep(700 * time.Millisecond)

	return nil
}

func (h *Dev) PowerOff() error {
	buf := []byte{byte(TCS3472_POWER_OFF)}
	err := h.writeBytes(TCS3472_ENABLE, buf)
	if err != nil {
		return err
	}

	return nil
}

func (h *Dev) GetColor() (Color, error) {
	var c Color
	data, err := h.readBytes(TCS3472_CLEAR_LOW, 2)
	if err != nil {
		return c, err
	}
	c.Clear = uint16(data[1])<<8 + uint16(data[0])

	data, err = h.readBytes(TCS3472_RED_LOW, 2)
	if err != nil {
		return c, err
	}
	c.Red = uint16(data[1])<<8 + uint16(data[0])

	data, err = h.readBytes(TCS3472_GREEN_LOW, 2)
	if err != nil {
		return c, err
	}
	c.Green = uint16(data[1])<<8 + uint16(data[0])

	data, err = h.readBytes(TCS3472_BLUE_LOW, 2)
	if err != nil {
		return c, err
	}
	c.Blue = uint16(data[1])<<8 + uint16(data[0])
	time.Sleep(time.Duration((256-h.ITime)*12/5+1) * time.Millisecond)
	return c, nil
}

func (h *Dev) GetRGB() (Color, error) {
	c, err := h.GetColor()
	if err != nil {
		return c, err
	}

	sum := float32(c.Clear)

	// Avoid divide by zero errors ... if clear = 0 return black
	if sum == 0 {
		return Color{Red: 0, Green: 0, Blue: 0}, nil
	}

	c.Red = uint16(float32(c.Red) / sum * 255.0)
	c.Green = uint16(float32(c.Green) / sum * 255.0)
	c.Blue = uint16(float32(c.Blue) / sum * 255.0)
	return c, nil
}

func (h *Dev) GetId() (byte, error) {
	data, err := h.readBytes(TCS3472_ID, 1)
	if err != nil {
		return 0, err
	}
	return data[0], nil
}
