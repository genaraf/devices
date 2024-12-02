package tcs3472

import (
	"log"
	"testing"

	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

func TestDev_tcs3472(t *testing.T) {
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
	id, err := m.GetId()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID:%d\n", id)
	err = m.PowerOn()
	if err != nil {
		log.Fatal(err)
	}
	c, err := m.GetColor()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Clear:%d, Red:%d, Green:%d, Blue:%d\n", c.Clear, c.Red, c.Green, c.Blue)
	c, err = m.GetRGB()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Red:%d, Green:%d, Blue:%d\n", c.Red, c.Green, c.Blue)
	err = m.PowerOff()
	if err != nil {
		log.Fatal(err)
	}
}
