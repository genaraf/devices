package ws15364

import (
	"fmt"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/pca9685"
)

// I2CAddr is the default I2C address for the ws15364 components.
const I2CAddr uint16 = 0x40

const (
	_PWMA_CHANNEL = int(0)
	_AIN1_CHANNEL = int(1)
	_AIN2_CHANNEL = int(2)
	_PWMB_CHANNEL = int(5)
	_BIN1_CHANNEL = int(3)
	_BIN2_CHANNEL = int(4)
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
	CW  Direction = 0x01 // clockwise
	CCW Direction = 0x02 // countclockwise
)

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

// Dev is an handle to an DFR0592 Motors driver.
type Dev struct {
	c i2c.Dev
	d *pca9685.Dev
}

// New creates a new driver for CCS811 VOC sensor.
func New(bus i2c.Bus, opts *Opts) (*Dev, error) {
	if opts.I2cAddress < 0x01 || opts.I2cAddress > 0x70 {
		return nil, fmt.Errorf("invalid device address")
	}

	dev := &Dev{c: i2c.Dev{Bus: bus, Addr: opts.I2cAddress}}
	var err error
	dev.d, err = pca9685.NewI2C(bus, dev.c.Addr)
	if err != nil {
		return nil, err
	}

	if err := dev.SetMoterPwmFrequency(opts.PwmFreq); err != nil {
		return nil, err
	}

	// init channels
	if err := dev.d.SetAllPwm(0, 0); err != nil {
		return nil, err
	}
	dev.MotorStop(M1)
	dev.MotorStop(M2)
	return dev, nil
}

func (dev *Dev) Close() {
	if dev != nil {
		dev.MotorStop(M1)
		dev.MotorStop(M2)
	}
}

func (dev *Dev) SetMoterPwmFrequency(frequency int16) error {
	if frequency < 50 || frequency > 1526 {
		return fmt.Errorf("frequency out of range: 50-1526")
	}
	if err := dev.d.SetPwmFreq(physic.Frequency(frequency) * physic.Hertz); err != nil {
		return err
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
	if speed < 0 || speed > 100 {
		return fmt.Errorf("speed out of range: 0-100")
	}
	s := gpio.Duty((4095.0 / 100.0) * speed)
	if id == M1 {
		dev.d.SetPwm(_PWMA_CHANNEL, 0, s)
		//		dev.d.SetFullOn(_PWMA_CHANNEL)
		if direction == CW {
			dev.d.SetFullOff(_AIN1_CHANNEL)
			dev.d.SetFullOn(_AIN2_CHANNEL)
		} else {
			dev.d.SetFullOn(_AIN1_CHANNEL)
			dev.d.SetFullOff(_AIN2_CHANNEL)
		}
	} else if id == M2 {
		dev.d.SetPwm(_PWMB_CHANNEL, 0, s)
		//		dev.d.SetFullOn(_PWMB_CHANNEL)
		if direction == CW {
			dev.d.SetFullOff(_BIN1_CHANNEL)
			dev.d.SetFullOn(_BIN2_CHANNEL)
		} else {
			dev.d.SetFullOn(_BIN1_CHANNEL)
			dev.d.SetFullOff(_BIN2_CHANNEL)
		}
	} else {
		return fmt.Errorf("wrong motor id")
	}
	return nil
}

// Motor stop
// id: MotorId          Motor Id M1 or M2
func (dev *Dev) MotorStop(id MotorId) error {

	if id == M1 {
		dev.d.SetFullOff(_AIN1_CHANNEL)
		dev.d.SetFullOff(_AIN2_CHANNEL)
	} else if id == M2 {
		dev.d.SetFullOff(_BIN1_CHANNEL)
		dev.d.SetFullOff(_BIN2_CHANNEL)
	} else {
		return fmt.Errorf("wrong motor id")
	}
	return nil
}
