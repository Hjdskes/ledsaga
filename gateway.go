package sladdfri

import (
	"fmt"
)

const (
	uriGatewayIdent        = "/15011/9063"
	uriGatewayInfo         = "/15011/15012"
	uriGatewayReboot       = "/15011/9030"
	uriGatewayFactoryReset = "/15011/9031"
)

// The Gateway struct holds all information related to the Trådfri
// gateway.
type Gateway struct {
	// The identifier of this gateway
	ID string `json:"9081"`

	// The NTP server the gateway uses
	NTPServer string `json:"9023"`

	// The firmware version of the gateway
	FirmwareVersion string `json:"9029"`

	// The current time as a Unix timestamp
	CurrentTimestamp int64 `json:"9059"`

	// The current time of the gateway in the format YYYY-MM-DDTHH:MM:SS.MMM
	CurrentTimestampUTC string `json:"9060"`

	// The amount of seconds in which this gateway accepts pairing requests from
	// new devices. A value of 0 means this gateway is not in commissioning mode.
	CommissioningMode uint32 `json:"9061"`

	// URL pointing to the release notes of the latest (?) update.
	ReleaseNotesURL string `json:"9056"`

	// The name of the gateway.
	Name string `json:"9035"`

	// All of the following fields have been reverse-engineered through the Android APK file.
	// Their naming and type matches the Java source code, but their function is unknown. It is
	// also likely that we may be able to use more precise types (e.g. uint8) for many of these.
	TimeSource              int    `json:"9071"`
	OtaUpdateState          int    `json:"9054"`
	UpdateProgress          int    `json:"9055"`
	UpdatePriority          int    `json:"9066"`
	UpdateAcceptedTimestamp int    `json:"9069"`
	ForceOtaUpdateCheck     string `json:"9032"`
	DstTimeOffset           int    `json:"9080"`
	DstStartMonth           int    `json:"9072"`
	DstStartDay             int    `json:"9073"`
	DstStartHour            int    `json:"9074"`
	DstStartMinute          int    `json:"9075"`
	DstEndMonth             int    `json:"9076"`
	DstEndDay               int    `json:"9077"`
	DstEndHour              int    `json:"9078"`
	DstEndMinute            int    `json:"9079"`
	GoogleHomePairStatus    int    `json:"9105"`
	AlexaPairStatus         int    `json:"9093"`
	CertificateProvisioned  int    `json:"9092"`
}

func (g *Gateway) String() string {
	return fmt.Sprintf("ID: %s\n"+
		"NTP server: %s\n"+
		"Firmware version: %s\n"+
		"Current time: %s\n"+
		"Commissioning: %d seconds\n", g.ID, g.NTPServer, g.FirmwareVersion, g.CurrentTimestampUTC, g.CommissioningMode)
}
