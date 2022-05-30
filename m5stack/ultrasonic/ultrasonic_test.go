package ultrasonic

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

	s, err := New(b, &DefaultOpts)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i <= 20; i++ {
		fmt.Printf("Distance: %.2f\n", s.GetDistance())
		time.Sleep(2 * time.Second)
	}

	s.Close()

}
