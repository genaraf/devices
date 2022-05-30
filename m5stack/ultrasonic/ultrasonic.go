package ultrasonic

import (
	"fmt"
	"time"

	"periph.io/x/conn/v3/i2c"
)

// I2CAddr is the default I2C address for the m5stack ultrasnic.
const I2CAddr uint16 = 0x57

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

	return &Dev{c: i2c.Dev{Bus: bus, Addr: opts.I2cAddress}}, nil
}

func (dev *Dev) Close() {

}

func (dev *Dev) GetDistance() float64 {
	b := make([]byte, 1)
	b[0] = 1
	_, err := dev.c.Write(b)
	if err != nil {
		return -1
	}
	r := make([]byte, 3)
	time.Sleep(20 * time.Millisecond)
	err = dev.c.Tx(nil, r)
	if err != nil {
		return -1
	}
	d := float64(uint32(r[0])<<16+uint32(r[1])<<8+uint32(r[2])) / 1000
	if d > 4500.0 {
		return 4500.0
	}

	return d
}
