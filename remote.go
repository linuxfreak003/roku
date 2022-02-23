// Package roku implements a library for
// interacting with Roku devices using the
// External Control Protocol (ECP)
// Example can be found at http://jonwillia.ms/roku/roku
package roku

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	ssdp "github.com/koron/go-ssdp"
)

// Remote is the base type for interacting with the Roku
// Device is populated when a `NewRemote` is called
// to create a new remote
type Remote struct {
	Addr   string
	Device *DeviceInfo
}

// RokuDevice is returned by FindRokuDevices
// For convenience finding Roku's on the network
type RokuDevice struct {
	Addr string
	Name string
}

// FindRokuDevices will do and SSDP search to
// find all Roku devices on the network
func FindRokuDevices() ([]*RokuDevice, error) {
	devices, err := ssdp.Search("roku:ecp", 3, "")
	if err != nil {
		log.Printf("could not find any devices on network: %v", err)
	}

	res := []*RokuDevice{}
	for _, d := range devices {
		u, err := url.Parse(d.Location)
		if err != nil {
			return nil, err
		}
		res = append(res, &RokuDevice{
			Addr: u.Host,
			Name: d.Server,
		})
	}

	return res, nil
}

// NewRemote sets up a remote to the given ip.
func NewRemote(addr string) (*Remote, error) {
	if addr == "" {
		return nil, fmt.Errorf("no address given: %v", addr)
	}

	r := &Remote{Addr: addr}

	info, err := r.DeviceInfo()
	if err != nil {
		return nil, err
	}

	r.Device = info

	return r, nil
}

// Refresh reloads the devince info on the remote
func (r *Remote) Refresh() error {
	info, err := r.DeviceInfo()
	if err != nil {
		return err
	}

	r.Device = info

	return nil
}

// ActiveApp returns the app currently running
func (r *Remote) ActiveApp() (*App, error) {
	b, err := r.query("active-app")
	if err != nil {
		return nil, err
	}

	type ActiveApp struct {
		XMLName xml.Name `xml:"active-app"`
		App     *App     `xml:"app"`
	}

	var app = &ActiveApp{}
	err = xml.Unmarshal(b, app)

	return app.App, err
}

// Launch will launch a given App
func (r *Remote) Launch(app *App) error {
	return r.launch(app.Id, nil)
}

// LaunchWithValues will launch a given App with exta arguments
func (r *Remote) LaunchWithValues(app *App, values url.Values) error {
	return r.launch(app.Id, values)
}

// Install will install the given app.
// Thie requres already knowing the App Id
// of the App you want to install
func (r *Remote) Install(app *App) error {
	return r.install(app.Id)
}

// Input sends a string of input
// (useful for things like filling a search box)
func (r *Remote) InputString(in string) error {
	return r.literalInput(in)
}

func (r *Remote) InputRune(rn rune) error {
	return r.literalInput(string(rn))
}

// Apps will get all installed apps fromt he device
func (r *Remote) Apps() ([]*App, error) {
	b, err := r.query("apps")
	if err != nil {
		return nil, err
	}

	type Apps struct {
		XMLName xml.Name `xml:"apps"`
		Apps    []*App   `xml:"app"`
	}

	apps := &Apps{}
	err = xml.Unmarshal(b, apps)

	return apps.Apps, err
}

// DeviceInfo shows the device info for connected device.
func (r *Remote) DeviceInfo() (*DeviceInfo, error) {
	b, err := r.query("device-info")
	if err != nil {
		return nil, err
	}

	info := &DeviceInfo{}
	err = xml.Unmarshal(b, info)

	return info, err
}

// PlayerStatus returns the media player state
func (r *Remote) PlayerStatus() (*PlayerStatus, error) {
	b, err := r.query("media-player")
	if err != nil {
		return nil, err
	}

	status := &PlayerStatus{}
	err = xml.Unmarshal(b, status)

	return status, err
}

