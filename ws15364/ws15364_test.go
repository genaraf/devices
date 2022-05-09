package ws15364

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
	m.SetMoterPwmFrequency(1000)

	for i := int(10); i <= 100; i += 10 {
		m.MotorMovement(M1, CW, i)
		m.MotorMovement(M2, CCW, i)
		time.Sleep(2 * time.Second)
		fmt.Printf("M1 speed: %d\n", i)
		fmt.Printf("M2 speed: %d\n", i)
	}
	m.MotorStop(M1)
	m.MotorStop(M2)
	m.Close()

}
