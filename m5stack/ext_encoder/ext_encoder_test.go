package ext_encoder

import (
	"fmt"
	"log"
	"testing"

	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

func TestDev_ExtEncoder(t *testing.T) {
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
	fmt.Printf("ext-encoder version:%d\n", v)
}