// Equivalent of pressing Home button on the remote
func (r *Remote) Home() error { return r.keypress("Home") }

// Equivalent of pressing Rev button on the remote
func (r *Remote) Rev() error { return r.keypress("Rev") }

// Equivalent of pressing Fwd button on the remote
func (r *Remote) Fwd() error { return r.keypress("Fwd") }

// Equivalent of pressing Play button on the remote
func (r *Remote) Play() error { return r.keypress("Play") }

// Equivalent of pressing Select button on the remote
func (r *Remote) Select() error { return r.keypress("Select") }

// Equivalent of pressing Left button on the remote
func (r *Remote) Left() error { return r.keypress("Left") }

// Equivalent of pressing Right button on the remote
func (r *Remote) Right() error { return r.keypress("Right") }

// Equivalent of pressing Down button on the remote
func (r *Remote) Down() error { return r.keypress("Down") }

// Equivalent of pressing Up button on the remote
func (r *Remote) Up() error { return r.keypress("Up") }

// Equivalent of pressing Back button on the remote
func (r *Remote) Back() error { return r.keypress("Back") }

// Equivalent of pressing Instant Replay button on the remote
func (r *Remote) InstantReplay() error { return r.keypress("InstantReplay") }

// Equivalent of pressing Info button on the remote
func (r *Remote) Info() error { return r.keypress("Info") }

// Equivalent of pressing Backspace button on the remote
func (r *Remote) Backspace() error { return r.keypress("Backspace") }

// Equivalent of pressing Search button on the remote
func (r *Remote) Search() error { return r.keypress("Search") }

// Equivalent of pressing Enter button on the remote
func (r *Remote) Enter() error { return r.keypress("Enter") }

// Equivalent of pressing Volume Down button on the remote
// not available on all devices
func (r *Remote) VolumeDown() error { return r.keypress("VolumeDown") }

// Equivalent of pressing Mute button on the remote
// not available on all devices
func (r *Remote) VolumeMute() error { return r.keypress("VolumeMute") }

// Equivalent of pressing Volume Up button on the remote
// not available on all devices
func (r *Remote) VolumeUp() error { return r.keypress("VolumeUp") }

// Equivalent of pressing Power button on the remote
// not available on all devices
func (r *Remote) PowerOff() error { return r.keypress("PowerOff") }

// Equivalent of pressing Power button on the remote
// not available on all devices
func (r *Remote) PowerOn() error { return r.keypress("PowerOn") }

// Equivalent of pressing Channel Up button on the remote
// not available on all devices
func (r *Remote) ChannelUp() error { return r.keypress("ChannelUp") }

// Equivalent of pressing Channel Down button on the remote
// not available on all devices
func (r *Remote) ChannelDown() error { return r.keypress("ChannelDown") }

// helper method for hitting the `install` endpoint
func (r *Remote) install(appID string) error {
	URL := fmt.Sprintf("http://%s/install/%s", r.Addr, appID)

	resp, err := http.Post(URL, "", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// helper method for hitting the `launch` endpoint
func (r *Remote) launch(appID string, values url.Values) error {
	URL := fmt.Sprintf("http://%s/launch/%s?%s", r.Addr, appID, values.Encode())

	resp, err := http.Post(URL, "", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// helper method for hitting the `query` endpoint
func (r *Remote) query(cmd string) ([]byte, error) {
	URL := fmt.Sprintf("http://%s/query/%s", r.Addr, cmd)

	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// helper method for entering literal input
func (r *Remote) literalInput(input string) error {
	for _, c := range input {
		escaped := url.QueryEscape(string(c))

		err := r.keypress("Lit_" + escaped)
		if err != nil {
			return err
		}
	}

	return nil
}

// helper method for hitting the `keypress` endpoint
func (r *Remote) keypress(cmd string) error {
	URL := fmt.Sprintf("http://%s/keypress/%s", r.Addr, cmd)

	resp, err := http.Post(URL, "", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
