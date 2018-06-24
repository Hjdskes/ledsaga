package sladdfri

import (
	"fmt"
	"time"
)

const (
	uriMoods = "/15005"
)

// The Mood struct holds all information related to a mood on the
// Trådfri gateway.
type Mood struct {
	// Numeric identifier of this mood.
	ID uint32 `json:"9003"`

	// The time at which this mood was created.
	CreatedAt int64 `json:"9002"`

	// The name of this mood, as given by the user.
	Name string `json:"9001"`

	// Whether this mood was predefined by Ikea, or created by the user.
	IsPredefined uint8 `json:"9068"`

	//
	Index int32 `json:"9057"`

	//
	IsActive uint8 `json:"9058"`

	//
	LightControls []LightControl `json:"15013"`

	//
	UseCurrentLightSettings uint8 `json:"9070"`
}

func (m *Mood) String() string {
	createdAt := time.Unix(m.CreatedAt, 0)
	s := fmt.Sprintf("ID: %d Name: %q Created: %s\n", m.ID, m.Name, createdAt.Format(time.RFC1123))
	s += fmt.Sprintf("Predefined: %d Index: %d Active: %d\n", m.IsPredefined, m.Index, m.IsActive)

	d := "[ "
	for _, control := range m.LightControls {
		d += fmt.Sprintf("%d ", control.ID)
	}
	d += "]"
	s += fmt.Sprintf("Devices: %s Using current light settings: %d\n", d, m.UseCurrentLightSettings)
	return s
}

// The data sent to the gateway in a request to add a new mood.
type AddMoodRequest struct {
	// The name of the new mood.
	Name string `json:"9001"`

	//
	IsActive uint8 `json:"9058"`
}

func (c *Client) moodParent() (*uint32, error) {
	parent := make([]uint32, 2)
	err := c.getRequest(uriMoods, &parent)
	if err != nil {
		return nil, err
	}
	return &parent[0], nil
}
