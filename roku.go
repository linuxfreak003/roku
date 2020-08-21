// Package roku implements a library for
// interacting with Roku devices using the
// External Control Protocol (ECP)
// Example can be found at http://github.com/linuxfreak003/roku
package roku

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	ssdp "github.com/bcurren/go-ssdp"
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

// DeviceInfo is the information about the Roku
// The most useful fields are probably `PowerMode`
// (To know if device is On/Ready/etc) and maybe
// the Name and network information
type DeviceInfo struct {
	XMLName                     xml.Name `xml:"device-info"`
	Udn                         string   `xml:"udn"`
	SerialNumber                string   `xml:"serial-number"`
	DeviceId                    string   `xml:"device-id"`
	AdvertisingId               string   `xml:"advertising-id"`
	VendorName                  string   `xml:"vendor-name"`
	ModelName                   string   `xml:"model-name"`
	ModelNumber                 string   `xml:"model-number"`
	ModelRegion                 string   `xml:"model-region"`
	IsTv                        bool     `xml:"is-tv"`
	IsStick                     bool     `xml:"is-stick"`
	ScreenSize                  string   `xml:"screen-size"`
	PanelId                     string   `xml:"panel-id"`
	TunerType                   string   `xml:"tuner-type"`
	SuuportsEthernet            bool     `xml:"supports-ethernet"`
	WifiMac                     string   `xml:"wifi-mac"`
	WifiDriver                  string   `xml:"wifi-driver"`
	HasWifiExtender             bool     `xml:"has-wifi-extender"`
	HasWifi5GSupport            bool     `xml:"has-wifi-5G-support"`
	CanUseWifiExtender          bool     `xml:"can-use-wifi-extender"`
	EthernetMac                 string   `xml:"ethernet-mac"`
	NetworkType                 string   `xml:"network-type"`
	NetworkName                 string   `xml:"network-name"`
	FriendlyDeviceName          string   `xml:"friendly-device-name"`
	FriendlyModelName           string   `xml:"friendly-model-name"`
	DefaultDeviceName           string   `xml:"default-device-name"`
	UserDeviceName              string   `xml:"user-device-name"`
	UserDeviceLocation          string   `xml:"user-device-location"`
	BuildNumber                 string   `xml:"build-number"`
	SoftwareVersion             string   `xml:"software-version"`
	SoftwareBuild               string   `xml:"software-build"`
	SecureDevice                bool     `xml:"secure-device"`
	Language                    string   `xml:"language"`
	Country                     string   `xml:"country"`
	Locale                      string   `xml:"locale"`
	TimeZoneAuto                bool     `xml:"time-zone-auto"`
	TimeZone                    string   `xml:"time-zone"`
	TimeZoneName                string   `xml:"time-zone-name"`
	TimeZoneTz                  string   `xml:"time-zone-tz"`
	TimeZoneOffset              int      `xml:"time-zone-offset"`
	ClockFormat                 string   `xml:"clock-format"`
	Uptime                      int      `xml:"uptime"`
	PowerMode                   string   `xml:"power-mode"`
	SupportsSuspend             bool     `xml:"supports-suspend"`
	SupportsFindRemote          bool     `xml:"supports-find-remote"`
	FindRemoteIsPossible        bool     `xml:"find-remote-is-possible"`
	SupportsAudioGuide          bool     `xml:"supports-audio-guide"`
	SupportsRva                 bool     `xml:"supports-rva"`
	DeveloperEnabled            bool     `xml:"developer-enabled"`
	SearchEnabled               bool     `xml:"search-enabled"`
	SearchChannelsEnabled       bool     `xml:"search-channels-enabled"`
	VoiceSearchEnabled          bool     `xml:"voice-search-enabled"`
	NotificationsEnabled        bool     `xml:"notifications-enabled"`
	NotificationsFirstUse       bool     `xml:"notifications-first-use"`
	SupportsPrivateListening    bool     `xml:"supports-private-listening"`
	SupportsPrivateListeningDtv bool     `xml:"supports-private-listening-dtv"`
	SupportsWarmStandby         bool     `xml:"supports-warm-standby"`
	HeadphonesConnected         bool     `xml:"headphones-connected"`
	ExpertPqEnabled             string   `xml:"expert-pq-enabled"`
	SupportsEcsTextedit         bool     `xml:"supports-ecs-textedit"`
	SupportsEcsMicrophone       bool     `xml:"supports-ecs-microphone"`
	SupportsWakeOnWlan          bool     `xml:"supports-wake-on-wlan"`
	HasPlayOnRoku               bool     `xml:"has-play-on-roku"`
	HasMobileScreensaver        bool     `xml:"has-mobile-screensaver"`
	SupportUrl                  string   `xml:"support-url"`
	GrandcentralVersion         string   `xml:"grandcentral-version"`
	TrcVersion                  string   `xml:"trc-version"`
	TrcChannelVersion           string   `xml:"trc-channel-version"`
	DavinciVersion              string   `xml:"davinci-version"`
}

