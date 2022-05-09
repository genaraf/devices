package ws15364

import (
	"fmt"

	"periph.io/x/conn/v3/i2c"
)

// I2CAddr is the default I2C address for the ws15364 components.
const I2CAddr uint16 = 0x40

const (
	_PWMA_CHANNEL = 0
	_AIN1_CHANNEL = 1
	_AIN2_CHANNEL = 2
	_PWMB_CHANNEL = 5
	_BIN1_CHANNEL = 3
	_BIN2_CHANNEL = 4
)

// Enum motor ID
type MotorId byte

const (
	M1 MotorId = 0x01
	M2 MotorId = 0x02
)

// Orientation
type Direction byte

const (
	CW   Direction = 0x01 // clockwise
	CCW  Direction = 0x02 // countclockwise
	STOP Direction = 0x05 // stop
)

// Opts holds the configuration options.
type Opts struct {
	I2cAddress uint16
}

// DefaultOpts are the recommended default options.
var DefaultOpts = Opts{
	I2cAddress: I2CAddr,
}

// Dev is an handle to an DFR0592 Motors driver.
type Dev struct {
	c i2c.Dev
}

// New creates a new driver for CCS811 VOC sensor.
func New(bus i2c.Bus, opts *Opts) (*Dev, error) {
	if opts.I2cAddress < 0x01 || opts.I2cAddress > 0x70 {
		return nil, fmt.Errorf("invalid device address")
	}

	dev := &Dev{c: i2c.Dev{Bus: bus, Addr: opts.I2cAddress}}

	dev.MotorStop(M1)
	dev.MotorStop(M2)
	return dev, nil
}

func (drv *Dev) Close() {
	if drv != nil {
		drv.MotorStop(M1)
		drv.MotorStop(M2)
	}
}

func (dev *Dev) SetMoterPwmFrequency(frequency int) error {
	if frequency < 100 || frequency > 12750 {
		return fmt.Errorf("frequency out of range: 100-12750")
	}
	return nil
}

// Motor movement
// id: MotorId          Motor Id M1 or M2
// direction: Direction Motor orientation, CW (clockwise) or CCW (counterclockwise)
// speed: float         Motor pwm duty cycle, in range 0 to 100, otherwise no effective
func (dev *Dev) MotorMovement(id MotorId, direction Direction, speed float32) error {
	if direction != CW && direction != CCW {
		return fmt.Errorf("wrond direction parameter")
	}
	if speed < 0.0 || speed > 100.0 {
		return fmt.Errorf("speed out of range: 0.0-100.0")
	}
	return nil
}

// Motor stop
// id: MotorId          Motor Id M1 or M2
func (dev *Dev) MotorStop(id MotorId) error {
	return nil
}
