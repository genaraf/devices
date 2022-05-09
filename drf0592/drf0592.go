package drf0592

import (
	"fmt"
	"time"

	"periph.io/x/conn/v3/i2c"
)

// I2CAddr is the default I2C address for the drf0592 components.
const I2CAddr uint16 = 0x10

// Register number
type Register byte

const (
	_REG_SLAVE_ADDR               Register = 0x00
	_REG_PID                               = 0x01
	_REG_PVD                               = 0x02
	_REG_CTRL_MODE                         = 0x03
	_REG_ENCODER1_EN                       = 0x04
	_REG_ENCODER1_SPPED                    = 0x05
	_REG_ENCODER1_REDUCTION_RATIO          = 0x07
	_REG_ENCODER2_EN                       = 0x09
	_REG_ENCODER2_SPEED                    = 0x0a
	_REG_ENCODER2_REDUCTION_RATIO          = 0x0c
	_REG_MOTOR_PWM                         = 0x0e
	_REG_MOTOR1_ORIENTATION                = 0x0f
	_REG_MOTOR1_SPEED                      = 0x10
	_REG_MOTOR2_ORIENTATION                = 0x12
	_REG_MOTOR2_SPEED                      = 0x13

	_REG_DEF_PID = 0xdf
	_REG_DEF_VID = 0x10
)

// Enum motor ID
type MotorId byte

const (
	M1 MotorId = 0x01
	M2 MotorId = 0x02
)

// Board status
type Status byte

const (
	STA_OK                      Status = 0x00
	STA_ERR                     Status = 0x01
	STA_ERR_DEVICE_NOT_DETECTED Status = 0x02
	STA_ERR_SOFT_VERSION        Status = 0x03
	STA_ERR_PARAMETER           Status = 0x04
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

func checkBoard(bus i2c.Bus, addr uint16) bool {
	i2cbus := i2c.Dev{Bus: bus, Addr: addr}
	r := make([]byte, 1)
	err := i2cbus.Tx([]byte{_REG_PID}, r)
	if err != nil {
		return false
	}
	pid := r[0]
	err = i2cbus.Tx([]byte{_REG_PVD}, r)
	if err != nil {
		return false
	}
	vid := r[0]
	if pid != _REG_DEF_PID || vid != _REG_DEF_VID {
		return false
	}
	return true
}

func Detecte(bus i2c.Bus) []byte {
	var addrList []byte
	for addr := uint16(1); addr < 128; addr++ {
		if checkBoard(bus, addr) {
			addrList = append(addrList, byte(addr))
		}
	}
	return addrList
}

// New creates a new driver for CCS811 VOC sensor.
func New(bus i2c.Bus, opts *Opts) (*Dev, error) {
	if opts.I2cAddress < 0x01 || opts.I2cAddress > 0x70 {
		return nil, fmt.Errorf("invalid device address")
	}

	if !checkBoard(bus, opts.I2cAddress) {
		return nil, fmt.Errorf("device not detected")
	}

	dev := &Dev{c: i2c.Dev{Bus: bus, Addr: opts.I2cAddress}}

	// set DC motor mode
	err := dev.c.Tx([]byte{_REG_CTRL_MODE, 0}, nil)
	if err != nil {
		return nil, fmt.Errorf("error set DC motor mode")
	}
	dev.MotorStop(M1)
	dev.MotorStop(M2)
	dev.SetEncoderDisable(M1)
	dev.SetEncoderDisable(M2)
	return dev, nil
}

func (drv *Dev) Close() {
	if drv != nil {
		drv.MotorStop(M1)
		drv.MotorStop(M2)
	}
}

func (dev *Dev) SetEncoderEnable(id MotorId) error {
	err := dev.c.Tx([]byte{byte(_REG_ENCODER1_EN + 5*(id-1)), 0x01}, nil)
	if err != nil {
		return fmt.Errorf("error SetEncoderEnable: %v", err)
	}
	return nil
}

func (dev *Dev) SetEncoderDisable(id MotorId) error {
	err := dev.c.Tx([]byte{byte(_REG_ENCODER1_EN + 5*(id-1)), 0x0}, nil)
	if err != nil {
		return fmt.Errorf("error SetEncoderDisable: %v", err)
	}
	return nil
}

func (dev *Dev) SetEncoderReductionRatio(id MotorId, reductionRatio uint16) error {
	if reductionRatio < 1 || reductionRatio > 2000 {
		return fmt.Errorf("reductionRatio out of range: 1-2000")
	}
	err := dev.c.Tx([]byte{byte(_REG_ENCODER1_REDUCTION_RATIO + 5*(id-1)), byte(reductionRatio >> 8), byte(reductionRatio & 0xFF)}, nil)
	if err != nil {
		return fmt.Errorf("error SetEncoderReductionRatio: %v", err)
	}
	return nil
}

func (dev *Dev) GetEncoderSpeed(id MotorId) (int32, error) {
	r := make([]byte, 2)
	err := dev.c.Tx([]byte{byte(_REG_ENCODER1_SPPED + 5*(id-1))}, r)
	if err != nil {
		return 0, fmt.Errorf("error GetEncoderSpeed: %v", err)
	}
	s := int32(r[0]*0xFF + r[1])
	if s&0x8000 > 0 {
		return int32(-(0x10000 - uint32(s))), nil
	}
	return s, nil
}

func (dev *Dev) SetMoterPwmFrequency(frequency int) error {
	if frequency < 100 || frequency > 12750 {
		return fmt.Errorf("frequency out of range: 100-12750")
	}
	err := dev.c.Tx([]byte{byte(_REG_MOTOR_PWM), byte(frequency / 50)}, nil)
	if err != nil {
		return fmt.Errorf("error SetMoterPwmFrequency: %v", err)
	}
	time.Sleep(100 * time.Millisecond)
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
	reg := byte(_REG_MOTOR1_ORIENTATION + (id-1)*3)
	err := dev.c.Tx([]byte{byte(reg), byte(direction)}, nil)
	if err != nil {
		return fmt.Errorf("error set orientation: %v", err)
	}
	err = dev.c.Tx([]byte{byte(reg + 1), byte(speed), byte(uint16(speed*10.0) % 10)}, nil)
	if err != nil {
		return fmt.Errorf("error set speed: %v", err)
	}
	return nil
}

// Motor stop
// id: MotorId          Motor Id M1 or M2
func (dev *Dev) MotorStop(id MotorId) error {
	err := dev.c.Tx([]byte{byte(_REG_MOTOR1_ORIENTATION + 3*(id-1)), byte(STOP)}, nil)
	if err != nil {
		return fmt.Errorf("error MotorStop: %v", err)
	}
	return nil
}

//  Set board controler address, reboot module to make it effective
//  param address: byte    Address to set, range in 1 to 127
func (dev *Dev) SetAddr(addr byte) error {
	if addr < 1 || addr > 127 {
		return fmt.Errorf("addres out of range (1..127)")
	}
	err := dev.c.Tx([]byte{byte(_REG_SLAVE_ADDR), byte(addr)}, nil)
	if err != nil {
		return fmt.Errorf("error SetAddr: %v", err)
	}
	return nil
}
