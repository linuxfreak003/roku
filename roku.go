package roku

import (
	"fmt"
	"net"
	"net/http"
)

type Remote struct {
	IP string
}

func NewRemote(ip string) (*Remote, error) {
	if net.ParseIP(ip) == nil {
		return nil, fmt.Errorf("invalid ip: %s", ip)
	}
	return &Remote{IP: ip}, nil
}

func (r *Remote) keypress(cmd string) error {
	URL := fmt.Sprintf("http://%s:8060/keypress/%s", r.IP, cmd)
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
