package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/eiannone/keyboard"
	"github.com/linuxfreak003/roku"
)

var UsageMessage = ` +----------------------------------+-------------------------------+
  | Back           B, u, or <Backsp>| Replay          R             |
  | Home           H                | Info/Settings   i             |
  | Left           h or <Left>      | Rewind          r             |
  | Down           j or <Down>      | Fast-Fwd        f             |
  | Up             k or <Up>        | Play/Pause      <Space>       |
  | Right          l or <Right>     | Enter Text      /             |
  | Ok/Enter       <Enter>          | Volume Up       + or <Ctrl-K> |
  | Volume Mute    m                | Volume Down     - or <Ctrl-J> |
  | Power Off/On   p                |                               |
  +---------------------------------+-------------------------------+
  (press q, Esc, or Ctrl-C to exit)`

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

	devices, err := roku.FindRokuDevices()
	if err != nil {
		return "", fmt.Errorf("could not find roku devices: %v", err)
	}
	fmt.Println("Devices:")
	for i, device := range devices {
		fmt.Printf("[%d] %s (%s)\n", i, device.Addr, device.Name)
	}
	char, _, err := keyboard.GetSingleKey()
	index := int(char - 48)
	if index < 0 || index >= len(devices) {
		return "", fmt.Errorf("invalid choice: %d", index)
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

	Usage()

	// Open Keyborad
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()

	// Command loop
CommandLoop:
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
		case keyboard.KeyEsc, keyboard.KeyCtrlC, 'q':
			break CommandLoop
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
		case 'p':
			if r.Device.PowerMode == "PowerOn" {
				err = r.PowerOff()
			} else {
				err = r.PowerOn()
			}
			r.Refresh()
		default:
			log.Printf("key does not match any command")
		}
		LogIf(err)
	}
}
