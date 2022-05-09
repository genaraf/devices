package drf0592

import (
	"fmt"
	"log"
	"testing"
	"time"

	"periph.io/x/conn/v3/i2c/i2creg"
	host "periph.io/x/host/v3"
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

	dl := Detecte(b)
	if len(dl) == 0 {
		t.Errorf("device not found")
	}
	fmt.Printf("t: founded device on %+q\n", dl)

	m, err := New(b, &DefaultOpts)
	if err != nil {
		log.Fatal(err)
	}
	m.SetMoterPwmFrequency(3000)
	m.SetEncoderEnable(M1)
	m.SetEncoderEnable(M2)
	m.SetEncoderReductionRatio(M1, 49)
	m.SetEncoderReductionRatio(M2, 49)
	for i := float32(10.0); i <= 100.0; i += 10.0 {
		m.MotorMovement(M1, CW, i)
		m.MotorMovement(M2, CCW, i)
		time.Sleep(2 * time.Second)
		s, _ := m.GetEncoderSpeed(M1)
		fmt.Printf("M1 speed: %.2f, encoder speed:%d\n", i, s)
		s, _ = m.GetEncoderSpeed(M2)
		fmt.Printf("M2 speed: %.2f, encoder speed:%d\n", i, s)
	}
	m.MotorStop(M1)
	m.MotorStop(M2)
	m.Close()
}
