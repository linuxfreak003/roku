package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/eiannone/keyboard"
	"jonwillia.ms/roku"
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
  | Enter Input    / + <Text>       | Launch (By ID)  <Ctrl-L> + ID    |
  +---------------------------------+----------------------------------+
  (press q, Esc, or Ctrl-C to exit)`

var ErrNoDevices = fmt.Errorf(`no devices found`)

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
		return "", ErrNoDevices
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
			return "", fmt.Errorf("could not get selection: %w", err)
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
		if errors.Is(err, ErrNoDevices) {
			fmt.Println(`Could not find any roku devices.
Please Try again, or enter IP address manually with '-ip' flag`)

			return
		} else if err != nil {
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
	fmt.Printf("> ")

	for {
		var err error

		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatalf("error getting key event: %v", err)
		}
		// This simplifies the switch statement substantially
		if key == 0 {
			key = keyboard.Key(char)
		}

		switch key {
		case keyboard.KeyEsc, keyboard.KeyCtrlC, keyboard.KeyCtrlD, 'q':
			return
		case '/':
			fmt.Printf("Enter Text: ")

			var s string
			s, err = GetInput()

			if err != nil {
				break
			}

			err = r.InputString(s + "\n")

			fmt.Printf("\n> ")
		case keyboard.KeyCtrlL:
			var s string
			s, err = GetInput()

			if err != nil {
				break
			}

			err = r.Launch(&roku.App{
				Id: s,
			})
		case keyboard.KeyArrowLeft, 'h':
			err = r.Left()
		case keyboard.KeyArrowDown, 'j':
			err = r.Down()
		case keyboard.KeyArrowUp, 'k':
			err = r.Up()
		case keyboard.KeyArrowRight, 'l':
			err = r.Right()
		case 'i':
			err = r.Info()
		case 'H':
			err = r.Home()
		case keyboard.KeySpace:
			err = r.Play()
		case keyboard.KeyEnter:
			err = r.Select()
		case keyboard.KeyBackspace, keyboard.KeyBackspace2, 'B', 'u':
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

				fmt.Printf("> ")
			}
		case 's':
			var ps *roku.PlayerStatus

			ps, err = r.PlayerStatus()
			if err == nil {
				fmt.Printf("Media Player\nApp: [%s] %s\n", ps.Plugin.Id, ps.Plugin.Name)
				fmt.Printf("Error: %v State: %s\n", ps.Error, ps.State)
				fmt.Printf("Bandwidth: %s\n", ps.Plugin.Bandwidth)
				fmt.Printf("Position: %s\n", ps.Position)
				fmt.Printf("Live: %v\n", ps.IsLive)
				fmt.Printf("> ")
			}
		case 'p':
			// Refresh DeviceInfo current state
			_ = r.Refresh()

			switch r.Device.PowerMode {
			case "PowerOn":
				err = r.PowerOff()
				r.Device.PowerMode = "Ready"
			case "Ready":
				err = r.PowerOn()
				r.Device.PowerMode = "PowerOn"
			default:
				err = fmt.Errorf("Unrecognized power mode: %s", r.Device.PowerMode)
			}
		default:
			log.Printf("'%s' key does not match any command", string(key))
			fmt.Printf("> ")
		}

		LogIf(err)
	}
}

func GetInput() (string, error) {
	var s string

	for char, key, err := keyboard.GetKey(); key != keyboard.KeyEnter; char, key, err = keyboard.GetKey() {
		if err != nil {
			return "", fmt.Errorf("coult not get key: %w", err)
		}

		val := string(key)
		if key == 0 {
			val = string(char)
		}

		fmt.Printf("%s", val)
		s += val
	}

	return s, nil
}
