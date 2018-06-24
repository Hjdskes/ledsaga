package sladdfri

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/zubairhamed/canopus"
)

const (
	tradfriPort     = 5684
	preauthIdentity = "Client_identity"
)

// Client represent the connection to a Trådfri gateway. Any and all
// communication goes through this struct's methods.
type Client struct {
	// Hostname or IP address of the gateway for this client
	Gateway string

	// Gateway code at the bottom of your gateway; used for authentication
	Key string

	// Preshared key to use when communicating with the gateway
	psk string

	// CoAP connection with the gateway
	connection canopus.Connection
}

// A PSKRequest is sent to the gateway in an authentication request.
type PSKRequest struct {
	// Identifier sent in an authentication request.
	Ident string `json:"9090"`
}

// A PSKResponse is received from the gateway following an authentication request.
type PSKResponse struct {
	// Preshared key as returned by the gateway following an authentication request.
	PSK string `json:"9091"`
}

// Creates a new Client, connecting to the given gateway using the given authentication.
func NewClient(gateway, key string) *Client {
	return &Client{
		Gateway: gateway,
		Key:     key,
	}
}

// Connects the client to its gateway using the given identifier.
func (c *Client) Connect(ident string) error {
	address := fmt.Sprintf("%s:%d", c.Gateway, tradfriPort)
	log.Printf("Connecting to gateway: %s\n", address)

	if c.psk == "" {
		err := c.generatePSK(address, ident)
		if err != nil {
			return err
		}
	}

	var err error
	c.connection, err = canopus.DialDTLS(address, ident, c.psk)
	return err
}

func (c *Client) generatePSK(address, ident string) error {
	log.Printf("Requesting PSK...\n")

	conn, err := canopus.DialDTLS(address, preauthIdentity, c.Key)
	if err != nil {
		return err
	}

	payload := PSKRequest{Ident: ident}
	// TODO: cannot use c.postRequest because we need to process the status code of the reply.
	req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Post)
	req.SetRequestURI(uriGatewayIdent)
	data := canopus.NewJSONPayload(payload).GetBytes()
	req.SetPayload(data)

	resp, err := conn.Send(req)
	if err != nil {
		return err
	}

	if resp.GetMessage().GetCode() == canopus.CoapCodeCreated {
		var pskResp PSKResponse
		err := json.Unmarshal(resp.GetMessage().GetPayload().GetBytes(), &pskResp)
		if err == nil {
			c.psk = pskResp.PSK
			log.Printf("PSK: %s\n", c.psk)
		}
		return nil
	} else {
		return errors.New("Unable to get PSK")
	}
}

func (c *Client) observe(uri string) error {
	_, err := c.connection.ObserveResource(uri)
	return err
}

func (c *Client) request(uri string, messageMethod canopus.CoapCode, payload interface{}) ([]byte, error) {
	req := canopus.NewRequest(canopus.MessageConfirmable, messageMethod)

	switch messageMethod {
	case canopus.Put, canopus.Post:
		if payload != nil {
			data := canopus.NewJSONPayload(payload).GetBytes()
			req.SetPayload(data)
		}
	case canopus.Get, canopus.Delete:
		// Do nothing.
	default:
		error := fmt.Sprintf("Invalid CoAP message type: %d\n", messageMethod)
		return nil, errors.New(error)
	}

	req.SetRequestURI(uri)
	resp, err := c.connection.Send(req)
	if err != nil {
		log.Printf("<- error: %+v", err)
		return nil, err
	}
	rdata := resp.GetMessage().GetPayload().GetBytes()
	log.Printf("<- %s", string(rdata))
	return rdata, nil
}

func (c *Client) putRequest(uri string, payload interface{}) error {
	log.Printf("PUT %s payload %s", uri, payload)
	_, err := c.request(uri, canopus.Put, payload)
	return err
}

func (c *Client) postRequest(uri string, payload interface{}) error {
	log.Printf("POST %s", uri)
	_, err := c.request(uri, canopus.Post, payload)
	return err
}

func (c *Client) getRequest(uri string, out interface{}) error {
	log.Printf("GET %s", uri)
	data, err := c.request(uri, canopus.Get, nil)
	if err == nil {
		err = json.Unmarshal(data, out)
	}
	return err
}

func (c *Client) deleteRequest(uri string) error {
	log.Printf("DELETE %s", uri)
	_, err := c.request(uri, canopus.Delete, nil)
	return err
}

