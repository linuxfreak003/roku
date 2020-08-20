package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/eiannone/keyboard"
	"github.com/linuxfreak003/roku"
)

var UsageMessage = ` +----------------------------------+----------------------------------+
  | Left           h or <Left>      | Rewind          r                |
  | Down           j or <Down>      | Fast-Fwd        f                |
  | Up             k or <Up>        | Replay          R                |
  | Right          l or <Right>     | Play/Pause      <Space>          |
  | Volume Up      + or <Ctrl-K>    | Ok/Enter        <Enter>          |
  | Volume Down    - or <Ctrl-J>    | Back            B, u, or <Backsp>|
  | Volume Mute    m                | Home            H                |
  | Power Off/On   p                | Info/Settings   i                |
  | List Apps      a                | Player Status   s                |
  +---------------------------------+----------------------------------+
  (press q, Esc, or Ctrl-C to exit)`

var NoDevicesError = fmt.Errorf(`Could not find any roku devices.
Please Try again, or enter IP address manually with '-ip' flag`)

func Usage() {
	fmt.Println(UsageMessage)
}

func LogIf(err error) {
	if err != nil {
		log.Printf("error: %v", err)
	}
}

func GetRokuAddress() (string, error) {
	fmt.Println("Searching for Roku devices...")

	// Get devices
	devices, err := roku.FindRokuDevices()
	if err != nil {
		return "", fmt.Errorf("Could not find roku devices: %v", err)
	}

	if len(devices) == 0 {
		return "", NoDevicesError
	}

	var index int

	// Have user select which device if more than 1
	if len(devices) > 1 {
		fmt.Println("Roku Devices:")

		for i, device := range devices {
			fmt.Printf("[%d] %s (%s)\n", i, device.Addr, device.Name)
		}

		fmt.Println("Select a Device:")

		// Using this method, a user won't actually be able
		// to select any options higher than '9'
		char, _, err := keyboard.GetSingleKey()
		if err != nil {
			return "", fmt.Errorf("Could not get selection: %v", err)
		}

		index := int(char - 48)
		if index < 0 || index >= len(devices) {
			return "", fmt.Errorf("invalid choice: %d", index)
		}
	}

	return devices[index].Addr, nil
}

func main() {
	var ip string
	var port int

	flag.StringVar(&ip, "ip", "", "IP address of roku device")
	flag.IntVar(&port, "port", 8060, "port to use for roku device")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", ip, port)

	var err error
	if ip == "" {
		addr, err = GetRokuAddress()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}

	// Create Remote
	r, err := roku.NewRemote(addr)
	if err != nil {
		log.Fatalf("could not create remote: %v", err)
	}

	fmt.Printf("Connected to %s in %s\n", r.Device.UserDeviceName, r.Device.UserDeviceLocation)
	fmt.Printf("  on %s network '%s'\n", r.Device.NetworkType, r.Device.NetworkName)
	fmt.Printf("Mode: %s\n", r.Device.PowerMode)

	active, err := r.ActiveApp()
	if err != nil {
		log.Printf("could not get apps: %v", err)
	}

	fmt.Printf("Active App: %v\n", active.Name)

	Usage()

	// Open Keyboard
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()

	CommandLoop(r)
	fmt.Printf("\nShutting down...\n")
}

func CommandLoop(r *roku.Remote) {
	for {
		var err error

		fmt.Printf("> ")

		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatalf("error getting key event: %v", err)
		}
		// This simplifies the switch statement substantially
		if key == 0 {
			key = keyboard.Key(char)
		}

		switch key {
		case keyboard.KeyEsc, keyboard.KeyCtrlC, 'q':
			return
		case keyboard.KeyArrowLeft, 'h':
			err = r.Left()
		case keyboard.KeyArrowDown, 'j':
			err = r.Down()
		case keyboard.KeyArrowUp, 'k':
			err = r.Up()
		case keyboard.KeyArrowRight, 'l':
			err = r.Right()
		case 'H':
			err = r.Home()
		case keyboard.KeySpace:
			err = r.Play()
		case keyboard.KeyEnter:
			err = r.Select()
		case keyboard.KeyBackspace, 'B', 'u':
			err = r.Back()
		case '+', keyboard.KeyCtrlK:
			err = r.VolumeUp()
		case '-', keyboard.KeyCtrlJ:
			err = r.VolumeDown()
		case 'm':
			err = r.VolumeMute()
		case 'R':
			err = r.InstantReplay()
		case 'r':
			err = r.Rev()
		case 'f':
			err = r.Fwd()
		case 'a':
			var apps []*roku.App

			apps, err = r.Apps()
			if err == nil {
				fmt.Printf("Installed apps:\n")

				for _, app := range apps {
					fmt.Printf("[%s]\t%s\n", app.Id, app.Name)
				}
			}
		case 's':
			var ps *roku.PlayerStatus

			ps, err = r.PlayerStatus()
			if err == nil {
				fmt.Printf("App: [%s] %s\n", ps.Plugin.Id, ps.Plugin.Name)
				fmt.Printf("Error: %v State: %s\n", ps.Error, ps.State)
				fmt.Printf("Bandwidth: %s\n", ps.Plugin.Bandwidth)
				fmt.Printf("Position: %s\n", ps.Position)
				fmt.Printf("Live: %v\n", ps.IsLive)
			}
		case 'p':
			if r.Device.PowerMode == "PowerOn" {
				err = r.PowerOff()
			} else {
				err = r.PowerOn()
			}
			// Refresh DeviceInfo
			_ = r.Refresh()
		default:
			log.Printf("'%s' key does not match any command", string(key))
		}

		LogIf(err)
	}
}
