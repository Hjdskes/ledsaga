package sladdfri

import (
	"fmt"
	"time"
)

const (
	uriGroups      = "/15004"
	uriGroupAdd    = "/15004/add"
	uriGroupRemove = "/15004/remove"
)

// The Group struct holds all information relating to a group on the
// Trådfri gateway.
type Group struct {
	// Whether the light bulbs in this group are on or off. Read-write. Defined in IPSO 3311.
	Power uint8 `json:"5850"`

	// Dimmer value, i.e. how bright the bulbs in this group are. Valid values are in the range [0,254].
	// Read-write. This code is defined in IPSO 3311, but note that the values in Ikea's implementation
	// are not a percentage.
	Dim uint8 `json:"5851"`

	// The name of this group, as given by the user.
	Name string `json:"9001"`

	// The time at which this group was created.
	CreatedAt int64 `json:"9002"`

	// Numeric identifier of this group.
	ID uint32 `json:"9003"`

	//
	AccessoryLink struct {
		//
		LinkedItems struct {
			// Numeric identifier of the light bulbs in this group.
			DeviceIDs []uint32 `json:"9003"`
		} `json:"15002"`
	} `json:"9018,omitempty"`

	// The identifier of the currently active mood, if any.
	MoodID uint32 `json:"9039"`
}

func (g *Group) String() string {
	createdAt := time.Unix(g.CreatedAt, 0)
	s := fmt.Sprintf("ID: %d Name: %q Created: %s\n", g.ID, g.Name, createdAt.Format(time.RFC1123))
	s += fmt.Sprintf("Power: %d Dim: %d\n", g.Power, g.Dim)
	s += fmt.Sprintf("Linked devices: %v\n", g.AccessoryLink.LinkedItems.DeviceIDs)
	return s
}

type AddGroupRequest struct {
	// Numeric identifiers of the elements in the new group.
	ID []uint32 `json:"9003,omitempty"`

	// The name of the group, as given by the user.
	Name string `json:"9001"`
}