// Sets the NTP server used by the gateway.
func (c *Client) SetNTP(NTPServer string) error {
	payload := Gateway{
		NTPServer: NTPServer,
	}
	return c.putRequest(uriGatewayInfo, payload)
}

// Sets the gateway into commissioning mode for the given duration in
// seconds.
func (c *Client) SetCommissioningMode(seconds uint32) error {
	payload := Gateway{
		CommissioningMode: seconds,
	}
	return c.putRequest(uriGatewayInfo, payload)
}

// Reboots the gateway.
func (c *Client) Reboot() error {
	return c.postRequest(uriGatewayReboot, nil)
}

// Resets the gateway to factory defaults.
func (c *Client) FactoryReset() error {
	return c.postRequest(uriGatewayFactoryReset, nil)
}

// Gets the gateway information, see Gateway.
func (c *Client) GetGateway() (*Gateway, error) {
	var gatewayInfo Gateway
	err := c.getRequest(uriGatewayInfo, &gatewayInfo)
	if err != nil {
		return nil, err
	}
	return &gatewayInfo, nil
}

// Gets the given group's information, see Group.
func (c *Client) GetGroup(id uint32) (*Group, error) {
	uri := fmt.Sprintf("%s/%d", uriGroups, id)
	var desc Group
	err := c.getRequest(uri, &desc)
	if err != nil {
		return nil, err
	}
	return &desc, nil
}

// Gets the given mood's information, see Mood.
func (c *Client) GetMood(id uint32, parent *uint32) (*Mood, error) {
	if parent == nil {
		var err error
		parent, err = c.moodParent()
		if err != nil {
			return nil, err
		}
	}
	uri := fmt.Sprintf("%s/%d/%d", uriMoods, *parent, id)
	var desc Mood
	err := c.getRequest(uri, &desc)
	if err != nil {
		return nil, err
	}
	return &desc, nil
}

// Gets the given device's information, see Device.
func (c *Client) GetDevice(id uint32) (*Device, error) {
	uri := fmt.Sprintf("%s/%d", uriDevices, id)
	var desc Device
	err := c.getRequest(uri, &desc)
	if err != nil {
		return nil, err
	}
	return &desc, nil
}

// Adds a new group to the gateway, consisting of the given devices
// using the given name.
func (c *Client) AddGroup(ids []uint32, name string) error {
	log.Printf("ID: %v\n", ids)

	existingIds, err := c.ListDeviceIds()
	if err != nil {
		return err
	}

	// The gateway happily accepts groups consisting of non-existing device identifiers.
	// This loop scans for non-existing identifiers to prevent this.
	for _, id := range ids {
		var found bool
		for _, existingId := range existingIds {
			if id == existingId {
				found = true
				break
			}
		}
		if !found {
			return errors.New("All identifiers must exist")
		}
	}

	payload := AddGroupRequest{
		ID:   ids,
		Name: name,
	}
	return c.putRequest(uriGroupAdd, payload)
}

// Changes the group's, whose identifier matches the one from the
// given Group, settings to that of the given Group.
func (c *Client) SetGroup(g Group) error {
	uri := fmt.Sprintf("%s/%d", uriGroups, g.ID)
	return c.putRequest(uri, g)
}

// Removes the given group from the gateway.
func (c *Client) RemoveGroup(id uint32) error {
	// TODO: why does this not have to use /15004/remove?
	return c.deleteRequest(fmt.Sprintf("%s/%d", uriGroups, id))
}

// Adds a mood of the given name to the gateway.
func (c *Client) AddMood(name string) error {
	parent, err := c.moodParent()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s/%d", uriMoods, *parent)
	payload := AddMoodRequest{
		Name:     name,
		IsActive: 1,
	}
	return c.postRequest(uri, payload)
}

// Removes the given mood from the gateway.
func (c *Client) RemoveMood(id uint32) error {
	parent, err := c.moodParent()
	if err != nil {
		return err
	}
	uri := fmt.Sprintf("%s/%d/%d", uriMoods, *parent, id)
	return c.deleteRequest(uri)
}

// Changes the given device's settings to that of the given
// LightControl.
func (c *Client) SetDevice(id uint32, change LightControl) error {
	payload := DeviceSet{
		[]LightControl{change},
	}
	uri := fmt.Sprintf("%s/%d", uriDevices, id)
	return c.putRequest(uri, payload)
}

