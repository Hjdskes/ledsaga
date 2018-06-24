package sladdfri

// The LightControl struct holds all settings to control a given
// Trådfri light bulb.
type LightControl struct {
	// The hex color string of this light bulb. Read-write. Defined in IPSO 3311, 3335.
	Color string `json:"5706"`

	// The hue of this light bulb, only for RGB bulbs.
	ColorHue int `json:"5707"`

	// The saturation of this light bulb, only for RGB bulbs.
	ColorSat int `json:"5708"`

	//
	ColorX int `json:"5709"`

	//
	ColorY int `json:"5710"`

	// Whether this light bulb is on or off. Read-write. Defined in IPSO 3311.
	Power uint8 `json:"5850"`

	// Dimmer value, i.e. how bright this bulb is. Valid values are in the range [0,254].
	// Read-write. This code is defined in IPSO 3311, but note that the values in Ikea's implementation
	// are not a percentage.
	Dim uint8 `json:"5851"`

	// Current color temperature in mired. Valid values are in the range [250,454],
	// which corresponds to [4000K,2200K].
	Mireds int `json:"5711"`

	// The duration of a transition in tenths of a second, only for RGB bulbs.
	TransitionDuration int `json:"5712"`

	// The total power in Wh that the light has used. Read-only. Defined in IPSO 3311.
	CumulativeActivePower float64 `json:"5805"`

	// The time in seconds that the light has been on. Writing a value of 0 resets the counter. Read-write.
	// Defined in IPSO 3311, 3342.
	OnTime uint32 `json:"5852"`

	// The power factor of the light. Read-only. Defined in IPSO 3311.
	PowerFactor float64 `json:"5820"`

	// If present, the type of sensor defined as the UCUM Unit Definition. Read-only.
	// Defined in IPSO 3311, 3335.
	SensorUnit string `json:"5701"`

	// Numeric identifier of this bulb.
	ID uint32 `json:"9003"`
}

// The DeviceSet struct is used in a request to change a Trådfri light
// bulb's settings.
type DeviceSet struct {
	LightControl []LightControl `json:"3311"`
}
