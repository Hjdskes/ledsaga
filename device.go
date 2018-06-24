package sladdfri

import (
	"fmt"
	"time"
)

const (
	uriDevices = "/15001"
)

// The application type of the sensor or actuator. Defined in IPSO 3311, 3335, 3342.
type DeviceType uint8

const (
	// The Trådfri wireless remote.
	Remote DeviceType = 0

	// The Trådfri wireless dimmer.
	Dimmer DeviceType = 1

	// Any Trådfri light bulb.
	Light DeviceType = 2

	// The Trådfri motion sensor.
	Sensor DeviceType = 4
)

func (t DeviceType) String() string {
	switch t {
	case Remote:
		return "Remote"
	case Dimmer:
		return "Dimmer"
	case Light:
		return "Light"
	case Sensor:
		return "Sensor"
	default:
		return "Unknown"
	}
}

// Available power sources. Defined in IPSO 3.
type PowerSource uint8

const (
	DC          PowerSource = 0
	InternalBat PowerSource = 1
	ExternalBat PowerSource = 2
	Battery     PowerSource = 3
	PoE         PowerSource = 4
	USB         PowerSource = 5
	AC          PowerSource = 6
	Solar       PowerSource = 7
)

func (p PowerSource) String() string {
	switch p {
	case DC:
		return "DC"
	case InternalBat:
		return "Internal Battery"
	case ExternalBat:
		return "External Battery"
	case Battery:
		return "Battery"
	case PoE:
		return "Power over Ethernet"
	case USB:
		return "USB"
	case AC:
		return "AC"
	case Solar:
		return "Solar"
	default:
		return "Unknown"
	}
}

// The Device struct holds all information related to a Trådfri
// device.
type Device struct {
	// Device related information according to IPSO 3.
	Device struct {
		// Read-only. Defined in IPSO 3.
		Manufacturer string `json:"0"`
		// Read-only. Defined in IPSO 3.
		ModelNumber string `json:"1"`
		// Read-only. Defined in IPSO 3.
		Serial string `json:"2"`
		// Read-only. Defined in IPSO 3.
		FirmwareVersion string `json:"3"`
		// See PowerSource. Read-only. Defined in IPSO 3.
		AvailablePowerSource PowerSource `json:"6"`
		// Battery level as a percentage. Read-only. Defined in IPSO 3.
		BatteryLevel uint8 `json:"9"`
	} `json:"3"`

	// A list of light source controls, according to IPSO 3311. See LightControl.
	LightControl []LightControl `json:"3311"`

	// The application type of this device, see DeviceType. Read-write. Defined in IPSO 3311, 3335, 3342.
	Type DeviceType `json:"5750"`

	// The name of this device, as given by the user.
	Name string `json:"9001"`

	// The time at which this bulb was paired with the gateway.
	CreatedAt int64 `json:"9002"`

	// Numeric identifier of this device.
	ID uint32 `json:"9003"`

	// Whether this device is reachable or not.
	Reachable uint8 `json:"9019"`

	//
	LastSeen int64 `json:"9020"`

	//
	OtaUpdateState int `json:"9054"`
}

func (d *Device) String() string {
	s := fmt.Sprintf("ID: %d Name: %q\nType: %s Model: %q\n", d.ID, d.Name, d.Type, d.Device.ModelNumber)
	s += fmt.Sprintf("Firmware: %s Manufacturer: %q\n", d.Device.FirmwareVersion, d.Device.Manufacturer)
	s += fmt.Sprintf("Power: %s\n", d.Device.AvailablePowerSource)

	createdAt := time.Unix(d.CreatedAt, 0)
	lastSeen := time.Unix(d.LastSeen, 0)
	s += fmt.Sprintf("Created at: %s ", createdAt.Format(time.RFC1123))
	s += fmt.Sprintf("Last seen: %s\n", lastSeen.Format(time.RFC1123))

	if d.Type == Light {
		for count, entry := range d.LightControl {
			power := "off"
			if entry.Power == 1 {
				power = "on"
			}
			pc := DimToPercentage(entry.Dim)
			s += fmt.Sprintf("Light Control Set %d, Power: %s, Dim: %d%%\n", count, power, pc)
			s += "Color: "
			s += fmt.Sprintf("%dK ", MiredToKelvin(entry.Mireds))
			s += fmt.Sprintf("#%s ", entry.Color)
			s += fmt.Sprintf("X:%d/Y:%d ", entry.ColorX, entry.ColorY)
			s += fmt.Sprintf("Hue: %d Sat: %d ", entry.ColorHue, entry.ColorSat)
			s += "\n"
		}
	} else if d.Type == Remote || d.Type == Dimmer {
		s += fmt.Sprintf("Level: %v%%\n", d.Device.BatteryLevel)
	}

	return s
}
