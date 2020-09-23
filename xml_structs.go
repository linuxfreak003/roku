package roku

import "encoding/xml"

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
