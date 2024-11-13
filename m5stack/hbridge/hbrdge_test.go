package hbridge

import (
	"fmt"
	"log"
	"testing"
	"time"

	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

func TestDev_MotorMovement(t *testing.T) {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Use i2creg I²C bus registry to find the first available I²C bus.
	b, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer b.Close()

	m, err := New(b, &DefaultOpts)
	if err != nil {
		log.Fatal(err)
	}
	v, err := m.GetFirmwareVersion()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("hbridge version:%d", v)
	m.SetDriverPWMFreq(1500)
	m.SetDriverDirection(HBRIDGE_FORWARD)
	for i := uint8(10); i <= 200; i += 10 {
		m.SetDriverSpeed8Bits(i)
		fmt.Printf("speed: %.d`\n", i)
		time.Sleep(1 * time.Second)
	}

	time.Sleep(2 * time.Second)
	m.SetDriverDirection(HBRIDGE_STOP)
	m.Close()

}
