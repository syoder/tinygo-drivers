// Package captouch provides a simple driver for turning a GPIO input into a
// a capacitive touch input. The input needs some sort of high impedence drain
// so that it is possible to charge the pin and time how long it takes to drain
// enough so that the input goes from high to low.
//
package captouch

import (
	"machine"
	"time"
)

const (
	chargeTime            = 10 * time.Microsecond
	timeout               = 10000
	samples               = 5
	calibrationMultiplyer = 1.5
)

// Device is a capacitive touch input
type Device struct {
	pin       machine.Pin
	threshold int
}

// New returns a new capacitive touch input, given a GPIO input pin
func New(pin machine.Pin) Device {
	return Device{pin: pin}
}

// Configure calibrates this device by measuring the input's capacitance
func (d *Device) Configure() {
	d.threshold = int(float64(d.Read()) * calibrationMultiplyer)
}

// SetThreshold sets the threshold above which a capacitance reading should be
// considered a touch. This is normally set automatically by Configure.
func (d *Device) SetThreshold(threshold int) {
	d.threshold = threshold
}

// Read returns a measure of capacitance of the input
func (d *Device) Read() int {
	count := 0

	for i := 0; i < samples; i++ {
		// charge the pin
		d.pin.Configure(machine.PinConfig{machine.PinOutput})
		d.pin.High()
		time.Sleep(chargeTime)

		// time how long it takes to drain
		d.pin.Configure(machine.PinConfig{machine.PinInput})
		for d.pin.Get() {
			count++
			if count >= timeout {
				return timeout
			}
		}
	}
	return count
}

// Get returns true if the input is being touched
func (d *Device) Get() bool {
	return d.Read() > d.threshold
}
