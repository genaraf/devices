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
	m.SetMoterPwmFrequency(1500)

	for i := float32(10); i <= 100; i += 10 {
		m.MotorMovement(M1, CW, i)
		m.MotorMovement(M2, CCW, i)
		fmt.Printf("M1 speed: %.2f\n", i)
		fmt.Printf("M2 speed: %.2f\n", i)
		time.Sleep(2 * time.Second)
	}

	time.Sleep(2 * time.Second)
	m.MotorStop(M1)
	m.MotorStop(M2)
	m.Close()

}