// Removes the given device from the gateway.
func (c *Client) RemoveDevice(id uint32) error {
	return c.deleteRequest(fmt.Sprintf("%s/%d", uriDevices, id))
}

// Lists the identifiers of all devices connected to the gateway.
func (c *Client) ListDeviceIds() (deviceIds []uint32, err error) {
	err = c.getRequest(uriDevices, &deviceIds)
	return deviceIds, err
}

// Lists the group settings of all devices connected to the gateway.
func (c *Client) ListGroups() ([]*Group, error) {
	log.Println("Requesting groups... ")
	var groupIds []uint32
	err := c.getRequest(uriGroups, &groupIds)
	if err != nil {
		return nil, err
	}

	log.Println("Enumerating...")
	groups := make([]*Group, len(groupIds))
	for i, group := range groupIds {
		var desc *Group
		desc, err = c.GetGroup(group)
		if err != nil {
			return nil, err
		}
		log.Printf("Found group: %+v\n", desc)
		groups[i] = desc

		// sleep for a while to avoid flood protection
		time.Sleep(100 * time.Millisecond)
	}

	return groups, nil
}

// Lists the mood settings of all the moods on the gateway.
func (c *Client) ListMoods() ([]*Mood, error) {
	log.Println("Requesting moods... ")
	parent, err := c.moodParent()
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("%s/%d", uriMoods, *parent)
	var moodIds []uint32
	err = c.getRequest(uri, &moodIds)
	if err != nil {
		return nil, err
	}

	log.Println("Enumerating...")
	moods := make([]*Mood, len(moodIds))
	for i, mood := range moodIds {
		var desc *Mood
		desc, err = c.GetMood(mood, parent)
		if err != nil {
			return nil, err
		}
		log.Printf("Found mood: %+v\n", desc)
		moods[i] = desc

		// sleep for a while to avoid flood protection
		time.Sleep(100 * time.Millisecond)
	}

	return moods, nil
}

// Lists the device settings of all the devices connected to the
// gateway.
func (c *Client) ListDevices() (devices []*Device, err error) {
	deviceIds, err := c.ListDeviceIds()
	if err != nil {
		return
	}

	log.Println("Enumerating...")
	for _, device := range deviceIds {
		var desc *Device
		desc, err = c.GetDevice(device)
		if err != nil {
			return
		}
		log.Printf("Found device: %s\n", desc)
		devices = append(devices, desc)

		// sleep for a while to avoid flood protection
		time.Sleep(100 * time.Millisecond)
	}

	return
}

func (c *Client) observerGateway(in chan canopus.ObserveMessage, out chan *Gateway) {
	for msg := range in {
		value := msg.GetValue()
		if value, ok := value.(canopus.MessagePayload); ok {
			gi := &Gateway{}
			err := json.Unmarshal(value.GetBytes(), gi)
			if err == nil {
				out <- gi
			}
		}
	}
}

func (c *Client) observerDevices(in chan canopus.ObserveMessage, out chan *Device) {
	for msg := range in {
		value := msg.GetValue()
		if value, ok := value.(canopus.MessagePayload); ok {
			dd := &Device{}
			err := json.Unmarshal(value.GetBytes(), dd)
			if err == nil {
				out <- dd
			}
		}
	}
}

// Observe the gateway for changes. These changes will be sent over
// the channel returned by GatewayEvents, which must be called first.
func (c *Client) ObserveGateway() error {
	return c.observe(uriGatewayInfo)
}

// Returns a channel over which any updates to any devices will be
// sent, see ObserveDevice.
func (c *Client) GatewayEvents() <-chan *Gateway {
	out := make(chan *Gateway)
	in := make(chan canopus.ObserveMessage)
	go c.connection.Observe(in)
	go c.observerGateway(in, out)
	return out
}

// Observe the given device, i.e., any changes through other channels
// (such as a remote) will be sent over the channel returned by
// DeviceEvents, which must be called first.
func (c *Client) ObserveDevice(deviceId uint32) error {
	uri := fmt.Sprintf("%s/%d", uriDevices, deviceId)
	return c.observe(uri)
}

// Returns a channel over which any updates to any devices will be
// sent, see ObserveDevice.
func (c *Client) DeviceEvents() <-chan *Device {
	out := make(chan *Device)
	in := make(chan canopus.ObserveMessage)
	go c.connection.Observe(in)
	go c.observerDevices(in, out)
	return out
}
