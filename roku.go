package roku

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	ssdp "github.com/bcurren/go-ssdp"
)

type Remote struct {
	Addr   string
	Device *DeviceInfo
}

type RokuDevice struct {
	Addr string
	Name string
}

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

type App struct {
	Name string `xml:",chardata"`
	Id   string `xml:"id,attr"`
}

// FindRokuDevices will serach the network and
// return the ip addresses and server name
// of any roku devices
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

// NewRemote sets up a remote to the given ip
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
	fmt.Printf("Connected to %s in %s\n", r.Device.UserDeviceName, r.Device.UserDeviceLocation)
	apps, err := r.ActiveApp()
	if err != nil {
		return nil, fmt.Errorf("could not get apps: %v", err)
	}
	fmt.Println(apps)
	// Check to make sure addr leads to a roku
	return r, nil
}

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

// Apps return all the available apps for the device
func (r *Remote) Apps() ([]*App, error) {
	b, err := r.query("apps")
	if err != nil {
		return nil, fmt.Errorf("could not query apps: %v", err)
	}

	type Apps struct {
		XMLName xml.Name `xml:"apps"`
		Apps    []*App   `xml:"app"`
	}
	var apps = &Apps{}

	xml.Unmarshal(b, apps)
	return apps.Apps, nil
}

// DeviceInfo shows the device info for
func (r *Remote) DeviceInfo() (*DeviceInfo, error) {
	b, err := r.query("device-info")
	if err != nil {
		return nil, fmt.Errorf("could not query: %v", err)
	}
	var info DeviceInfo
	xml.Unmarshal(b, &info)
	return &info, nil
}

// MediaPlayer
// func (r *Remote) MediaPlayer() (*MediaPlayer, error) {
// }

// Each Function here maps to a button press
func (r *Remote) Home() error          { return r.keypress("Home") }
func (r *Remote) Rev() error           { return r.keypress("Rev") }
func (r *Remote) Fwd() error           { return r.keypress("Fwd") }
func (r *Remote) Play() error          { return r.keypress("Play") }
func (r *Remote) Select() error        { return r.keypress("Select") }
func (r *Remote) Left() error          { return r.keypress("Left") }
func (r *Remote) Right() error         { return r.keypress("Right") }
func (r *Remote) Down() error          { return r.keypress("Down") }
func (r *Remote) Up() error            { return r.keypress("Up") }
func (r *Remote) Back() error          { return r.keypress("Back") }
func (r *Remote) InstantReplay() error { return r.keypress("InstantReplay") }
func (r *Remote) Info() error          { return r.keypress("Info") }
func (r *Remote) Backspace() error     { return r.keypress("Backspace") }
func (r *Remote) Search() error        { return r.keypress("Search") }
func (r *Remote) Enter() error         { return r.keypress("Enter") }

// Only Available on some Devices
func (r *Remote) VolumeDown() error  { return r.keypress("VolumeDown") }
func (r *Remote) VolumeMute() error  { return r.keypress("VolumeMute") }
func (r *Remote) VolumeUp() error    { return r.keypress("VolumeUp") }
func (r *Remote) PowerOff() error    { return r.keypress("PowerOff") }
func (r *Remote) ChannelUp() error   { return r.keypress("ChannelUp") }
func (r *Remote) ChannelDown() error { return r.keypress("ChannelDown") }

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
	return ioutil.ReadAll(resp.Body)
}

// keypress sends the actual keypress event to the Roku ECP API
func (r *Remote) keypress(cmd string) error {
	URL := fmt.Sprintf("http://%s/keypress/%s", r.Addr, cmd)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, URL, nil)
	if err != nil {
		return fmt.Errorf("could not build HTTP request: %v", err)
	}
	_, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("could not send HTTP request: %v", err)
	}
	return nil
}