// App holds the app name and corresponding id
// The ID is needed to install/open and app
type App struct {
	Name string `xml:",chardata"`
	Id   string `xml:"id,attr"`
}

// PlayerStatus holds the media status
// of the app currently running
type PlayerStatus struct {
	XMLName  xml.Name `xml:"player"`
	Error    bool     `xml:"error,attr"`
	State    string   `xml:"state,attr"`
	Plugin   Plugin   `xml:"plugin"`
	Format   Format   `xml:"format"`
	Position string   `xml:"position"`
	IsLive   bool     `xml:"is_live"`
}
type Plugin struct {
	Bandwidth string `xml:"bandwidth,attr"`
	Id        string `xml:"id,attr"`
	Name      string `xml:"name,attr"`
}
type Format struct {
	Audio    string `xml:"audio,attr"`
	Captions string `xml:"captions,attr"`
	Drm      string `xml:"drm,attr"`
	Video    string `xml:"video,attr"`
}

// FindRokuDevices will do and SSDP search to
// find all Roku devices on the network
func FindRokuDevices() ([]*RokuDevice, error) {
	devices, err := ssdp.Search("roku:ecp", 3*time.Second)
	if err != nil {
		log.Printf("could not find any devices on network: %v", err)
	}
	res := []*RokuDevice{}
	for _, d := range devices {
		res = append(res, &RokuDevice{
			Addr: d.Location.Host,
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
		return nil, fmt.Errorf("could not connect to roku: %v", err)
	}

	r.Device = info

	return r, nil
}

// Refresh reloads the devince info on the remote
func (r *Remote) Refresh() error {
	info, err := r.DeviceInfo()
	if err != nil {
		return fmt.Errorf("could not connect to roku: %v", err)
	}

	r.Device = info

	return nil
}

// ActiveApp returns the app currently running
func (r *Remote) ActiveApp() (*App, error) {
	b, err := r.query("active-app")
	if err != nil {
		return nil, fmt.Errorf("could not query active-app: %v", err)
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
	return r.launch(app.Id)
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
	return r.literal_input(in)
}

func (r *Remote) InputRune(rn rune) error {
	return r.literal_input(string(rn))
}

// Apps will get all installed apps fromt he device
func (r *Remote) Apps() ([]*App, error) {
	b, err := r.query("apps")
	if err != nil {
		return nil, fmt.Errorf("could not query apps: %v", err)
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
		return nil, fmt.Errorf("could not query: %v", err)
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
func (r *Remote) install(appId string) error {
	URL := fmt.Sprintf("http://%s/install/%s", r.Addr, appId)
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, URL, nil)
	if err != nil {
		return fmt.Errorf("could not build HTTP request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	return nil
}

// helper method for hitting the `launch` endpoint
func (r *Remote) launch(appId string) error {
	URL := fmt.Sprintf("http://%s/launch/%s", r.Addr, appId)
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, URL, nil)
	if err != nil {
		return fmt.Errorf("could not build HTTP request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	return nil
}

// helper method for hitting the `query` endpoint
func (r *Remote) query(cmd string) ([]byte, error) {
	URL := fmt.Sprintf("http://%s/query/%s", r.Addr, cmd)
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, fmt.Errorf("could not build HTTP request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// helper method for entering literal input
func (r *Remote) literal_input(input string) error {
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
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, URL, nil)
	if err != nil {
		return fmt.Errorf("could not build HTTP request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	return nil
}
